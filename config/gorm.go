package config

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewDatabase(viperCfg *viper.Viper, log *logrus.Logger) *gorm.DB {
	dbHost := viperCfg.GetString("DB_HOST")
	dbUser := viperCfg.GetString("DB_USER")
	dbPassword := viperCfg.GetString("DB_PASSWORD")
	dbName := viperCfg.GetString("DB_NAME")
	dbPort := viperCfg.GetInt("DB_PORT")
	dbSsl := viperCfg.GetString("DB_SSL_MODE")
	dbTimezone := viperCfg.GetString("DB_TIMEZONE")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		dbHost,
		dbUser,
		dbPassword,
		dbName,
		dbPort,
		dbSsl,
		dbTimezone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.New(&logrusWriter{
			Logger: log,
		}, logger.Config{
			SlowThreshold:        5 * time.Second,
			Colorful:             true,
			ParameterizedQueries: true,
			LogLevel:             logger.Info,
		}),
	})

	if err != nil {
		log.Fatalf("failed to create database : %+v", err)
	}

	return db
}

type logrusWriter struct {
	Logger *logrus.Logger
}

func (w *logrusWriter) Printf(message string, args ...interface{}) {
	w.Logger.Tracef(message, args...)
}
