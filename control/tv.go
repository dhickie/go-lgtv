package control

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

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
	var respPayload connection.GetVolumeResponsePayload
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
	var respPayload connection.GetMuteResponsePayload
	err := tv.doRequest(uriGetMute, nil, &respPayload)
	return respPayload.Mute, err
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

// GetChannelList returns a slice of available TV channels
func (tv *LgTv) GetChannelList() ([]Channel, error) {
	var respPayload connection.GetChannelListResponsePayload
	err := tv.doRequest(uriGetChannelList, nil, &respPayload)
	if err != nil {
		return nil, err
	}

	channels := make([]Channel, len(respPayload.ChannelList))
	for i, v := range respPayload.ChannelList {
		channelNum, err := strconv.Atoi(v.ChannelNumber)
		if err != nil {
			return nil, err
		}

		channels[i] = Channel{
			ChannelName:   v.ChannelName,
			ChannelNumber: channelNum,
			IsHdtv:        v.HDTV,
			IsScrambled:   v.Scrambled,
			tv:            tv,
		}
	}

	return channels, nil
}

// GetCurrentChannel returns the channel the TV is currently set to
func (tv *LgTv) GetCurrentChannel() (Channel, error) {
	var respPayload connection.GetCurrentChannelResponsePayload
	err := tv.doRequest(uriGetCurrentChannel, nil, &respPayload)
	if err != nil {
		return Channel{}, err
	}

	channelNum, err := strconv.Atoi(respPayload.ChannelNumber)
	if err != nil {
		return Channel{}, err
	}

	return Channel{
		ChannelName:   respPayload.ChannelName,
		ChannelNumber: channelNum,
		IsHdtv:        false,
		IsScrambled:   respPayload.IsScrambled,
		tv:            tv,
	}, nil
}

// GetChannelProgramList gets the list of programs broadcast on the current channel
func (tv *LgTv) GetChannelProgramList() (ChannelProgramList, error) {
	var respPayload connection.GetChannelProgramInfoResponsePayload
	err := tv.doRequest(uriGetChannelProgramInfo, nil, &respPayload)
	if err != nil {
		return ChannelProgramList{}, err
	}

	channelNum, err := strconv.Atoi(respPayload.Channel.ChannelNumber)
	if err != nil {
		return ChannelProgramList{}, err
	}

	programList := ChannelProgramList{
		Channel: Channel{
			ChannelName:   respPayload.Channel.ChannelName,
			ChannelNumber: channelNum,
			IsHdtv:        respPayload.Channel.HDTV,
			IsScrambled:   respPayload.Channel.Scrambled,
			tv:            tv,
		},
		Programs: make([]Program, len(respPayload.ProgramList)),
	}

	for i, v := range respPayload.ProgramList {
		duration, err := time.ParseDuration(fmt.Sprintf("%vs", v.Duration))
		if err != nil {
			return ChannelProgramList{}, err
		}

		stringTimes := []string{v.StartTime, v.EndTime}
		times := make([]time.Time, 2)
		for i, v := range stringTimes {
			t, err := parseTime(v)
			if err != nil {
				return ChannelProgramList{}, err
			}

			times[i] = t
		}

		programList.Programs[i] = Program{
			Name:      v.ProgramName,
			Genre:     v.Genre,
			StartTime: times[0],
			EndTime:   times[1],
			Duration:  duration,
		}
	}

	return programList, nil
}

// SwitchInput switches the input of the TV to the one with the specified input ID
func (tv *LgTv) SwitchInput(inputID string) error {
	payload := connection.SwitchInputPayload{
		InputID: inputID,
	}
	return tv.doRequest(uriSwitchInput, payload, nil)
}

// ListExternalInputs lists the external input devices for the TV
func (tv *LgTv) ListExternalInputs() ([]Input, error) {
	var respPayload connection.GetExternalInputListResponsePayload
	err := tv.doRequest(uriGetExternalInputList, nil, &respPayload)
	if err != nil {
		return nil, err
	}

	inputs := make([]Input, len(respPayload.Devices))
	for i, v := range respPayload.Devices {
		inputs[i] = Input{
			ID:    v.ID,
			Label: v.Label,
			tv:    tv,
		}
	}

	return inputs, nil
}

// ListInstalledApps lists the apps currently installed on the TV
func (tv *LgTv) ListInstalledApps() ([]App, error) {
	var respPayload connection.GetInstalledAppsResponsePayload
	err := tv.doRequest(uriListApps, nil, &respPayload)
	if err != nil {
		return nil, err
	}

	apps := make([]App, len(respPayload.Apps))
	for i, v := range respPayload.Apps {
		apps[i] = App{
			Name: v.Title,
			ID:   v.ID,
			tv:   tv,
		}
	}

	return apps, nil
}

// LaunchApp launches the app with the provided ID. If successfully launched,
// it returns the ID of the new session
func (tv *LgTv) LaunchApp(appID string) (string, error) {
	payload := connection.LaunchAppPayload{
		ID: appID,
	}
	var respPayload connection.LaunchAppResponsePayload
	err := tv.doRequest(uriLaunchApp, payload, &respPayload)
	if err != nil {
		return "", err
	}

	return respPayload.SessionID, nil
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

func parseTime(strTime string) (time.Time, error) {
	loc, err := time.LoadLocation("UTC")

	if err != nil {
		return time.Time{}, err
	}

	elements := strings.Split(strTime, ",")
	intElements := make([]int, 6)

	for i, v := range elements {
		num, err := strconv.Atoi(v)
		if err != nil {
			return time.Time{}, err
		}

		intElements[i] = num
	}

	return time.Date(
		intElements[0],
		time.Month(intElements[1]),
		intElements[2],
		intElements[3],
		intElements[4],
		intElements[5],
		0,
		loc), nil
}
