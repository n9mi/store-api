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

func TestCustomerAddressGetAll(t *testing.T) {
	testItems := map[string]TestItem{
		"get_all_success": {
			"response_code":   fiber.StatusOK,
			"response_status": "SUCCESS",
		},
	}

	for testName, testItem := range testItems {
		t.Run(testName, func(t *testing.T) {
			request := NewRequestWithToken(fiber.MethodGet, CustomerAddressURL, "", ExistingCustomer["token"])

			response, err := TestCfg.App.Test(request)
			require.Nil(t, err)
			require.Equal(t, testItem["response_code"], response.StatusCode)

			responseBodyByte, err := io.ReadAll(response.Body)
			require.Nil(t, err)

			responseBody := new(dto.Response[[]dto.CustomerAddressResponse])
			err = json.Unmarshal(responseBodyByte, responseBody)
			require.Nil(t, err)

			require.Equal(t, testItem["response_status"], responseBody.Status)
			require.True(t, len(responseBody.Messages["_success"]) > 0)
		})
	}
}

func TestCustomerAddressCreate(t *testing.T) {
	t.Run("create_success", func(t *testing.T) {
		testItem := TestItem{
			"request_street":     "street",
			"request_city":       "city",
			"request_province":   "province",
			"request_is_default": "true",
			"response_code":      fiber.StatusOK,
			"response_status":    "SUCCESS",
		}
		requestBody := fmt.Sprintf(`{"street":"%s","city":"%s","province":"%s","is_default":%s}`,
			testItem["request_street"],
			testItem["request_city"],
			testItem["request_province"],
			testItem["request_is_default"])
		request := NewRequestWithToken(fiber.MethodPost, CustomerAddressURL, requestBody, ExistingCustomer["token"])

		response, err := TestCfg.App.Test(request)
		require.Nil(t, err)
		require.Equal(t, testItem["response_code"], response.StatusCode)

		responseBodyByte, err := io.ReadAll(response.Body)
		require.Nil(t, err)

		responseBody := new(dto.Response[*dto.CustomerAddressResponse])
		err = json.Unmarshal(responseBodyByte, responseBody)
		require.Nil(t, err)

		require.Equal(t, testItem["response_status"], responseBody.Status)
		require.True(t, len(responseBody.Messages["_success"]) > 0)

		var custAddrEntity *entity.CustomerAddress
		TestCfg.DB.First(&custAddrEntity, "id = ?", responseBody.Data.ID)
		require.Equal(t, custAddrEntity.ID, responseBody.Data.ID)
		require.Equal(t, custAddrEntity.Street, responseBody.Data.Street)
		require.Equal(t, custAddrEntity.City, responseBody.Data.City)
		require.Equal(t, custAddrEntity.Province, responseBody.Data.Province)
		require.True(t, responseBody.Data.IsDefault)

		// make sure that only one default address
		var countDefault int64
		TestCfg.DB.Model(new(entity.CustomerAddress)).
			Where("user_id = ?", custAddrEntity.UserID).
			Where("is_default = ?", true).
			Count(&countDefault)
		require.Equal(t, int64(1), countDefault)
	})

	testItems := map[string]TestItem{
		"create_bad_request": {
			"request_street":     "",
			"request_city":       "",
			"request_province":   "",
			"request_is_default": "true",
			"request_token":      ExistingCustomer["token"],
			"response_code":      fiber.StatusBadRequest,
			"response_status":    "BAD_REQUEST",
		},
		"create_forbidden": {
			"request_street":     "street",
			"request_city":       "city",
			"request_province":   "province",
			"request_is_default": "true",
			"request_token":      ExistingMerchant["token"],
			"response_code":      fiber.StatusForbidden,
			"response_status":    "FORBIDDEN",
		},
		"create_unauthorized": {
			"request_street":     "street",
			"request_city":       "city",
			"request_province":   "province",
			"request_is_default": "true",
			"request_token":      "randomToken",
			"response_code":      fiber.StatusUnauthorized,
			"response_status":    "UNAUTHORIZED",
		},
	}

	for testName, testItem := range testItems {
		t.Run(testName, func(t *testing.T) {
			requestBody := fmt.Sprintf(`{"street":"%s","city":"%s","province":"%s","is_default":%s}`,
				testItem["request_street"],
				testItem["request_city"],
				testItem["request_province"],
				testItem["request_is_default"])
			request := NewRequestWithToken(fiber.MethodPost, CustomerAddressURL, requestBody, testItem["request_token"].(string))

			response, err := TestCfg.App.Test(request)
			require.Nil(t, err)
			require.Equal(t, testItem["response_code"], response.StatusCode)

			responseBodyByte, err := io.ReadAll(response.Body)
			require.Nil(t, err)

			responseBody := new(dto.Response[any])
			err = json.Unmarshal(responseBodyByte, responseBody)
			require.Nil(t, err)
			fmt.Println(responseBody)

			require.Equal(t, testItem["response_status"], responseBody.Status)
			require.True(t, len(responseBody.Messages) > 0)
		})
	}
}
