package routers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/kongwoojin/ipTIME-API/cmd/structs"
)

func GetMacAuthList(client *http.Client, router *structs.Router) []structs.MacAuth {
	raw, err := serviceCall(client, router, "wireless/mac/show", nil)
	if err != nil {
		return nil
	}

	var list []structs.MacAuth
	if err := json.Unmarshal(raw, &list); err != nil {
		return nil
	}
	return list
}

func ChangeMacAuthPolicy(client *http.Client, router *structs.Router, bsstag string, policy structs.MacAuthPolicy) (bool, error) {
	params := map[string]string{
		"bsstag": bsstag,
		"policy": string(policy),
	}

	_, err := serviceCall(client, router, "wireless/mac/policy", params)
	if err != nil {
		return false, fmt.Errorf("failed to change MAC auth policy for %s: %w", bsstag, err)
	}
	return true, nil
}

func AddMacAuth(client *http.Client, router *structs.Router, bsstag string, macAddress string) (bool, error) {
	params := map[string]string{
		"bsstag": bsstag,
		"mac":    macAddress,
	}

	_, err := serviceCall(client, router, "wireless/mac/add", params)
	if err != nil {
		return false, fmt.Errorf("failed to add MAC %s to %s: %w", macAddress, bsstag, err)
	}
	return true, nil
}

func RemoveMacAuth(client *http.Client, router *structs.Router, bsstag string, macAddress string) (bool, error) {
	params := map[string]string{
		"bsstag": bsstag,
		"mac":    macAddress,
	}

	_, err := serviceCall(client, router, "wireless/mac/del", params)
	if err != nil {
		return false, fmt.Errorf("failed to remove MAC %s from %s: %w", macAddress, bsstag, err)
	}
	return true, nil
}
