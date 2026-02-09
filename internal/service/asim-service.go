package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"reconTool/global"
	"reconTool/internal/dto"

	"github.com/google/uuid"
)

var asimHttpClient = &http.Client{Timeout: 120 * time.Second}

func CallAsimApi[TReq any, TResp any](method, endpoint string, body TReq) (*dto.AsimBaseResponse[TResp], error) {
	cfg, err := GetGlobalConfig()
	if err != nil {
		return nil, fmt.Errorf("Can not get config: %v", err)
	}

	url := fmt.Sprintf("%s/lcs-new/api/partner/%s", cfg.PartnerUrl, endpoint)

	var bodyBytes []byte
	if any(body) != nil {
		bodyBytes, _ = json.Marshal(body)
	}

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))

	requestID := uuid.New().String()
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("requestId", requestID)
	req.Header.Set("deviceId", cfg.DeviceId)
	req.Header.Set("accept", "*/*")
	req.Header.Set("Authorization", "Bearer "+cfg.Token)

	start := time.Now()
	resp, err := asimHttpClient.Do(req)

	var respBodyStr string
	if err == nil {
		b, _ := io.ReadAll(resp.Body)
		respBodyStr = string(b)
		resp.Body = io.NopCloser(bytes.NewBuffer(b))
	}

	global.Logger.LogPartnerCall(req, respBodyStr, time.Since(start))

	if err != nil {
		return nil, fmt.Errorf("Error connect third party: %v", err)
	}
	defer resp.Body.Close()

	var baseRes dto.AsimBaseResponse[TResp]
	if err := json.Unmarshal([]byte(respBodyStr), &baseRes); err != nil {
		return &dto.AsimBaseResponse[TResp]{
			IsSucceeded:  false,
			ErrorMessage: "Error Formatter json from third party",
		}, nil
	}

	return &baseRes, nil
}
