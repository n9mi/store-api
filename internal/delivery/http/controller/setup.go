package controller

import (
	"store-api/internal/delivery/http/controller/auth"
	"store-api/internal/service"

	"github.com/sirupsen/logrus"
)

type ControllerSetup struct {
	AuthController *auth.AuthController
}

func Setup(logger *logrus.Logger, services *service.Services) *ControllerSetup {
	return &ControllerSetup{
		AuthController: auth.NewAuthController(logger, services),
	}
}
