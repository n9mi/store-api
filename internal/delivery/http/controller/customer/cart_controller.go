package customer

import (
	"store-api/internal/dto"
	"store-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type CartController struct {
	Logger          *logrus.Logger
	CartItemService service.CartItemService
}

func NewCartController(logger *logrus.Logger, services *service.Services) *CartController {
	return &CartController{
		Logger:          logger,
		CartItemService: services.CartItemService,
	}
}

func (ct *CartController) Insert(c *fiber.Ctx) error {
	request := new(dto.CartItemRequest)
	if err := c.BodyParser(request); err != nil {
		ct.Logger.Warnf("failed to parse request : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	userData, ok := c.Locals("AUTH_DATA").(*dto.AuthDTO)
	if !ok {
		ct.Logger.Warnf("failed to parse auth data")
		return fiber.NewError(fiber.StatusInternalServerError, "something wrong")
	}

	if err := ct.CartItemService.Create(c.UserContext(), request, userData.UserID); err != nil {
		return err
	}

	resp := dto.Response[any]{
		Status: "SUCCESS",
		Messages: map[string]string{
			"_success": "success to add product into cart",
		},
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}
