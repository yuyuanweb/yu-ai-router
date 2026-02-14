package strategy

import (
	"sort"

	"github.com/yupi/airouter/go-backend/internal/constant"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
)

type LatencyFirstRoutingStrategy struct{}

func NewLatencyFirstRoutingStrategy() *LatencyFirstRoutingStrategy {
	return &LatencyFirstRoutingStrategy{}
}

func (s *LatencyFirstRoutingStrategy) GetStrategyType() string {
	return constant.RoutingStrategyLatencyFirst
}

func (s *LatencyFirstRoutingStrategy) SelectModel(models []entity.Model, _ string) *entity.Model {
	sorted := sortByLatency(models)
	if len(sorted) == 0 {
		return nil
	}
	selected := sorted[0]
	return &selected
}

func (s *LatencyFirstRoutingStrategy) GetFallbackModels(models []entity.Model, _ string) []entity.Model {
	sorted := sortByLatency(models)
	if len(sorted) <= 1 {
		return make([]entity.Model, 0)
	}
	return sorted[1:]
}

func sortByLatency(models []entity.Model) []entity.Model {
	cloned := append([]entity.Model(nil), models...)
	sort.Slice(cloned, func(i, j int) bool {
		li := latencyOrderValue(cloned[i])
		lj := latencyOrderValue(cloned[j])
		if li == lj {
			return cloned[i].ID < cloned[j].ID
		}
		return li < lj
	})
	return cloned
}

func latencyOrderValue(model entity.Model) int {
	if model.AvgLatency <= 0 {
		return 999999
	}
	return model.AvgLatency
}
