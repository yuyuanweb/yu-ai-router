package task

import (
	"context"
	"log"
	"time"

	"github.com/yupi/airouter/go-backend/internal/service"
)

const healthCheckInterval = 30 * time.Second

type HealthCheckTask struct {
	healthCheckService *service.HealthCheckService
}

func NewHealthCheckTask(healthCheckService *service.HealthCheckService) *HealthCheckTask {
	return &HealthCheckTask{healthCheckService: healthCheckService}
}

func (t *HealthCheckTask) Start(ctx context.Context) {
	go func() {
		t.executeHealthCheck()
		ticker := time.NewTicker(healthCheckInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				t.executeHealthCheck()
			}
		}
	}()
}

func (t *HealthCheckTask) executeHealthCheck() {
	log.Printf("execute scheduled health check task")
	defer func() {
		if r := recover(); r != nil {
			log.Printf("health check task panic: %v", r)
		}
	}()
	t.healthCheckService.CheckAllProviders()
}
