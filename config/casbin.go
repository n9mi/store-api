package config

import (
	"github.com/casbin/casbin/v2"
	"github.com/sirupsen/logrus"
)

func NewEnforcer(logger *logrus.Logger) *casbin.Enforcer {
	enforcer, err := casbin.NewEnforcer("./internal/casbin/model.conf", "./internal/casbin/policy.csv")
	if err != nil {
		logger.Fatalf("error load casbin configuration : %+v", err)
	}

	return enforcer
}

func NewTestEnforcer(logger *logrus.Logger) *casbin.Enforcer {
	enforcer, err := casbin.NewEnforcer("../internal/casbin/model.conf", "../internal/casbin/policy.csv")
	if err != nil {
		logger.Fatalf("error load casbin configuration : %+v", err)
	}

	return enforcer
}
