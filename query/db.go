package query

import (
	permify "github.com/Permify/go-role"
	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/pkg/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"moul.io/zapgorm2"
)

var (
	DB    *gorm.DB
	Perms *permify.Permify
)

func SetupDB(logger *zap.Logger) error {
	dbLogger := zapgorm2.New(logger.Named("db"))
	dbLogger.LogLevel = gormlogger.Info
	dbLogger.SetAsDefault()
	db, err := gorm.Open(mysql.Open(config.C.Database.DSN), &gorm.Config{Logger: dbLogger})
	if err != nil {
		return err
	}

	// Need to use gorm's AutoMigrate for our "non-existing" (at least on a basic ESX FiveM server) models
	if err := db.AutoMigrate(
		// User related
		&model.Account{},
		&model.AccountUser{},
		&model.UserProps{},
		// User location
		model.UserLocation{},
		// Document related
		&model.Document{},
		&model.DocumentMentions{},
		&model.DocumentJobAccess{},
		&model.DocumentUserAccess{},
	); err != nil {
		return err
	}

	// Initialize Permify go-role
	tablePrefix := "arpanet_"
	Perms, err = permify.New(permify.Options{
		Migrate:     true,
		DB:          db,
		TablePrefix: &tablePrefix,
	})
	if err != nil {
		return err
	}

	// Set the DB var and default for the query package
	DB = db
	SetDefault(DB)

	return nil
}
