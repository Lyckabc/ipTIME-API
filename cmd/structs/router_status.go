package structs

type RouterStatus struct {
	SystemInfo SystemInfo
	LANInfo    LANInfo
	PortLinks  []PortLinkStatus
}
