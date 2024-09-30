package config

import (
	customhandler "store-api/internal/custom_handler"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func NewFiber(viperCfg *viper.Viper) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:      viperCfg.GetString("APP_NAME"),
		ErrorHandler: customhandler.NewCustomErrorHandler(),
	})

	return app
}
