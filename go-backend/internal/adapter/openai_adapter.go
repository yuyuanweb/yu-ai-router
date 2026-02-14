package adapter

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
)

type OpenAIAdapter struct {
	base *BaseOpenAICompatibleAdapter
}

func NewOpenAIAdapter() *OpenAIAdapter {
	return &OpenAIAdapter{
		base: NewBaseOpenAICompatibleAdapter("/v1/chat/completions"),
	}
}

func (a *OpenAIAdapter) Supports(providerName string) bool {
	provider := strings.ToLower(strings.TrimSpace(providerName))
	switch provider {
	case "openai", "gpt", "qwen", "dashscope", "tongyi", "aliyun", "deepseek":
		return true
	default:
		return false
	}
}

func (a *OpenAIAdapter) Invoke(model *entity.Model, provider *entity.ModelProvider, chatRequest dto.ChatRequest) ([]byte, error) {
	resp, body, err := a.base.invoke(model, provider, chatRequest, false)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("upstream status=%d body=%s", resp.StatusCode, string(body))
	}
	return body, nil
}

func (a *OpenAIAdapter) InvokeStream(model *entity.Model, provider *entity.ModelProvider, chatRequest dto.ChatRequest) (*http.Response, error) {
	resp, body, err := a.base.invoke(model, provider, chatRequest, true)
	if err != nil {
		if len(body) > 0 {
			return nil, fmt.Errorf("upstream stream request failed: body=%s", string(body))
		}
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		defer resp.Body.Close()
		_, _ = io.ReadAll(resp.Body)
		return nil, fmt.Errorf("upstream status=%d", resp.StatusCode)
	}
	return resp, nil
}
