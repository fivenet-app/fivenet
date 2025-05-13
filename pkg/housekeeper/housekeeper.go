package housekeeper

import (
	"context"
	"database/sql"
	"fmt"
	"slices"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	jet "github.com/go-jet/jet/v2/mysql"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
)

var Module = fx.Module("db_housekeeper",
	fx.Provide(
		New,
	),
)

const DefaultDeleteLimit = 500

type Housekeeper struct {
	logger *zap.Logger
	tracer trace.Tracer

	db *sql.DB

	getTablesListFn func() map[string]*Table
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	DB     *sql.DB
	TP     *tracesdk.TracerProvider

	Cron         croner.ICron
	CronHandlers *croner.Handlers
}

const (
	lastTableMapIndex = "last_key"
	lastJobName       = "last_job_name"
)

func New(p Params) *Housekeeper {
	h := &Housekeeper{
		logger: p.Logger.Named("housekeeper"),
		tracer: p.TP.Tracer("housekeeper"),
		db:     p.DB,
		getTablesListFn: func() map[string]*Table {
			return tablesList
		},
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := p.Cron.RegisterCronjob(ctx, &cron.Cronjob{
			Name:     "housekeeper.run",
			Schedule: "*/3 * * * *", // Every 3 minutes
			Timeout:  durationpb.New(1 * time.Minute),
		}); err != nil {
			return err
		}

		if err := p.Cron.RegisterCronjob(ctx, &cron.Cronjob{
			Name:     "housekeeper.job_delete",
			Schedule: "@everysecond", // Every second
		}); err != nil {
			return err
		}

		return nil
	}))

	/* p.CronHandlers.Add("housekeeper.job_delete", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := h.tracer.Start(ctx, "housekeeper.job_delete")
		defer span.End()

		dest := &cron.GenericCronData{
			Attributes: map[string]string{},
		}
		if data.Data == nil {
			data.Data = &anypb.Any{}
		}

		if err := data.Data.UnmarshalTo(dest); err != nil {
			h.logger.Warn("failed to unmarshal housekeeper cron data", zap.Error(err))
		}

		tJobProps := table.FivenetJobProps

		stmt := tJobProps.
			SELECT(tJobProps.Job).
			WHERE(tJobProps.DeletedAt.IS_NOT_NULL()).
			LIMIT(10)

		var jobs []string
		if err := stmt.QueryContext(ctx, h.db, &jobs); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return fmt.Errorf("failed to query jobs: %w", err)
			}
		}

		tableListsMu.Lock()
		defer tableListsMu.Unlock()

		tablesList := h.getTablesListFn()

		keys := []string{}
		for key := range tablesList {
			keys = append(keys, key)
		}
		slices.Sort(keys)
		if len(keys) == 0 {
			return nil
		}

		lastTblKey, ok := dest.Attributes[lastTableMapIndex]
		if !ok {
			// Take first table
			lastTblKey = keys[0]
		} else {
			idx := slices.Index(keys, lastTblKey)
			if idx == -1 || len(keys) <= idx+1 {
				h.logger.Debug("last table key not found in keys, starting from the beginning again")
				lastTblKey = keys[0]
			} else {
				lastTblKey = keys[idx+1]
			}
		}

		tbl, ok := tablesList[lastTblKey]
		if !ok {
			return nil
		}

		jobName := "usms"

		markedRows, err := h.SoftDeleteJobData(ctx, tbl, jobName)
		if err != nil {
			return fmt.Errorf("failed to soft delete rows from main table fivenet_wiki_pages: %w", err)
		}
		fmt.Println("markedRows", markedRows)

		dest.Attributes[lastTableMapIndex] = lastTblKey
		dest.Attributes[lastJobName] = jobName

		if err := data.Data.MarshalFrom(dest); err != nil {
			return fmt.Errorf("failed to marshal updated housekeeper cron data. %w", err)
		}

		return nil
	}) */

	p.CronHandlers.Add("housekeeper.run", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := h.tracer.Start(ctx, "housekeeper.run")
		defer span.End()

		dest := &cron.GenericCronData{
			Attributes: map[string]string{},
		}
		if data.Data == nil {
			data.Data = &anypb.Any{}
		}

		if err := data.Data.UnmarshalTo(dest); err != nil {
			h.logger.Warn("failed to unmarshal housekeeper cron data", zap.Error(err))
		}

		if err := h.runHousekeeper(ctx, dest); err != nil {
			return fmt.Errorf("error during docstore workflow handling. %w", err)
		}

		if err := data.Data.MarshalFrom(dest); err != nil {
			return fmt.Errorf("failed to marshal updated housekeeper cron data. %w", err)
		}

		return nil
	})

	return h
}

