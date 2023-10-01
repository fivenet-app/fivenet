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

var OwnedVehicles = newOwnedVehiclesTable("", "owned_vehicles", "")

type ownedVehiclesTable struct {
	mysql.Table

	// Columns
	Owner     mysql.ColumnString
	Plate     mysql.ColumnString
	Model     mysql.ColumnString
	Vehicle   mysql.ColumnString
	Type      mysql.ColumnString
	Stored    mysql.ColumnBool
	Carseller mysql.ColumnInteger
	Owners    mysql.ColumnString
	Trunk     mysql.ColumnString
	Glovebox  mysql.ColumnString

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
}

type OwnedVehiclesTable struct {
	ownedVehiclesTable

	NEW ownedVehiclesTable
}

// AS creates new OwnedVehiclesTable with assigned alias
func (a OwnedVehiclesTable) AS(alias string) *OwnedVehiclesTable {
	return newOwnedVehiclesTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new OwnedVehiclesTable with assigned schema name
func (a OwnedVehiclesTable) FromSchema(schemaName string) *OwnedVehiclesTable {
	return newOwnedVehiclesTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new OwnedVehiclesTable with assigned table prefix
func (a OwnedVehiclesTable) WithPrefix(prefix string) *OwnedVehiclesTable {
	return newOwnedVehiclesTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new OwnedVehiclesTable with assigned table suffix
func (a OwnedVehiclesTable) WithSuffix(suffix string) *OwnedVehiclesTable {
	return newOwnedVehiclesTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newOwnedVehiclesTable(schemaName, tableName, alias string) *OwnedVehiclesTable {
	return &OwnedVehiclesTable{
		ownedVehiclesTable: newOwnedVehiclesTableImpl(schemaName, tableName, alias),
		NEW:                newOwnedVehiclesTableImpl("", "new", ""),
	}
}

func newOwnedVehiclesTableImpl(schemaName, tableName, alias string) ownedVehiclesTable {
	var (
		OwnerColumn     = mysql.StringColumn("owner")
		PlateColumn     = mysql.StringColumn("plate")
		ModelColumn     = mysql.StringColumn("model")
		VehicleColumn   = mysql.StringColumn("vehicle")
		TypeColumn      = mysql.StringColumn("type")
		StoredColumn    = mysql.BoolColumn("stored")
		CarsellerColumn = mysql.IntegerColumn("carseller")
		OwnersColumn    = mysql.StringColumn("owners")
		TrunkColumn     = mysql.StringColumn("trunk")
		GloveboxColumn  = mysql.StringColumn("glovebox")
		allColumns      = mysql.ColumnList{OwnerColumn, PlateColumn, ModelColumn, VehicleColumn, TypeColumn, StoredColumn, CarsellerColumn, OwnersColumn, TrunkColumn, GloveboxColumn}
		mutableColumns  = mysql.ColumnList{OwnerColumn, ModelColumn, VehicleColumn, TypeColumn, StoredColumn, CarsellerColumn, OwnersColumn, TrunkColumn, GloveboxColumn}
	)

	return ownedVehiclesTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		Owner:     OwnerColumn,
		Plate:     PlateColumn,
		Model:     ModelColumn,
		Vehicle:   VehicleColumn,
		Type:      TypeColumn,
		Stored:    StoredColumn,
		Carseller: CarsellerColumn,
		Owners:    OwnersColumn,
		Trunk:     TrunkColumn,
		Glovebox:  GloveboxColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
