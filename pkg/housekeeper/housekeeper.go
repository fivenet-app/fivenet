package housekeeper

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
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
const (
	lastJobName       = "last_job_name"
	lastTableMapIndex = "last_key"
)

type Housekeeper struct {
	logger *zap.Logger
	tracer trace.Tracer

	db *sql.DB

	getTablesListFn func() map[string]*Table

	dryRun bool
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	DB     *sql.DB
	TP     *tracesdk.TracerProvider
}

type Result struct {
	fx.Out

	Housekeeper  *Housekeeper
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

func New(p Params) Result {
	h := &Housekeeper{
		logger: p.Logger.Named("housekeeper"),
		tracer: p.TP.Tracer("housekeeper"),

		db: p.DB,

		getTablesListFn: func() map[string]*Table {
			return tablesList
		},

		dryRun: false,
	}

	return Result{
		Housekeeper:  h,
		CronRegister: h,
	}
}

func (s *Housekeeper) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "housekeeper.run",
		Schedule: "* * * * *", // Every minute
		Timeout:  durationpb.New(2 * time.Second),
	}); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "housekeeper.job_delete",
		Schedule: "* * * * *", // Every minute
		Timeout:  durationpb.New(45 * time.Second),
	}); err != nil {
		return err
	}

	return nil
}

func (h *Housekeeper) RegisterCronjobHandlers(hand *croner.Handlers) error {
	hand.Add("housekeeper.job_delete", func(ctx context.Context, data *cron.CronjobData) error {
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

		if err := h.runJobSoftDelete(ctx, dest); err != nil {
			return fmt.Errorf("error during housekeeper (soft delete). %w", err)
		}

		if err := data.Data.MarshalFrom(dest); err != nil {
			return fmt.Errorf("failed to marshal updated housekeeper (soft delete) cron data. %w", err)
		}

		return nil
	})

	hand.Add("housekeeper.run", func(ctx context.Context, data *cron.CronjobData) error {
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

		if err := h.runHardDelete(ctx, dest); err != nil {
			return fmt.Errorf("error during housekeeper (hard delete). %w", err)
		}

		if err := data.Data.MarshalFrom(dest); err != nil {
			return fmt.Errorf("failed to marshal updated housekeeper (hard delete) cron data. %w", err)
		}

		return nil
	})

	return nil
}

func (h *Housekeeper) runJobSoftDelete(ctx context.Context, data *cron.GenericCronData) error {
	tJobProps := table.FivenetJobProps

	stmt := tJobProps.
		SELECT(tJobProps.Job).
		WHERE(jet.AND(
			tJobProps.DeletedAt.IS_NOT_NULL(),
			tJobProps.DeletedAt.LT_EQ(jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(1, jet.DAY))),
		)).
		ORDER_BY(tJobProps.Job.ASC())

	var jobs []string
	if err := stmt.QueryContext(ctx, h.db, &jobs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query jobs. %w", err)
		}
	}

	if len(jobs) == 0 {
		h.logger.Debug("no jobs found to soft delete")
		return nil
	}

	tablesList := h.getTablesListFn()

	keys := []string{}
	for key := range tablesList {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	if len(keys) == 0 {
		return nil
	}

	jobName := jobs[0]
	lastJob, ok := data.Attributes[lastJobName]
	if ok {
		if lastJob != jobName {
			h.logger.Debug("job name changed, starting from the beginning again")
			data.Attributes[lastJobName] = jobName
			data.Attributes[lastTableMapIndex] = keys[0]
		}
	}

	lastTblKey, ok := data.Attributes[lastTableMapIndex]
	if !ok {
		// Take first table
		lastTblKey = keys[0]
	} else {
		idx := slices.Index(keys, lastTblKey)
		if idx == -1 || len(keys) <= idx+1 {
			h.logger.Debug("next table key not found in keys, starting from the beginning again")
			lastTblKey = keys[0]

			// All job's data should be deleted now, delete the job props
			stmt := tJobProps.
				DELETE().
				WHERE(tJobProps.Job.EQ(jet.String(jobName))).
				LIMIT(1)

			if !h.dryRun {
				if _, err := stmt.ExecContext(ctx, h.db); err != nil {
					return fmt.Errorf("failed to delete job %s. %w", jobName, err)
				}
			} else {
				h.logger.Debug("dry run delete job props statement", zap.String("query", stmt.DebugSql()))
			}
		} else {
			lastTblKey = keys[idx+1]
		}
	}

	tbl, ok := tablesList[lastTblKey]
	if !ok {
		return nil
	}

	rowsAffected, err := h.SoftDeleteJobData(ctx, tbl, jobName)
	if err != nil {
		return fmt.Errorf("failed to soft delete rows for table %s (job: %s). %w", tbl.Table.TableName(), jobName, err)
	}

	// Only update the last table key if no rows were affected (empty table)
	if rowsAffected == 0 {
		data.Attributes[lastTableMapIndex] = lastTblKey
	}

	return nil
}

