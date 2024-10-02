package service

import (
	"context"
	"store-api/internal/dto"
	"store-api/internal/repository"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CustomerAddressService interface {
	FindAll(ctx context.Context, userID string) ([]dto.CustomerAddressResponse, error)
}

type CustomerAddressServiceImpl struct {
	DB                       *gorm.DB
	Logger                   *logrus.Logger
	CustomerAddressRepsitory *repository.CustomerAddressRepository
}

func NewCustomerAddressService(db *gorm.DB, logger *logrus.Logger, repositories repository.Repositories) *CustomerAddressServiceImpl {
	return &CustomerAddressServiceImpl{
		DB:                       db,
		Logger:                   logger,
		CustomerAddressRepsitory: repositories.CustomerAddressRepository,
	}
}

func (s *CustomerAddressServiceImpl) FindAll(ctx context.Context, userID string) ([]dto.CustomerAddressResponse, error) {
	custAddress, err := s.CustomerAddressRepsitory.FindByUserID(s.DB, userID)
	if err != nil {
		return nil, err
	}

	res := make([]dto.CustomerAddressResponse, len(custAddress))
	for i, e := range custAddress {
		res[i] = dto.CustomerAddressResponse{
			ID:        e.ID,
			Street:    e.Street,
			City:      e.City,
			Province:  e.Province,
			IsDefault: e.IsDefault,
		}
	}

	return res, nil
}
