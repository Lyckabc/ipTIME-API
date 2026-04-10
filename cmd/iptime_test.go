package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/Lyckabc/ipTIME-API/cmd/routers"
	"github.com/Lyckabc/ipTIME-API/cmd/structs"
)

const testPFName = "api_test_rule"
const testWOLMac = "AA:BB:CC:DD:EE:01"
const testWOLName = "api_test_device"

func setup(t *testing.T) (interface{ Do(interface{}) (interface{}, error) }, *structs.Router) {
	t.Helper()
	router := RouterInfo()
	client := CreateClient()
	if !routers.Login(client, router) {
		t.Fatal("Login failed — check .env credentials")
	}
	return nil, router
}

// TestLogin verifies credentials in .env work.
func TestLogin(t *testing.T) {
	router := RouterInfo()
	client := CreateClient()
	if !routers.Login(client, router) {
		t.Fatal("FAIL: login failed")
	}
	t.Log("PASS: login succeeded")
}

// TestRouterStatus checks that system/info, lan/info, port/link/status return valid data.
func TestRouterStatus(t *testing.T) {
	router := RouterInfo()
	client := CreateClient()
	if !routers.Login(client, router) {
		t.Fatal("login failed")
	}

	status, err := routers.RouterStatus(client, router)
	if err != nil {
		t.Fatalf("FAIL RouterStatus: %v", err)
	}

	if status.SystemInfo.Version == "" {
		t.Error("FAIL: firmware version is empty")
	} else {
		t.Logf("PASS: firmware=%s uptime=%ds wan=%s",
			status.SystemInfo.Version, status.SystemInfo.Uptime, status.SystemInfo.WanLink)
	}

	if status.LANInfo.IP == "" {
		t.Error("FAIL: LAN IP is empty")
	} else {
		t.Logf("PASS: LAN ip=%s mac=%s", status.LANInfo.IP, status.LANInfo.MAC)
	}

	if len(status.PortLinks) == 0 {
		t.Error("FAIL: port link list is empty")
	} else {
		t.Logf("PASS: port links=%d", len(status.PortLinks))
	}
}

// TestConnectedClients checks that at least one client is visible.
func TestConnectedClients(t *testing.T) {
	router := RouterInfo()
	client := CreateClient()
	if !routers.Login(client, router) {
		t.Fatal("login failed")
	}

	clients := routers.GetConnectedClientList(client, router)
	if len(clients) == 0 {
		t.Error("FAIL: no connected clients returned")
		return
	}
	t.Logf("PASS: %d client(s) connected", len(clients))
	for i, c := range clients {
		connType := "wired"
		if c.ConnectionType == structs.Wireless {
			connType = fmt.Sprintf("wifi(rssi=%d)", c.RSSI)
		}
		t.Logf("  [%d] %s %-16s %-15s %s/%s", i, c.MAC, c.IP, c.Name, connType, c.IPType)
	}
}

// TestPortForwardCRUD adds, verifies, and removes a test rule.
func TestPortForwardCRUD(t *testing.T) {
	router := RouterInfo()
	client := CreateClient()
	if !routers.Login(client, router) {
		t.Fatal("login failed")
	}

	// cleanup 먼저 (이전 테스트 잔여물 제거)
	routers.RemovePortForward(client, router, &structs.PortForward{Name: testPFName})

	// Add
	ok, err := routers.AddPortForward(client, router, &structs.PortForward{
		Name:     testPFName,
		Active:   false,
		Protocol: "TCP",
		Target:   "192.168.0.253",
		Src:      structs.PortRange{Start: 19998},
		Dst:      structs.PortRange{Start: 19998},
	})
	if err != nil || !ok {
		t.Fatalf("FAIL AddPortForward: %v", err)
	}
	t.Logf("PASS: port forward %q added", testPFName)

	// Verify it exists
	found := false
	for _, pf := range routers.GetPortForwardList(client, router) {
		if pf.Name == testPFName {
			found = true
			t.Logf("PASS: verified rule name=%s proto=%s src=%d dst=%d target=%s",
				pf.Name, pf.Protocol, pf.Src.Start, pf.Dst.Start, pf.Target)
			break
		}
	}
	if !found {
		t.Error("FAIL: added rule not found in list")
	}

	// Remove — RemovePortForward verifies absence in the del response (remaining list)
	ok, err = routers.RemovePortForward(client, router, &structs.PortForward{Name: testPFName})
	if err != nil || !ok {
		t.Fatalf("FAIL RemovePortForward: %v", err)
	}
	t.Logf("PASS: port forward %q removed and confirmed absent", testPFName)
}

// TestWOLCRUD adds, lists, sends signal, and removes a WOL entry.
func TestWOLCRUD(t *testing.T) {
	router := RouterInfo()
	client := CreateClient()
	if !routers.Login(client, router) {
		t.Fatal("login failed")
	}

	// cleanup
	routers.RemoveWOL(client, router, testWOLMac)

	// Add
	ok, err := routers.AddWOL(client, router, testWOLMac, testWOLName)
	if err != nil || !ok {
		t.Fatalf("FAIL AddWOL: %v", err)
	}
	t.Logf("PASS: WOL entry %s (%s) added", testWOLMac, testWOLName)

	// Verify
	found := false
	for _, w := range routers.GetWOLList(client, router) {
		if strings.EqualFold(w.MAC, testWOLMac) {
			found = true
			t.Logf("PASS: WOL entry verified mac=%s name=%s", w.MAC, w.Name)
			break
		}
	}
	if !found {
		t.Error("FAIL: WOL entry not found after add")
	}

	// Signal (fire-and-forget, just check no panic)
	routers.Wake(client, router, testWOLMac)
	t.Log("PASS: WOL signal sent (no error)")

	// Remove
	ok, err = routers.RemoveWOL(client, router, testWOLMac)
	if err != nil || !ok {
		t.Fatalf("FAIL RemoveWOL: %v", err)
	}
	t.Logf("PASS: WOL entry %s removed", testWOLMac)

	// Verify gone
	for _, w := range routers.GetWOLList(client, router) {
		if strings.EqualFold(w.MAC, testWOLMac) {
			t.Errorf("FAIL: WOL entry %s still exists after removal", testWOLMac)
		}
	}
	t.Log("PASS: WOL entry confirmed removed")
}

// TestMacAuth reads MAC auth policy (read-only, no modification).
func TestMacAuth(t *testing.T) {
	router := RouterInfo()
	client := CreateClient()
	if !routers.Login(client, router) {
		t.Fatal("login failed")
	}

	list := routers.GetMacAuthList(client, router)
	if len(list) == 0 {
		t.Error("FAIL: MAC auth list is empty")
		return
	}
	t.Logf("PASS: %d BSS band(s) returned", len(list))
	for _, m := range list {
		t.Logf("  bsstag=%s policy=%s macs=%d", m.BSSTag, m.Policy, len(m.MacList))
	}
}
