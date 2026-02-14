package adapter

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
)

const invokeTimeout = 60 * time.Second

type BaseOpenAICompatibleAdapter struct {
	httpClient *http.Client
	path       string
}

func NewBaseOpenAICompatibleAdapter(path string) *BaseOpenAICompatibleAdapter {
	return &BaseOpenAICompatibleAdapter{
		httpClient: &http.Client{Timeout: invokeTimeout},
		path:       path,
	}
}

func (a *BaseOpenAICompatibleAdapter) invoke(model *entity.Model, provider *entity.ModelProvider, chatRequest dto.ChatRequest, stream bool) (*http.Response, []byte, error) {
	payload := map[string]any{
		"model":    model.ModelKey,
		"messages": chatRequest.Messages,
		"stream":   stream,
	}
	if chatRequest.Temperature != nil {
		payload["temperature"] = *chatRequest.Temperature
	}
	if chatRequest.MaxTokens != nil {
		payload["max_tokens"] = *chatRequest.MaxTokens
	}
	if chatRequest.EnableReasoning != nil && *chatRequest.EnableReasoning {
		// 与 Java 版行为对齐：仅在模型声明支持深度思考时开启。
		if model.SupportReasoning == 1 {
			payload["enable_reasoning"] = true
			payload["enable_thinking"] = true
		} else {
			log.Printf("reasoning requested but model does not support: model=%s provider=%s", model.ModelKey, provider.ProviderName)
		}
	}
	if stream {
		payload["stream_options"] = map[string]any{"include_usage": true}
	}

	rawPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, nil, err
	}
	url := strings.TrimRight(provider.BaseURL, "/") + a.path
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(rawPayload))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+provider.APIKey)

	resp, err := a.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	if stream {
		if resp.StatusCode >= http.StatusBadRequest {
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			return nil, body, io.EOF
		}
		return resp, nil, nil
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}
	return resp, body, nil
}
