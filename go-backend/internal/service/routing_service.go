package service

import (
	"strings"

	"github.com/yupi/airouter/go-backend/internal/constant"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"github.com/yupi/airouter/go-backend/internal/repository"
	"github.com/yupi/airouter/go-backend/internal/strategy"
)

const maxRoutingFallbackSize = 3

type RoutingService struct {
	modelRepo   *repository.ModelRepository
	strategyMap map[string]strategy.RoutingStrategy
}

func NewRoutingService(modelRepo *repository.ModelRepository, strategies []strategy.RoutingStrategy) *RoutingService {
	m := make(map[string]strategy.RoutingStrategy, len(strategies))
	for _, item := range strategies {
		m[item.GetStrategyType()] = item
	}
	return &RoutingService{
		modelRepo:   modelRepo,
		strategyMap: m,
	}
}

func (s *RoutingService) SelectModel(strategyType, modelType, requestedModel string) (*entity.Model, []entity.Model, error) {
	finalStrategy := s.getStrategy(strategyType)
	models, err := s.modelRepo.ListRoutable(modelType)
	if err != nil {
		return nil, nil, err
	}
	selected := finalStrategy.SelectModel(models, requestedModel)
	fallbacks := finalStrategy.GetFallbackModels(models, requestedModel)
	if len(fallbacks) > maxRoutingFallbackSize {
		fallbacks = fallbacks[:maxRoutingFallbackSize]
	}
	return selected, fallbacks, nil
}

func (s *RoutingService) DetermineStrategyType(requestedStrategy, requestedModel string) string {
	if strings.TrimSpace(requestedStrategy) != "" {
		return requestedStrategy
	}
	if strings.TrimSpace(requestedModel) != "" {
		return constant.RoutingStrategyFixed
	}
	return constant.RoutingStrategyAuto
}

func (s *RoutingService) getStrategy(strategyType string) strategy.RoutingStrategy {
	if item, ok := s.strategyMap[strategyType]; ok {
		return item
	}
	if autoStrategy, ok := s.strategyMap[constant.RoutingStrategyAuto]; ok {
		return autoStrategy
	}
	return strategy.NewAutoRoutingStrategy()
}
