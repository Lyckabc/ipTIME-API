package structs

type ConnectionType int

const (
	UnknownConnectionType ConnectionType = iota
	Wired
	Wireless
)

type Client struct {
	MAC            string         `json:"mac"`
	IP             string         `json:"ip"`
	Name           string         `json:"name"`
	IPType         string         `json:"ipType"`   // "dhcp" or "static"
	ConnectionType ConnectionType `json:"connectionType"`
	Port           string         `json:"port,omitempty"`    // wired: "lan1"..
	RSSI           int            `json:"rssi,omitempty"`    // wireless: signal strength
	DownSpeed      int            `json:"downSpeed,omitempty"`
	UpSpeed        int            `json:"upSpeed,omitempty"`
}
