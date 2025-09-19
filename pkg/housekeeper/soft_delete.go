package housekeeper

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

// runJobSoftDelete executes the soft delete cronjob logic for the housekeeper.
// It processes jobs marked for deletion, iterates over tables, and updates the cron data.
func (h *Housekeeper) runJobSoftDelete(ctx context.Context, data *cron.GenericCronData) error {
	tJobProps := table.FivenetJobProps

	// Query jobs that are marked for deletion
	stmt := tJobProps.
		SELECT(tJobProps.Job).
		WHERE(mysql.AND(
			tJobProps.DeletedAt.IS_NOT_NULL(),
			tJobProps.DeletedAt.LT_EQ(mysql.CURRENT_TIMESTAMP().SUB(mysql.INTERVAL(1, mysql.DAY))),
		)).
		ORDER_BY(tJobProps.Job.ASC())

	var jobs []string
	if err := stmt.QueryContext(ctx, h.db, &jobs); err != nil && !errors.Is(err, qrm.ErrNoRows) {
		return fmt.Errorf("failed to query jobs. %w", err)
	}

	if len(jobs) == 0 {
		h.logger.Debug("no jobs found to soft delete")
		return nil
	}

	jobName := jobs[0]

	tables := h.getTablesListFn()
	if len(tables) == 0 {
		return nil
	}

	keys := make([]string, 0, len(tables))
	for k := range tables {
		keys = append(keys, k)
	}
	slices.Sort(keys)

	// Reset table progress if job changed
	if prevJob := data.GetAttribute(lastJobName); prevJob != jobName {
		h.logger.Debug("job name changed, starting from beginning", zap.String("job", jobName))
		data.SetAttribute(lastJobName, jobName)
		data.SetAttribute(lastTableMapIndex, keys[0])
	}

	lastTblKey := data.GetAttribute(lastTableMapIndex)
	currentIdx := slices.Index(keys, lastTblKey)
	if currentIdx == -1 {
		currentIdx = 0
	}

	nextTblKey := keys[currentIdx]
	tbl := tables[nextTblKey]
	if tbl == nil {
		return fmt.Errorf("no table found for key: %s", nextTblKey)
	}

	rowsAffected, err := h.SoftDeleteJobData(ctx, tbl, jobName)
	if err != nil {
		return fmt.Errorf(
			"failed to soft delete rows for table %s (job: %s). %w",
			tbl.Table.TableName(),
			jobName,
			err,
		)
	}

	metricSoftDeleteAffectedRows.Set(float64(rowsAffected))

	// Advance table key only if this table had no more rows to delete
	if rowsAffected < DefaultDeleteLimit {
		// If we've reached the end of the table list, delete the job prop
		if currentIdx+1 >= len(keys) {
			delStmt := tJobProps.
				DELETE().
				WHERE(tJobProps.Job.EQ(mysql.String(jobName))).
				LIMIT(1)

			if !h.dryRun {
				if _, err := delStmt.ExecContext(ctx, h.db); err != nil {
					return fmt.Errorf("failed to delete job %s. %w", jobName, err)
				}
			} else {
				h.logger.Debug("dry run: delete job props", zap.String("query", delStmt.DebugSql()))
			}
			// Clear progress
			data.SetAttribute(lastJobName, "")
			data.SetAttribute(lastTableMapIndex, "")
		} else {
			data.SetAttribute(lastTableMapIndex, keys[currentIdx+1])
		}
	}

	return nil
}

// SoftDeleteJobData marks rows as deleted in the main table and its dependant tables
// by setting the DeletedAt column to the current timestamp for the given job.
func (h *Housekeeper) SoftDeleteJobData(
	ctx context.Context,
	table *Table,
	jobName string,
) (int64, error) {
	h.logger.Debug(
		"starting soft delete",
		zap.String("table", table.Table.TableName()),
		zap.String("job", jobName),
	)

	rowsAffected := int64(0)

	if table.DeletedAtColumn != nil && table.JobColumn != nil {
		// Mark rows as deleted in the current table for the given job
		r, err := h.markRowsAsDeletedInJob(ctx, table, jobName)
		if err != nil {
			return rowsAffected, fmt.Errorf(
				"failed to soft delete rows from main table %s. %w",
				table.Table.TableName(),
				err,
			)
		}
		rowsAffected += r
	}

	// Traverse dependencies and soft delete in dependant tables.
	for _, dep := range table.DependantTables {
		if dep.DeletedAtColumn == nil && table.JobColumn == nil {
			continue
		}

		r, err := h.softDeleteJobData(ctx, table, dep, jobName)
		if err != nil {
			return rowsAffected, fmt.Errorf(
				"failed to soft delete rows from dependant table %s. %w",
				dep.Table.TableName(),
				err,
			)
		}
		rowsAffected += r
	}

	h.logger.Debug(
		"soft delete completed",
		zap.String("table", table.Table.TableName()),
		zap.String("job", jobName),
		zap.Int64("rows", rowsAffected),
	)
	return rowsAffected, nil
}

