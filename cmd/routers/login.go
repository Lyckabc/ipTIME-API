package routers

import (
	"encoding/json"
	"net/http"

	"github.com/kongwoojin/ipTIME-API/cmd/structs"
)

func Login(client *http.Client, router *structs.Router) bool {
	params := map[string]string{
		"id": router.Username,
		"pw": router.Password,
	}

	raw, err := serviceCall(client, router, "session/login", params)
	if err != nil {
		return false
	}

	var result string
	if err := json.Unmarshal(raw, &result); err != nil {
		return false
	}

	return result == "done"
}
