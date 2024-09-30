package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

func NewValidator(viperCfg *viper.Viper) *validator.Validate {
	validator := validator.New()

	return validator
}
