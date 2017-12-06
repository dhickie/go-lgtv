package control

import (
	"errors"
	"net"

	"github.com/dhickie/go-lgtv/connection"
)

// ErrNotConnected is returned if an request is attempted to a TV which is not connected to the client
var ErrNotConnected = errors.New("Client is not connected to TV")

// LgTv represents the TV being controlled
type LgTv struct {
	ip          net.IP
	conn        *connection.Connection
	isConnected bool
}

// NewTV returns a new LgTv object with the specified IP and pin
func NewTV(ip net.IP) LgTv {
	return LgTv{ip, nil, false}
}

// Connect connects to the tv using the provided client key. If an empty client key
// is provided, a new one will be provisioned
func (tv *LgTv) Connect(clientKey string) (string, error) {
	conn, err := connection.NewConnection(tv.ip)
	if err != nil {
		return "", err
	}

	tv.conn = conn
	clientKey, err = tv.conn.Register(clientKey)
	if err == nil {
		tv.isConnected = true
	}

	return clientKey, err
}

// VolumeUp increases the volume by 1
func (tv *LgTv) VolumeUp() error {
	return tv.doRequest(uriVolumeUp)
}

// VolumeDown decreases the volume by 1
func (tv *LgTv) VolumeDown() error {
	return tv.doRequest(uriVolumeDown)
}

// Play plays the current media
func (tv *LgTv) Play() error {
	return tv.doRequest(uriPlay)
}

// Pause pauses the current media
func (tv *LgTv) Pause() error {
	return tv.doRequest(uriPause)
}

// Stop stops the current media
func (tv *LgTv) Stop() error {
	return tv.doRequest(uriStop)
}

// Rewind rewinds the current media
func (tv *LgTv) Rewind() error {
	return tv.doRequest(uriRewind)
}

// FastForward fast forwards the current media
func (tv *LgTv) FastForward() error {
	return tv.doRequest(uriFastForward)
}

// ChannelUp changes the current channel up by 1
func (tv *LgTv) ChannelUp() error {
	return tv.doRequest(uriChannelUp)
}

// ChannelDown changes the current channel down by 1
func (tv *LgTv) ChannelDown() error {
	return tv.doRequest(uriChannelDown)
}

// TurnOff turns the tv off
func (tv *LgTv) TurnOff() error {
	return tv.doRequest(uriTurnOff)
}

func (tv *LgTv) doRequest(uri string) error {
	if tv.isConnected {
		return tv.conn.Request(uri)
	}

	return ErrNotConnected
}
