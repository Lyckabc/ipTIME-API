package routers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kongwoojin/ipTIME-API/cmd/structs"
)

func GetWOLList(client *http.Client, router *structs.Router) []structs.WOL {
	raw, err := serviceCall(client, router, "wol/show", nil)
	if err != nil {
		return nil
	}

	var list []structs.WOL
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil
	}
	return list
}

func AddWOL(client *http.Client, router *structs.Router, macAddress string, name string) (bool, error) {
	params := map[string]string{
		"mac":    macAddress,
		"pcname": name,
	}

	_, err := serviceCall(client, router, "wol/add", params)
	if err != nil {
		return false, fmt.Errorf("failed to add %s to Wake on LAN list: %w", macAddress, err)
	}
	return true, nil
}

func RemoveWOL(client *http.Client, router *structs.Router, macAddress string) (bool, error) {
	// wol/del expects an array of MAC address strings
	params := []string{macAddress}

	_, err := serviceCall(client, router, "wol/del", params)
	if err != nil {
		return false, fmt.Errorf("failed to remove %s from Wake on LAN list: %w", macAddress, err)
	}
	return true, nil
}

func Wake(client *http.Client, router *structs.Router, macAddress string) {
	params := map[string]string{"mac": macAddress}
	serviceCall(client, router, "wol/signal", params) //nolint:errcheck
}
