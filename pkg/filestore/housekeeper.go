package filestore

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/file"
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/pkg/storage"
	jet "github.com/go-jet/jet/v2/mysql"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/durationpb"
)

type Housekeeper struct {
	logger  *zap.Logger
	tracer  trace.Tracer
	db      *sql.DB
	storage storage.IStorage

	// joinTables describes each “live‐reference” table:
	//   the Jet table, the column for file_id, and (optionally) a column
	//   for parent_id if you ever want to filter by namespace.
	getTablesListFn func() []joinInfo

	// gracePeriod for soft-deleted rows (e.g. 7 days).
	gracePeriod time.Duration

	batchSize int64
}

// joinInfo holds information about one join table.
type joinInfo struct {
	Table   jet.Table         // e.g. table.FivenetDocumentsFiles
	FileCol jet.ColumnInteger // e.g. table.FivenetDocumentsFiles.FileID
}

type Result struct {
	fx.Out

	Housekeeper  *Housekeeper
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

type HousekeeperParams struct {
	fx.In

	Logger *zap.Logger
	TP     *tracesdk.TracerProvider
	DB     *sql.DB
}

func NewHousekeeper(p HousekeeperParams) (Result, error) {
	h := &Housekeeper{
		logger: p.Logger.Named("housekeeper"),
		tracer: p.TP.Tracer("mstlystcdata-cache"),
		db:     p.DB,

		getTablesListFn: func() []joinInfo {
			tableListsMu.Lock()
			defer tableListsMu.Unlock()

			return tablesList
		},

		gracePeriod: 7 * 24 * time.Hour, // 7 days

		batchSize: 100,
	}
	return Result{
		Housekeeper:  h,
		CronRegister: h,
	}, nil
}

func (s *Housekeeper) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "filestore.housekeeper",
		Schedule: "*/2 * * * *", // Every 2 minutes
		Timeout:  durationpb.New(60 * time.Second),
	}); err != nil {
		return err
	}

	return nil
}

func (h *Housekeeper) RegisterCronjobHandlers(hand *croner.Handlers) error {
	hand.Add("filestore.housekeeper", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := h.tracer.Start(ctx, "filestore.housekeeper")
		defer span.End()

		if err := h.Run(ctx); err != nil {
			h.logger.Error("failed to run filestore housekeeper job", zap.Error(err))
			return err
		}

		return nil
	})

	return nil
}

func (h *Housekeeper) Run(ctx context.Context) error {
	cutoff := time.Now().Add(-h.gracePeriod)
	maxDeletes := 200
	deletes := 0

	joinTables := h.getTablesListFn()

	// 1) Build the WHERE clause for this batch:
	//
	//   WHERE
	//   (deleted_at IS NOT NULL AND deleted_at < cutoff)
	//   OR
	//   (deleted_at IS NULL
	//     AND NOT EXISTS (SELECT 1 FROM fivenet_documents_files WHERE file_id = f.id)
	//     AND NOT EXISTS (SELECT 1 FROM fivenet_wiki_pages_files WHERE file_id = f.id)
	//     AND … for each join table …
	//   )
	//

	// A) soft-deleted & expired
	softExpired := tFiles.DeletedAt.IS_NOT_NULL().AND(
		tFiles.DeletedAt.LT(jet.TimestampT(cutoff)),
	)

	// B) orphaned & not yet soft-deleted
	candidateWhere := softExpired

	if len(joinTables) > 0 {
		orphanCond := tFiles.DeletedAt.IS_NULL()

		for _, ji := range joinTables {
			// NOT EXISTS (SELECT 1 FROM ji.table WHERE ji.fileCol = files.id)
			orphanCond = orphanCond.AND(
				jet.NOT(jet.EXISTS(
					ji.Table.SELECT(jet.RawInt("1")).
						FROM(ji.Table).
						WHERE(ji.FileCol.EQ(tFiles.ID)),
				)),
			)
		}

		candidateWhere = jet.OR(softExpired, orphanCond)
	}

	for {
		// 2) Select up to batchSize candidates
		var files []*file.File
		err := tFiles.
			SELECT(
				tFiles.ID,
				tFiles.FilePath,
			).
			FROM(tFiles).
			WHERE(candidateWhere).
			ORDER_BY(tFiles.ID). // deterministic ordering
			LIMIT(h.batchSize).
			QueryContext(ctx, h.db, &files)
		if err != nil {
			return fmt.Errorf("failed to select candidates %w", err)
		}

		// If no more candidates, we’re done
		if len(files) == 0 {
			break
		}

		// 3) Process each candidate in this batch in one transaction
		tx, err := h.db.BeginTx(ctx, nil)
		if err != nil {
			return fmt.Errorf("failed to begin tx. %w", err)
		}

		deletes += len(files)

		for _, c := range files {
			// A) Delete from storage
			if err := h.storage.Delete(ctx, c.FilePath); err != nil {
				h.logger.Error("delete storage key failed (skipping DB removal)", zap.String("file_path", c.FilePath), zap.Error(err))
				continue
			}

			// B) Delete row from fivenet_files
			if _, err := tFiles.
				DELETE().
				WHERE(tFiles.ID.EQ(jet.Uint64(c.Id))).
				ExecContext(ctx, tx); err != nil {
				tx.Rollback()
				return fmt.Errorf("delete fivenet_files id=%d. %w", c.Id, err)
			}
			h.logger.Debug("deleted file", zap.Uint64("file_id", c.Id), zap.String("file_path", c.FilePath))
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit tx. %w", err)
		}

		// If we haven't reached the maximum yet loop again to pick up the next batch. The cutoff and conditions remain the same.
		if deletes >= maxDeletes {
			break
		}
	}

	metricDeletedFiles.Set(float64(deletes))

	return nil
}
