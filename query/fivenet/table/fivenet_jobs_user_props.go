//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package table

import (
	"github.com/go-jet/jet/v2/mysql"
)

var FivenetJobsUserProps = newFivenetJobsUserPropsTable("", "fivenet_jobs_user_props", "")

type fivenetJobsUserPropsTable struct {
	mysql.Table

	// Columns
	UserID       mysql.ColumnInteger
	Job          mysql.ColumnString
	AbsenceBegin mysql.ColumnDate
	AbsenceEnd   mysql.ColumnDate
	Note         mysql.ColumnString
	NamePrefix   mysql.ColumnString
	NameSuffix   mysql.ColumnString

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type FivenetJobsUserPropsTable struct {
	fivenetJobsUserPropsTable

	NEW fivenetJobsUserPropsTable
}

// AS creates new FivenetJobsUserPropsTable with assigned alias
func (a FivenetJobsUserPropsTable) AS(alias string) *FivenetJobsUserPropsTable {
	return newFivenetJobsUserPropsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetJobsUserPropsTable with assigned schema name
func (a FivenetJobsUserPropsTable) FromSchema(schemaName string) *FivenetJobsUserPropsTable {
	return newFivenetJobsUserPropsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetJobsUserPropsTable with assigned table prefix
func (a FivenetJobsUserPropsTable) WithPrefix(prefix string) *FivenetJobsUserPropsTable {
	return newFivenetJobsUserPropsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetJobsUserPropsTable with assigned table suffix
func (a FivenetJobsUserPropsTable) WithSuffix(suffix string) *FivenetJobsUserPropsTable {
	return newFivenetJobsUserPropsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetJobsUserPropsTable(schemaName, tableName, alias string) *FivenetJobsUserPropsTable {
	return &FivenetJobsUserPropsTable{
		fivenetJobsUserPropsTable: newFivenetJobsUserPropsTableImpl(schemaName, tableName, alias),
		NEW:                       newFivenetJobsUserPropsTableImpl("", "new", ""),
	}
}

func newFivenetJobsUserPropsTableImpl(schemaName, tableName, alias string) fivenetJobsUserPropsTable {
	var (
		UserIDColumn       = mysql.IntegerColumn("user_id")
		JobColumn          = mysql.StringColumn("job")
		AbsenceBeginColumn = mysql.DateColumn("absence_begin")
		AbsenceEndColumn   = mysql.DateColumn("absence_end")
		NoteColumn         = mysql.StringColumn("note")
		NamePrefixColumn   = mysql.StringColumn("name_prefix")
		NameSuffixColumn   = mysql.StringColumn("name_suffix")
		allColumns         = mysql.ColumnList{UserIDColumn, JobColumn, AbsenceBeginColumn, AbsenceEndColumn, NoteColumn, NamePrefixColumn, NameSuffixColumn}
		mutableColumns     = mysql.ColumnList{UserIDColumn, JobColumn, AbsenceBeginColumn, AbsenceEndColumn, NoteColumn, NamePrefixColumn, NameSuffixColumn}
		defaultColumns     = mysql.ColumnList{}
	)

	return fivenetJobsUserPropsTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		UserID:       UserIDColumn,
		Job:          JobColumn,
		AbsenceBegin: AbsenceBeginColumn,
		AbsenceEnd:   AbsenceEndColumn,
		Note:         NoteColumn,
		NamePrefix:   NamePrefixColumn,
		NameSuffix:   NameSuffixColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
