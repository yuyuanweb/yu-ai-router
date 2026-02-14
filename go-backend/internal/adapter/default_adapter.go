package adapter

import (
	"fmt"
	"net/http"

	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
)

type DefaultAdapter struct {
	base *BaseOpenAICompatibleAdapter
}

func NewDefaultAdapter() *DefaultAdapter {
	return &DefaultAdapter{
		base: NewBaseOpenAICompatibleAdapter("/v1/chat/completions"),
	}
}

func (a *DefaultAdapter) Supports(_ string) bool {
	return true
}

func (a *DefaultAdapter) Invoke(model *entity.Model, provider *entity.ModelProvider, chatRequest dto.ChatRequest) ([]byte, error) {
	resp, body, err := a.base.invoke(model, provider, chatRequest, false)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("upstream status=%d body=%s", resp.StatusCode, string(body))
	}
	return body, nil
}

func (a *DefaultAdapter) InvokeStream(model *entity.Model, provider *entity.ModelProvider, chatRequest dto.ChatRequest) (*http.Response, error) {
	resp, body, err := a.base.invoke(model, provider, chatRequest, true)
	if err != nil {
		if len(body) > 0 {
			return nil, fmt.Errorf("upstream stream request failed: body=%s", string(body))
		}
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("upstream status=%d", resp.StatusCode)
	}
	return resp, nil
}
