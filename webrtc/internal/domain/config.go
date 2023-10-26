package domain

type ICEServer struct {
	URLs       []string `json:"urls"`
	UserName   string   `json:"username,omitempty"`
	Credential string   `json:"credential,omitempty"`
}

type RTCConfig struct {
	LifetimeDuration   string      `json:"lifetimeDuration,omitempty"`
	IceServers         []ICEServer `json:"iceServers"`
	BlockStatus        string      `json:"blockStatus,omitempty"`
	IceTransportPolicy string      `json:"iceTransportPolicy,omitempty"`
}
