package test

import (
	"encoding/json"
	"fmt"
	"io"
	"store-api/internal/dto"
	"store-api/internal/entity"
	"testing"
	"time"

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
			request := NewRequestWithToken(fiber.MethodPost, CustomerCartURL, requestBody, ExistingCustomer["token"])

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

func deleteAllCartItems() {
	TestCfg.DB.Where("1 = 1").Delete(new(entity.CartItem))
}

func TestCustomerGetItemsCart(t *testing.T) {
	t.Run("cart_get_items_success", func(t *testing.T) {
		deleteAllCartItems()

		numProductLim := 5

		var products []entity.Product
		err := TestCfg.DB.Limit(numProductLim).Find(&products).Error
		require.Nil(t, err)

		cartItems := make([]entity.CartItem, numProductLim)
		for i := range 5 {
			cartItems[i] = entity.CartItem{
				UserID:    ExistingCustomer["userID"],
				ProductID: products[i].ID,
				Quantity:  products[i].Stock,
				CreatedAt: time.Now(),
			}
		}
		err = TestCfg.DB.Create(cartItems).Error
		require.Nil(t, err)

		request := NewRequestWithToken(fiber.MethodGet, CustomerCartURL, "", ExistingCustomer["token"])

		response, err := TestCfg.App.Test(request)
		require.Nil(t, err)
		require.Equal(t, fiber.StatusOK, response.StatusCode)

		responseBodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)

		responseBody := new(dto.Response[[]dto.CartItemResponse])
		err = json.Unmarshal(responseBodyByte, responseBody)
		require.Nil(t, err)

		require.Equal(t, "SUCCESS", responseBody.Status)

		require.True(t, len(responseBody.Messages["_success"]) > 0)
		require.Equal(t, len(responseBody.Data), numProductLim)
		for i := range numProductLim {
			require.Equal(t, responseBody.Data[i].Product.ID, cartItems[i].ProductID)
			require.Equal(t, responseBody.Data[i].Quantity, cartItems[i].Quantity)
			require.Equal(t, responseBody.Data[i].TotalPrice, float64(cartItems[i].Quantity)*products[i].PriceIdr)
		}

		deleteAllCartItems()
	})
}

func TestCustomerDeleteItemCart(t *testing.T) {
	var product entity.Product
	err := TestCfg.DB.First(&product).Error
	require.Nil(t, err)

	cartItem := entity.CartItem{
		UserID:    ExistingCustomer["userID"],
		ProductID: product.ID,
		Quantity:  1,
		CreatedAt: time.Now(),
	}
	err = TestCfg.DB.Create(&cartItem).Error
	require.Nil(t, err)

	testItems := map[string]TestItem{
		"cart_delete_success": {
			"request_product_id": product.ID,
			"response_code":      fiber.StatusOK,
			"response_status":    "SUCCESS",
		},
		"cart_delete_product_not_found": {
			"request_product_id": product.ID,
			"response_code":      fiber.StatusNotFound,
			"response_status":    "NOT_FOUND",
		},
	}
	for testName, testItem := range testItems {
		t.Run(testName, func(t *testing.T) {
			request := NewRequestWithToken(fiber.MethodDelete,
				CustomerCartURL+"/"+testItem["request_product_id"].(string), "",
				ExistingCustomer["token"])

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
				res := TestCfg.DB.Where("user_id = ?", ExistingCustomer["userID"]).
					Where("product_id = ?", testItem["request_product_id"]).Find(new(entity.CartItem))
				require.Nil(t, res.Error)
				require.Equal(t, 0, int(res.RowsAffected))
			} else {
				require.True(t, len(responseBody.Messages["_error"]) > 0)
			}
		})
	}
}
