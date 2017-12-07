package connection

// Represents a request made to the Web OS
type request struct {
	ID      int         `json:"id"`
	Type    string      `json:"type"`
	URI     string      `json:"uri"`
	Payload interface{} `json:"payload"`
}

// Represents an payload sent with a request to register with the Web OS
type registerReqPayload struct {
	PairingType string   `json:"pairingType"`
	Manifest    manifest `json:"manifest"`
	ClientKey   string   `json:"client-key"`
}

// Represents a payload which doesn't hold any data
type emptyPayload struct {
}

// Manifest represent an optional manifest sent in the payload
type manifest struct {
	Permissions []string `json:"permissions"`
}

// SetVolumePayload is the payload sent with a SetVolume request
type SetVolumePayload struct {
	Volume int `json:"volume"`
}

// SetMutePayload is the payload sent with a SetMute request
type SetMutePayload struct {
	Mute bool `json:"mute"`
}

// SetChannelPayload is the payload sent with an OpenChannel request
type SetChannelPayload struct {
	ChannelNumber string `json:"channelNumber"`
}

// SwitchInputPayload is the payload sent with a SwitchInput request
type SwitchInputPayload struct {
	InputID string `json:"inputId"`
}
