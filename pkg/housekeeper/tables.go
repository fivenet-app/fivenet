package housekeeper

import (
	"sync"

	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	tablesMu   = sync.Mutex{}
	tablesList = map[string]*Table{}
)

type Table struct {
	Table           jet.Table
	TimestampColumn jet.ColumnTimestamp
	Condition       jet.BoolExpression
	MinDays         int
}

func AddTable(tbl *Table) {
	tablesMu.Lock()
	defer tablesMu.Unlock()

	if tbl.MinDays < 30 {
		tbl.MinDays = 30
	}

	tablesList[tbl.Table.TableName()] = tbl
}
