package controller

import (
	"store-api/internal/delivery/http/controller/auth"
	"store-api/internal/delivery/http/controller/customer"
	"store-api/internal/service"

	"github.com/sirupsen/logrus"
)

type Controllers struct {
	AuthController            *auth.AuthController
	ProductController         *customer.ProductController
	CustomerAddressController *customer.AddressController
	CartController            *customer.CartController
	OrderController           *customer.OrderController
}

func Setup(logger *logrus.Logger, services *service.Services) *Controllers {
	return &Controllers{
		AuthController:            auth.NewAuthController(logger, services),
		ProductController:         customer.NewProductController(logger, services),
		CustomerAddressController: customer.NewAddressController(logger, services),
		CartController:            customer.NewCartController(logger, services),
		OrderController:           customer.NewOrderController(logger, services),
	}
}
