//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package tables

import (
	"github.com/go-jet/jet/v2/mysql"
)

var FivenetOwnedVehicles = newFivenetOwnedVehiclesTable("", "fivenet_owned_vehicles", "")

type fivenetOwnedVehiclesTable struct {
	mysql.Table

	// Columns
	Owner mysql.ColumnString
    Job   mysql.ColumnString
	Plate mysql.ColumnString
	Model mysql.ColumnString
	Type  mysql.ColumnString
    Data  mysql.ColumnString

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
}

type FivenetOwnedVehiclesTable struct {
	fivenetOwnedVehiclesTable

	NEW fivenetOwnedVehiclesTable
}

// AS creates new FivenetOwnedVehiclesTable with assigned alias
func (a FivenetOwnedVehiclesTable) AS(alias string) *FivenetOwnedVehiclesTable {
	return newFivenetOwnedVehiclesTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetOwnedVehiclesTable with assigned schema name
func (a FivenetOwnedVehiclesTable) FromSchema(schemaName string) *FivenetOwnedVehiclesTable {
	return newFivenetOwnedVehiclesTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetOwnedVehiclesTable with assigned table prefix
func (a FivenetOwnedVehiclesTable) WithPrefix(prefix string) *FivenetOwnedVehiclesTable {
	return newFivenetOwnedVehiclesTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetOwnedVehiclesTable with assigned table suffix
func (a FivenetOwnedVehiclesTable) WithSuffix(suffix string) *FivenetOwnedVehiclesTable {
	return newFivenetOwnedVehiclesTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetOwnedVehiclesTable(schemaName, tableName, alias string) *FivenetOwnedVehiclesTable {
	return &FivenetOwnedVehiclesTable{
		fivenetOwnedVehiclesTable: newFivenetOwnedVehiclesTableImpl(schemaName, tableName, alias),
		NEW:                       newFivenetOwnedVehiclesTableImpl("", "new", ""),
	}
}

func newFivenetOwnedVehiclesTableImpl(schemaName, tableName, alias string) fivenetOwnedVehiclesTable {
	var (
		OwnerColumn    = mysql.StringColumn("owner")
        JobColumn      = mysql.StringColumn("job")
		PlateColumn    = mysql.StringColumn("plate")
		ModelColumn    = mysql.StringColumn("model")
		TypeColumn     = mysql.StringColumn("type")
        DataColumn     = mysql.StringColumn("data")
		allColumns     = mysql.ColumnList{OwnerColumn, JobColumn, PlateColumn, ModelColumn, TypeColumn, DataColumn}
		mutableColumns = mysql.ColumnList{OwnerColumn, JobColumn, ModelColumn, TypeColumn, DataColumn}
	)

	return fivenetOwnedVehiclesTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		Owner: OwnerColumn,
        Job:   JobColumn,
		Plate: PlateColumn,
		Model: ModelColumn,
		Type:  TypeColumn,
        Data:  DataColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
