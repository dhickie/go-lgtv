package control

import (
	"errors"
	"net"
	"strconv"

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
	return tv.doRequest(uriVolumeUp, nil, nil)
}

// VolumeDown decreases the volume by 1
func (tv *LgTv) VolumeDown() error {
	return tv.doRequest(uriVolumeDown, nil, nil)
}

// SetVolume sets the volume to the specified value
func (tv *LgTv) SetVolume(value int) error {
	payload := connection.SetVolumePayload{
		Volume: value,
	}
	return tv.doRequest(uriSetVolume, payload, nil)
}

// GetVolume returns the current volume of the TV
func (tv *LgTv) GetVolume() (int, error) {
	respPayload := connection.GetVolumeResponsePayload{}
	err := tv.doRequest(uriGetVolume, nil, &respPayload)
	if err != nil {
		return 0, err
	}

	return respPayload.Volume, nil
}

// SetMute sets the mute status of the TV
func (tv *LgTv) SetMute(isMute bool) error {
	payload := connection.SetMutePayload{
		Mute: isMute,
	}
	return tv.doRequest(uriSetMute, payload, nil)
}

// GetMute gets the mute status of the TV
func (tv *LgTv) GetMute() (bool, error) {
	respPayload := connection.GetMuteResponsePayload{}
	err := tv.doRequest(uriGetMute, nil, &respPayload)
	if err != nil {
		return false, err
	}

	return respPayload.Mute, nil
}

// Play plays the current media
func (tv *LgTv) Play() error {
	return tv.doRequest(uriPlay, nil, nil)
}

// Pause pauses the current media
func (tv *LgTv) Pause() error {
	return tv.doRequest(uriPause, nil, nil)
}

// Stop stops the current media
func (tv *LgTv) Stop() error {
	return tv.doRequest(uriStop, nil, nil)
}

// Rewind rewinds the current media
func (tv *LgTv) Rewind() error {
	return tv.doRequest(uriRewind, nil, nil)
}

// FastForward fast forwards the current media
func (tv *LgTv) FastForward() error {
	return tv.doRequest(uriFastForward, nil, nil)
}

// ChannelUp changes the current channel up by 1
func (tv *LgTv) ChannelUp() error {
	return tv.doRequest(uriChannelUp, nil, nil)
}

// ChannelDown changes the current channel down by 1
func (tv *LgTv) ChannelDown() error {
	return tv.doRequest(uriChannelDown, nil, nil)
}

// SetChannel sets the current viewed channel to the specified number
func (tv *LgTv) SetChannel(channelNumber int) error {
	payload := connection.SetChannelPayload{
		ChannelNumber: strconv.Itoa(channelNumber),
	}
	return tv.doRequest(uriSetChannel, payload, nil)
}

func (tv *LgTv) GetChannelList() error {
	return tv.doRequest(uriGetChannelList, nil, nil)
}

// SwitchInput switches the input of the TV to the one with the specified input ID
func (tv *LgTv) SwitchInput(inputID string) error {
	payload := connection.SwitchInputPayload{
		InputID: inputID,
	}
	return tv.doRequest(uriSwitchInput, payload, nil)
}

// TurnOff turns the tv off
func (tv *LgTv) TurnOff() error {
	return tv.doRequest(uriTurnOff, nil, nil)
}

func (tv *LgTv) doRequest(uri string, reqPayload interface{}, respPayload interface{}) error {
	if tv.isConnected {
		return tv.conn.Request(uri, reqPayload, respPayload)
	}

	return ErrNotConnected
}
