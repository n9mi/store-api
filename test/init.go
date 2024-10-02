package test

import (
	"encoding/json"
	"fmt"
	"io"
	"store-api/config"
	"store-api/internal/dto"

	"github.com/gofiber/fiber/v2"
)

var TestCfg config.ConfigBootstrap

type TestItem map[string]interface{}

const (
	RegisterURL        = "/api/v1/auth/register"
	LoginURL           = "/api/v1/auth/login"
	CustomerProductURL = "/api/v1/customer/products"
	CustomerAddressURL = "/api/v1/customer/addresses"
)

var ValidCustomer = map[string]any{
	"name":     "test customer",
	"email":    "customer@test.com",
	"password": "password",
}

var ExistingCustomer = map[string]string{
	"name":     "Customer 1",
	"email":    "customer1@test.com",
	"password": "password",
	"token":    "",
}

func SetupAuthenticatedCustomer() {
	requestBody := fmt.Sprintf(`{"email":"%s","password":"%s","as_role":"%s"}`,
		ExistingCustomer["email"],
		ExistingCustomer["password"],
		"customer")
	request := NewRequest(fiber.MethodPost, LoginURL, requestBody)

	response, _ := TestCfg.App.Test(request)
	responseBodyByte, _ := io.ReadAll(response.Body)
	responseBody := new(dto.Response[dto.LoginDTO])
	json.Unmarshal(responseBodyByte, responseBody)

	ExistingCustomer["token"] = responseBody.Data.AccessToken
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

	TestCfg = config.ConfigBootstrap{
		App:       app,
		Logger:    logger,
		DB:        db,
		Validator: validator,
		ViperCfg:  viperCfg,
		Enforcer:  enforcer,
	}
	config.Bootstrap(&TestCfg)

	SetupAuthenticatedCustomer()
}
