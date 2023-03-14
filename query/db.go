package query

import (
	"database/sql"
	"embed"
	"errors"

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
	db           *sql.DB
)

func SetupDB(logger *zap.Logger) (*sql.DB, error) {
	if err := MigrateDB(logger, config.C.Database.DSN); err != nil {
		return nil, err
	}

	// Connect to database
	var err error
	db, err = sql.Open("mysql", config.C.Database.DSN)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func MigrateDB(logger *zap.Logger, dsn string) error {
	logger.Info("starting database migrations")
	// Connect to database
	db, err := sql.Open("mysql", dsn+"&multiStatements=true")
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
		if !errors.Is(migrate.ErrNoChange, err) {
			return err
		} else {
			logger.Info("database migration have caused no changes")
		}
	} else {
		logger.Info("completed database migrations changes have been made")
	}

	return db.Close()
}
