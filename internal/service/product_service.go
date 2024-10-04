package service

import (
	"context"
	"store-api/internal/dto"
	"store-api/internal/repository"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ProductService interface {
	FindAll(ctx context.Context, request *dto.FindAndSearchProductRequest) ([]dto.ProductDTO, *dto.Pagination, error)
}

type ProductServiceImpl struct {
	DB                *gorm.DB
	Logger            *logrus.Logger
	ProductRepository *repository.ProductRepository
}

func NewProductService(db *gorm.DB, logger *logrus.Logger, repositories *repository.Repositories) *ProductServiceImpl {
	return &ProductServiceImpl{
		DB:                db,
		Logger:            logger,
		ProductRepository: repositories.ProductRepository,
	}
}

func (s *ProductServiceImpl) FindAll(ctx context.Context, request *dto.FindAndSearchProductRequest) ([]dto.ProductDTO, *dto.Pagination, error) {
	res, pagination, err := s.ProductRepository.FindAll(s.DB, *request)
	if err != nil {
		s.Logger.Warnf("error getting products : %+v", err)
		return nil, nil, err
	}

	return res, pagination, nil
}
