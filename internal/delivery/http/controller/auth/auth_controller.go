package auth

import (
	"store-api/internal/dto"
	"store-api/internal/service"
	"time"

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

func (ct *AuthController) Login(c *fiber.Ctx) error {
	request := new(dto.LoginRequest)
	if err := c.BodyParser(request); err != nil {
		ct.Logger.Warnf("failed to parse request : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	loginDTO, err := ct.AuthService.Login(c.UserContext(), request)
	if err != nil {
		return err
	}

	cookie := fiber.Cookie{
		Name:     "REFRESH_TOKEN",
		Value:    loginDTO.RefreshToken,
		Expires:  time.Unix(loginDTO.RefreshTokenExpiredUnix, 0),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)

	response := dto.Response[*dto.LoginDTO]{
		Status: "SUCCESS",
		Messages: map[string]string{
			"_success": "user login success",
		},
		Data: loginDTO,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
