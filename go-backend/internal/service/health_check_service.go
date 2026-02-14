package service

import (
	"log"
	"math"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/yupi/airouter/go-backend/internal/model/entity"
	"github.com/yupi/airouter/go-backend/internal/repository"
)

const (
	maxHealthHistorySize   = 100
	healthCheckTimeout     = 10 * time.Second
	modelMetricsStatsHours = 24

	costWeight       = 0.3
	latencyWeight    = 0.3
	successRateWeight = 0.2
	priorityWeight   = 0.2

	healthThresholdHealthy  = 80.0
	healthThresholdDegraded = 50.0
)

type HealthCheckService struct {
	providerRepo   *repository.ProviderRepository
	modelRepo      *repository.ModelRepository
	requestLogRepo *repository.RequestLogRepository
	httpClient     *http.Client

	historyMu       sync.RWMutex
	healthHistory   map[int64][]bool
}

func NewHealthCheckService(
	providerRepo *repository.ProviderRepository,
	modelRepo *repository.ModelRepository,
	requestLogRepo *repository.RequestLogRepository,
) *HealthCheckService {
	return &HealthCheckService{
		providerRepo:   providerRepo,
		modelRepo:      modelRepo,
		requestLogRepo: requestLogRepo,
		httpClient: &http.Client{
			Timeout: healthCheckTimeout,
		},
		healthHistory: make(map[int64][]bool),
	}
}

func (s *HealthCheckService) CheckAllProviders() {
	providers, err := s.providerRepo.ListActive()
	if err != nil {
		log.Printf("health check list active providers failed: err=%v", err)
		return
	}
	if len(providers) == 0 {
		log.Printf("health check skipped: no active providers")
		return
	}

	log.Printf("health check start: providerCount=%d", len(providers))
	var wg sync.WaitGroup
	for _, provider := range providers {
		providerCopy := provider
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.checkProviderHealth(&providerCopy)
		}()
	}
	wg.Wait()
	s.syncModelMetricsFromRequestLog()
	log.Printf("health check finished")
}

func (s *HealthCheckService) checkProviderHealth(provider *entity.ModelProvider) {
	if provider == nil {
		return
	}
	start := time.Now()
	ok := s.sendHealthCheckRequest(provider)
	latency := int(time.Since(start).Milliseconds())

	if !ok {
		latency = 0
	}
	s.recordHealth(provider.ID, ok)
	successRate := s.calculateSuccessRate(provider.ID)
	healthStatus := determineHealthStatus(successRate)

	_, err := s.providerRepo.UpdateByID(provider.ID, map[string]any{
		"healthStatus": healthStatus,
		"avgLatency":   latency,
		"successRate":  round2(successRate),
	})
	if err != nil {
		log.Printf("update provider health failed: providerId=%d err=%v", provider.ID, err)
		return
	}
	log.Printf("provider health checked: providerId=%d healthy=%t latency=%d successRate=%.2f healthStatus=%s", provider.ID, ok, latency, successRate, healthStatus)
}

