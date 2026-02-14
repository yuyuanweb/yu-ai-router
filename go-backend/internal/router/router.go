package router

import (
	"strconv"

	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"

	"github.com/yupi/airouter/go-backend/internal/config"
	"github.com/yupi/airouter/go-backend/internal/controller"
	"github.com/yupi/airouter/go-backend/internal/middleware"
	"github.com/yupi/airouter/go-backend/internal/service"
)

const sessionPoolSize = 10

func New(
	cfg *config.Config,
	healthController *controller.HealthController,
	userController *controller.UserController,
	apiKeyController *controller.ApiKeyController,
	providerController *controller.ProviderController,
	modelController *controller.ModelController,
	chatController *controller.ChatController,
	internalChatController *controller.InternalChatController,
	statsController *controller.StatsController,
	userService *service.UserService,
) (*gin.Engine, error) {
	engine := gin.New()
	engine.Use(gin.Logger())
	engine.Use(middleware.Recovery())
	engine.Use(middleware.CORS())

	store, err := redisStore.NewStoreWithDB(
		sessionPoolSize,
		"tcp",
		cfg.RedisAddr,
		cfg.RedisUsername,
		cfg.RedisPassword,
		strconv.Itoa(cfg.RedisDB),
		[]byte(cfg.SessionSecret),
	)
	if err != nil {
		return nil, err
	}
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   cfg.SessionMaxAge,
		HttpOnly: true,
	})
	engine.Use(sessions.Sessions(cfg.SessionName, store))

	apiGroup := engine.Group(cfg.ContextPath)
	{
		healthGroup := apiGroup.Group("/health")
		healthGroup.GET("/", healthController.HealthCheck)

		userGroup := apiGroup.Group("/user")
		userGroup.POST("/register", userController.UserRegister)
		userGroup.POST("/login", userController.UserLogin)
		userGroup.GET("/get/login", userController.GetLoginUser)
		userGroup.POST("/logout", userController.UserLogout)

		userGroup.POST("/add", middleware.RequireAdmin(userService), userController.AddUser)
		userGroup.GET("/get", middleware.RequireAdmin(userService), userController.GetUserByID)
		userGroup.GET("/get/vo", userController.GetUserVOByID)
		userGroup.POST("/delete", middleware.RequireAdmin(userService), userController.DeleteUser)
		userGroup.POST("/update", middleware.RequireAdmin(userService), userController.UpdateUser)
		userGroup.POST("/list/page/vo", middleware.RequireAdmin(userService), userController.ListUserVOByPage)

		apiKeyGroup := apiGroup.Group("/api/key")
		apiKeyGroup.Use(middleware.RequireLogin(userService))
		apiKeyGroup.POST("/create", apiKeyController.CreateApiKey)
		apiKeyGroup.GET("/list/my", apiKeyController.ListMyApiKeys)
		apiKeyGroup.POST("/revoke", apiKeyController.RevokeApiKey)

		providerGroup := apiGroup.Group("/provider")
		providerGroup.POST("/add", middleware.RequireAdmin(userService), providerController.AddProvider)
		providerGroup.POST("/delete", middleware.RequireAdmin(userService), providerController.DeleteProvider)
		providerGroup.POST("/update", middleware.RequireAdmin(userService), providerController.UpdateProvider)
		providerGroup.GET("/get/vo", providerController.GetProviderVOByID)
		providerGroup.POST("/list/page/vo", providerController.ListProviderVOByPage)
		providerGroup.GET("/list/vo", providerController.ListProviderVO)
		providerGroup.GET("/list/healthy", providerController.ListHealthyProviders)

		modelGroup := apiGroup.Group("/model")
		modelGroup.POST("/add", middleware.RequireAdmin(userService), modelController.AddModel)
		modelGroup.POST("/delete", middleware.RequireAdmin(userService), modelController.DeleteModel)
		modelGroup.POST("/update", middleware.RequireAdmin(userService), modelController.UpdateModel)
		modelGroup.GET("/get/vo", modelController.GetModelVOByID)
		modelGroup.POST("/list/page/vo", modelController.ListModelVOByPage)
		modelGroup.GET("/list/vo", modelController.ListModelVO)
		modelGroup.GET("/list/active", modelController.ListActiveModels)
		modelGroup.GET("/list/active/provider/:providerId", modelController.ListActiveModelsByProvider)
		modelGroup.GET("/list/active/type/:modelType", modelController.ListActiveModelsByType)

		statsGroup := apiGroup.Group("/stats")
		statsGroup.Use(middleware.RequireLogin(userService))
		statsGroup.GET("/my/tokens", statsController.GetMyTokenStats)
		statsGroup.GET("/my/logs", statsController.GetMyLogs)

		internalChatGroup := apiGroup.Group("/internal/chat")
		internalChatGroup.Use(middleware.RequireLogin(userService))
		internalChatGroup.POST("/completions", internalChatController.ChatCompletions)

		chatGroup := apiGroup.Group("/v1/chat")
		chatGroup.POST("/completions", chatController.ChatCompletions)
	}

	return engine, nil
}
