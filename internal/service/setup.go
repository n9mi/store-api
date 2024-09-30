package service

import (
	"store-api/internal/repository"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type Services struct {
	JWTService  JWTService
	AuthService AuthService
}

func Setup(viperCfg *viper.Viper,
	db *gorm.DB,
	validator *validator.Validate,
	logger *logrus.Logger,
	repositories *repository.Repositories) *Services {
	return &Services{
		JWTService:  NewJwtService(viperCfg, db, logger, repositories),
		AuthService: NewAuthService(db, validator, logger, repositories),
	}
}
