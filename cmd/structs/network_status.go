package structs

type LANInfo struct {
	Type            string `json:"type"`
	Ifname          string `json:"ifname"`
	PhyIfname       string `json:"phy_ifname"`
	Link            bool   `json:"link"`
	Status          string `json:"status"`
	ProtocolStatus  string `json:"protocol_status"`
	MAC             string `json:"mac"`
	IP              string `json:"ip"`
	Mask            string `json:"mask"`
	MTU             int    `json:"mtu"`
	ConnectedPeriod int    `json:"connected_period"`
}

type PortLinkStatus struct {
	Type string `json:"type"`
	Port int    `json:"port"`
	Link string `json:"link"`
}
