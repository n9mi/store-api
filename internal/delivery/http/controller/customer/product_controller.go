package customer

import (
	"store-api/internal/dto"

	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
}

func NewProductController() *ProductController {
	return &ProductController{}
}

func (ct *ProductController) GetAll(c *fiber.Ctx) error {

	response := dto.Response[any]{
		Status: "SUCCESS",
	}
	return c.Status(fiber.StatusOK).JSON(response)
}
