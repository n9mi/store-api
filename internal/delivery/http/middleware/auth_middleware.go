package middleware

import (
	"store-api/internal/service"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewAuthMiddleware(viperCfg *viper.Viper,
	validate *validator.Validate,
	logger *logrus.Logger,
	enforcer *casbin.Enforcer,
	JWTService service.JWTService,
	authService service.AuthService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization", "")

		if !strings.Contains(header, "Bearer ") {
			log.Warnf("authorization token hasn't been provided")
			return fiber.NewError(fiber.StatusUnauthorized, "empty_authorization")
		}

		accessToken := strings.Replace(header, "Bearer ", "", -1)
		userAuthData, err := JWTService.ParseAccessToken(accessToken)
		if err != nil {
			return err
		}

		if err := authService.ValidateAuthData(userAuthData); err != nil {
			return err
		}

		enforce, err := enforcer.Enforce(userAuthData.UserCurrentRole, c.Path(), c.Method())
		if err != nil {
			log.Warnf("failed to enforce : %+v", err)
			return fiber.NewError(fiber.StatusInternalServerError, "something wrong")
		}
		if !enforce {
			log.Warn("forbidden : current user role does'nt match")
			return fiber.NewError(fiber.StatusForbidden, "invalid_user_role")
		}

		c.Locals("AUTH_DATA", userAuthData)
		return c.Next()
	}
}
