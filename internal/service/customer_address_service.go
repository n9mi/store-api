package service

import (
	"context"
	"store-api/internal/dto"
	"store-api/internal/entity"
	"store-api/internal/repository"
	"store-api/util"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CustomerAddressService interface {
	FindAll(ctx context.Context, userID string) ([]dto.CustomerAddressResponse, error)
	Create(ctx context.Context, request *dto.CustomerAddressRequest, userID string) (*dto.CustomerAddressResponse, error)
}

type CustomerAddressServiceImpl struct {
	DB                       *gorm.DB
	Logger                   *logrus.Logger
	Validator                *validator.Validate
	CustomerAddressRepsitory *repository.CustomerAddressRepository
}

func NewCustomerAddressService(db *gorm.DB, logger *logrus.Logger, validator *validator.Validate, repositories repository.Repositories) *CustomerAddressServiceImpl {
	return &CustomerAddressServiceImpl{
		DB:                       db,
		Logger:                   logger,
		Validator:                validator,
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

func (s *CustomerAddressServiceImpl) Create(ctx context.Context, request *dto.CustomerAddressRequest,
	userID string) (*dto.CustomerAddressResponse, error) {
	if err := s.Validator.Struct(request); err != nil {
		s.Logger.Warnf("validation error : %+v", err)
		return nil, err
	}

	tx := s.DB.WithContext(ctx).Begin()
	defer util.RecoverRollback(tx)

	if request.IsDefault {
		if err := s.CustomerAddressRepsitory.SetDefaultToFalseByUserID(tx, userID); err != nil {
			tx.Rollback()

			s.Logger.Warnf("failed to create user address, can't set default to default : %+v", err)
			return nil, err
		}
	}

	e := entity.CustomerAddress{
		ID:        uuid.NewString(),
		Street:    request.Street,
		City:      request.City,
		Province:  request.Province,
		IsDefault: request.IsDefault,
		UserID:    userID,
	}
	if err := s.CustomerAddressRepsitory.Create(s.DB, &e); err != nil {
		tx.Rollback()

		s.Logger.Warnf("failed to create user address : %+v", err)
		return nil, err
	}

	res := &dto.CustomerAddressResponse{
		ID:        e.ID,
		Street:    e.Street,
		City:      e.City,
		Province:  e.Province,
		IsDefault: e.IsDefault,
	}
	return res, tx.Commit().Error
}
