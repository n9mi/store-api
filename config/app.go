package config

import (
	"store-api/database"
	"store-api/internal/delivery/http/controller"
	"store-api/internal/delivery/http/route"
	"store-api/internal/repository"
	"store-api/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ConfigBootstrap struct {
	ViperCfg  *viper.Viper
	App       *fiber.App
	Logger    *logrus.Logger
	DB        *gorm.DB
	Validator *validator.Validate
}

func Bootstrap(cfg *ConfigBootstrap) {
	repositories := repository.Setup()
	services := service.Setup(cfg.ViperCfg, cfg.DB, cfg.Validator, cfg.Logger, repositories)
	controllers := controller.Setup(cfg.Logger, services)
	router := route.NewRouter(cfg.App, controllers)

	router.Setup()

	if err := database.DropTable(cfg.DB); err != nil {
		cfg.Logger.Fatalf("failed to drop tables : %+v", err)
	}

	if err := database.CreateTable(cfg.DB); err != nil {
		cfg.Logger.Fatalf("failed to create tables : %+v", err)
	}

	if err := database.Seed(cfg.DB, repositories); err != nil {
		cfg.Logger.Fatalf("failed to seed into database : %+v", err)
	}
}
