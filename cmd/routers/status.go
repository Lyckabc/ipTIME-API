package routers

import (
	"encoding/json"
	"net/http"

	"github.com/Lyckabc/ipTIME-API/cmd/structs"
)

func RouterStatus(client *http.Client, router *structs.Router) (*structs.RouterStatus, error) {
	status := &structs.RouterStatus{}

	sysRaw, err := serviceCall(client, router, "system/info", nil)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(sysRaw, &status.SystemInfo); err != nil {
		return nil, err
	}

	lanRaw, err := serviceCall(client, router, "network/interface/lan/info", nil)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(lanRaw, &status.LANInfo); err != nil {
		return nil, err
	}

	portRaw, err := serviceCall(client, router, "port/link/status", nil)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(portRaw, &status.PortLinks); err != nil {
		return nil, err
	}

	return status, nil
}
