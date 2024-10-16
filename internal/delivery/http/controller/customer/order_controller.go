package customer

import (
	"store-api/internal/dto"
	"store-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type OrderController struct {
	Logger       *logrus.Logger
	OrderService service.OrderService
}

func NewOrderController(logger *logrus.Logger, service *service.Services) *OrderController {
	return &OrderController{
		Logger:       logger,
		OrderService: service.OrderService,
	}
}

func (ct *OrderController) Create(c *fiber.Ctx) error {
	request := new(dto.OrderCreateRequest)
	if err := c.BodyParser(request); err != nil {
		ct.Logger.Warnf("failed to parse request : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	userData, ok := c.Locals("AUTH_DATA").(*dto.AuthDTO)
	if !ok {
		ct.Logger.Warnf("failed to parse auth data")
		return fiber.NewError(fiber.StatusInternalServerError, "something wrong")
	}

	result, err := ct.OrderService.Create(c.UserContext(), request, userData.UserID)
	if err != nil {
		return err
	}

	resp := dto.Response[*dto.OrderCreateResponse]{
		Status: "SUCCESS",
		Data:   result,
		Messages: map[string]string{
			"_success": "success to create an order",
		},
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}
