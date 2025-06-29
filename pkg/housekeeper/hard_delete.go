package housekeeper

import (
	"context"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

// runHardDelete executes the hard delete cronjob logic for the housekeeper.
// It determines which table to process next, runs the hard delete, and updates the cron data.
func (h *Housekeeper) runHardDelete(ctx context.Context, data *cron.GenericCronData) error {
	tablesList := h.getTablesListFn()

	// Collect and sort table keys for deterministic processing order.
	keys := []string{}
	for key := range tablesList {
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		return nil
	}
	slices.Sort(keys)

	// Determine next table to process
	var currentIdx int
	lastKey := data.GetAttribute(lastTableMapIndex)

	if lastKey != "" {
		if idx := slices.Index(keys, lastKey); idx != -1 {
			currentIdx = (idx + 1) % len(keys)
		} else {
			// Unknown key, start from the beginning
			currentIdx = 0
		}
	}

	nextKey := keys[currentIdx]

	tbl, ok := tablesList[nextKey]
	if !ok {
		return fmt.Errorf("no table config for key: %s", nextKey)
	}

	rowsAffected, err := h.HardDelete(ctx, tbl)
	if err != nil {
		return err
	}

	metricHardDeleteAffectedRows.WithLabelValues(tbl.Table.TableName()).Set(float64(rowsAffected))

	// Only update the last table key if less than the limit rows were affected
	if rowsAffected < DefaultDeleteLimit {
		data.SetAttribute(lastTableMapIndex, nextKey)
	}

	return nil
}

// HardDelete performs a hard delete operation on the given table and its dependant tables.
// It traverses dependencies, deletes rows, and returns the total number of affected rows.
func (h *Housekeeper) HardDelete(ctx context.Context, table *Table) (int64, error) {
	h.logger.Debug("starting hard delete", zap.String("table", table.Table.TableName()))

	rowsAffected := int64(0)

	// Traverse dependencies and perform hard delete on each dependant table.
	for _, dep := range table.DependantTables {
		r, err := h.hardDelete(ctx, table, dep, table.MinDays)
		if err != nil {
			return rowsAffected, fmt.Errorf("failed to hard delete rows from dependant table %s. %w", dep.Table.TableName(), err)
		}
		rowsAffected += r
	}

	// Delete rows in the current table.
	r, err := h.deleteRows(ctx, nil, table, table.MinDays)
	if err != nil {
		return rowsAffected, fmt.Errorf("failed to hard delete rows from main table %s. %w", table.Table.TableName(), err)
	}
	rowsAffected += r

	h.logger.Debug("hard delete completed", zap.String("table", table.Table.TableName()), zap.Int64("rows", rowsAffected))
	return rowsAffected, nil
}

// hardDelete recursively performs hard delete operations on dependant tables and then on the current table.
func (h *Housekeeper) hardDelete(ctx context.Context, parent *Table, table *Table, minDays int) (int64, error) {
	rowsAffected := int64(0)

	// Traverse dependencies and delete rows from child tables first.
	for _, child := range table.DependantTables {
		r, err := h.deleteRows(ctx, table, child, minDays)
		if err != nil {
			return rowsAffected, fmt.Errorf("failed to hard delete dependant rows from dependant table %s. %w", child.Table.TableName(), err)
		}
		rowsAffected += r
	}

	// Delete rows in the current dependant table.
	r, err := h.deleteRows(ctx, parent, table, minDays)
	if err != nil {
		return rowsAffected, fmt.Errorf("failed to hard delete rows from dependant table %s. %w", parent.Table.TableName(), err)
	}
	rowsAffected += r

	return rowsAffected, nil
}

// deleteRows deletes rows from the specified table (optionally filtered by parent) that are older than minDays.
// Returns the number of affected rows. If dryRun is enabled, no rows are actually deleted.
func (h *Housekeeper) deleteRows(ctx context.Context, parent *Table, table *Table, minDays int) (rowsAffected int64, err error) {
	var condition jet.BoolExpression

	if parent != nil {
		// If a parent is specified, only delete rows whose foreign key matches deleted/old rows in the parent.
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
			// Use DateColumn if available for filtering.
			condition = jet.AND(
				table.DateColumn.IS_NOT_NULL(),
				table.DateColumn.LT_EQ(
					jet.CAST(
						jet.CURRENT_DATE().SUB(jet.INTERVAL(minDays, jet.DAY)),
					).AS_DATE(),
				),
			)
		} else {
			// Otherwise, use TimestampColumn or DeletedAtColumn.
			if table.TimestampColumn != nil {
				condition = jet.AND(
					table.TimestampColumn.IS_NOT_NULL(),
					table.TimestampColumn.LT_EQ(
						jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(minDays, jet.DAY)),
					),
				)
			} else {
				condition = jet.AND(
					table.DeletedAtColumn.IS_NOT_NULL(),
					table.DeletedAtColumn.LT_EQ(
						jet.CURRENT_DATE().SUB(jet.INTERVAL(minDays, jet.DAY)),
					),
				)
			}
		}
	}

	stmt := table.Table.
		DELETE().
		WHERE(condition).
		LIMIT(DefaultDeleteLimit)

	if !h.dryRun {
		// Execute the delete statement if not in dry run mode.
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
