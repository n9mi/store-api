package service

import (
	"errors"
	"store-api/internal/dto"
	"store-api/internal/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type JWTService interface {
	GenerateAccessToken(authDto *dto.AuthDTO) (string, int64, error)
	GenerateRefreshToken(authDto *dto.AuthDTO) (string, int64, error)
	ParseAccessToken(token string) (*dto.AuthDTO, error)
	ParseRefreshToken(token string) (*dto.AuthDTO, error)
}

type JWTServiceImpl struct {
	ViperCfg *viper.Viper
	Logger   *logrus.Logger
}

func NewJwtService(viperCfg *viper.Viper, db *gorm.DB, logger *logrus.Logger, repositories *repository.Repositories) *JWTServiceImpl {
	return &JWTServiceImpl{
		ViperCfg: viperCfg,
		Logger:   logger,
	}
}

func (s *JWTServiceImpl) GenerateAccessToken(authDto *dto.AuthDTO) (string, int64, error) {
	accessTokenKey := s.ViperCfg.GetString("ACCESS_TOKEN_KEY")
	accessExpireMin := s.ViperCfg.GetInt("ACCESS_TOKEN_EXPIRE_MINUTES")

	return s.generateJWTToken(accessTokenKey, accessExpireMin, authDto)
}

func (s *JWTServiceImpl) GenerateRefreshToken(authDto *dto.AuthDTO) (string, int64, error) {
	refreshTokenKey := s.ViperCfg.GetString("REFRESH_TOKEN_KEY")
	refreshExpireMin := s.ViperCfg.GetInt("REFRESH_TOKEN_EXPIRE_MINUTES")

	return s.generateJWTToken(refreshTokenKey, refreshExpireMin, authDto)
}

func (s *JWTServiceImpl) ParseAccessToken(token string) (*dto.AuthDTO, error) {
	accessTokenKey := s.ViperCfg.GetString("ACCESS_TOKEN_KEY")

	return s.parseJWTToken(accessTokenKey, token)
}

func (s *JWTServiceImpl) ParseRefreshToken(token string) (*dto.AuthDTO, error) {
	refreshTokenKey := s.ViperCfg.GetString("REFRESH_TOKEN_KEY")

	return s.parseJWTToken(refreshTokenKey, token)
}

func (s *JWTServiceImpl) generateJWTToken(key string, expMinutes int, authDTO *dto.AuthDTO) (string, int64, error) {
	timeMin := time.Duration(expMinutes) * time.Minute
	expAtUnix := time.Now().Add(timeMin).Unix()

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp":  expAtUnix,
		"data": authDTO,
	})
	sg, err := t.SignedString([]byte(key))
	if err != nil {
		s.Logger.Warnf("failed to generate token : %+v", err)
		return "", 0, err
	}

	return sg, expAtUnix, nil
}

func (s *JWTServiceImpl) parseJWTToken(key string, token string) (*dto.AuthDTO, error) {
	t, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			s.Logger.Warnf("expired token : %+v", err)
			return nil, fiber.NewError(fiber.StatusUnauthorized, "expired_token")
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			s.Logger.Warnf("malformed token : %+v", err)
			return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid_token")
		}

		s.Logger.Warnf("failed to parse token : %+v", err)
		return nil, err
	}

	claims := t.Claims.(jwt.MapClaims)
	authDTOMap, ok := claims["data"].(map[string]interface{})
	if !ok {
		s.Logger.Warnf("failed to parse token : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid_token")
	}

	var authDTO dto.AuthDTO
	authDTO.UserID, ok = authDTOMap["UserID"].(string)
	if !ok {
		s.Logger.Warnf("failed to parse UserID from token : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid_token")
	}
	authDTO.UserEmail, ok = authDTOMap["UserEmail"].(string)
	if !ok {
		s.Logger.Warnf("failed to parse UserEmail from token : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid_token")
	}
	authDTO.UserName, ok = authDTOMap["UserName"].(string)
	if !ok {
		s.Logger.Warnf("failed to parse UserName from token : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid_token")
	}
	authDTO.UserCurrentRole, ok = authDTOMap["UserCurrentRole"].(string)
	if !ok {
		s.Logger.Warnf("failed to parse UserCurrentRole from token : %+v", err)
		return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid_token")
	}

	return &authDTO, nil
}
