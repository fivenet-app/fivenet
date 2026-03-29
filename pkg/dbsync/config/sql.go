package dbsyncconfig

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
)

func getWhereCondition(
	table DBSyncTable,
	state *TableSyncState,
	cursorIDColumn string,
) string {
	if state == nil {
		return ""
	}

	lastCheck := state.GetLastCheck()
	lastID := state.GetLastID()

	hasCursorID := cursorIDColumn != "" && lastID != nil && strings.TrimSpace(*lastID) != ""
	hasUpdatedAt := table.UpdatedTimeColumn != nil &&
		*table.UpdatedTimeColumn != "" &&
		lastCheck != nil &&
		!lastCheck.IsZero()

	// Only quote the cursor ID column if it doesn't already contain backticks, to avoid double quoting.
	if !strings.Contains(cursorIDColumn, "`") {
		cursorIDColumn = fmt.Sprintf("%#q", cursorIDColumn)
	}

	if hasUpdatedAt {
		updatedAtCol := *table.UpdatedTimeColumn
		updatedAtValue := lastCheck.Format("2006-01-02 15:04:05.000")

		// Only quote the updated time column if it doesn't already contain backticks, to avoid double quoting.
		if !strings.Contains(updatedAtCol, "`") {
			updatedAtCol = fmt.Sprintf("%#q", updatedAtCol)
		}

		if hasCursorID {
			return fmt.Sprintf(
				"(%s > '%s' OR (%s = '%s' AND %s > %s))\n",
				updatedAtCol, updatedAtValue,
				updatedAtCol, updatedAtValue,
				cursorIDColumn, sqlLiteral(*lastID),
			)
		}

		return fmt.Sprintf("%s >= '%s'\n", updatedAtCol, updatedAtValue)
	}

	if hasCursorID {
		return fmt.Sprintf("%s > %s\n", cursorIDColumn, sqlLiteral(*lastID))
	}

	return ""
}

func sqlLiteral(value string) string {
	if _, err := strconv.ParseInt(value, 10, 64); err == nil {
		return value
	}
	if _, err := strconv.ParseUint(value, 10, 64); err == nil {
		return value
	}

	return fmt.Sprintf("'%s'", strings.ReplaceAll(value, "'", "''"))
}

func prepareStringQuery(
	q string,
	table DBSyncTable,
	state *TableSyncState,
	limit int64,
	cursorIDColumn string,
) string {
	limitStr := strconv.FormatInt(limit, 10)

	q = strings.ReplaceAll(q, "$limit", limitStr)

	where := getWhereCondition(table, state, cursorIDColumn)
	if where != "" {
		// Prepend "WHERE " if there is a condition
		where = "WHERE " + where
	}
	q = strings.ReplaceAll(q, "$whereCondition", where)

	return q
}

func prepareStringQueryWithWhereCondition(
	q string,
	limit int64,
	whereCondition []string,
) string {
	limitStr := strconv.FormatInt(limit, 10)

	q = strings.ReplaceAll(q, "$limit", limitStr)

	whereCondition = slices.DeleteFunc(whereCondition, func(c string) bool {
		return strings.TrimSpace(c) == ""
	})

	where := ""
	if len(whereCondition) > 0 {
		where = "WHERE " + strings.Join(whereCondition, " AND ") + "\n"
	}

	q = strings.ReplaceAll(q, "$whereCondition", where)

	return q
}

func buildQueryFromColumns(
	tableName string,
	columns map[string]string,
	whereCondition []string,
	limit int64,
	orderByColumns []string,
) string {
	columnsList := []string{}
	keys := make([]string, 0, len(columns))
	for k := range columns {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	for _, alias := range keys {
		column := columns[alias]
		if column == "" || column == "-" {
			continue
		}

		columnsList = append(columnsList, fmt.Sprintf("%#q AS %#q", column, alias))
	}

	q := fmt.Sprintf("SELECT %s\nFROM %#q\n",
		strings.Join(columnsList, ", "),
		tableName,
	)
	whereCondition = slices.DeleteFunc(whereCondition, func(c string) bool {
		return strings.TrimSpace(c) == ""
	})
	if len(whereCondition) > 0 {
		q += "WHERE " + strings.Join(whereCondition, " AND ") + "\n"
	}

	if len(orderByColumns) > 0 {
		q += "ORDER BY " + strings.Join(orderByColumns, ", ") + "\n"
	}

	if limit > 0 {
		q += fmt.Sprintf("LIMIT %d;", limit)
	} else {
		q += ";"
	}

	return q
}
