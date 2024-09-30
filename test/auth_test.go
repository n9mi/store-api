package test

import (
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

const (
	registerURL = "/api/v1/auth/register"
)

func TestRegister(t *testing.T) {
	testItems := map[string]testItem{
		"register_success": {
			"request_name":     "test customer",
			"request_email":    "customer@test.com",
			"request_password": "password",
			"request_as_role":  "customer",
			"response_status":  "success",
		},
	}

	for testName, testItem := range testItems {
		t.Run(testName, func(t *testing.T) {
			requestBody := fmt.Sprintf(`{"name":"%s","email":"%s","password":"%s","as_role":"%s"}`,
				testItem["request_name"],
				testItem["request_email"],
				testItem["request_password"],
				testItem["request_as_role"])
			request := newRequest(fiber.MethodPost, registerURL, requestBody)

			response, err := testCfg.App.Test(request)
			require.Nil(t, err)
			require.Equal(t, fiber.StatusOK, response.StatusCode)

			responseBodyByte, err := io.ReadAll(response.Body)
			require.Nil(t, err)

			responseBody := new(testResponse[any])
			err = json.Unmarshal(responseBodyByte, responseBody)
			require.Nil(t, err)

			require.Equal(t, testItem["response_status"], responseBody.Status)
			require.True(t, len(responseBody.Messages) > 0)
		})
	}
}
