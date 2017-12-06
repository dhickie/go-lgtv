package connection

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"encoding/json"

	"github.com/gorilla/websocket"
)

const (
	wsPort                 = 3000
	registerTimeoutSeconds = 60
)

var (
	errUnknownResponseType = errors.New("Unknown response type recieved")

	// ErrRegisterTimeout is returned if a registered response is not recieved from the TV
	// before a timeout
	ErrRegisterTimeout = errors.New("Timeout waiting for registered response")
)

// Connection represents a web socket connection to the TV
type Connection struct {
	conn          *websocket.Conn
	connOpen      *bool
	idLock        sync.Mutex
	lastRequestID int
	respChans     map[int]chan response
}

// NewConnection creates a new web socket connection to the TV at the given IP address
func NewConnection(ip net.IP) (*Connection, error) {
	url := fmt.Sprintf("ws://%v:%v", ip, wsPort)

	// Dial the websocket connection
	c, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return nil, err
	}

	// Set the routine going to get responses
	connOpen := true
	connection := &Connection{c,
		&connOpen,
		sync.Mutex{},
		0,
		make(map[int]chan response),
	}

	go connection.respWorker()

	return connection, nil
}

// Register registers with the TV using the provided client key.
// If no client key is provided, the TV will generate a new one
func (c *Connection) Register(clientKey string) (string, error) {
	// Create the request
	requestID := c.getID()
	request := request{
		ID:   requestID,
		Type: reqTypeRegister,
		Payload: registerReqPayload{
			PairingType: pairTypePrompt,
			Manifest: manifest{
				Permissions: getPermissions(),
			},
			ClientKey: clientKey,
		},
	}

	// Marshal the request
	message, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	// Create the channel to recieve the response
	respChan := c.addRespChannel(requestID)
	defer c.removeRespChan(requestID)

	// Send the message to the websocket
	err = c.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return "", err
	}

	// For register requests, multiple responses are recieved, wait until we get either
	// an error response, or a registered response (or we timeout)
	ticker := time.NewTicker(registerTimeoutSeconds * time.Second)
	for {
		select {
		case <-ticker.C:
			return "", ErrRegisterTimeout
		case resp := <-respChan:
			if resp.Type == respTypeRegistered {
				return resp.Payload.(registerRespPayload).ClientKey, nil
			} else if resp.Type == respTypeError {
				return "", errors.New(resp.Error)
			}
		}
	}
}

func (c *Connection) addRespChannel(reqID int) chan response {
	respChan := make(chan response)
	c.respChans[reqID] = respChan

	return respChan
}

func (c *Connection) removeRespChan(reqID int) {
	if val, ok := c.respChans[reqID]; ok {
		close(val)
		delete(c.respChans, reqID)
	}
}

func (c *Connection) respWorker() {
	for *c.connOpen {
		// Read a message from the connection
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			continue
		}

		// Unmarshal the response
		resp, err := unmarshalResponse(message)
		if err != nil {
			continue
		}

		// Send the response to the appropriate channel
		if val, ok := c.respChans[resp.ID]; ok {
			val <- resp
		}
	}
}

func unmarshalResponse(message []byte) (response, error) {
	// Unmarshal the common properties first
	var resp response
	err := json.Unmarshal(message, &resp)
	if err != nil {
		return resp, err
	}

	// Use json.RawMessage to extract the JSON for the payload
	var propertyMap map[string]*json.RawMessage
	err = json.Unmarshal(message, &propertyMap)
	if err != nil {
		return resp, err
	}

	// Unmarshal the payload based on the type of the response
	if resp.Type == respTypeResponse {
		var payload responsePayload
		err = json.Unmarshal(*propertyMap["payload"], &payload)
		resp.Payload = payload
		return resp, err
	} else if resp.Type == respTypeRegistered {
		var payload registerRespPayload
		err = json.Unmarshal(*propertyMap["payload"], &payload)
		resp.Payload = payload
		return resp, err
	} else if resp.Type == respTypeError {
		var payload emptyPayload
		err = json.Unmarshal(*propertyMap["payload"], &payload)
		resp.Payload = payload
		return resp, err
	}

	return response{}, errUnknownResponseType
}

func (c *Connection) getID() int {
	c.idLock.Lock()
	defer c.idLock.Unlock()

	c.lastRequestID++
	return c.lastRequestID
}

func getPermissions() []string {
	return []string{
		permissionLaunch,
		permissionControlAudio,
		permissionControlPower,
	}
}
