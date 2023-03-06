package query

import (
	"database/sql"
	"embed"

	"github.com/galexrt/arpanet/pkg/config"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"go.uber.org/zap"
)

var (
	//go:embed migrations/*.sql
	migrationsFS embed.FS

	DB *sql.DB
)

func SetupDB(logger *zap.Logger) error {
	if err := migrateDB(logger); err != nil {
		return err
	}

	// Connect to database
	db, err := sql.Open("mysql", config.C.Database.DSN)
	if err != nil {
		return err
	}

	// Set the DB var and default for the query package
	DB = db

	return nil
}

func migrateDB(logger *zap.Logger) error {
	logger.Info("starting database migrations")
	// Connect to database
	db, err := sql.Open("mysql", config.C.Database.DSN+"&multiStatements=true")
	if err != nil {
		return err
	}

	// Setup migrate source and driver
	source, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		return err
	}
	driver, err := mysql.WithInstance(db, &mysql.Config{
		MigrationsTable: "arpanet_zschema_migrations",
	})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance(
		"iofs", source,
		"mysql", driver)
	if err != nil {
		return err
	}
	// Run migrations
	if err := m.Up(); err != nil {
		if err != migrate.ErrNoChange {
			return err
		} else {
			logger.Info("database migration have caused no changes")
		}
	}

	return db.Close()
}
