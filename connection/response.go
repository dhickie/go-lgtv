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

// GetChannelListResponsePayload is the payload returned to "GetChannelList" requests
type GetChannelListResponsePayload struct {
	ReturnValue          bool                `json:"returnValue"`
	ValueList            string              `json:"valueList"`
	DataSource           int                 `json:"dataSource"`
	DataType             int                 `json:"dataType"`
	CableAnalogSkipped   bool                `json:"cableAnalogSkipped"`
	ScannedChannelCount  ScannedChannelCount `json:"scannedChannelCount"`
	DeviceSourceIndex    int                 `json:"deviceSourceIndex"`
	ChannelListCount     int                 `json:"channelListCount"`
	ChannelLogoServerURL string              `json:"channelLogoServerUrl"`
	IPChanInteractiveURL string              `json:"ipChanInteractiveUrl"`
	ChannelList          []Channel           `json:"channelList"`
}

// GetCurrentChannelResponsePayload is the payload returned to "GetCurrentChannel" requests
type GetCurrentChannelResponsePayload struct {
	ReturnValue      bool        `json:"returnValue"`
	ChannelID        string      `json:"channelId"`
	PhysicalNumber   int         `json:"physicalNumber"`
	IsScrambled      bool        `json:"isScrambled"`
	ChannelTypeName  string      `json:"channelTypeName"`
	IsLocked         bool        `json:"isLocked"`
	DualChannel      DualChannel `json:"dualChannel"`
	IsChannelChanged bool        `json:"isChannelChanged"`
	ChannelModeName  string      `json:"channelModeName"`
	ChannelNumber    string      `json:"channelNumber"`
	IsFineTuned      bool        `json:"isFineTuned"`
	ChannelTypeID    int         `json:"channelTypeId"`
	IsDescrambled    bool        `json:"isDescrambled"`
	IsSkipped        bool        `json:"isSkipped"`
	IsHevcChannel    bool        `json:"isHEVCChannel"`
	HybridTvType     string      `json:"hybridtvType"`
	IsInvisible      bool        `json:"isInvisible"`
	FavoriteGroup    string      `json:"favoriteGroup"`
	ChannelName      string      `json:"channelName"`
	ChannelModeID    int         `json:"channelModeId"`
	SignalChannelID  string      `json:"signalChannelId"`
}

// GetChannelProgramInfoResponsePayload is the payload returned to "GetChannelProgramInfo" requests
type GetChannelProgramInfoResponsePayload struct {
	ReturnValue bool      `json:"returnValue"`
	Channel     Channel   `json:"channel"`
	ProgramList []Program `json:"programList"`
}

// Program represents a program on TV
type Program struct {
	ChannelList     string   `json:"channelId"`
	Duration        int      `json:"duration"`
	EndTime         string   `json:"endTime"`
	LocalEndTime    string   `json:"localEndTime"`
	LocalStartTime  string   `json:"localStartTime"`
	Genre           string   `json:"genre"`
	ProgramID       string   `json:"programId"`
	ProgramName     string   `json:"programName"`
	Rating          []Rating `json:"rating"`
	SignalChannelID string   `json:"signalChannelId"`
	StartTime       string   `json:"startTime"`
	TableID         int      `json:"tableId"`
}

// Rating represents the rating of a TV program
type Rating struct {
	RatingString string `json:"ratingString"`
	RatingValue  int    `json:"ratingValue"`
	Region       int    `json:"region"`
	ID           string `json:"_id"`
}

// ScannedChannelCount represents count details of the returned channels
type ScannedChannelCount struct {
	TerrestrialAnalogCount  int `json:"terrestrialAnalogCount"`
	TerrestrialDigitalCount int `json:"terrestrialDigitalCount"`
	CableAnalogCount        int `json:"cableAnalogCount"`
	CableDigitalCount       int `json:"cableDigitalCount"`
	SatelliteDigitalCount   int `json:"satelliteDigitalCount"`
}

