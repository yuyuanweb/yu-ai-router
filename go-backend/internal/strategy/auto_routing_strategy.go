package strategy

import (
	"math"
	"sort"

	"github.com/yupi/airouter/go-backend/internal/constant"
	"github.com/yupi/airouter/go-backend/internal/model/entity"
)

type AutoRoutingStrategy struct{}

func NewAutoRoutingStrategy() *AutoRoutingStrategy {
	return &AutoRoutingStrategy{}
}

func (s *AutoRoutingStrategy) GetStrategyType() string {
	return constant.RoutingStrategyAuto
}

func (s *AutoRoutingStrategy) SelectModel(models []entity.Model, _ string) *entity.Model {
	sorted := sortByScore(models)
	if len(sorted) == 0 {
		return nil
	}
	selected := sorted[0]
	return &selected
}

func (s *AutoRoutingStrategy) GetFallbackModels(models []entity.Model, _ string) []entity.Model {
	sorted := sortByScore(models)
	if len(sorted) <= 1 {
		return make([]entity.Model, 0)
	}
	return sorted[1:]
}

func sortByScore(models []entity.Model) []entity.Model {
	cloned := append([]entity.Model(nil), models...)
	sort.Slice(cloned, func(i, j int) bool {
		si := scoreOrderValue(cloned[i])
		sj := scoreOrderValue(cloned[j])
		if si == sj {
			return cloned[i].ID < cloned[j].ID
		}
		return si < sj
	})
	return cloned
}

func scoreOrderValue(model entity.Model) float64 {
	if model.Score <= 0 {
		return math.MaxFloat64
	}
	return model.Score
}
