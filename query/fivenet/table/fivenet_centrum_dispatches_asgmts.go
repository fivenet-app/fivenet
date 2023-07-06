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

var FivenetCentrumDispatchesAsgmts = newFivenetCentrumDispatchesAsgmtsTable("", "fivenet_centrum_dispatches_asgmts", "")

type fivenetCentrumDispatchesAsgmtsTable struct {
	mysql.Table

	// Columns
	DispatchID mysql.ColumnInteger
	UnitID     mysql.ColumnInteger

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
}

type FivenetCentrumDispatchesAsgmtsTable struct {
	fivenetCentrumDispatchesAsgmtsTable

	NEW fivenetCentrumDispatchesAsgmtsTable
}

// AS creates new FivenetCentrumDispatchesAsgmtsTable with assigned alias
func (a FivenetCentrumDispatchesAsgmtsTable) AS(alias string) *FivenetCentrumDispatchesAsgmtsTable {
	return newFivenetCentrumDispatchesAsgmtsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetCentrumDispatchesAsgmtsTable with assigned schema name
func (a FivenetCentrumDispatchesAsgmtsTable) FromSchema(schemaName string) *FivenetCentrumDispatchesAsgmtsTable {
	return newFivenetCentrumDispatchesAsgmtsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetCentrumDispatchesAsgmtsTable with assigned table prefix
func (a FivenetCentrumDispatchesAsgmtsTable) WithPrefix(prefix string) *FivenetCentrumDispatchesAsgmtsTable {
	return newFivenetCentrumDispatchesAsgmtsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetCentrumDispatchesAsgmtsTable with assigned table suffix
func (a FivenetCentrumDispatchesAsgmtsTable) WithSuffix(suffix string) *FivenetCentrumDispatchesAsgmtsTable {
	return newFivenetCentrumDispatchesAsgmtsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetCentrumDispatchesAsgmtsTable(schemaName, tableName, alias string) *FivenetCentrumDispatchesAsgmtsTable {
	return &FivenetCentrumDispatchesAsgmtsTable{
		fivenetCentrumDispatchesAsgmtsTable: newFivenetCentrumDispatchesAsgmtsTableImpl(schemaName, tableName, alias),
		NEW:                                 newFivenetCentrumDispatchesAsgmtsTableImpl("", "new", ""),
	}
}

func newFivenetCentrumDispatchesAsgmtsTableImpl(schemaName, tableName, alias string) fivenetCentrumDispatchesAsgmtsTable {
	var (
		DispatchIDColumn = mysql.IntegerColumn("dispatch_id")
		UnitIDColumn     = mysql.IntegerColumn("unit_id")
		allColumns       = mysql.ColumnList{DispatchIDColumn, UnitIDColumn}
		mutableColumns   = mysql.ColumnList{}
	)

	return fivenetCentrumDispatchesAsgmtsTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		DispatchID: DispatchIDColumn,
		UnitID:     UnitIDColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
