package middleware

import (
	"store-api/internal/service"

	"github.com/casbin/casbin/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type Middlewares struct {
	AuthMiddleware func(c *fiber.Ctx) error
}

func Setup(viperCfg *viper.Viper,
	validate *validator.Validate,
	logger *logrus.Logger,
	enforcer *casbin.Enforcer,
	services *service.Services) *Middlewares {
	return &Middlewares{
		AuthMiddleware: NewAuthMiddleware(viperCfg, validate, logger, enforcer, services.JWTService, services.AuthService),
	}
}
