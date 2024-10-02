package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NewRequest(method string, url string, requestBody string) *http.Request {
	request := httptest.NewRequest(method, url, strings.NewReader(requestBody))

	request.Header.Add(fiber.HeaderContentType, fiber.MIMEApplicationJSON)

	return request
}

func NewRequestWithToken(method string, url string, requestBody string, token string) *http.Request {
	request := NewRequest(method, url, requestBody)
	bearerAuth := fmt.Sprintf("Bearer %s", token)

	request.Header.Add("Authorization", bearerAuth)

	return request
}
