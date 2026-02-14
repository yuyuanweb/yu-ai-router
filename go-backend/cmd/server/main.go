package main

import (
	"context"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/yupi/airouter/go-backend/internal/adapter"
	"github.com/yupi/airouter/go-backend/internal/config"
	"github.com/yupi/airouter/go-backend/internal/controller"
	"github.com/yupi/airouter/go-backend/internal/repository"
	"github.com/yupi/airouter/go-backend/internal/router"
	"github.com/yupi/airouter/go-backend/internal/service"
	"github.com/yupi/airouter/go-backend/internal/strategy"
	"github.com/yupi/airouter/go-backend/internal/task"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config failed: %v", err)
	}

	db, err := gorm.Open(mysql.Open(cfg.MySQLDSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("open mysql with gorm failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("get sql db failed: %v", err)
	}
	defer sqlDB.Close()

	if err = sqlDB.Ping(); err != nil {
		log.Fatalf("ping mysql failed: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	apiKeyRepo := repository.NewApiKeyRepository(db)
	requestLogRepo := repository.NewRequestLogRepository(db)
	providerRepo := repository.NewProviderRepository(db)
	modelRepo := repository.NewModelRepository(db)

	userService := service.NewUserService(userRepo)
	apiKeyService := service.NewApiKeyService(apiKeyRepo)
	requestLogService := service.NewRequestLogService(requestLogRepo, apiKeyService)
	providerService := service.NewProviderService(providerRepo)
	modelService := service.NewModelService(modelRepo, providerRepo)
	healthCheckService := service.NewHealthCheckService(providerRepo, modelRepo, requestLogRepo)

	routingStrategies := []strategy.RoutingStrategy{
		strategy.NewAutoRoutingStrategy(),
		strategy.NewFixedRoutingStrategy(),
		strategy.NewCostFirstRoutingStrategy(),
		strategy.NewLatencyFirstRoutingStrategy(),
	}
	routingService := service.NewRoutingService(modelRepo, routingStrategies)
	adapterFactory := adapter.NewModelAdapterFactory(
		[]adapter.ModelAdapter{
			adapter.NewZhipuAdapter(),
			adapter.NewOpenAIAdapter(),
		},
		adapter.NewDefaultAdapter(),
	)
	modelInvokeService := service.NewModelInvokeService(adapterFactory)
	chatService := service.NewChatService(requestLogService, routingService, modelInvokeService, providerService)

	healthController := controller.NewHealthController()
	userController := controller.NewUserController(userService)
	apiKeyController := controller.NewApiKeyController(apiKeyService, userService)
	providerController := controller.NewProviderController(providerService)
	modelController := controller.NewModelController(modelService)
	chatController := controller.NewChatController(chatService, apiKeyService)
	internalChatController := controller.NewInternalChatController(chatService, apiKeyService, userService)
	statsController := controller.NewStatsController(requestLogService, userService)
	healthCheckTask := task.NewHealthCheckTask(healthCheckService)

	engine, err := router.New(
		cfg,
		healthController,
		userController,
		apiKeyController,
		providerController,
		modelController,
		chatController,
		internalChatController,
		statsController,
		userService,
	)
	if err != nil {
		log.Fatalf("build router failed: %v", err)
	}

	taskCtx, cancelTask := context.WithCancel(context.Background())
	defer cancelTask()
	healthCheckTask.Start(taskCtx)

	if err = engine.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server run failed: %v", err)
	}
}
