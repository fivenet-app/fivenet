package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/fivenet-app/fivenet/v2025/cmd/envs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/content"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/storage"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"go.uber.org/fx"
)

type MigrationsHTMLToJSONCmd struct {
	db *sql.DB

	DryRun bool `default:"true" negatable:"" help:"By default migrations will only be simulated, use --no-dry-run to disable dry run mode and apply changes."`
}

func (c *MigrationsHTMLToJSONCmd) Run() error {
	fxOpts := getFxBaseOpts(12*time.Hour, false, true)

	if err := os.Setenv(envs.SkipDBMigrationsEnv, "true"); err != nil {
		return err
	}

	fxOpts = append(
		fxOpts,
		fx.Invoke(
			func(lifecycle fx.Lifecycle, cfg *config.Config, db *sql.DB, storage storage.IStorage, shutdowner fx.Shutdowner) {
				lifecycle.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							exitCode := 0
							c.db = db
							if err := c.run(ctx); err != nil {
								// handle error, set non-zero exit code so caller knows the job failed
								exitCode = 1
								fmt.Println("Error running filestore command:", err)
							}
							_ = shutdowner.Shutdown(fx.ExitCode(exitCode))
						}()
						return nil
					},
				})
			},
		),
	)

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}

func (c *MigrationsHTMLToJSONCmd) run(ctx context.Context) error {
	tables := []struct {
		table            mysql.Table
		idCol            mysql.ColumnInteger
		legacyContentCol mysql.ColumnString
	}{
		{
			table:            table.FivenetDocuments,
			idCol:            table.FivenetDocuments.ID,
			legacyContentCol: table.FivenetDocuments.ContentJSON,
		},
		{
			table:            table.FivenetWikiPages,
			idCol:            table.FivenetWikiPages.ID,
			legacyContentCol: table.FivenetWikiPages.Content,
		},
		{
			table:            table.FivenetQualifications,
			idCol:            table.FivenetQualifications.ID,
			legacyContentCol: table.FivenetQualifications.Content,
		},
		{
			table:            table.FivenetMailerMessages,
			idCol:            table.FivenetMailerMessages.ID,
			legacyContentCol: table.FivenetMailerMessages.Content,
		},
		{
			table:            table.FivenetDocumentsComments,
			idCol:            table.FivenetDocumentsComments.ID,
			legacyContentCol: table.FivenetDocumentsComments.Content,
		},
		{
			table:            table.FivenetCalendarEntries,
			idCol:            table.FivenetCalendarEntries.ID,
			legacyContentCol: table.FivenetCalendarEntries.Content,
		},
		{
			table:            table.FivenetCalendar,
			idCol:            table.FivenetCalendar.ID,
			legacyContentCol: table.FivenetCalendar.Description,
		},
	}

	if c.DryRun {
		for _, tbl := range tables {
			count, err := c.countTable(ctx, tbl.table, tbl.idCol, tbl.legacyContentCol)
			if err != nil {
				return err
			}

			fmt.Printf(
				"Dry run enabled, no changes will be applied. Counted %d rows in table %s that need to be migrated from HTML to custom JSON format.\n",
				count,
				tbl.table.TableName(),
			)
		}

		fmt.Println(
			"Run the migration command again with `--no-dry-run` to disable dry run mode and apply the changes to the tables.",
		)

		return nil
	}

	total := 0
	for _, tbl := range tables {
		t, err := c.migrateTable(ctx, tbl.table, tbl.idCol, tbl.legacyContentCol)
		if err != nil {
			return err
		}
		total += t
	}

	tDocuments := table.FivenetDocuments
	stmt := tDocuments.
		UPDATE().
		SET(
			tDocuments.ContentText.SET(tDocuments.Summary),
		).
		WHERE(mysql.AND(
			tDocuments.Summary.NOT_EQ(mysql.String("")),
			tDocuments.ContentText.EQ(mysql.String("")),
		))

	if _, err := stmt.ExecContext(ctx, c.db); err != nil {
		return fmt.Errorf(
			"failed to backfill summary in fivenet_documents (for existing documents). %w",
			err,
		)
	}

	fmt.Printf("Migration applied to total of %d rows (tables %d).\n", total, len(tables))

	return nil
}

func (c *MigrationsHTMLToJSONCmd) countTable(
	ctx context.Context,
	table mysql.Table,
	idCol mysql.ColumnInteger,
	legacyContentCol mysql.ColumnString,
) (int64, error) {
	stmt := table.
		SELECT(
			mysql.COUNT(idCol).AS("data_count.total"),
		).
		FROM(table).
		WHERE(mysql.AND(
			legacyContentCol.LIKE(mysql.String("<%")),
		))

	var dest database.DataCount
	if err := stmt.QueryContext(ctx, c.db, &dest); err != nil {
		return 0, fmt.Errorf(
			"failed to count rows of table %s (for dry-run estimation). %w",
			table.TableName(),
			err,
		)
	}

	return dest.Total, nil
}

func (c *MigrationsHTMLToJSONCmd) migrateTable(
	ctx context.Context,
	table mysql.Table,
	idCol mysql.ColumnInteger,
	contentCol mysql.ColumnString,
) (int, error) {
	total := 0
	for {
		stmt := table.
			SELECT(
				idCol.AS("id"),
				contentCol.AS("content"),
			).
			FROM(table).
			WHERE(mysql.AND(
				contentCol.LIKE(mysql.String("<%")),
			)).
			ORDER_BY(idCol.ASC()).
			LIMIT(100)

		var dest []struct {
			ID      int64
			Content *content.Content
		}
		if err := stmt.QueryContext(ctx, c.db, &dest); err != nil {
			return 0, fmt.Errorf("failed to query rows of table %s. %w", table.TableName(), err)
		}

		if len(dest) == 0 {
			break
		}
		for _, row := range dest {
			fmt.Printf("Migrating %s ID: %d\n", table.TableName(), row.ID)

			insertStmt := table.
				UPDATE(
					contentCol,
				).
				SET(
					row.Content,
				).
				WHERE(idCol.EQ(mysql.Int64(row.ID))).
				LIMIT(1)

			_, err := insertStmt.ExecContext(ctx, c.db)
			if err != nil {
				return 0, fmt.Errorf(
					"failed to update ID %d in table %s. %w",
					row.ID,
					table.TableName(),
					err,
				)
			}
		}

		total += len(dest)

		select {
		case <-ctx.Done():
			return 0, ctx.Err()

		case <-time.After(150 * time.Millisecond):
		}
	}

	fmt.Printf("Migration applied to total of %d of table %s.\n", total, table.TableName())

	return total, nil
}
