package main

import (
	"fmt"
	"log"

	"github.com/Lyckabc/ipTIME-API/cmd/routers"
	"github.com/Lyckabc/ipTIME-API/cmd/structs"
)

func main() {
	router := RouterInfo()
	client := CreateClient()

	if !routers.Login(client, router) {
		log.Fatal("Login failed")
	}
	fmt.Println("Login successful")

	// Router status
	status, err := routers.RouterStatus(client, router)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Firmware: %s | Uptime: %ds | WAN: %s\n",
		status.SystemInfo.Version, status.SystemInfo.Uptime, status.SystemInfo.WanLink)
	fmt.Printf("LAN: %s (%s)\n", status.LANInfo.IP, status.LANInfo.MAC)

	// Connected clients
	fmt.Println("\n--- Connected Clients ---")
	for i, c := range routers.GetConnectedClientList(client, router) {
		connType := "wired"
		if c.ConnectionType == structs.Wireless {
			connType = fmt.Sprintf("wifi rssi=%d", c.RSSI)
		}
		fmt.Printf("[%d] %s %-16s %-20s (%s/%s)\n", i, c.MAC, c.IP, c.Name, connType, c.IPType)
	}

	// Port forwards
	fmt.Println("\n--- Port Forward Rules ---")
	for i, pf := range routers.GetPortForwardList(client, router) {
		active := "on"
		if !pf.Active {
			active = "off"
		}
		fmt.Printf("[%d] %-30s %s:%d -> %s:%d [%s]\n",
			i, pf.Name, pf.Protocol, pf.Src.Start, pf.Target, pf.Dst.Start, active)
	}

	// Add and remove a test port forward
	_, addErr := routers.AddPortForward(client, router, &structs.PortForward{
		Name:     "api_test",
		Active:   false,
		Protocol: "TCP",
		Target:   "192.168.0.253",
		Src:      structs.PortRange{Start: 19999},
		Dst:      structs.PortRange{Start: 19999},
	})
	if addErr != nil {
		log.Printf("AddPortForward: %v", addErr)
	} else {
		fmt.Println("\nPort forward 'api_test' added")
		routers.RemovePortForward(client, router, &structs.PortForward{Name: "api_test"})
		fmt.Println("Port forward 'api_test' removed")
	}

	// WOL list
	fmt.Println("\n--- Wake on LAN ---")
	for i, w := range routers.GetWOLList(client, router) {
		fmt.Printf("[%d] %s (%s)\n", i, w.MAC, w.Name)
	}

	// MAC auth
	fmt.Println("\n--- Wireless MAC Auth ---")
	for _, m := range routers.GetMacAuthList(client, router) {
		fmt.Printf("[%s] policy=%s macs=%d\n", m.BSSTag, m.Policy, len(m.MacList))
	}
}
