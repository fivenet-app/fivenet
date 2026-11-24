package housekeeper

import (
	"context"
	"fmt"
	"slices"
	"strconv"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	"github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

// runHardDelete executes the hard delete cronjob logic for the housekeeper.
// It determines which table to process next, runs the hard delete, and updates the cron data.
func (h *Housekeeper) runHardDelete(
	ctx context.Context,
	data *cron.GenericCronData,
) (int64, error) {
	tablesList := h.getTablesListFn()

	// Collect and sort table keys for deterministic processing order.
	keys := []string{}
	for key := range tablesList {
		keys = append(keys, key)
	}
	if len(keys) == 0 {
		data.SetAttribute(currentTable, "")
		data.SetAttribute(backlogRemaining, "false")
		return 0, nil
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
		return 0, fmt.Errorf("no table config for key: %s", nextKey)
	}
	data.SetAttribute(currentTable, tbl.Table.TableName())

	rowsAffected, err := h.HardDelete(ctx, tbl)
	if err != nil {
		return 0, err
	}

	metricHardDeleteRowsAffected.WithLabelValues(tbl.Table.TableName()).Set(float64(rowsAffected))

	// A full batch means the cursor stays on this table for the next run. This
	// is intentionally a cheap progress signal: a full batch can only be
	// produced while there was at least one eligible row at the start of it.
	backlog := rowsAffected >= DefaultDeleteLimit
	data.SetAttribute(backlogRemaining, strconv.FormatBool(backlog))

	// Only update the last table key if less than the limit rows were affected.
	if !backlog {
		data.SetAttribute(lastTableMapIndex, nextKey)
	}

	return rowsAffected, nil
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
			return rowsAffected, fmt.Errorf(
				"failed to hard delete rows from dependant table %s. %w",
				dep.Table.TableName(),
				err,
			)
		}
		rowsAffected += r
	}

	// Delete rows in the current table.
	r, err := h.deleteRows(ctx, nil, table, table.MinDays)
	if err != nil {
		return rowsAffected, fmt.Errorf(
			"failed to hard delete rows from main table %s. %w",
			table.Table.TableName(),
			err,
		)
	}
	rowsAffected += r

	h.logger.Debug(
		"hard delete completed",
		zap.String("table", table.Table.TableName()),
		zap.Int64("rows", rowsAffected),
	)
	return rowsAffected, nil
}

// hardDelete recursively performs hard delete operations on dependant tables and then on the current table.
func (h *Housekeeper) hardDelete(
	ctx context.Context,
	parent *Table,
	table *Table,
	minDays int,
) (int64, error) {
	rowsAffected := int64(0)

	// Traverse dependencies and delete rows from child tables first.
	for _, child := range table.DependantTables {
		r, err := h.deleteRows(ctx, table, child, minDays)
		if err != nil {
			return rowsAffected, fmt.Errorf(
				"failed to hard delete dependant rows from dependant table %s. %w",
				child.Table.TableName(),
				err,
			)
		}
		rowsAffected += r
	}

	// Delete rows in the current dependant table.
	r, err := h.deleteRows(ctx, parent, table, minDays)
	if err != nil {
		return rowsAffected, fmt.Errorf(
			"failed to hard delete rows from dependant table %s. %w",
			parent.Table.TableName(),
			err,
		)
	}
	rowsAffected += r

	return rowsAffected, nil
}

// deleteRows deletes rows from the specified table (optionally filtered by parent) that are older than minDays.
// Returns the number of affected rows. If dryRun is enabled, no rows are actually deleted.
func (h *Housekeeper) deleteRows(
	ctx context.Context,
	parent *Table,
	table *Table,
	minDays int,
) (int64, error) {
	var condition mysql.BoolExpression

	if parent != nil {
		// If a parent is specified, only delete rows whose foreign key matches deleted/old rows in the parent.
		if table.ForeignKeyString != nil && parent.KeyColumnString != nil {
			parentKeys := parent.Table.
				SELECT(parent.KeyColumnString).
				WHERE(mysql.AND(
					parent.DeletedAtColumn.IS_NOT_NULL(),
					parent.DeletedAtColumn.LT_EQ(
						mysql.CURRENT_DATE().SUB(mysql.INTERVAL(minDays, mysql.DAY)),
					),
				))
			condition = mysql.AND(table.ForeignKeyString.IN(parentKeys))
		} else {
			parentKeys := parent.Table.
				SELECT(parent.IDColumn).
				WHERE(mysql.AND(
					parent.DeletedAtColumn.IS_NOT_NULL(),
					parent.DeletedAtColumn.LT_EQ(
						mysql.CURRENT_DATE().SUB(mysql.INTERVAL(minDays, mysql.DAY)),
					),
				))
			condition = mysql.AND(table.ForeignKey.IN(parentKeys))
		}
	} else {
		if table.DateColumn != nil {
			// Use DateColumn if available for filtering.
			condition = mysql.AND(
				table.DateColumn.IS_NOT_NULL(),
				table.DateColumn.LT_EQ(
					mysql.CAST(
						mysql.CURRENT_DATE().SUB(mysql.INTERVAL(minDays, mysql.DAY)),
					).AS_DATE(),
				),
			)
		} else {
			// Otherwise, use TimestampColumn or DeletedAtColumn.
			if table.TimestampColumn != nil {
				condition = mysql.AND(
					table.TimestampColumn.IS_NOT_NULL(),
					table.TimestampColumn.LT_EQ(
						mysql.CURRENT_TIMESTAMP().SUB(mysql.INTERVAL(minDays, mysql.DAY)),
					),
				)
			} else {
				condition = mysql.AND(
					table.DeletedAtColumn.IS_NOT_NULL(),
					table.DeletedAtColumn.LT_EQ(
						mysql.CURRENT_DATE().SUB(mysql.INTERVAL(minDays, mysql.DAY)),
					),
				)
			}
		}
	}

	stmt := table.Table.
		DELETE().
		WHERE(condition)

	// Retention deletes can be ordered by their eligibility timestamp. This
	// makes batches progress from the oldest rows first without requiring each
	// table registration to repeat an ordering declaration. Dependant-table
	// deletes are instead driven by a parent-table relationship, so an age
	// column on the child is not relevant to their eligibility.
	if parent == nil {
		if orderBy := retentionDeleteOrderBy(table); len(orderBy) > 0 {
			stmt = stmt.ORDER_BY(orderBy...)
		}
	}
	stmt = stmt.LIMIT(DefaultDeleteLimit)

	var rowsAffected int64

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

	h.logger.Debug(
		"hard deleted rows",
		zap.String("table", table.Table.TableName()),
		zap.Int64("rows", rowsAffected),
	)
	return rowsAffected, nil
}

// retentionDeleteOrderBy returns an oldest-first order that matches the age
// column selected by deleteRows. IDColumn is used as a deterministic tie-breaker
// when the registration provides one.
func retentionDeleteOrderBy(table *Table) []mysql.OrderByClause {
	orderBy := make([]mysql.OrderByClause, 0, 2)

	switch {
	case table.DateColumn != nil:
		orderBy = append(orderBy, table.DateColumn.ASC())
	case table.TimestampColumn != nil:
		orderBy = append(orderBy, table.TimestampColumn.ASC())
	case table.DeletedAtColumn != nil:
		orderBy = append(orderBy, table.DeletedAtColumn.ASC())
	default:
		return nil
	}

	if table.IDColumn != nil {
		orderBy = append(orderBy, table.IDColumn.ASC())
	}

	return orderBy
}
