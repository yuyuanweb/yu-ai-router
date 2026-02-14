package adapter

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
)

type ZhipuAdapter struct {
	base *BaseOpenAICompatibleAdapter
}

func NewZhipuAdapter() *ZhipuAdapter {
	return &ZhipuAdapter{
		base: NewBaseOpenAICompatibleAdapter("/chat/completions"),
	}
}

func (a *ZhipuAdapter) Supports(providerName string) bool {
	provider := strings.ToLower(strings.TrimSpace(providerName))
	return provider == "zhipu" || provider == "zhipuai" || provider == "glm"
}

func (a *ZhipuAdapter) Invoke(model *entity.Model, provider *entity.ModelProvider, chatRequest dto.ChatRequest) ([]byte, error) {
	resp, body, err := a.base.invoke(model, provider, chatRequest, false)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("upstream status=%d body=%s", resp.StatusCode, string(body))
	}
	return body, nil
}

func (a *ZhipuAdapter) InvokeStream(model *entity.Model, provider *entity.ModelProvider, chatRequest dto.ChatRequest) (*http.Response, error) {
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
