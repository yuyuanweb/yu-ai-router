package service

import (
	"net/http"

	"github.com/yupi/airouter/go-backend/internal/adapter"
	"github.com/yupi/airouter/go-backend/internal/model/dto"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
)

type ModelInvokeService struct {
	adapterFactory *adapter.ModelAdapterFactory
}

func NewModelInvokeService(adapterFactory *adapter.ModelAdapterFactory) *ModelInvokeService {
	return &ModelInvokeService{adapterFactory: adapterFactory}
}

func (s *ModelInvokeService) Invoke(model *entity.Model, provider *entity.ModelProvider, chatRequest dto.ChatRequest) ([]byte, error) {
	modelAdapter, err := s.adapterFactory.GetAdapter(provider.ProviderName)
	if err != nil {
		return nil, err
	}
	return modelAdapter.Invoke(model, provider, chatRequest)
}

func (s *ModelInvokeService) InvokeStream(model *entity.Model, provider *entity.ModelProvider, chatRequest dto.ChatRequest) (*http.Response, error) {
	modelAdapter, err := s.adapterFactory.GetAdapter(provider.ProviderName)
	if err != nil {
		return nil, err
	}
	return modelAdapter.InvokeStream(model, provider, chatRequest)
}
