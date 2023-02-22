package cmd

import (
	"database/sql"
	"fmt"

	"github.com/galexrt/rphub/model"
	"github.com/galexrt/rphub/pkg/config"
	"github.com/galexrt/rphub/pkg/sync"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
	"moul.io/zapgorm2"
)

var syncCmd = &cobra.Command{
	Use: "sync",
	RunE: func(cmd *cobra.Command, args []string) error {
		// refer to https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			config.C.FiveM.Database.Username, config.C.FiveM.Database.Password,
			config.C.FiveM.Database.Host, config.C.FiveM.Database.Port,
			config.C.FiveM.Database.DBName)
		dbLogger := zapgorm2.New(logger.Named("db"))
		dbLogger.SetAsDefault()
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return err
		}

		sy := sync.New(logger.Named("sync"))
		sy.SyncJobs(db)
		sy.SyncJobGrades(db)
		sy.SyncUsers(db)

		// TODO Sync jobs and job_grades tables from FiveM to our database

		return nil
	},
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return model.SetupDB(logger)
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
