package service

import (
	"context"
	"errors"
	"store-api/internal/dto"
	"store-api/internal/entity"
	"store-api/internal/repository"
	"store-api/util"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CartItemService interface {
	Create(ctx context.Context, request *dto.CartItemRequest, userID string) error
	FindAll(ctx context.Context, userID string) ([]dto.CartItemResponse, error)
}

type CartItemServiceImpl struct {
	DB                 *gorm.DB
	Logger             *logrus.Logger
	Validator          *validator.Validate
	CartItemRepository *repository.CartItemRepository
	ProductRepository  *repository.ProductRepository
}

func NewCartServiceImpl(db *gorm.DB, logger *logrus.Logger, validator *validator.Validate, repositories *repository.Repositories) *CartItemServiceImpl {
	return &CartItemServiceImpl{
		DB:                 db,
		Logger:             logger,
		Validator:          validator,
		CartItemRepository: repositories.CartItemRepository,
		ProductRepository:  repositories.ProductRepository,
	}
}

func (s *CartItemServiceImpl) FindAll(ctx context.Context, userID string) ([]dto.CartItemResponse, error) {
	return s.CartItemRepository.FindByUserIDWithAssociation(s.DB, userID)
}

func (s *CartItemServiceImpl) Create(ctx context.Context, request *dto.CartItemRequest, userID string) error {
	if err := s.Validator.Struct(request); err != nil {
		s.Logger.Warnf("validation error : %+v", err)
		return err
	}

	var product entity.Product
	if err := s.ProductRepository.FindByID(s.DB, &product, request.ProductID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Warnf("product with id %s not found : %+v", request.ProductID, err)
			return fiber.NewError(fiber.StatusNotFound, "product not found")
		}

		s.Logger.Warnf("can't find product: %+v", err)
		return err
	}

	if product.Stock < request.Quantity {
		s.Logger.Warnf("requested quantity is exceeding product stock with product ID %s", product.ID)
		return fiber.NewError(fiber.StatusInternalServerError, "requested quantity is exceeding product stock")
	}

	tx := s.DB.WithContext(ctx).Begin()
	defer util.RecoverRollback(tx)

	e := entity.CartItem{
		UserID:    userID,
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
		CreatedAt: time.Now(),
	}
	if err := s.CartItemRepository.Save(tx, &e); err != nil {
		tx.Rollback()

		s.Logger.Warnf("failed to insert product with id %s into cart : %+v", product.ID, err)
		return err
	}

	return tx.Commit().Error
}
