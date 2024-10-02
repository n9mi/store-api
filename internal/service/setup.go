package service

import (
	"store-api/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Services struct {
	AuthService            AuthService
	JWTService             JWTService
	ProductService         ProductService
	CustomerAddressService CustomerAddressService
}

func Setup(viperCfg *viper.Viper,
	db *gorm.DB,
	validator *validator.Validate,
	logger *logrus.Logger,
	repositories *repository.Repositories) *Services {
	JWTService := NewJwtService(viperCfg, db, logger, repositories)

	return &Services{
		AuthService:            NewAuthService(db, validator, logger, JWTService, repositories),
		JWTService:             JWTService,
		ProductService:         NewProductService(db, logger, repositories),
		CustomerAddressService: NewCustomerAddressService(db, logger, *repositories),
	}
}
