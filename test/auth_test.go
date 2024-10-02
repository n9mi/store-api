package test

import (
	"encoding/json"
	"fmt"
	"io"
	"store-api/internal/dto"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	testItems := map[string]TestItem{
		"register_success": {
			"request_name":     ValidCustomer["name"],
			"request_email":    ValidCustomer["email"],
			"request_password": ValidCustomer["password"],
			"request_as_role":  "customer",
			"response_code":    fiber.StatusOK,
			"response_status":  "SUCCESS",
		},
		"register_validation": {
			"request_name":     "",
			"request_email":    "",
			"request_password": "",
			"request_as_role":  "",
			"response_code":    fiber.StatusBadRequest,
			"response_status":  "BAD_REQUEST",
		},
		"register_duplicate": {
			"request_name":     ExistingCustomer["name"],
			"request_email":    ExistingCustomer["email"],
			"request_password": ExistingCustomer["password"],
			"request_as_role":  "customer",
			"response_code":    fiber.StatusConflict,
			"response_status":  "CONFLICT",
		},
	}

	for testName, testItem := range testItems {
		t.Run(testName, func(t *testing.T) {
			requestBody := fmt.Sprintf(`{"name":"%s","email":"%s","password":"%s","as_role":"%s"}`,
				testItem["request_name"],
				testItem["request_email"],
				testItem["request_password"],
				testItem["request_as_role"])
			request := NewRequest(fiber.MethodPost, RegisterURL, requestBody)

			response, err := TestCfg.App.Test(request)
			require.Nil(t, err)
			require.Equal(t, testItem["response_code"], response.StatusCode)

			responseBodyByte, err := io.ReadAll(response.Body)
			require.Nil(t, err)

			responseBody := new(dto.Response[any])
			err = json.Unmarshal(responseBodyByte, responseBody)
			require.Nil(t, err)

			require.Equal(t, testItem["response_status"], responseBody.Status)
			require.True(t, len(responseBody.Messages) > 0)
		})
	}
}

func TestLogin(t *testing.T) {
	testItems := map[string]TestItem{
		"login_success": {
			"request_email":    ValidCustomer["email"],
			"request_password": ValidCustomer["password"],
			"request_as_role":  "customer",
			"response_code":    fiber.StatusOK,
			"response_status":  "SUCCESS",
		},
		"login_bad_request__empty_email_password": {
			"request_email":    "",
			"request_password": "",
			"request_as_role":  "customer",
			"response_code":    fiber.StatusBadRequest,
			"response_status":  "BAD_REQUEST",
		},
		"login_unauthorized__empty_as_role": {
			"request_email":    ValidCustomer["email"],
			"request_password": ValidCustomer["password"],
			"request_as_role":  "",
			"response_code":    fiber.StatusUnauthorized,
			"response_status":  "UNAUTHORIZED",
		},
		"login_unauthorized__invalid_credentials": {
			"request_email":    "email@email.com",
			"request_password": "passwords",
			"request_as_role":  "customer",
			"response_code":    fiber.StatusUnauthorized,
			"response_status":  "UNAUTHORIZED",
		},
	}

	for testName, testItem := range testItems {
		t.Run(testName, func(t *testing.T) {
			requestBody := fmt.Sprintf(`{"email":"%s","password":"%s","as_role":"%s"}`,
				testItem["request_email"],
				testItem["request_password"],
				testItem["request_as_role"])
			request := NewRequest(fiber.MethodPost, LoginURL, requestBody)

			response, err := TestCfg.App.Test(request)
			require.Nil(t, err)
			require.Equal(t, testItem["response_code"], response.StatusCode)

			responseBodyByte, err := io.ReadAll(response.Body)
			require.Nil(t, err)

			responseBody := new(dto.Response[dto.LoginDTO])
			err = json.Unmarshal(responseBodyByte, responseBody)
			require.Nil(t, err)

			require.Equal(t, testItem["response_status"], responseBody.Status)

			switch testItem["response_code"] {
			case fiber.StatusOK:
				require.True(t, len(responseBody.Messages["_success"]) > 0)
				require.True(t, len(responseBody.Data.AccessToken) > 0)
			case fiber.StatusBadRequest:
				require.True(t, len(responseBody.Messages["email"]) > 0)
				require.True(t, len(responseBody.Messages["password"]) > 0)
			case fiber.StatusUnauthorized:
				require.True(t, len(responseBody.Messages["_error"]) > 0)
			}
		})
	}
}
