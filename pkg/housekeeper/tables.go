package housekeeper

import (
	"sync"

	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	// Mutexes to protect access to the table list maps
	tableListsMu   = sync.Mutex{}
	fkTableListsMu = sync.Mutex{}

	tablesList   = map[string]*Table{}
	fkTablesList = map[string]*JobTable{}
)

type Table struct {
	Table           jet.Table
	TimestampColumn jet.ColumnTimestamp
	DateColumn      jet.ColumnDate
	Condition       jet.BoolExpression
	MinDays         int
}

func AddTable(tbls ...*Table) {
	tableListsMu.Lock()
	defer tableListsMu.Unlock()

	for _, tbl := range tbls {
		if tbl.MinDays < 30 {
			tbl.MinDays = 30
		}

		tablesList[tbl.Table.TableName()] = tbl
	}
}

type JobTable struct {
	Source *JobTableSource

	// Table which relies on the source ID column
	TargetTable jet.Table

	TargetSourceIDColumn jet.ColumnInteger
	// To be able to (soft) delete the target rows (if nil will "hard" delete)
	TargetDeletedAtColumn jet.ColumnTimestamp
	// Target job column (required when TargetSourceIDColumn not set)
	TargetJobColumn jet.ColumnString

	TargetSourceIDValue jet.Expression
}

type JobTableSource struct {
	SourceTable jet.Table
	// To find rows that are part of the job to delete
	SourceJobColumn jet.ColumnString
	// So the resource can be set as "(soft) deleted"
	SourceDeletedAtColumn jet.ColumnTimestamp
	// Neded to have ID(s) for results
	SourceIDColumn jet.ColumnInteger
}

func AddJobTable(tbls ...*JobTable) {
	fkTableListsMu.Lock()
	defer fkTableListsMu.Unlock()

	for _, tbl := range tbls {
		fkTablesList[tbl.TargetTable.TableName()] = tbl
	}
}