// softDeleteJobData recursively marks rows as deleted in dependant tables for the given job.
func (h *Housekeeper) softDeleteJobData(
	ctx context.Context,
	parent *Table,
	table *Table,
	jobName string,
) (int64, error) {
	rowsAffected := int64(0)

	if table.DeletedAtColumn != nil && parent.JobColumn != nil {
		// Mark rows as deleted in the current table for the given job
		r, err := h.markRowsAsDeleted(ctx, parent, table, jobName)
		if err != nil {
			return rowsAffected, fmt.Errorf(
				"failed to soft delete rows from dependant table %s. %w",
				parent.Table.TableName(),
				err,
			)
		}
		rowsAffected += r
	}

	// Traverse dependencies and soft delete in child tables.
	for _, child := range table.DependantTables {
		if table.JobColumn == nil || child.DeletedAtColumn == nil {
			continue
		}

		r, err := h.markRowsAsDeleted(ctx, table, child, jobName)
		if err != nil {
			return rowsAffected, fmt.Errorf(
				"failed to soft delete dependant rows from dependant table %s. %w",
				child.Table.TableName(),
				err,
			)
		}
		rowsAffected += r
	}

	return rowsAffected, nil
}

// markRowsAsDeletedInJob sets the DeletedAt column to the current timestamp for all rows in the table
// matching the given job name and not already deleted.
func (h *Housekeeper) markRowsAsDeletedInJob(
	ctx context.Context,
	table *Table,
	jobName string,
) (int64, error) {
	condition := mysql.AND(
		table.JobColumn.EQ(mysql.String(jobName)),
		table.DeletedAtColumn.IS_NULL(),
	)

	stmt := table.Table.
		UPDATE().
		SET(
			table.DeletedAtColumn.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(condition).
		LIMIT(DefaultDeleteLimit)

	var rowsAffected int64
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

	h.logger.Debug(
		"soft deleted rows",
		zap.String("table", table.Table.TableName()),
		zap.Int64("rows", rowsAffected),
	)
	return rowsAffected, nil
}

// markRowsAsDeleted sets the DeletedAt column to the current timestamp for all rows in the table
// matching the given job name and not already deleted, filtered by parent table.
func (h *Housekeeper) markRowsAsDeleted(
	ctx context.Context,
	parentTable *Table,
	table *Table,
	jobName string,
) (int64, error) {
	var condition mysql.BoolExpression
	if table.JobColumn != nil {
		condition = table.JobColumn.EQ(mysql.String(jobName))
	} else {
		condition = parentTable.JobColumn.EQ(mysql.String(jobName))
	}

	condition = condition.AND(mysql.AND(
		table.DeletedAtColumn.IS_NULL(),
		table.ForeignKey.IN(
			parentTable.Table.
				SELECT(parentTable.IDColumn).
				WHERE(mysql.AND(
					parentTable.JobColumn.EQ(mysql.String(jobName)),
					parentTable.DeletedAtColumn.IS_NULL(),
				)),
		),
	))

	stmt := table.Table.
		UPDATE().
		SET(
			table.DeletedAtColumn.SET(mysql.CURRENT_TIMESTAMP()),
		).
		WHERE(condition).
		LIMIT(DefaultDeleteLimit)

	var rowsAffected int64
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

	h.logger.Debug(
		"soft deleted rows",
		zap.String("table", table.Table.TableName()),
		zap.Int64("rows", rowsAffected),
	)
	return rowsAffected, nil
}
