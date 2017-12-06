package connection

const (
	// Request types
	reqTypeRegister = "register"
	reqTypeRequest  = "request"

	// Response types
	respTypeRegistered = "registered"
	respTypeResponse   = "response"
	respTypeError      = "error"

	// Pairing type
	pairTypePrompt = "PROMPT"

	// Permissions
	permissionLaunch       = "LAUNCH"
	permissionControlAudio = "CONTROL_AUDIO"
	permissionControlPower = "CONTROL_POWER"
)

// Represents a request made to the Web OS
type request struct {
	ID      int         `json:"id"`
	Type    string      `json:"type"`
	URI     string      `json:"uri"`
	Payload interface{} `json:"payload"`
}

// Represents a response from the Web OS made to a request
type response struct {
	ID      int         `json:"id"`
	Type    string      `json:"type"`
	Error   string      `json:"error"`
	Payload interface{} `json:"payload"`
}

// Temporary structure used to determine the type of response before unmarshalling
// the correct model
type tempResponse struct {
	ID    int    `json:"id"`
	Type  string `json:"type"`
	Error string `json:"error"`
}

// Represents an payload sent with a request to register with the Web OS
type registerReqPayload struct {
	PairingType string   `json:"pairingType"`
	Manifest    manifest `json:"manifest"`
	ClientKey   string   `json:"client-key"`
}

// Represents a "registered" response payload to a request to register
type registerRespPayload struct {
	ClientKey string `json:"client-key"`
}

// Represents a response payload to a request to register
type responsePayload struct {
	ReturnValue bool `json:"returnValue"`
}

// Represents a payload which doesn't hold any data
type emptyPayload struct {
}

// Manifest represent an optional manifest sent in the payload
type manifest struct {
	Permissions []string `json:"permissions"`
}
