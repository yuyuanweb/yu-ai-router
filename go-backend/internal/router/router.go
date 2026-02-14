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
	}

	return engine, nil
}
