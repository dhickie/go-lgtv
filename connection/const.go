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
	permissionLaunch          = "LAUNCH"
	permissionControlAudio    = "CONTROL_AUDIO"
	permissionControlPower    = "CONTROL_POWER"
	permissionControlInputTv  = "CONTROL_INPUT_TV"
	permissionReadChannelList = "READ_TV_CHANNEL_LIST"
)
