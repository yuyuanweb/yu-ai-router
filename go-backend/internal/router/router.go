package router

import (
	"strconv"
	"time"

	"github.com/gin-contrib/sessions"
	redisStore "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"

	"github.com/yupi/airouter/go-backend/internal/config"
	"github.com/yupi/airouter/go-backend/internal/controller"
	"github.com/yupi/airouter/go-backend/internal/middleware"
	"github.com/yupi/airouter/go-backend/internal/service"
)

const sessionPoolSize = 10
const (
	apiChatRateLimitPerSecond      = 60
	internalChatRateLimitPerSecond = 30
	imageRateLimitPerMinute        = 10
)

func New(
	cfg *config.Config,
	healthController *controller.HealthController,
	userController *controller.UserController,
	apiKeyController *controller.ApiKeyController,
	providerController *controller.ProviderController,
	modelController *controller.ModelController,
	blacklistController *controller.BlacklistController,
	chatController *controller.ChatController,
	internalChatController *controller.InternalChatController,
	statsController *controller.StatsController,
	pluginController *controller.PluginController,
	userProviderKeyController *controller.UserProviderKeyController,
	rechargeController *controller.RechargeController,
	stripeWebhookController *controller.StripeWebhookController,
	balanceController *controller.BalanceController,
	imageController *controller.ImageController,
	userService *service.UserService,
	blacklistService *service.BlacklistService,
	rateLimitService *service.RateLimitService,
) (*gin.Engine, error) {
	engine := gin.New()
	engine.Use(middleware.IPBlacklistFilter(blacklistService))
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
		userGroup.GET("/quota/my", middleware.RequireLogin(userService), userController.GetMyQuota)
		userGroup.POST("/quota/set", middleware.RequireAdmin(userService), userController.SetUserQuota)
		userGroup.POST("/quota/reset", middleware.RequireAdmin(userService), userController.ResetUserQuota)
		userGroup.POST("/disable", middleware.RequireAdmin(userService), userController.DisableUser)
		userGroup.POST("/enable", middleware.RequireAdmin(userService), userController.EnableUser)
		userGroup.GET("/analysis", middleware.RequireAdmin(userService), userController.GetUserAnalysis)

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

		blacklistGroup := apiGroup.Group("/admin/blacklist")
		blacklistGroup.Use(middleware.RequireAdmin(userService))
		blacklistGroup.GET("/list", blacklistController.List)
		blacklistGroup.POST("/add", blacklistController.Add)
		blacklistGroup.POST("/remove", blacklistController.Remove)
		blacklistGroup.GET("/check", blacklistController.Check)
		blacklistGroup.GET("/count", blacklistController.Count)

		statsGroup := apiGroup.Group("/stats")
		statsGroup.Use(middleware.RequireLogin(userService))
		statsGroup.GET("/my/tokens", statsController.GetMyTokenStats)
		statsGroup.GET("/my/logs", statsController.GetMyLogs)
		statsGroup.GET("/my/cost", statsController.GetMyCostStats)
		statsGroup.GET("/my/summary", statsController.GetMySummaryStats)
		statsGroup.GET("/my/daily", statsController.GetMyDailyStats)
		statsGroup.POST("/history/my/page", statsController.PageMyHistory)
		statsGroup.GET("/history/detail", statsController.GetHistoryDetail)
		statsGroup.POST("/history/page", middleware.RequireAdmin(userService), statsController.PageHistory)

		pluginGroup := apiGroup.Group("/plugin")
		pluginGroup.GET("/list/enabled", pluginController.ListEnabledPlugins)
		pluginGroup.GET("/get", pluginController.GetPlugin)
		pluginGroup.POST("/execute", pluginController.ExecutePlugin)
		pluginGroup.GET("/list", middleware.RequireAdmin(userService), pluginController.ListPlugins)
		pluginGroup.POST("/update", middleware.RequireAdmin(userService), pluginController.UpdatePlugin)
		pluginGroup.POST("/enable", middleware.RequireAdmin(userService), pluginController.EnablePlugin)
		pluginGroup.POST("/disable", middleware.RequireAdmin(userService), pluginController.DisablePlugin)
		pluginGroup.POST("/reload", middleware.RequireAdmin(userService), pluginController.ReloadPlugin)
		pluginGroup.POST("/reload/all", middleware.RequireAdmin(userService), pluginController.ReloadAllPlugins)

		byokGroup := apiGroup.Group("/byok")
		byokGroup.Use(middleware.RequireLogin(userService))
		byokGroup.POST("/add", userProviderKeyController.AddUserProviderKey)
		byokGroup.POST("/update", userProviderKeyController.UpdateUserProviderKey)
		byokGroup.POST("/delete", userProviderKeyController.DeleteUserProviderKey)
		byokGroup.GET("/my/list", userProviderKeyController.ListMyProviderKeys)

		rechargeGroup := apiGroup.Group("/recharge")
		rechargeGroup.GET("/stripe/success", rechargeController.StripeSuccess)
		rechargeGroup.GET("/stripe/cancel", rechargeController.StripeCancel)
		rechargeGroup.Use(middleware.RequireLogin(userService))
		rechargeGroup.POST("/stripe/create", rechargeController.CreateStripeRecharge)
		rechargeGroup.GET("/list/my", rechargeController.GetMyRechargeRecords)

		balanceGroup := apiGroup.Group("/balance")
		balanceGroup.Use(middleware.RequireLogin(userService))
		balanceGroup.GET("/my", balanceController.GetMyBalance)
		balanceGroup.GET("/billing/my", balanceController.GetMyBillingRecords)

		webhookGroup := apiGroup.Group("/webhook")
		webhookGroup.POST("/stripe", stripeWebhookController.HandleStripeWebhook)

		imageGroup := apiGroup.Group("/v1/images")
		imageGroup.Use(middleware.RateLimit(rateLimitService, middleware.RateLimitTypeIP, imageRateLimitPerMinute, time.Minute))
		imageGroup.POST("/generations", imageController.GenerateImage)
		imageGroup.GET("/my/records", middleware.RequireLogin(userService), imageController.GetMyRecords)

		internalChatGroup := apiGroup.Group("/internal/chat")
		internalChatGroup.Use(middleware.RequireLogin(userService))
		internalChatGroup.Use(middleware.RateLimit(rateLimitService, middleware.RateLimitTypeIP, internalChatRateLimitPerSecond, time.Second))
		internalChatGroup.POST("/completions", internalChatController.ChatCompletions)
		internalChatGroup.POST("/completions/upload", internalChatController.ChatCompletionsWithFile)

		chatGroup := apiGroup.Group("/v1/chat")
		chatGroup.Use(middleware.RateLimit(rateLimitService, middleware.RateLimitTypeAPIKey, apiChatRateLimitPerSecond, time.Second))
		chatGroup.POST("/completions", chatController.ChatCompletions)
	}

	return engine, nil
}
