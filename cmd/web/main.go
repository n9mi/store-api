package main

import (
	"fmt"
	"store-api/config"
)

func main() {
	viperCfg, err := config.NewViper()
	if err != nil {
		panic(err)
	}
	app := config.NewFiber(viperCfg)
	logger := config.NewLogrus(viperCfg)
	db := config.NewDatabase(viperCfg, logger)
	validator := config.NewValidator(viperCfg)
	enforcer := config.NewEnforcer(viperCfg, logger)

	cfg := config.ConfigBootstrap{
		App:       app,
		Logger:    logger,
		DB:        db,
		Validator: validator,
		ViperCfg:  viperCfg,
		Enforcer:  enforcer,
	}
	config.Bootstrap(&cfg)

	webPort := viperCfg.GetInt("APP_PORT")
	if err := app.Listen(fmt.Sprintf(":%d", webPort)); err != nil {
		logger.Fatalf("failed to run the app : %+v", err)
	}
}
