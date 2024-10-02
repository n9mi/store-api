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
		return fiber.NewError(fiber.StatusInternalServerError, "something wrong")
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

func (ct *AddressController) Create(c *fiber.Ctx) error {
	userData, ok := c.Locals("AUTH_DATA").(*dto.AuthDTO)
	if !ok {
		ct.Logger.Warnf("failed to parse auth data")
		return fiber.NewError(fiber.StatusInternalServerError, "something wrong")
	}

	request := new(dto.CustomerAddressRequest)
	if err := c.BodyParser(request); err != nil {
		ct.Logger.Warnf("failed to parse request : %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request")
	}

	result, err := ct.CustomerAddressService.Create(c.UserContext(), request, userData.UserID)
	if err != nil {
		return err
	}

	response := dto.Response[*dto.CustomerAddressResponse]{
		Status: "SUCCESS",
		Messages: map[string]string{
			"_success": "success creating user address",
		},
		Data: result,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
