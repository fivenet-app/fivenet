package housekeeper

import (
	"context"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

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

	rowsAffected, err := h.HardDelete(ctx, tbl)
	if err != nil {
		return err
	}

	metricHardDeleteAffectedRows.WithLabelValues(tbl.Table.TableName()).Set(float64(rowsAffected))

	// Only update the last table key if less than the limit rows were affected
	if rowsAffected < DefaultDeleteLimit {
		data.Attributes[lastTableMapIndex] = lastTblKey
	}

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
