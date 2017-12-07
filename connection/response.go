package connection

// Represents a response from the Web OS made to a request
type response struct {
	ID      int         `json:"id"`
	Type    string      `json:"type"`
	Error   string      `json:"error"`
	Payload interface{} `json:"payload"`
}

// Represents a "registered" response payload to a request to register
type registerRespPayload struct {
	ClientKey string `json:"client-key"`
}

// Represents a response payload to a request to register
type responsePayload struct {
	ReturnValue bool `json:"returnValue"`
}

// GetVolumeResponsePayload is the payload returned to "GetVolume" requests
type GetVolumeResponsePayload struct {
	ReturnValue bool   `json:"returnValue"`
	Scenario    string `json:"scenario"`
	Volume      int    `json:"volume"`
	Muted       bool   `json:"muted"`
	VolumeMax   int    `json:"volumeMax"`
}

// GetMuteResponsePayload is the payload returned to "GetMute" requests
type GetMuteResponsePayload struct {
	ReturnValue bool `json:"returnValue"`
	Mute        bool `json:"mute"`
}