func (h *Housekeeper) runHousekeeper(ctx context.Context, data *cron.GenericCronData) error {
	tableListsMu.Lock()
	defer tableListsMu.Unlock()

	tablesList := h.getTablesListFn()

	keys := []string{}
	for key := range tablesList {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	if len(keys) == 0 {
		return nil
	}

	lastTblKey, ok := data.Attributes[lastTableMapIndex]
	if !ok {
		// Take first table
		lastTblKey = keys[0]
	} else {
		idx := slices.Index(keys, lastTblKey)
		if idx == -1 || len(keys) <= idx+1 {
			h.logger.Debug("last table key not found in keys, starting from the beginning again")
			lastTblKey = keys[0]
		} else {
			lastTblKey = keys[idx+1]
		}
	}

	tbl, ok := tablesList[lastTblKey]
	if !ok {
		return nil
	}

	var condition jet.BoolExpression
	if tbl.DeletedAtColumn != nil {
		condition = jet.AND(
			tbl.DeletedAtColumn.IS_NOT_NULL(),
			tbl.DeletedAtColumn.LT_EQ(
				jet.CURRENT_DATE().SUB(jet.INTERVAL(tbl.MinDays, jet.DAY)),
			),
		)
	} else if tbl.DateColumn != nil {
		condition = jet.AND(
			tbl.DateColumn.IS_NOT_NULL(),
			tbl.DateColumn.LT_EQ(
				jet.CAST(
					jet.CURRENT_DATE().SUB(jet.INTERVAL(tbl.MinDays, jet.DAY)),
				).AS_DATE(),
			),
		)
	} else {
		condition = jet.AND(
			tbl.TimestampColumn.IS_NOT_NULL(),
			tbl.TimestampColumn.LT_EQ(
				jet.CURRENT_DATE().SUB(jet.INTERVAL(tbl.MinDays, jet.DAY)),
			),
		)
	}

	if tbl.Condition != nil {
		condition = condition.AND(tbl.Condition)
	}

	stmt := tbl.Table.
		DELETE().
		WHERE(condition).
		LIMIT(DefaultDeleteLimit)

	res, err := stmt.ExecContext(ctx, h.db)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected > 0 {
		h.logger.Info("housekeeper run deleted rows", zap.String("table", lastTblKey), zap.Int64("rows", rowsAffected))
	}

	data.Attributes[lastTableMapIndex] = lastTblKey

	return nil
}

// SoftDeleteJobData marks rows as deleted in the main table and its dependent tables
// by setting the DeletedAt column to the current timestamp.
func (h *Housekeeper) SoftDeleteJobData(ctx context.Context, table *Table, jobName string) (int64, error) {
	h.logger.Info("starting soft delete", zap.String("table", table.Table.TableName()), zap.String("job", jobName))

	markedRows := int64(0)

	if table.DeletedAtColumn != nil {
		// Mark rows as deleted in the current table for the given job
		if err := h.markRowsAsDeletedInJob(ctx, table, jobName); err != nil {
			return markedRows, fmt.Errorf("failed to soft delete rows from main table %s: %w", table.Table.TableName(), err)
		}
	}

	// Traverse dependencies
	for _, dep := range table.DependentTables {
		if dep.DeletedAtColumn == nil {
			continue
		}

		r, err := h.softDeleteJobData(ctx, table, dep, jobName)
		if err != nil {
			return markedRows, fmt.Errorf("failed to soft delete rows from dependent table %s: %w", dep.Table.TableName(), err)
		}
		markedRows += r
	}

	h.logger.Info("soft delete completed", zap.String("table", table.Table.TableName()), zap.String("job", jobName))
	return markedRows, nil
}

func (h *Housekeeper) softDeleteJobData(ctx context.Context, parent *Table, table *Table, jobName string) (int64, error) {
	h.logger.Info("starting soft delete", zap.String("table", parent.Table.TableName()), zap.String("job", jobName))

	markedRows := int64(0)

	if table.DeletedAtColumn != nil {
		// Mark rows as deleted in the current table for the given job
		r, err := h.markRowsAsDeleted(ctx, parent, table, jobName)
		if err != nil {
			return markedRows, fmt.Errorf("failed to soft delete rows from dependent table %s: %w", parent.Table.TableName(), err)
		}
		markedRows += r
	}

	// Traverse dependencies
	for _, child := range table.DependentTables {
		if child.DeletedAtColumn == nil {
			continue
		}

		r, err := h.markRowsAsDeleted(ctx, table, child, jobName)
		if err != nil {
			return markedRows, fmt.Errorf("failed to soft delete dependent rows from dependent table %s: %w", child.Table.TableName(), err)
		}
		markedRows += r
	}

	h.logger.Info("soft delete dependent completed", zap.String("table", parent.Table.TableName()), zap.String("job", jobName))
	return markedRows, nil
}

func (h *Housekeeper) markRowsAsDeletedInJob(ctx context.Context, table *Table, jobName string) error {
	condition := jet.AND(
		table.JobColumn.EQ(jet.String(jobName)),
		table.DeletedAtColumn.IS_NULL(),
	)

	stmt := table.Table.
		UPDATE().
		SET(
			table.DeletedAtColumn.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(condition).
		LIMIT(DefaultDeleteLimit)

	fmt.Println("stmt", stmt.DebugSql())
	return nil

	res, err := stmt.ExecContext(ctx, h.db)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	h.logger.Info("marked rows as deleted", zap.String("table", table.Table.TableName()), zap.Int64("rows", rowsAffected))
	return nil
}

func (h *Housekeeper) markRowsAsDeleted(ctx context.Context, parentTable *Table, table *Table, jobName string) (int64, error) {
	var condition jet.BoolExpression
	if table.JobColumn != nil {
		condition = table.JobColumn.EQ(jet.String(jobName))
	} else {
		condition = parentTable.JobColumn.EQ(jet.String(jobName))
	}

	condition = condition.AND(
		table.DeletedAtColumn.IS_NULL(),
	)
	condition = condition.AND(
		table.ForeignKey.IN(
			parentTable.Table.
				SELECT(parentTable.IDColumn).
				WHERE(jet.AND(
					parentTable.JobColumn.EQ(jet.String(jobName)),
					parentTable.DeletedAtColumn.IS_NULL(),
				)),
		),
	)

	stmt := table.Table.
		UPDATE().
		SET(
			table.DeletedAtColumn.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(condition).
		LIMIT(DefaultDeleteLimit)

	fmt.Println("stmt", stmt.DebugSql())
	return 0, nil

	res, err := stmt.ExecContext(ctx, h.db)
	if err != nil {
		return 0, err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	h.logger.Info("marked rows as deleted", zap.String("table", table.Table.TableName()), zap.Int64("rows", rowsAffected))
	return rowsAffected, nil
}

func (h *Housekeeper) HardDelete(ctx context.Context, table *Table) error {
	h.logger.Info("starting hard delete", zap.String("table", table.Table.TableName()))

	// Traverse dependencies
	for _, dep := range table.DependentTables {
		if err := h.hardDelete(ctx, table, dep); err != nil {
			return fmt.Errorf("failed to hard delete rows from dependent table %s: %w", dep.Table.TableName(), err)
		}
	}

	// Mark rows as deleted in the current table for the given job
	if err := h.deleteRows(ctx, table, nil, table.MinDays); err != nil {
		return fmt.Errorf("failed to hard delete rows from main table %s: %w", table.Table.TableName(), err)
	}

	h.logger.Info("hard delete completed", zap.String("table", table.Table.TableName()))
	return nil
}

func (h *Housekeeper) hardDelete(ctx context.Context, parent *Table, table *Table) error {
	h.logger.Info("starting hard delete", zap.String("table", parent.Table.TableName()))

	// Traverse dependencies
	for _, child := range table.DependentTables {
		if child.DeletedAtColumn == nil {
			continue
		}

		if err := h.deleteRows(ctx, table, child, table.MinDays); err != nil {
			return fmt.Errorf("failed to hard delete dependent rows from dependent table %s: %w", child.Table.TableName(), err)
		}
	}

	// Mark rows as deleted in the current table for the given job
	if err := h.deleteRows(ctx, parent, table, table.MinDays); err != nil {
		return fmt.Errorf("failed to hard delete rows from dependent table %s: %w", parent.Table.TableName(), err)
	}

	h.logger.Info("hard delete dependent completed", zap.String("table", parent.Table.TableName()))
	return nil
}

func (h *Housekeeper) deleteRows(ctx context.Context, table *Table, parent *Table, minDays int) error {
	var condition jet.BoolExpression
	if parent != nil {
		condition = jet.AND()
	}
	// TODO

	stmt := table.Table.
		DELETE().
		WHERE(condition).
		LIMIT(DefaultDeleteLimit)

	fmt.Println("stmt", stmt.DebugSql())
	return nil

	res, err := stmt.ExecContext(ctx, h.db)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	h.logger.Info("deleted rows", zap.String("table", table.Table.TableName()), zap.Int64("rows", rowsAffected))
	return nil
}
