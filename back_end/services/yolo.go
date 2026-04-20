package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultYoloInferURL     = "http://127.0.0.1:8001/predict/fish-species"
	defaultYoloTimeoutInSec = 5
)

type yoloInferenceResponse struct {
	TypeConfidence map[string]float64 `json:"type_confidence"`
	Result         map[string]float64 `json:"result"`
	Msg            string             `json:"msg"`
	Type           string             `json:"type"`
	Confidence     float64            `json:"confidence"`
	Predictions    []struct {
		Type       string  `json:"type"`
		Label      string  `json:"label"`
		Confidence float64 `json:"confidence"`
		Score      float64 `json:"score"`
	} `json:"predictions"`
}

// PredictFishTypeConfidence 调用独立推理服务，返回鱼种-置信度映射。
// 约定优先读取 type_confidence 字段；兼容 result / predictions 两种格式。
func PredictFishTypeConfidence(imageData []byte) (map[string]float64, error) {
	endpoint := strings.TrimSpace(os.Getenv("YOLO_INFER_URL"))
	if endpoint == "" {
		endpoint = defaultYoloInferURL
	}

	timeoutSec := defaultYoloTimeoutInSec
	if s := strings.TrimSpace(os.Getenv("YOLO_INFER_TIMEOUT_SEC")); s != "" {
		if v, err := strconv.Atoi(s); err == nil && v > 0 {
			timeoutSec = v
		}
	}

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("fish_pic", "fish.jpg")
	if err != nil {
		return nil, fmt.Errorf("create form file: %w", err)
	}
	if _, err := part.Write(imageData); err != nil {
		return nil, fmt.Errorf("write image data: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("close multipart writer: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, endpoint, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: time.Duration(timeoutSec) * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("call yolo inference service: %w", err)
	}
	defer resp.Body.Close()

	payload, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read yolo response: %w", err)
	}

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("yolo inference service returned status %d", resp.StatusCode)
	}

	var parsed yoloInferenceResponse
	if err := json.Unmarshal(payload, &parsed); err != nil {
		return nil, fmt.Errorf("decode yolo response: %w", err)
	}

	typeConfidence := normalizeTypeConfidence(parsed)
	if len(typeConfidence) == 0 {
		return nil, fmt.Errorf("empty yolo result")
	}

	return typeConfidence, nil
}

func normalizeTypeConfidence(resp yoloInferenceResponse) map[string]float64 {
	if strings.TrimSpace(resp.Type) != "" {
		return map[string]float64{resp.Type: resp.Confidence}
	}

	if len(resp.TypeConfidence) > 0 {
		return resp.TypeConfidence
	}
	if len(resp.Result) > 0 {
		return resp.Result
	}

	out := make(map[string]float64)
	for _, item := range resp.Predictions {
		label := strings.TrimSpace(item.Type)
		if label == "" {
			label = strings.TrimSpace(item.Label)
		}
		if label == "" {
			continue
		}

		confidence := item.Confidence
		if confidence == 0 {
			confidence = item.Score
		}

		if old, exists := out[label]; !exists || confidence > old {
			out[label] = confidence
		}
	}
	return out
}