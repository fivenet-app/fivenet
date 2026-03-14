//nolint:forbidigo // Migration command prints progress to stdout.
package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/cmd/envs"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents"
	documentsdata "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/data"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	docstats "github.com/fivenet-app/fivenet/v2026/pkg/stats"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/timeutils"
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

func (c *MigrationsStatsBackfillCmd) Run() error {
	fxOpts := getFxBaseOpts(12*time.Hour, false, true)

	if err := os.Setenv(envs.SkipDBMigrationsEnv, "true"); err != nil {
		return err
	}

	fxOpts = append(
		fxOpts,
		fx.Invoke(
			func(lifecycle fx.Lifecycle, _ *config.Config, db *sql.DB, shutdowner fx.Shutdowner) {
				lifecycle.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
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
					},
				})
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
		var count int64
		if err := c.db.QueryRowContext(
			ctx,
			`SELECT COUNT(*) FROM fivenet_documents WHERE created_at >= ?`,
			start.UTC(),
		).Scan(&count); err != nil {
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
		rows, err := c.db.QueryContext(
			ctx,
			`SELECT id, creator_job, created_at, deleted_at, draft, data
FROM fivenet_documents
WHERE created_at >= ?
  AND id > ?
  AND data IS NOT NULL
ORDER BY id ASC
LIMIT ?`,
			start.UTC(),
			lastID,
			c.BatchSize,
		)
		if err != nil {
			return totalProcessed, fmt.Errorf(
				"failed querying documents for stats backfill. %w",
				err,
			)
		}

		batchCount := 0
		for rows.Next() {
			var (
				id         int64
				creatorJob string
				createdAt  time.Time
				deletedAt  sql.NullTime
				draft      bool
				dataRaw    sql.NullString
			)
			if err := rows.Scan(&id, &creatorJob, &createdAt, &deletedAt, &draft, &dataRaw); err != nil {
				rows.Close()
				return totalProcessed, fmt.Errorf(
					"failed scanning document row for stats backfill. %w",
					err,
				)
			}

			var data *documentsdata.DocumentData
			if dataRaw.Valid && strings.TrimSpace(dataRaw.String) != "" {
				data = &documentsdata.DocumentData{}
				if err := protojson.Unmarshal([]byte(dataRaw.String), data); err != nil {
					rows.Close()
					return totalProcessed, fmt.Errorf(
						"failed to decode document.data for document %d. %w",
						id,
						err,
					)
				}
			}

			doc := &documents.Document{
				Id:         id,
				CreatorJob: creatorJob,
				CreatedAt:  timestamp.New(createdAt),
				Meta: &documents.DocumentMeta{
					Draft: draft,
				},
				Data: data,
			}
			if deletedAt.Valid {
				doc.DeletedAt = timestamp.New(deletedAt.Time)
			}

			if err := c.stats.RebuildDocumentMetrics(ctx, doc); err != nil {
				rows.Close()
				return totalProcessed, fmt.Errorf(
					"failed to rebuild metrics for document %d. %w",
					id,
					err,
				)
			}

			lastID = id
			batchCount++
			totalProcessed++
		}

		if err := rows.Err(); err != nil {
			rows.Close()
			return totalProcessed, fmt.Errorf("row iteration error during stats backfill. %w", err)
		}
		rows.Close()

		if batchCount == 0 {
			break
		}

		fmt.Printf("Processed %d documents so far (last id: %d)\n", totalProcessed, lastID)
	}

	return totalProcessed, nil
}

func parseBackfillStart(input string) (time.Time, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return time.Time{}, fmt.Errorf("start is required")
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
