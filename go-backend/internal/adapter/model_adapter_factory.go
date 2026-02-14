package adapter

import (
	"strings"

	"github.com/yupi/airouter/go-backend/internal/errno"
)

type ModelAdapterFactory struct {
	adapters       []ModelAdapter
	defaultAdapter ModelAdapter
}

func NewModelAdapterFactory(adapters []ModelAdapter, defaultAdapter ModelAdapter) *ModelAdapterFactory {
	return &ModelAdapterFactory{
		adapters:       adapters,
		defaultAdapter: defaultAdapter,
	}
}

func (f *ModelAdapterFactory) GetAdapter(providerName string) (ModelAdapter, error) {
	if strings.TrimSpace(providerName) == "" {
		return nil, errno.NewWithMessage(errno.ParamsError, "提供者名称不能为空")
	}
	for _, item := range f.adapters {
		if item.Supports(providerName) {
			return item, nil
		}
	}
	return f.defaultAdapter, nil
}
