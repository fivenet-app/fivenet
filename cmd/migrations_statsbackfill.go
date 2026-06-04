//nolint:forbidigo // Migration command prints progress to stdout.
package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/cmd/envs"
	"github.com/fivenet-app/fivenet/v2026/cmd/fxopts"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	documentsdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/data"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	docstats "github.com/fivenet-app/fivenet/v2026/pkg/stats"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/timeutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"go.uber.org/fx"
	"google.golang.org/protobuf/encoding/protojson"
)

type MigrationsStatsBackfillCmd struct {
	db    *sql.DB
	stats *docstats.Service

	Start     string `help:"Start date/time (inclusive), e.g. 2026-01-01 or 2026-01-01T00:00:00Z."             required:""`
	BatchSize int    `help:"Number of documents processed per batch."                                                      default:"250"`
	DryRun    bool   `help:"By default this command only estimates impact. Use --no-dry-run to apply changes."             default:"true" negatable:""`
}

func (c *MigrationsStatsBackfillCmd) Run(cli *CLI) error {
	fxOpts := fxopts.GetFxBaseOpts(12*time.Hour, false, true)

	if err := os.Setenv(envs.SkipDBMigrationsEnv, "true"); err != nil {
		return err
	}

	fxOpts = append(
		fxOpts,
		fx.Invoke(
			func(lifecycle fx.Lifecycle, _ *config.Config, db *sql.DB, shutdowner fx.Shutdowner) {
				lifecycle.Append(fx.StartHook(func(ctx context.Context) error {
					go func() {
						exitCode := 0
						c.db = db
						c.stats = docstats.NewService(
							db,
							nil,
						)
						if err := c.run(ctx); err != nil {
							exitCode = 1
							fmt.Println("Error running stats backfill command:", err)
						}
						_ = shutdowner.Shutdown(fx.ExitCode(exitCode))
					}()
					return nil
				}))
			},
		),
	)

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}

func (c *MigrationsStatsBackfillCmd) run(ctx context.Context) error {
	start, err := parseBackfillStart(c.Start)
	if err != nil {
		return err
	}
	if c.BatchSize <= 0 {
		c.BatchSize = 250
	}

	fmt.Printf(
		"Starting stats backfill from %s (batch size: %d)\n",
		start.Format(time.RFC3339),
		c.BatchSize,
	)

	startDay := timeutils.StartOfDay(start.UTC())
	endDay := timeutils.StartOfDay(time.Now().UTC())

	if c.DryRun {
		tDocuments := table.FivenetDocuments
		stmt := tDocuments.
			SELECT(
				mysql.COUNT(mysql.STAR).AS("data_count.total"),
			).
			WHERE(tDocuments.CreatedAt.GT_EQ(mysql.TimestampT(start.UTC())))

		var count database.DataCount
		if err := stmt.QueryContext(ctx, c.db, &count); err != nil {
			return fmt.Errorf("failed to estimate stats backfill document count. %w", err)
		}

		fmt.Printf(
			"Dry run enabled. Would rebuild metrics for %d documents and rebuild rollups from %s to %s.\n",
			count,
			startDay.Format(time.DateOnly),
			endDay.Format(time.DateOnly),
		)
		fmt.Println("Run again with --no-dry-run to apply changes.")
		return nil
	}

	totalProcessed, err := c.rebuildDocumentMetrics(ctx, start)
	if err != nil {
		return err
	}
	fmt.Printf("Finished metrics rebuild for %d documents\n", totalProcessed)

	fmt.Printf(
		"Rebuilding document-column rollups from %s to %s\n",
		startDay.Format(time.DateOnly),
		endDay.Format(time.DateOnly),
	)
	if err := c.stats.RebuildDocumentColumnRollups(ctx, startDay, endDay); err != nil {
		return fmt.Errorf("failed rebuilding document-column rollups. %w", err)
	}

	fmt.Printf(
		"Rebuilding document-metric rollups from %s to %s\n",
		startDay.Format(time.DateOnly),
		endDay.Format(time.DateOnly),
	)
	if err := c.stats.RebuildDocumentMetricRollups(ctx, startDay, endDay); err != nil {
		return fmt.Errorf("failed rebuilding document-metric rollups. %w", err)
	}

	fmt.Println("Stats backfill complete")
	return nil
}

