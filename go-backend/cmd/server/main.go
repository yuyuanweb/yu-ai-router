package main

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/yupi/airouter/go-backend/internal/config"
	"github.com/yupi/airouter/go-backend/internal/controller"
	"github.com/yupi/airouter/go-backend/internal/repository"
	"github.com/yupi/airouter/go-backend/internal/router"
	"github.com/yupi/airouter/go-backend/internal/service"
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
	userService := service.NewUserService(userRepo)
	healthController := controller.NewHealthController()
	userController := controller.NewUserController(userService)

	engine, err := router.New(cfg, healthController, userController, userService)
	if err != nil {
		log.Fatalf("build router failed: %v", err)
	}

	if err = engine.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server run failed: %v", err)
	}
}
