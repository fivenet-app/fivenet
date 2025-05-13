package housekeeper

import (
	"sync"

	jet "github.com/go-jet/jet/v2/mysql"
)

const DefaultDeleteLimit = 500

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

func (t *JobTable) GetQuery(job string) jet.Statement {
	return t.GetQueryWithLimit(job, DefaultDeleteLimit)
}

func (t *JobTable) GetQueryWithLimit(job string, limit int64) jet.Statement {
	if t.Source == nil {
		if t.TargetDeletedAtColumn == nil {
			return t.TargetTable.
				DELETE().
				WHERE(jet.AND(
					t.TargetJobColumn.EQ(jet.String(job)),
				)).
				LIMIT(limit)
		} else {
			return t.TargetTable.
				UPDATE().
				SET(t.TargetDeletedAtColumn.SET(jet.CURRENT_TIMESTAMP())).
				WHERE(jet.AND(
					t.TargetJobColumn.EQ(jet.String(job)),
					t.TargetDeletedAtColumn.IS_NULL(),
				)).
				LIMIT(limit)
		}
	} else {
		// Handle case with Source using USING AND FROM
		joinCondition := t.Source.SourceIDColumn.EQ(t.TargetSourceIDColumn)
		whereCondition := t.Source.SourceJobColumn.EQ(jet.String(job))

		if t.Source.SourceDeletedAtColumn != nil {
			whereCondition = whereCondition.AND(t.Source.SourceDeletedAtColumn.IS_NOT_NULL())
		}
		if t.TargetDeletedAtColumn != nil {
			whereCondition = whereCondition.AND(t.TargetDeletedAtColumn.IS_NULL())
		}

		if t.TargetDeletedAtColumn == nil {
			return t.TargetTable.
				DELETE().
				USING(t.Source.SourceTable).
				WHERE(jet.AND(
					whereCondition,
					joinCondition,
				)).
				LIMIT(limit)
		} else {
			return t.TargetTable.
				INNER_JOIN(t.Source.SourceTable, joinCondition).
				UPDATE().
				SET(
					t.TargetDeletedAtColumn.SET(jet.CURRENT_TIMESTAMP()),
				).
				WHERE(jet.AND(
					whereCondition,
					joinCondition,
				)).
				LIMIT(limit)
		}
	}
}
