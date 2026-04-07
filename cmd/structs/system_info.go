package structs

type SystemInfo struct {
	Version         string        `json:"version"`
	Uptime          int           `json:"uptime"`
	RemotePort      int           `json:"remote_port"`
	ConnectedPeriod int           `json:"connected_period"`
	WanLink         string        `json:"wan_link"`
	DHCP            DHCPInfo      `json:"dhcpd"`
	Wireless        []WirelessBand `json:"wireless"`
}

type DHCPInfo struct {
	Enable  bool   `json:"enable"`
	StartIP string `json:"startip"`
	EndIP   string `json:"endip"`
}

type WirelessBand struct {
	Band string       `json:"band"`
	BSS  []WirelessBS `json:"bss"`
}

type WirelessBS struct {
	Enable   bool   `json:"enable"`
	SSID     string `json:"ssid"`
	Password string `json:"password"`
}
