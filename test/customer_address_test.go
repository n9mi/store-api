package test

import (
	"encoding/json"
	"io"
	"store-api/internal/dto"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestCustomerProducts(t *testing.T) {
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
