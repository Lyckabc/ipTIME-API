package routers

import (
	"encoding/json"
	"net/http"

	"github.com/kongwoojin/ipTIME-API/cmd/structs"
)

type stationInfo struct {
	IP   string `json:"ip"`
	Type string `json:"type"`
	Name string `json:"name"`
}

type wiredConn struct {
	Port      string `json:"port"`
	Link      int    `json:"link"`
	DownBytes int    `json:"down_bytes"`
	UpBytes   int    `json:"up_bytes"`
}

type wirelessConn struct {
	BSS       string `json:"bss"`
	Duration  int    `json:"duration"`
	RSSI      int    `json:"rssi"`
	DownSpeed int    `json:"down_speed"`
	UpSpeed   int    `json:"up_speed"`
	DownBytes int    `json:"down_bytes"`
	UpBytes   int    `json:"up_bytes"`
}

type stationConnection struct {
	Type     string        `json:"type"`
	Wired    *wiredConn    `json:"wired"`
	Wireless *wirelessConn `json:"wireless"`
}

type station struct {
	MAC        string            `json:"mac"`
	Info       stationInfo       `json:"info"`
	Connection stationConnection `json:"connection"`
}

func GetConnectedClientList(client *http.Client, router *structs.Router) []structs.Client {
	raw, err := serviceCall(client, router, "network/interface/lan/stations", nil)
	if err != nil {
		return nil
	}

	var stations []station
	if err := json.Unmarshal(raw, &stations); err != nil {
		return nil
	}

	var clients []structs.Client
	for _, s := range stations {
		if s.Info.IP == "" {
			continue
		}
		c := structs.Client{
			MAC:    s.MAC,
			IP:     s.Info.IP,
			Name:   s.Info.Name,
			IPType: s.Info.Type,
		}
		switch s.Connection.Type {
		case "wired":
			c.ConnectionType = structs.Wired
			if s.Connection.Wired != nil {
				c.Port = s.Connection.Wired.Port
			}
		case "wireless":
			c.ConnectionType = structs.Wireless
			if s.Connection.Wireless != nil {
				c.RSSI = s.Connection.Wireless.RSSI
				c.DownSpeed = s.Connection.Wireless.DownSpeed
				c.UpSpeed = s.Connection.Wireless.UpSpeed
			}
		default:
			c.ConnectionType = structs.UnknownConnectionType
		}
		clients = append(clients, c)
	}
	return clients
}
