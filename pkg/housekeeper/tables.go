package housekeeper

import (
	"fmt"
	"sync"

	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	// Mutexes to protect access to the table list maps
	tableListsMu = sync.Mutex{}

	tablesList = map[string]*Table{}
)

type Table struct {
	Table           jet.Table
	JobColumn       jet.ColumnString
	DeletedAtColumn jet.ColumnTimestamp // Optional column for soft delete
	ForeignKey      jet.ColumnInteger   // Optional for first/main table
	IDColumn        jet.ColumnInteger

	TimestampColumn jet.ColumnTimestamp
	DateColumn      jet.ColumnDate
	MinDays         int

	DependantTables []*Table // Allow tables to have their dependants
}

func AddTable(tbl *Table) {
	tableListsMu.Lock()
	defer tableListsMu.Unlock()

	if tbl.DeletedAtColumn == nil && tbl.TimestampColumn == nil && tbl.DateColumn == nil {
		panic(fmt.Sprintf("table %s must have a DeletedAt, TimestampColumn, or DateColumn column set for soft delete!", tbl.Table.TableName()))
	}

	ensureMinDays(tbl)

	tablesList[tbl.Table.TableName()] = tbl
}

func ensureMinDays(tbl *Table) {
	if tbl.MinDays < 30 {
		tbl.MinDays = 30
	}

	for _, t := range tbl.DependantTables {
		ensureMinDays(t)
	}
}
