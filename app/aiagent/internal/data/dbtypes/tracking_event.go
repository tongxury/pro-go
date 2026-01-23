package dbtypes

type TrackingEvent struct {
	Ip          string
	CountryCode string
	AppId       string `json:"appId"`
	Build       struct {
		PackageName string `json:"packageName"`
		Version     string `json:"version"`
	} `json:"build"`
	ClientId string `json:"clientId"`
	Device   Device `json:"device"`
	DeviceTs int64  `json:"deviceTs"`
	EventId  string `json:"eventId"`
	Payload  struct {
		Event string `json:"event"`
	} `json:"payload"`
	Platform  string `json:"platform"`
	SendAt    int64  `json:"sendAt"`
	SessionId string `json:"sessionId"`
	Version   string `json:"version"`
}

type Device struct {
	Brand           string `json:"brand"`
	Carrier         string `json:"carrier"`
	DisplayLanguage string `json:"displayLanguage"`
	Idfa            string `json:"idfa"`
	Language        string `json:"language"`
	Manufacturer    string `json:"manufacturer"`
	Memory          struct {
		Available   int64 `json:"available"`
		IsLowMemory bool  `json:"isLowMemory"`
		Total       int64 `json:"total"`
	} `json:"memory"`
	Model       string `json:"model"`
	Network     string `json:"network"`
	Orientation int64  `json:"orientation"`
	OsVersion   string `json:"osVersion"`
	Screen      struct {
		Density    float64 `json:"density"`
		Resolution struct {
			Height int64 `json:"height"`
			Width  int64 `json:"width"`
		} `json:"resolution"`
	} `json:"screen"`
	SdkVersion int64 `json:"sdkVersion"`
	Storage    struct {
		Available int64 `json:"available"`
		Free      int64 `json:"free"`
		Total     int64 `json:"total"`
	} `json:"storage"`
	Timezone string `json:"timezone"`
}
