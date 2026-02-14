package strategy

import (
	"sort"

	"github.com/yupi/airouter/go-backend/internal/constant"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
)

type CostFirstRoutingStrategy struct{}

func NewCostFirstRoutingStrategy() *CostFirstRoutingStrategy {
	return &CostFirstRoutingStrategy{}
}

func (s *CostFirstRoutingStrategy) GetStrategyType() string {
	return constant.RoutingStrategyCostFirst
}

func (s *CostFirstRoutingStrategy) SelectModel(models []entity.Model, _ string) *entity.Model {
	sorted := sortByCost(models)
	if len(sorted) == 0 {
		return nil
	}
	selected := sorted[0]
	return &selected
}

func (s *CostFirstRoutingStrategy) GetFallbackModels(models []entity.Model, _ string) []entity.Model {
	sorted := sortByCost(models)
	if len(sorted) <= 1 {
		return make([]entity.Model, 0)
	}
	return sorted[1:]
}

func sortByCost(models []entity.Model) []entity.Model {
	cloned := append([]entity.Model(nil), models...)
	sort.Slice(cloned, func(i, j int) bool {
		ci := cloned[i].InputPrice + cloned[i].OutputPrice
		cj := cloned[j].InputPrice + cloned[j].OutputPrice
		if ci == cj {
			return cloned[i].ID < cloned[j].ID
		}
		return ci < cj
	})
	return cloned
}
