package service

import (
	"context"
	"errors"
	"fmt"
	"store-api/internal/dto"
	"store-api/internal/entity"
	"store-api/internal/repository"
	"store-api/util"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrderService interface {
	Create(ctx context.Context, request *dto.OrderCreateRequest, userID string) (*dto.OrderCreateResponse, error)
}

type OrderServiceImpl struct {
	DB                        *gorm.DB
	Logger                    *logrus.Logger
	Validator                 *validator.Validate
	OrderRepository           *repository.OrderRepository
	OrderItemRepository       *repository.OrderItemRepository
	PaymentMethodRepository   *repository.PaymentMethodRepository
	CustomerAddressRepository *repository.CustomerAddressRepository
	CartItemRepository        *repository.CartItemRepository
}

func NewOrderServiceImpl(db *gorm.DB, logger *logrus.Logger, validator *validator.Validate, repositories *repository.Repositories) *OrderServiceImpl {
	return &OrderServiceImpl{
		DB:                        db,
		Logger:                    logger,
		Validator:                 validator,
		OrderRepository:           repositories.OrderRepository,
		OrderItemRepository:       repositories.OrderItemRepository,
		PaymentMethodRepository:   repositories.PaymentMethodRepository,
		CustomerAddressRepository: repositories.CustomerAddressRepository,
		CartItemRepository:        repositories.CartItemRepository,
	}
}

func (s *OrderServiceImpl) Create(ctx context.Context, request *dto.OrderCreateRequest, userID string) (*dto.OrderCreateResponse, error) {
	if err := s.Validator.Struct(request); err != nil {
		s.Logger.Warnf("validation error : %+v", err)
		return nil, err
	}

	tx := s.DB.WithContext(ctx).Begin()
	defer util.RecoverRollback(tx)

	if ok, err := s.CustomerAddressRepository.IsExistByIDAndUserID(tx, request.CustomerAddressID, userID); err != nil {
		tx.Rollback()

		s.Logger.Warnf("failed to make order : %+v", err)
		return nil, err
	} else if !ok {
		s.Logger.Warnf("failed to make order - customer address doens't exists : %+v", err)
		return nil, fiber.NewError(fiber.StatusNotFound, "customer address doesn't exsits")
	}

	var paymentMethod entity.PaymentMethod
	if err := s.PaymentMethodRepository.FindByID(tx, &paymentMethod, request.PaymentMethodID); err != nil {
		tx.Rollback()

		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Warnf("failed to make order - payment doens't exists : %+v", err)
			return nil, fiber.NewError(fiber.StatusNotFound, "payment doesn't exsits")
		}

		s.Logger.Warnf("failed to make order : %+v", err)
		return nil, err
	}

	custCartItems, err := s.CartItemRepository.FindByUserID(tx, userID)
	if err != nil {
		tx.Rollback()

		s.Logger.Warnf("failed to make order : %+v", err)
		return nil, err
	}
	if len(custCartItems) == 0 {
		s.Logger.Warn("failed to make order : cart empty")
		return nil, fiber.NewError(fiber.StatusBadRequest, "empty cart")
	}

	order := entity.Order{
		ID:                   uuid.NewString(),
		UserID:               userID,
		CustomerAddressID:    request.CustomerAddressID,
		PaymentMethodID:      request.PaymentMethodID,
		PaymentCode:          fmt.Sprintf("%s-%s-%s", paymentMethod.Code, util.GenerateRandomString(4), util.GenerateRandomString(4)),
		PaymentCodeExpiredAt: time.Now().Add(time.Duration(24) * time.Hour),
	}
	if err := s.OrderRepository.Create(tx, &order); err != nil {
		tx.Rollback()

		s.Logger.Warnf("failed to make order : %+v", err)
		return nil, err
	}

	orderItem := make([]*entity.OrderItem, len(custCartItems))
	for i, cI := range custCartItems {
		orderItem[i] = &entity.OrderItem{
			ID:        uuid.NewString(),
			OrderID:   order.ID,
			ProductID: cI.ProductID,
			Quantity:  cI.Quantity,
		}
	}
	if err := s.OrderItemRepository.CreateMany(tx, orderItem); err != nil {
		tx.Rollback()

		s.Logger.Warnf("failed to make order : %+v", err)
		return nil, err
	}
	if err := s.CartItemRepository.Empty(tx, userID); err != nil {
		tx.Rollback()

		s.Logger.Warnf("failed to make order : %+v", err)
		return nil, err
	}
	orderItems, err := s.OrderRepository.FindByOrderIDWithAssociation(tx, order.ID)
	if err != nil {
		tx.Rollback()

		s.Logger.Warnf("failed to make order : %+v", err)
		return nil, err
	}
	var ordTotalPrice float64 = 0
	for _, oI := range orderItems {
		ordTotalPrice += (oI.TotalPrice)
	}

	res := &dto.OrderCreateResponse{
		PaymentCode:          order.PaymentCode,
		PaymentCodeExpiredAt: order.PaymentCodeExpiredAt,
		OrderItems:           orderItems,
		TotalPrice:           ordTotalPrice,
	}
	return res, tx.Commit().Error
}
