package model

import (
	"github.com/galexrt/rphub/pkg/config"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"moul.io/zapgorm2"
)

var DB *gorm.DB

func SetupDB(logger *zap.Logger) error {
	// refer to https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	dbLogger := zapgorm2.New(logger.Named("db"))
	dbLogger.SetAsDefault()
	db, err := gorm.Open(postgres.Open(config.C.Database.DSN), &gorm.Config{Logger: dbLogger})
	if err != nil {
		return err
	}

	DB = db

	return nil
}
