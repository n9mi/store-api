package customhandler

import (
	"fmt"
	"store-api/internal/dto"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func NewCustomErrorHandler() func(*fiber.Ctx, error) error {
	return func(c *fiber.Ctx, err error) error {
		if err != nil {
			defResp := dto.Response[any]{
				Status:     "INTERNAL_SERVER_ERROR",
				Messages:   make(map[string]string),
				Data:       nil,
				Pagination: nil,
			}
			code := fiber.StatusInternalServerError

			if errConv, ok := err.(validator.ValidationErrors); ok {
				code = fiber.StatusBadRequest
				defResp.Status = "BAD_REQUEST"
				for _, errItem := range errConv {
					errItemLow := strings.ToLower(errItem.Field())
					switch errItem.Tag() {
					case "required":
						defResp.Messages[errItemLow] = fmt.Sprintf("%s is required", errItemLow)
					case "min":
						defResp.Messages[errItemLow] = fmt.Sprintf("%s should be more than %s characters", errItemLow, errItem.Param())
					case "max":
						defResp.Messages[errItemLow] = fmt.Sprintf("%s should be less than %s characters", errItemLow, errItem.Param())
					case "gte":
						defResp.Messages[errItemLow] = fmt.Sprintf("%s should be more than %s", errItemLow, errItem.Param())
					case "lte":
						defResp.Messages[errItemLow] = fmt.Sprintf("%s should be less than %s", errItemLow, errItem.Param())
					case "email":
						defResp.Messages[errItemLow] = fmt.Sprintf("%s should be a valid email", errItemLow)
					}
				}
			} else if errConv, ok := err.(*fiber.Error); ok {
				code = errConv.Code
				switch code {
				case fiber.StatusBadRequest:
					defResp.Status = "BAD_REQUEST"
				case fiber.StatusConflict:
					defResp.Status = "CONFLICT"
				case fiber.StatusUnauthorized:
					defResp.Status = "UNAUTHORIZED"
				default:
					defResp.Status = "FAILED"
				}
				defResp.Messages["_error"] = errConv.Message
			} else {
				defResp.Messages["_error"] = "something wrong"
			}

			return c.Status(code).JSON(defResp)
		}

		return nil
	}
}
