package customhandler

import (
	"fmt"
	"store-api/internal/dto"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewCustomErrorHandler() func(*fiber.Ctx, error) error {
	return func(c *fiber.Ctx, err error) error {
		if err != nil {
			defResp := dto.Response[any]{
				Status:     "failed",
				Messages:   make(map[string]string),
				Data:       nil,
				Pagination: nil,
			}
			code := fiber.StatusInternalServerError

			if errConv, ok := err.(validator.ValidationErrors); ok {
				code = fiber.StatusBadRequest

				for _, errItem := range errConv {
					switch errItem.Tag() {
					case "required":
						defResp.Messages[errItem.Field()] = fmt.Sprintf("%s is required", errItem.Field())
					case "min":
						defResp.Messages[errItem.Field()] = fmt.Sprintf("%s should be more than %s characters", errItem.Field(), errItem.Param())
					case "max":
						defResp.Messages[errItem.Field()] = fmt.Sprintf("%s should be less than %s characters", errItem.Field(), errItem.Param())
					case "gte":
						defResp.Messages[errItem.Field()] = fmt.Sprintf("%s should be more than %s", errItem.Field(), errItem.Param())
					case "lte":
						defResp.Messages[errItem.Field()] = fmt.Sprintf("%s should be less than %s", errItem.Field(), errItem.Param())
					case "email":
						defResp.Messages[errItem.Field()] = fmt.Sprintf("%s should be a valid email", errItem.Field())
					}
				}
			} else if errConv, ok := err.(*fiber.Error); ok {
				code = errConv.Code

				defResp.Messages["_error"] = errConv.Message
			} else {
				defResp.Messages["_error"] = "something wrong"
			}

			return c.Status(code).JSON(defResp)
		}

		return nil
	}
}
