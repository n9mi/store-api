package auth

import (
	"store-api/internal/dto"
	"store-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AuthController struct {
	Logger      *logrus.Logger
	AuthService service.AuthService
}

func NewAuthController(logger *logrus.Logger, services *service.Services) *AuthController {
	return &AuthController{
		Logger:      logger,
		AuthService: services.AuthService,
	}
}

func (ct *AuthController) Register(c *fiber.Ctx) error {
	request := new(dto.RegisterRequest)
	if err := c.BodyParser(request); err != nil {
		ct.Logger.Warnf("failed to parse request : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	if err := ct.AuthService.Register(c.UserContext(), request); err != nil {
		return err
	}

	response := dto.Response[any]{
		Status: "SUCCESS",
		Messages: map[string]string{
			"_success": "user registration success",
		},
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
