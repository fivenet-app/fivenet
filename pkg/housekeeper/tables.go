// Package housekeeper provides table definitions and helpers for background cleanup jobs.
package housekeeper

import (
	"fmt"
	"sync"

	"github.com/go-jet/jet/v2/mysql"
)

var (
	// tableListsMu is a mutex to protect access to the tablesList map.
	tableListsMu = sync.Mutex{}

	// tablesList holds all registered cleanup policies. The key includes the
	// retention column so one physical table can have more than one policy.
	tablesList = map[string]*Table{}
)

// Table represents a database table and its housekeeping-relevant columns and relationships.
type Table struct {
	// Table is the Jet representation of the SQL table.
	Table mysql.Table
	// JobColumn is the column used to identify jobs for this table.
	JobColumn mysql.ColumnString
	// DeletedAtColumn is an optional timestamp column for soft deletes.
	DeletedAtColumn mysql.ColumnTimestamp
	// ForeignKey is an optional foreign key column for dependent tables.
	ForeignKey mysql.ColumnInteger
	// IDColumn is the primary key column for the table.
	IDColumn mysql.ColumnInteger
	// ForeignKeyString and KeyColumnString support dependencies joined by a
	// string key, such as property definitions and property values.
	ForeignKeyString mysql.ColumnString
	KeyColumnString  mysql.ColumnString

	// TimestampColumn is an optional timestamp column for filtering by age.
	TimestampColumn mysql.ColumnTimestamp
	// DateColumn is an optional date column for filtering by age.
	DateColumn mysql.ColumnDate
	// MinDays is the minimum age in days before a row is eligible for deletion.
	MinDays int

	// DependantTables lists tables that depend on this table (for cascading deletes).
	DependantTables []*Table
}

// AddTable registers a new table for housekeeping. Panics if the table is missing required columns for soft delete.
func AddTable(tbl *Table) {
	tableListsMu.Lock()
	defer tableListsMu.Unlock()

	// Ensure the table has at least one column for soft delete logic.
	if tbl.DeletedAtColumn == nil && tbl.TimestampColumn == nil && tbl.DateColumn == nil {
		panic(
			fmt.Sprintf(
				"table %s must have a DeletedAt, TimestampColumn, or DateColumn column set for soft delete!",
				tbl.Table.TableName(),
			),
		)
	}

	// Ensure minimum days is set for the table and its dependants.
	ensureMinDays(tbl)

	tablesList[tableRegistrationKey(tbl)] = tbl
}

// tableRegistrationKey identifies a cleanup policy by its physical table and
// the age column used to select rows. This mirrors deleteRows' precedence when
// more than one age column is set on a registration.
func tableRegistrationKey(tbl *Table) string {
	var columnName string
	switch {
	case tbl.DateColumn != nil:
		columnName = tbl.DateColumn.Name()
	case tbl.TimestampColumn != nil:
		columnName = tbl.TimestampColumn.Name()
	case tbl.DeletedAtColumn != nil:
		columnName = tbl.DeletedAtColumn.Name()
	default:
		return tbl.Table.TableName()
	}

	return tbl.Table.TableName() + ":" + columnName
}

// ensureMinDays sets a default MinDays if not set, and applies recursively to dependant tables.
func ensureMinDays(tbl *Table) {
	if tbl.MinDays <= 0 {
		tbl.MinDays = 30
	}

	for _, t := range tbl.DependantTables {
		ensureMinDays(t)
	}
}
