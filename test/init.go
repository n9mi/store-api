package test

import (
	"encoding/json"
	"fmt"
	"io"
	"store-api/config"
	"store-api/internal/dto"

	"github.com/gofiber/fiber/v2"
)

var testCfg config.ConfigBootstrap

type testItem map[string]interface{}

const (
	registerURL        = "/api/v1/auth/register"
	loginURL           = "/api/v1/auth/login"
	customerProductURL = "/api/v1/customer/products"
)

var validCustomer = map[string]any{
	"name":     "test customer",
	"email":    "customer@test.com",
	"password": "password",
}

var existingCustomer = map[string]string{
	"name":     "Customer 1",
	"email":    "customer1@test.com",
	"password": "password",
	"token":    "",
}

func SetupAuthenticatedCustomer() {
	requestBody := fmt.Sprintf(`{"email":"%s","password":"%s","as_role":"%s"}`,
		existingCustomer["email"],
		existingCustomer["password"],
		"customer")
	request := newRequest(fiber.MethodPost, loginURL, requestBody)

	response, _ := testCfg.App.Test(request)
	responseBodyByte, _ := io.ReadAll(response.Body)
	responseBody := new(dto.Response[dto.LoginDTO])
	json.Unmarshal(responseBodyByte, responseBody)

	existingCustomer["token"] = responseBody.Data.AccessToken
}

func init() {
	viperCfg, err := config.NewViper()
	if err != nil {
		panic(err)
	}
	app := config.NewFiber(viperCfg)
	logger := config.NewLogrus(viperCfg)
	db, err := config.NewDatabase(viperCfg, logger)
	if err != nil {
		logger.Fatalf("failed to create database : %+v", err)
	}
	validator := config.NewValidator(viperCfg)
	enforcer := config.NewTestEnforcer(logger)

	testCfg = config.ConfigBootstrap{
		App:       app,
		Logger:    logger,
		DB:        db,
		Validator: validator,
		ViperCfg:  viperCfg,
		Enforcer:  enforcer,
	}
	config.Bootstrap(&testCfg)

	SetupAuthenticatedCustomer()
}
