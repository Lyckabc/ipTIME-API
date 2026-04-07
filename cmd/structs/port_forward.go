package structs

type PortForward struct {
	Name     string `json:"name"`
	Active   bool   `json:"active"`
	Protocol string `json:"protocol"`
	Target   string `json:"target"`
	Priority int    `json:"priority"`
	Fixed    bool   `json:"fixed"`
	Src      PortRange `json:"src"`
	Dst      PortRange `json:"dst"`
}

type PortRange struct {
	Start int `json:"start"`
	End   int `json:"end,omitempty"`
}
