package connection

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"encoding/json"

	"github.com/gorilla/websocket"
)

const (
	wsPort                 = 3000
	registerTimeoutSeconds = 60
	requestTimeoutSeconds  = 10
)

var (
	errUnknownResponseType = errors.New("Unknown response type recieved")

	// ErrRegisterTimeout is returned if a registered response is not recieved from the TV
	// before a timeout
	ErrRegisterTimeout = errors.New("Timeout waiting for registered response")
	// ErrFailResponse is returned when the TV returns a fail response to a request
	ErrFailResponse = errors.New("TV returned fail response to request")
	// ErrRequestTimeout is returned when no response is recieved to a request before a timeout.
	// Note that it can also be returned if we do get a response, but an error occurs processing it
	// on our end.
	ErrRequestTimeout = errors.New("Timeout waiting for response to request")
	// ErrConnectionTimeout is returned when we fail to open the websocket connection before the timeout
	ErrConnectionTimeout = errors.New("Failed to connect to TV's websocket connection before timeout")
)

// Connection represents a web socket connection to the TV
type Connection struct {
	url           string
	conn          *websocket.Conn
	connOpen      bool
	idLock        *sync.Mutex
	lastRequestID int
	respChans     map[int]chan response
	respPayloads  map[int]interface{}
}

// NewConnection creates a new web socket connection to the TV at the given IP address. The timeout is in milliseconds.
func NewConnection(ip net.IP) *Connection {
	url := fmt.Sprintf("ws://%v:%v", ip, wsPort)

	return &Connection{
		url:           url,
		idLock:        new(sync.Mutex),
		connOpen:      false,
		lastRequestID: 0,
		respChans:     make(map[int]chan response),
		respPayloads:  make(map[int]interface{}),
	}
}

// Open opens the connection to the TV, and sets the response worker going
func (c *Connection) Open(timeout int) error {
	// Dial the websocket connection, with a timeout
	conChan := make(chan *websocket.Conn)
	errChan := make(chan error)
	go func(url string) {
		c, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			errChan <- err
			return
		}

		conChan <- c
	}(c.url)

	ticker := time.NewTicker(time.Duration(timeout) * time.Millisecond)
	select {
	case <-ticker.C:
		return ErrConnectionTimeout
	case err := <-errChan:
		return err
	case conn := <-conChan:
		c.conn = conn
		c.conn.SetCloseHandler(c.closeHandler)
		c.connOpen = true

		// Set the routine going to get responses
		go c.respWorker()

		return nil
	}
}

// IsOpen returns whether or not the connection to the TV is currently open
func (c *Connection) IsOpen() bool {
	return c.connOpen
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
				return resp.Payload.(*registerRespPayload).ClientKey, nil
			} else if resp.Type == respTypeError {
				return "", errors.New(resp.Error)
			}
		}
	}
}

// Request makes a request to the TV to perform an action
func (c *Connection) Request(uri string, reqPayload interface{}, respPayload interface{}) error {
	// Create the request
	requestID := c.getID()
	request := request{
		ID:   requestID,
		Type: reqTypeRequest,
		URI:  uri,
	}

	if reqPayload != nil {
		request.Payload = reqPayload
	}

	// Marshal the request
	message, err := json.Marshal(request)
	if err != nil {
		return err
	}

	// Create the channel to recieve the response
	respChan := c.addRespChannel(requestID)
	defer c.removeRespChan(requestID)

	// Add the response payload to the map of response payloads if one has been provided
	if respPayload != nil {
		c.respPayloads[requestID] = respPayload
		defer delete(c.respPayloads, requestID)
	}

	// Send the message to the websocket
	err = c.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		return err
	}

	// Wait for the response (or timeout)
	ticker := time.NewTicker(requestTimeoutSeconds * time.Second)
	defer ticker.Stop()

	select {
	case <-ticker.C:
		return ErrRequestTimeout
	case resp := <-respChan:
		if resp.Error != "" {
			return errors.New(resp.Error)
		}
	}

	return nil
}

// Close closes the connection to the TV
func (c *Connection) Close() error {
	c.connOpen = false
	return c.conn.Close()
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
	for c.connOpen {
		// Read a message from the connection
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			// Error messages are fatal, set the connection to closed and kill the worker
			c.connOpen = false
			return
		}

		// Get the ID of the response
		respID, err := getRespID(message)
		if err != nil {
			continue
		}

		// Unmarshal the response
		var payload interface{}
		if val, ok := c.respPayloads[respID]; ok {
			payload = val
		} else {
			payload = &responsePayload{}
		}
		resp, err := unmarshalResponse(message, payload)
		if err != nil {
			continue
		}

		// Send the response to the appropriate channel
		if val, ok := c.respChans[resp.ID]; ok {
			val <- resp
		}
	}
}

// closeHandler is called if the websocket connection receives a close message
func (c *Connection) closeHandler(code int, text string) error {
	// Set the connection as closed
	c.connOpen = false
	return nil
}

func (c *Connection) getID() int {
	c.idLock.Lock()
	defer c.idLock.Unlock()

	c.lastRequestID++
	return c.lastRequestID
}

func unmarshalResponse(message []byte, respPayload interface{}) (response, error) {
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
	var payload interface{}
	if resp.Type == respTypeResponse {
		// If this is a "response" type, then use the provided payload pointer
		payload = respPayload
	} else if resp.Type == respTypeRegistered {
		payload = &registerRespPayload{}
	} else if resp.Type == respTypeError {
		payload = &emptyPayload{}
	} else {
		return response{}, errUnknownResponseType
	}

	err = json.Unmarshal(*propertyMap["payload"], payload)
	resp.Payload = payload
	return resp, err
}

func getRespID(message []byte) (int, error) {
	var propertyMap map[string]*json.RawMessage
	err := json.Unmarshal(message, &propertyMap)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(string(*propertyMap["id"]))
}

func getPermissions() []string {
	return []string{
		permissionLaunch,
		permissionControlAudio,
		permissionControlPower,
		permissionControlPlayback,
		permissionControlInputTv,
		permissionReadChannelList,
		permissionReadCurrentChannel,
		permissionReadRunningApps,
		permissionReadInstalledApps,
		permissionReadInputList,
	}
}