// Channel represents all the details of a TV channel
type Channel struct {
	ChannelID           string         `json:"channelId"`
	ProgramID           string         `json:"programId"`
	SignalChannelID     string         `json:"signalChannelId"`
	ChanCode            string         `json:"chanCode"`
	ChannelMode         string         `json:"channelMode"`
	ChannelModeID       int            `json:"channelModeId"`
	ChannelType         string         `json:"channelType"`
	ChannelTypeID       int            `json:"channelTypeId"`
	ChannelNumber       string         `json:"channelNumber"`
	MajorNumber         int            `json:"majorNumber"`
	MinorNumber         int            `json:"minorNumber"`
	ChannelName         string         `json:"channelName"`
	Skipped             bool           `json:"skipped"`
	Locked              bool           `json:"locked"`
	Descrambled         bool           `json:"descrambled"`
	Scrambled           bool           `json:"scrambled"`
	ServiceType         int            `json:"serviceType"`
	FavoriteGroup       []string       `json:"favoriteGroup"`
	ImgURL              string         `json:"imgUrl"`
	Display             int            `json:"display"`
	SatelliteName       string         `json:"satelliteName"`
	FineTuned           bool           `json:"fineTuned"`
	Frequency           int            `json:"Frequency"`
	ShortCut            int            `json:"shortCut"`
	Bandwidth           int            `json:"Bandwidth"`
	HDTV                bool           `json:"HDTV"`
	Invisible           bool           `json:"Invisible"`
	TV                  bool           `json:"TV"`
	DTV                 bool           `json:"DTV"`
	ATV                 bool           `json:"ATV"`
	Data                bool           `json:"Data"`
	Radio               bool           `json:"Radio"`
	Numeric             bool           `json:"Numeric"`
	PrimaryCh           bool           `json:"PrimaryCh"`
	SpecialService      bool           `json:"specialService"`
	CASystemIDList      CASystemIDList `json:"CASystemIDList"`
	CASystemIDListCount int            `json:"CASystemIDListCount"`
	GroupIDList         []int          `json:"groupIdList"`
	ChannelGenreCode    string         `json:"channelGenreCode"`
	FavoriteIdxA        int            `json:"favoriteIdxA"`
	FavoriteIdxB        int            `json:"favoriteIdxB"`
	FavoriteIDxC        int            `json:"favoriteIdxC"`
	FavoriteIDxD        int            `json:"favoriteIdxD"`
	FavoriteIDxE        int            `json:"favoriteIdxE"`
	FavoriteIDxF        int            `json:"favoriteIdxF"`
	FavoriteIDxG        int            `json:"favoriteIdxG"`
	FavoriteIDxH        int            `json:"favoriteIdxH"`
	ImgURL2             string         `json:"imgUrl2"`
	ChannelLogoSize     string         `json:"channelLogoSize"`
	IPChanServerURL     string         `json:"ipChanServerUrl"`
	PayChan             bool           `json:"payChan"`
	IPChannelCode       string         `json:"IPChannelCode"`
	IPCallNumber        string         `json:"ipCallNumber"`
	UTOFlag             bool           `json:"otuFlag"`
	SatelliteLcn        bool           `json:"satelliteLcn"`
	WaterMarkURL        string         `json:"waterMarkUrl"`
	ChannelNameSortKey  string         `json:"channelNameSortKey"`
	IPChanType          string         `json:"ipChanType"`
	AdultFlag           int            `json:"adultFlag"`
	IPChanCategory      string         `json:"ipChanCategory"`
	IPChanInteractive   bool           `json:"ipChanInteractive"`
	CallSign            string         `json:"callSign"`
	AdFlag              int            `json:"adFlag"`
	Configured          bool           `json:"configured"`
	LastUpdated         string         `json:"lastUpdated"`
	IPChanCpID          string         `json:"ipChanCpId"`
	IsFreeviewPlay      int            `json:"isFreeviewPlay"`
	PlayerService       string         `json:"playerService"`
}

// CASystemIDList represents... something
type CASystemIDList struct {
}

// DualChannel represents dual channel details for dual channel TV channels
type DualChannel struct {
	DualChannelID       string `json:"dualChannelId"`
	DualChannelTypeName string `json:"dualChannelTypeName"`
	DualChannelTypeID   string `json:"dualChannelTypeId"`
	DualChannelNumber   string `json:"dualChannelNumber"`
}
