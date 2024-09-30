package service

import (
	"context"
	"errors"
	"fmt"
	"store-api/internal/dto"
	"store-api/internal/entity"
	"store-api/internal/repository"
	"store-api/util"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(ctx context.Context, request *dto.RegisterRequest) error
}

type AuthServiceImpl struct {
	DB             *gorm.DB
	Validator      *validator.Validate
	Logger         *logrus.Logger
	UserRepository *repository.UserRepository
	RoleRepository *repository.RoleRepository
}

func NewAuthService(db *gorm.DB, validator *validator.Validate, logger *logrus.Logger, repositories *repository.Repositories) *AuthServiceImpl {
	return &AuthServiceImpl{
		DB:             db,
		Validator:      validator,
		Logger:         logger,
		UserRepository: repositories.UserRepository,
		RoleRepository: repositories.RoleRepository,
	}
}

func (s *AuthServiceImpl) Register(ctx context.Context, request *dto.RegisterRequest) error {
	if err := s.Validator.Struct(request); err != nil {
		s.Logger.Warnf("request validation failed : %+v", err)
		return err
	}
	if request.AsRole == "" {
		s.Logger.Warnf("request validation failed : unspecified role")
		return fiber.NewError(fiber.StatusBadRequest, "unspecified role")
	}
	if !(request.AsRole == "merchant" || request.AsRole == "customer") {
		s.Logger.Warnf("request validation failed : invalid role")
		return fiber.NewError(fiber.StatusBadRequest, "invalid role")
	}

	// check if user already has role that specified in the request. if true, then return error.
	// if user already exists but has different role, then allow it to be assigned to the specified role.
	var user entity.User
	if err := s.UserRepository.FindByEmail(s.DB, &user, request.Email); len(user.ID) > 0 && err == nil {
		if isExists, err := s.UserRepository.HasRole(s.DB, user.ID, fmt.Sprintf("role_%s", request.AsRole)); isExists {
			s.Logger.Warnf("register conflict : user with role already exists")
			return fiber.NewError(fiber.StatusConflict, "user already exists")
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			s.Logger.Warnf("register failed : %+v", err)
			return fiber.NewError(fiber.StatusConflict, "something wrong")
		}
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.Logger.Warnf("register failed : %+v", err)
		return fiber.NewError(fiber.StatusConflict, "something wrong")
	}

	tx := s.DB.WithContext(ctx).Begin()
	defer util.RecoverRollback(tx)

	// if user doens't exists yet, create a new one. if user already exists, assigned to the specified role.
	if len(user.ID) == 0 {
		userPwd, err := util.NewUserPassword(request.Password)
		if err != nil {
			s.Logger.Warnf("register failed: failed to generate password")
			return fiber.ErrInternalServerError
		}

		user = entity.User{
			ID:       uuid.NewString(),
			Name:     request.Name,
			Email:    request.Email,
			Password: userPwd,
		}
		if err := s.UserRepository.Create(tx, &user); err != nil {
			tx.Rollback()

			s.Logger.Warnf("register failed: %+v", err)
			return fiber.NewError(fiber.StatusInternalServerError, "something wrong")
		}
	}
	var role entity.Role
	if err := s.RoleRepository.FindByID(tx, &role, fmt.Sprintf("role_%s", request.AsRole)); err != nil {
		tx.Rollback()

		s.Logger.Warnf("register failed: %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "something wrong")
	}
	if err := s.UserRepository.AssignRole(tx, &user, &role); err != nil {
		tx.Rollback()

		s.Logger.Warnf("register failed: %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "something wrong")
	}

	if err := tx.Commit().Error; err != nil {
		s.Logger.Warnf("register failed: %+v", err)
		return fiber.NewError(fiber.StatusInternalServerError, "something wrong")
	}

	return nil
}
