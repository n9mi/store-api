package customer

import (
	"store-api/internal/dto"
	"store-api/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type AddressController struct {
	Logger                 *logrus.Logger
	CustomerAddressService service.CustomerAddressService
}

func NewAddressController(logger *logrus.Logger, services *service.Services) *AddressController {
	return &AddressController{
		Logger:                 logger,
		CustomerAddressService: services.CustomerAddressService,
	}
}

func (ct *AddressController) GetAll(c *fiber.Ctx) error {
	userData, ok := c.Locals("AUTH_DATA").(*dto.AuthDTO)
	if !ok {
		ct.Logger.Warnf("failed to parse auth data")
	}

	result, err := ct.CustomerAddressService.FindAll(c.UserContext(), userData.UserID)
	if err != nil {
		return err
	}

	response := dto.Response[[]dto.CustomerAddressResponse]{
		Status: "SUCCESS",
		Messages: map[string]string{
			"_success": "success getting all user address",
		},
		Data: result,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
