package config

import (
	"github.com/casbin/casbin/v2"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func NewEnforcer(viperCfg *viper.Viper, logger *logrus.Logger) *casbin.Enforcer {
	modelConfPath := viperCfg.GetString("CASBIN_MODEL_PATH")
	modelPolicyPath := viperCfg.GetString("CASBIN_POLICY_PATH")

	enforcer, err := casbin.NewEnforcer(modelConfPath, modelPolicyPath)
	if err != nil {
		logger.Fatalf("error load casbin configuration : %+v", err)
	}

	return enforcer
}

func NewTestEnforcer(viperCfg *viper.Viper, logger *logrus.Logger) *casbin.Enforcer {
	modelConfPath := viperCfg.GetString("TEST_CASBIN_MODEL_PATH")
	modelPolicyPath := viperCfg.GetString("TEST_CASBIN_POLICY_PATH")

	enforcer, err := casbin.NewEnforcer(modelConfPath, modelPolicyPath)
	if err != nil {
		logger.Fatalf("error load casbin configuration : %+v", err)
	}

	return enforcer
}
