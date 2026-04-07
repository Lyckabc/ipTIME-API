package routers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/kongwoojin/ipTIME-API/cmd/structs"
)

type portForwardList struct {
	User []structs.PortForward `json:"user"`
	UPnP []structs.PortForward `json:"upnp"`
}

func GetPortForwardList(client *http.Client, router *structs.Router) []structs.PortForward {
	raw, err := serviceCall(client, router, "portforward/get", nil)
	if err != nil {
		return nil
	}

	var list portForwardList
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil
	}
	return list.User
}

func AddPortForward(client *http.Client, router *structs.Router, pf *structs.PortForward) (bool, error) {
	if pf.Target == router.Host {
		return false, fmt.Errorf("cannot add router IP to port forward rule")
	}
	if checkPortForwardExist(client, router, pf.Name) {
		return false, fmt.Errorf("port forward rule %q already exists", pf.Name)
	}

	proto := strings.ToLower(pf.Protocol)
	switch proto {
	case "tcp", "udp", "tcpudp", "gre":
	default:
		return false, fmt.Errorf("unknown protocol: %s", pf.Protocol)
	}

	params := map[string]interface{}{
		"name":     pf.Name,
		"active":   pf.Active,
		"protocol": proto,
		"target":   pf.Target,
		"src":      structs.PortRange{Start: pf.Src.Start, End: pf.Src.End},
		"dst":      structs.PortRange{Start: pf.Dst.Start, End: pf.Dst.End},
	}

	_, err := serviceCall(client, router, "portforward/add", params)
	if err != nil {
		return false, fmt.Errorf("failed to add port forward rule %q: %w", pf.Name, err)
	}
	return true, nil
}

func RemovePortForward(client *http.Client, router *structs.Router, pf *structs.PortForward) (bool, error) {
	if !checkPortForwardExist(client, router, pf.Name) {
		return false, fmt.Errorf("port forward rule %q not found", pf.Name)
	}

	// portforward/del takes {"type":"user","list":["name1","name2",...]}
	// and returns the remaining list after deletion.
	params := map[string]interface{}{
		"type": "user",
		"list": []string{pf.Name},
	}

	raw, err := serviceCall(client, router, "portforward/del", params)
	if err != nil {
		return false, fmt.Errorf("failed to remove port forward rule %q: %w", pf.Name, err)
	}

	// Verify the rule is absent from the returned remaining list.
	var remaining []structs.PortForward
	if err := json.Unmarshal(raw, &remaining); err == nil {
		for _, r := range remaining {
			if r.Name == pf.Name {
				return false, fmt.Errorf("port forward rule %q still present after deletion", pf.Name)
			}
		}
	}
	return true, nil
}

func checkPortForwardExist(client *http.Client, router *structs.Router, name string) bool {
	for _, pf := range GetPortForwardList(client, router) {
		if pf.Name == name {
			return true
		}
	}
	return false
}