// SoftDeleteJobData marks rows as deleted in the main table and its dependant tables
// by setting the DeletedAt column to the current timestamp.
func (h *Housekeeper) SoftDeleteJobData(ctx context.Context, table *Table, jobName string) (int64, error) {
	h.logger.Debug("starting soft delete", zap.String("table", table.Table.TableName()), zap.String("job", jobName))

	rowsAffected := int64(0)

	if table.DeletedAtColumn != nil && table.JobColumn != nil {
		// Mark rows as deleted in the current table for the given job
		r, err := h.markRowsAsDeletedInJob(ctx, table, jobName)
		if err != nil {
			return rowsAffected, fmt.Errorf("failed to soft delete rows from main table %s: %w", table.Table.TableName(), err)
		}
		rowsAffected += r
	}

	// Traverse dependencies
	for _, dep := range table.DependantTables {
		if dep.DeletedAtColumn == nil && table.JobColumn == nil {
			continue
		}

		r, err := h.softDeleteJobData(ctx, table, dep, jobName)
		if err != nil {
			return rowsAffected, fmt.Errorf("failed to soft delete rows from dependant table %s: %w", dep.Table.TableName(), err)
		}
		rowsAffected += r
	}

	h.logger.Debug("soft delete completed", zap.String("table", table.Table.TableName()), zap.String("job", jobName), zap.Int64("rows", rowsAffected))
	return rowsAffected, nil
}

func (h *Housekeeper) softDeleteJobData(ctx context.Context, parent *Table, table *Table, jobName string) (int64, error) {
	rowsAffected := int64(0)

	if table.DeletedAtColumn != nil && parent.JobColumn != nil {
		// Mark rows as deleted in the current table for the given job
		r, err := h.markRowsAsDeleted(ctx, parent, table, jobName)
		if err != nil {
			return rowsAffected, fmt.Errorf("failed to soft delete rows from dependant table %s: %w", parent.Table.TableName(), err)
		}
		rowsAffected += r
	}

	// Traverse dependencies
	for _, child := range table.DependantTables {
		if table.JobColumn == nil || child.DeletedAtColumn == nil {
			continue
		}

		r, err := h.markRowsAsDeleted(ctx, table, child, jobName)
		if err != nil {
			return rowsAffected, fmt.Errorf("failed to soft delete dependant rows from dependant table %s: %w", child.Table.TableName(), err)
		}
		rowsAffected += r
	}

	return rowsAffected, nil
}

func (h *Housekeeper) markRowsAsDeletedInJob(ctx context.Context, table *Table, jobName string) (rowsAffected int64, err error) {
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

	if !h.dryRun {
		res, err := stmt.ExecContext(ctx, h.db)
		if err != nil {
			return rowsAffected, err
		}

		rowsAffected, err = res.RowsAffected()
		if err != nil {
			return rowsAffected, err
		}
	} else {
		h.logger.Debug("dry run markRowsAsDeleted statement", zap.String("query", stmt.DebugSql()))
	}

	h.logger.Debug("soft deleted rows", zap.String("table", table.Table.TableName()), zap.Int64("rows", rowsAffected))
	return rowsAffected, nil
}

func (h *Housekeeper) markRowsAsDeleted(ctx context.Context, parentTable *Table, table *Table, jobName string) (rowsAffected int64, err error) {
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

	if !h.dryRun {
		res, err := stmt.ExecContext(ctx, h.db)
		if err != nil {
			return rowsAffected, err
		}

		rowsAffected, err = res.RowsAffected()
		if err != nil {
			return rowsAffected, err
		}
	} else {
		h.logger.Debug("dry run markRowsAsDeleted statement", zap.String("query", stmt.DebugSql()))
	}

	h.logger.Debug("soft deleted rows", zap.String("table", table.Table.TableName()), zap.Int64("rows", rowsAffected))
	return rowsAffected, nil
}

