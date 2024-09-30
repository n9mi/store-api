package test

import (
	"store-api/config"
	"store-api/internal/repository"
)

var testCfg config.ConfigBootstrap

type testItem map[string]interface{}

type testResponse[T any] struct {
	Status     string            `json:"status"`
	Messages   map[string]string `json:"messages"`
	Data       T                 `json:"data"`
	Pagination *pagination       `json:"pagination"`
}

type pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPage   int `json:"total_page"`
	PageSize    int `json:"page_size"`
}

var repositories *repository.Repositories

var validCustomer = map[string]any{
	"name":     "test customer",
	"email":    "customer@test.com",
	"password": "password",
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

	repositories = repository.Setup()
}
