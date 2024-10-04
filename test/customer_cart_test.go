package test

import (
	"encoding/json"
	"fmt"
	"io"
	"store-api/internal/dto"
	"store-api/internal/entity"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestCustomerCartInsert(t *testing.T) {
	var product entity.Product
	err := TestCfg.DB.First(&product).Error
	require.Nil(t, err)

	testItems := map[string]TestItem{
		"cart_insert_success": {
			"request_product_id": product.ID,
			"request_quantity":   product.Stock,
			"response_code":      fiber.StatusOK,
			"response_status":    "SUCCESS",
		},
		"cart_insert_update": {
			"request_product_id": product.ID,
			"request_quantity":   product.Stock,
			"response_code":      fiber.StatusOK,
			"response_status":    "SUCCESS",
		},
		"cart_insert_bad_request__product_id": {
			"request_product_id": "",
			"request_quantity":   product.Stock,
			"response_code":      fiber.StatusBadRequest,
			"response_status":    "BAD_REQUEST",
		},
		"cart_insert_bad_request__quantity": {
			"request_product_id": product.ID,
			"request_quantity":   -1,
			"response_code":      fiber.StatusBadRequest,
			"response_status":    "BAD_REQUEST",
		},
		"cart_insert_server_error__quantity": {
			"request_product_id": product.ID,
			"request_quantity":   100000,
			"response_code":      fiber.StatusInternalServerError,
			"response_status":    "INTERNAL_SERVER_ERROR",
		},
		"cart_insert_not_found": {
			"request_product_id": "someRandomId",
			"request_quantity":   product.Stock,
			"response_code":      fiber.StatusNotFound,
			"response_status":    "NOT_FOUND",
		},
	}
	for testName, testItem := range testItems {
		t.Run(testName, func(t *testing.T) {
			requestBody := fmt.Sprintf(`{"product_id":"%s","quantity":%d}`,
				testItem["request_product_id"],
				testItem["request_quantity"])
			request := NewRequestWithToken(fiber.MethodPost, CustomerInsertCartURL, requestBody, ExistingCustomer["token"])

			response, err := TestCfg.App.Test(request)
			require.Nil(t, err)
			require.Equal(t, testItem["response_code"], response.StatusCode)

			responseBodyByte, err := io.ReadAll(response.Body)
			require.Nil(t, err)

			responseBody := new(dto.Response[any])
			err = json.Unmarshal(responseBodyByte, responseBody)
			require.Nil(t, err)

			require.Equal(t, testItem["response_status"], responseBody.Status)

			if testItem["response_code"] == fiber.StatusOK {
				require.True(t, len(responseBody.Messages["_success"]) > 0)

				var cartItem entity.CartItem
				err := TestCfg.DB.Where("user_id = ?", ExistingCustomer["userID"]).
					Where("product_id = ?", product.ID).
					First(&cartItem).Error
				require.Nil(t, err)

				require.Equal(t, testItem["request_quantity"], cartItem.Quantity)
			}
		})
	}
}
