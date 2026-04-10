package routers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Lyckabc/ipTIME-API/cmd/structs"
)

type apiResponse struct {
	Result json.RawMessage `json:"result"`
	Error  *apiError       `json:"error"`
}

type apiError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func serviceCall(client *http.Client, router *structs.Router, method string, params interface{}) (json.RawMessage, error) {
	baseURL := fmt.Sprintf("http://%s:%d", router.Host, router.Port)

	payload := map[string]interface{}{"method": method}
	if params != nil {
		payload["params"] = params
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", baseURL+"/cgi/service.cgi", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Origin", "http://"+router.Host)
	req.Header.Set("Referer", "http://"+router.Host+"/ui/")
	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result apiResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Error != nil {
		return nil, fmt.Errorf("API error %d: %s", result.Error.Code, result.Error.Message)
	}

	return result.Result, nil
}