func (h *Housekeeper) runHardDelete(ctx context.Context, data *cron.GenericCronData) error {
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

	_, err := h.HardDelete(ctx, tbl)
	if err != nil {
		return err
	}

	data.Attributes[lastTableMapIndex] = lastTblKey

	return nil
}

func (h *Housekeeper) HardDelete(ctx context.Context, table *Table) (int64, error) {
	h.logger.Debug("starting hard delete", zap.String("table", table.Table.TableName()))

	rowsAffected := int64(0)

	// Traverse dependencies
	for _, dep := range table.DependantTables {
		r, err := h.hardDelete(ctx, table, dep, table.MinDays)
		if err != nil {
			return rowsAffected, fmt.Errorf("failed to hard delete rows from dependant table %s: %w", dep.Table.TableName(), err)
		}
		rowsAffected += r
	}

	// Mark rows as deleted in the current table for the given job
	r, err := h.deleteRows(ctx, nil, table, table.MinDays)
	if err != nil {
		return rowsAffected, fmt.Errorf("failed to hard delete rows from main table %s: %w", table.Table.TableName(), err)
	}
	rowsAffected += r

	h.logger.Debug("hard delete completed", zap.String("table", table.Table.TableName()), zap.Int64("rows", rowsAffected))
	return rowsAffected, nil
}

func (h *Housekeeper) hardDelete(ctx context.Context, parent *Table, table *Table, minDays int) (int64, error) {
	rowsAffected := int64(0)

	// Traverse dependencies
	for _, child := range table.DependantTables {
		r, err := h.deleteRows(ctx, table, child, minDays)
		if err != nil {
			return rowsAffected, fmt.Errorf("failed to hard delete dependant rows from dependant table %s: %w", child.Table.TableName(), err)
		}
		rowsAffected += r
	}

	// Mark rows as deleted in the current table for the given job
	r, err := h.deleteRows(ctx, parent, table, minDays)
	if err != nil {
		return rowsAffected, fmt.Errorf("failed to hard delete rows from dependant table %s: %w", parent.Table.TableName(), err)
	}
	rowsAffected += r

	return rowsAffected, nil
}

func (h *Housekeeper) deleteRows(ctx context.Context, parent *Table, table *Table, minDays int) (rowsAffected int64, err error) {
	var condition jet.BoolExpression

	if parent != nil {
		condition = jet.AND(
			table.ForeignKey.IN(
				parent.Table.
					SELECT(
						parent.IDColumn,
					).
					WHERE(jet.AND(
						parent.DeletedAtColumn.IS_NOT_NULL(),
						parent.DeletedAtColumn.LT_EQ(
							jet.CURRENT_DATE().SUB(jet.INTERVAL(minDays, jet.DAY)),
						),
					)),
			),
		)
	} else {
		if table.DateColumn != nil {
			condition = jet.AND(
				table.DateColumn.IS_NOT_NULL(),
				table.DateColumn.LT_EQ(
					jet.CAST(
						jet.CURRENT_DATE().SUB(jet.INTERVAL(minDays, jet.DAY)),
					).AS_DATE(),
				),
			)
		} else {
			var col jet.ColumnTimestamp
			if table.TimestampColumn != nil {
				col = table.TimestampColumn
			} else {
				col = table.DeletedAtColumn
			}
			condition = jet.AND(
				col.IS_NOT_NULL(),
				col.LT_EQ(
					jet.CURRENT_DATE().SUB(jet.INTERVAL(minDays, jet.DAY)),
				),
			)
		}
	}

	stmt := table.Table.
		DELETE().
		WHERE(condition).
		LIMIT(DefaultDeleteLimit)

	if !h.dryRun {
		res, err := stmt.ExecContext(ctx, h.db)
		if err != nil {
			return 0, err
		}

		rowsAffected, err = res.RowsAffected()
		if err != nil {
			return 0, err
		}
	} else {
		h.logger.Debug("dry run deleteRows statement", zap.String("query", stmt.DebugSql()))
	}

	h.logger.Debug("hard deleted rows", zap.String("table", table.Table.TableName()), zap.Int64("rows", rowsAffected))
	return rowsAffected, nil
}