func (s *HealthCheckService) sendHealthCheckRequest(provider *entity.ModelProvider) bool {
	baseURL := strings.TrimRight(provider.BaseURL, "/")
	if baseURL == "" {
		return false
	}
	req, err := http.NewRequest(http.MethodGet, baseURL+"/models", nil)
	if err != nil {
		return false
	}
	if strings.TrimSpace(provider.APIKey) != "" {
		req.Header.Set("Authorization", "Bearer "+provider.APIKey)
	}
	resp, err := s.httpClient.Do(req)
	if err != nil {
		log.Printf("provider health request failed: providerId=%d err=%v", provider.ID, err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusUnauthorized
}

func (s *HealthCheckService) recordHealth(providerID int64, ok bool) {
	s.historyMu.Lock()
	defer s.historyMu.Unlock()

	history := append(s.healthHistory[providerID], ok)
	if len(history) > maxHealthHistorySize {
		history = history[len(history)-maxHealthHistorySize:]
	}
	s.healthHistory[providerID] = history
}

func (s *HealthCheckService) calculateSuccessRate(providerID int64) float64 {
	s.historyMu.RLock()
	history := s.healthHistory[providerID]
	s.historyMu.RUnlock()
	if len(history) == 0 {
		return 100
	}
	var successCount int
	for _, ok := range history {
		if ok {
			successCount++
		}
	}
	return float64(successCount) * 100 / float64(len(history))
}

func (s *HealthCheckService) syncModelMetricsFromRequestLog() {
	activeModels, err := s.modelRepo.ListActive()
	if err != nil {
		log.Printf("sync model metrics failed: list active models err=%v", err)
		return
	}
	if len(activeModels) == 0 {
		return
	}

	statsRows, err := s.requestLogRepo.QueryModelStatsSince(time.Now().Add(-modelMetricsStatsHours * time.Hour))
	if err != nil {
		log.Printf("sync model metrics failed: query stats err=%v", err)
		return
	}
	statsMap := make(map[string]repository.ModelStatsRow, len(statsRows))
	for _, row := range statsRows {
		statsMap[row.ModelName] = row
	}
	norm := calculateNormParams(activeModels, statsMap)

	for _, model := range activeModels {
		stats, hasStats := statsMap[model.ModelKey]
		avgLatency := 0
		successRate := 100.0
		if hasStats {
			avgLatency = stats.AvgLatency
			successRate = stats.SuccessRate
		}
		healthStatus := determineHealthStatus(successRate)
		score := calculateScore(model, avgLatency, successRate, norm)

		_, updateErr := s.modelRepo.UpdateByID(model.ID, map[string]any{
			"healthStatus": healthStatus,
			"avgLatency":   avgLatency,
			"successRate":  round2(successRate),
			"score":        score,
		})
		if updateErr != nil {
			log.Printf("sync model metrics update failed: modelId=%d modelKey=%s err=%v", model.ID, model.ModelKey, updateErr)
		}
	}

	s.syncProviderMetrics(activeModels)
}

func (s *HealthCheckService) syncProviderMetrics(models []entity.Model) {
	if len(models) == 0 {
		return
	}
	type providerAgg struct {
		totalLatency int
		latencyCount int
		totalRate    float64
		rateCount    int
	}
	aggMap := make(map[int64]*providerAgg)
	for _, model := range models {
		item, ok := aggMap[model.ProviderID]
		if !ok {
			item = &providerAgg{}
			aggMap[model.ProviderID] = item
		}
		if model.AvgLatency > 0 {
			item.totalLatency += model.AvgLatency
			item.latencyCount++
		}
		if model.SuccessRate > 0 {
			item.totalRate += model.SuccessRate
			item.rateCount++
		}
	}
	for providerID, agg := range aggMap {
		avgLatency := 0
		successRate := 100.0
		if agg.latencyCount > 0 {
			avgLatency = agg.totalLatency / agg.latencyCount
		}
		if agg.rateCount > 0 {
			successRate = agg.totalRate / float64(agg.rateCount)
		}
		healthStatus := determineHealthStatus(successRate)
		_, err := s.providerRepo.UpdateByID(providerID, map[string]any{
			"healthStatus": healthStatus,
			"avgLatency":   avgLatency,
			"successRate":  round2(successRate),
		})
		if err != nil {
			log.Printf("sync provider metrics failed: providerId=%d err=%v", providerID, err)
		}
	}
}

type normParams struct {
	minCost     float64
	maxCost     float64
	minLatency  int
	maxLatency  int
	minPriority int
	maxPriority int
}

func calculateNormParams(models []entity.Model, stats map[string]repository.ModelStatsRow) normParams {
	p := normParams{
		minCost:     math.MaxFloat64,
		maxCost:     0,
		minLatency:  math.MaxInt,
		maxLatency:  0,
		minPriority: math.MaxInt,
		maxPriority: math.MinInt,
	}
	for _, model := range models {
		cost := model.InputPrice + model.OutputPrice
		if cost < p.minCost {
			p.minCost = cost
		}
		if cost > p.maxCost {
			p.maxCost = cost
		}
		if model.Priority < p.minPriority {
			p.minPriority = model.Priority
		}
		if model.Priority > p.maxPriority {
			p.maxPriority = model.Priority
		}
		if row, ok := stats[model.ModelKey]; ok && row.AvgLatency > 0 {
			if row.AvgLatency < p.minLatency {
				p.minLatency = row.AvgLatency
			}
			if row.AvgLatency > p.maxLatency {
				p.maxLatency = row.AvgLatency
			}
		}
	}
	if p.minCost == math.MaxFloat64 {
		p.minCost = 0
	}
	if p.minLatency == math.MaxInt {
		p.minLatency = 0
	}
	if p.minPriority == math.MaxInt {
		p.minPriority = 0
	}
	if p.maxPriority == math.MinInt {
		p.maxPriority = 100
	}
	if p.maxLatency == 0 {
		p.maxLatency = 10000
	}
	if p.maxCost == 0 {
		p.maxCost = 1
	}
	return p
}

func calculateScore(model entity.Model, avgLatency int, successRate float64, p normParams) float64 {
	costScore := normalizeFloat(model.InputPrice+model.OutputPrice, p.minCost, p.maxCost)
	if avgLatency <= 0 {
		avgLatency = 5000
	}
	latencyScore := normalizeInt(avgLatency, p.minLatency, p.maxLatency)
	successRateScore := 1.0 - successRate/100.0
	priorityScore := 1.0 - normalizeInt(model.Priority, p.minPriority, p.maxPriority)
	score := costScore*costWeight +
		latencyScore*latencyWeight +
		successRateScore*successRateWeight +
		priorityScore*priorityWeight
	return round4(score)
}

func normalizeFloat(value, min, max float64) float64 {
	if max <= min {
		return 0
	}
	return (value - min) / (max - min)
}

func normalizeInt(value, min, max int) float64 {
	if max <= min {
		return 0
	}
	return float64(value-min) / float64(max-min)
}

func determineHealthStatus(successRate float64) string {
	switch {
	case successRate >= healthThresholdHealthy:
		return "healthy"
	case successRate >= healthThresholdDegraded:
		return "degraded"
	default:
		return "unhealthy"
	}
}

func round2(v float64) float64 {
	return math.Round(v*100) / 100
}

func round4(v float64) float64 {
	return math.Round(v*10000) / 10000
}
