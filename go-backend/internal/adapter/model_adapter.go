package adapter

import (
	"net/http"

	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
)

type ModelAdapter interface {
	Supports(providerName string) bool
	Invoke(model *entity.Model, provider *entity.ModelProvider, chatRequest dto.ChatRequest) ([]byte, error)
	InvokeStream(model *entity.Model, provider *entity.ModelProvider, chatRequest dto.ChatRequest) (*http.Response, error)
}