func (c *MigrationsStatsBackfillCmd) rebuildDocumentMetrics(
	ctx context.Context,
	start time.Time,
) (int, error) {
	lastID := int64(0)
	totalProcessed := 0

	for {
		nextLastID, batchCount, err := c.rebuildDocumentMetricsBatch(ctx, start, lastID)
		if err != nil {
			return totalProcessed, fmt.Errorf(
				"failed rebuilding document metrics batch. %w",
				err,
			)
		}
		lastID = nextLastID
		totalProcessed += batchCount

		if batchCount == 0 {
			break
		}

		fmt.Printf("Processed %d documents so far (last id: %d)\n", totalProcessed, lastID)
	}

	return totalProcessed, nil
}

// rebuildDocumentMetricsBatch processes a batch of documents for metrics rebuilding, starting from the given lastID.
//
//nolint:nonamedreturns // Using named returns for clarity in deferred error handling.
func (c *MigrationsStatsBackfillCmd) rebuildDocumentMetricsBatch(
	ctx context.Context,
	start time.Time,
	lastID int64,
) (nextLastID int64, batchCount int, err error) {
	nextLastID = lastID

	tDocuments := table.FivenetDocuments

	stmt := tDocuments.
		SELECT(
			tDocuments.ID.AS("id"),
			tDocuments.CreatorJob.AS("creator_job"),
			tDocuments.CreatedAt.AS("created_at"),
			tDocuments.DeletedAt.AS("deleted_at"),
			tDocuments.Draft.AS("draft"),
			tDocuments.Data.AS("data"),
		).
		FROM(tDocuments).
		WHERE(mysql.AND(
			tDocuments.CreatedAt.GT_EQ(mysql.TimestampT(start.UTC())),
			tDocuments.ID.GT(mysql.Int64(lastID)),
			tDocuments.Data.IS_NOT_NULL(),
		)).
		ORDER_BY(tDocuments.ID.ASC()).
		LIMIT(int64(c.BatchSize))

	var rows []struct {
		DocID      int64          `alias:"id"`
		CreatorJob string         `alias:"creator_job"`
		CreatedAt  time.Time      `alias:"created_at"`
		DeletedAt  sql.NullTime   `alias:"deleted_at"`
		Draft      bool           `alias:"draft"`
		DataRaw    sql.NullString `alias:"data"`
	}
	if err := stmt.QueryContext(ctx, c.db, &rows); err != nil {
		return nextLastID, 0, fmt.Errorf("failed querying documents for stats backfill. %w", err)
	}

	for _, row := range rows {
		var data *documentsdata.DocumentData
		if row.DataRaw.Valid && strings.TrimSpace(row.DataRaw.String) != "" {
			data = &documentsdata.DocumentData{}
			if err := protojson.Unmarshal([]byte(row.DataRaw.String), data); err != nil {
				return nextLastID, batchCount, fmt.Errorf(
					"failed to decode document.data for document %d. %w",
					row.DocID,
					err,
				)
			}
		}

		doc := &documents.Document{
			Id:         row.DocID,
			CreatorJob: row.CreatorJob,
			CreatedAt:  timestamp.New(row.CreatedAt),
			Meta: &documents.DocumentMeta{
				Draft: row.Draft,
			},
			Data: data,
		}
		if row.DeletedAt.Valid {
			doc.DeletedAt = timestamp.New(row.DeletedAt.Time)
		}

		if err := c.stats.RebuildDocumentMetrics(ctx, doc); err != nil {
			return nextLastID, batchCount, fmt.Errorf(
				"failed to rebuild metrics for document %d. %w",
				row.DocID,
				err,
			)
		}

		nextLastID = row.DocID
		batchCount++
	}

	return nextLastID, batchCount, nil
}

func parseBackfillStart(input string) (time.Time, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return time.Time{}, errors.New("start is required")
	}

	layouts := []string{
		time.RFC3339,
		time.DateOnly,
	}

	for _, layout := range layouts {
		t, err := time.Parse(layout, input)
		if err == nil {
			if layout == time.DateOnly {
				return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC), nil
			}
			return t.UTC(), nil
		}
	}

	return time.Time{}, fmt.Errorf(
		"invalid start format %q (expected YYYY-MM-DD or RFC3339)",
		input,
	)
}
