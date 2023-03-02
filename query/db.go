package query

import (
	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/pkg/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

var DB *gorm.DB

func SetupDB(logger *zap.Logger) error {
	dbLogger := zapgorm2.New(logger.Named("db"))
	dbLogger.LogLevel = gormlogger.Info
	dbLogger.SetAsDefault()
	db, err := gorm.Open(mysql.Open(config.C.Database.DSN), &gorm.Config{Logger: dbLogger})
	if err != nil {
		return err
	}

	// Need to use gorm's AutoMigrate for our "non-existing" (at least on a basic ESX FiveM server) models
	db.AutoMigrate(
		&model.Account{},
		&model.Document{},
		&model.DocumentJobAccess{},
		&model.DocumentUserAccess{},
		&model.UserProps{},
	)

	// Set the DB var and default for the query package
	DB = db
	SetDefault(DB)

	return nil
}
