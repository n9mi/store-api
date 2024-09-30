package service

import (
	"store-api/internal/repository"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductService interface {
}

type ProductServiceImpl struct {
	DB                *gorm.DB
	Logger            *logrus.Logger
	ProductRepository *repository.ProductRepository
}
