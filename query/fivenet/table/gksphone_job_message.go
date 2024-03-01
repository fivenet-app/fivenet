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

var GksphoneJobMessage = newGksphoneJobMessageTable("", "gksphone_job_message", "")

type gksphoneJobMessageTable struct {
	mysql.Table

	// Columns
	ID      mysql.ColumnInteger
	Name    mysql.ColumnString
	Number  mysql.ColumnString
	Message mysql.ColumnString
	Photo   mysql.ColumnString
	Gps     mysql.ColumnString
	Owner   mysql.ColumnInteger
	Jobm    mysql.ColumnString
	Anon    mysql.ColumnString
	Time    mysql.ColumnTimestamp

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
}

type GksphoneJobMessageTable struct {
	gksphoneJobMessageTable

	NEW gksphoneJobMessageTable
}

// AS creates new GksphoneJobMessageTable with assigned alias
func (a GksphoneJobMessageTable) AS(alias string) *GksphoneJobMessageTable {
	return newGksphoneJobMessageTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new GksphoneJobMessageTable with assigned schema name
func (a GksphoneJobMessageTable) FromSchema(schemaName string) *GksphoneJobMessageTable {
	return newGksphoneJobMessageTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new GksphoneJobMessageTable with assigned table prefix
func (a GksphoneJobMessageTable) WithPrefix(prefix string) *GksphoneJobMessageTable {
	return newGksphoneJobMessageTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new GksphoneJobMessageTable with assigned table suffix
func (a GksphoneJobMessageTable) WithSuffix(suffix string) *GksphoneJobMessageTable {
	return newGksphoneJobMessageTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newGksphoneJobMessageTable(schemaName, tableName, alias string) *GksphoneJobMessageTable {
	return &GksphoneJobMessageTable{
		gksphoneJobMessageTable: newGksphoneJobMessageTableImpl(schemaName, tableName, alias),
		NEW:                     newGksphoneJobMessageTableImpl("", "new", ""),
	}
}

func newGksphoneJobMessageTableImpl(schemaName, tableName, alias string) gksphoneJobMessageTable {
	var (
		IDColumn       = mysql.IntegerColumn("id")
		NameColumn     = mysql.StringColumn("name")
		NumberColumn   = mysql.StringColumn("number")
		MessageColumn  = mysql.StringColumn("message")
		PhotoColumn    = mysql.StringColumn("photo")
		GpsColumn      = mysql.StringColumn("gps")
		OwnerColumn    = mysql.IntegerColumn("owner")
		JobmColumn     = mysql.StringColumn("jobm")
		AnonColumn     = mysql.StringColumn("anon")
		TimeColumn     = mysql.TimestampColumn("time")
		allColumns     = mysql.ColumnList{IDColumn, NameColumn, NumberColumn, MessageColumn, PhotoColumn, GpsColumn, OwnerColumn, JobmColumn, AnonColumn, TimeColumn}
		mutableColumns = mysql.ColumnList{NameColumn, NumberColumn, MessageColumn, PhotoColumn, GpsColumn, OwnerColumn, JobmColumn, AnonColumn, TimeColumn}
	)

	return gksphoneJobMessageTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:      IDColumn,
		Name:    NameColumn,
		Number:  NumberColumn,
		Message: MessageColumn,
		Photo:   PhotoColumn,
		Gps:     GpsColumn,
		Owner:   OwnerColumn,
		Jobm:    JobmColumn,
		Anon:    AnonColumn,
		Time:    TimeColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}