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

var FivenetRoleAttrs = newFivenetRoleAttrsTable("", "fivenet_role_attrs", "")

type fivenetRoleAttrsTable struct {
	mysql.Table

	// Columns
	RoleID    mysql.ColumnInteger
	CreatedAt mysql.ColumnTimestamp
	UpdatedAt mysql.ColumnTimestamp
	AttrID    mysql.ColumnInteger
	Value     mysql.ColumnString

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type FivenetRoleAttrsTable struct {
	fivenetRoleAttrsTable

	NEW fivenetRoleAttrsTable
}

// AS creates new FivenetRoleAttrsTable with assigned alias
func (a FivenetRoleAttrsTable) AS(alias string) *FivenetRoleAttrsTable {
	return newFivenetRoleAttrsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetRoleAttrsTable with assigned schema name
func (a FivenetRoleAttrsTable) FromSchema(schemaName string) *FivenetRoleAttrsTable {
	return newFivenetRoleAttrsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetRoleAttrsTable with assigned table prefix
func (a FivenetRoleAttrsTable) WithPrefix(prefix string) *FivenetRoleAttrsTable {
	return newFivenetRoleAttrsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetRoleAttrsTable with assigned table suffix
func (a FivenetRoleAttrsTable) WithSuffix(suffix string) *FivenetRoleAttrsTable {
	return newFivenetRoleAttrsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetRoleAttrsTable(schemaName, tableName, alias string) *FivenetRoleAttrsTable {
	return &FivenetRoleAttrsTable{
		fivenetRoleAttrsTable: newFivenetRoleAttrsTableImpl(schemaName, tableName, alias),
		NEW:                   newFivenetRoleAttrsTableImpl("", "new", ""),
	}
}

func newFivenetRoleAttrsTableImpl(schemaName, tableName, alias string) fivenetRoleAttrsTable {
	var (
		RoleIDColumn    = mysql.IntegerColumn("role_id")
		CreatedAtColumn = mysql.TimestampColumn("created_at")
		UpdatedAtColumn = mysql.TimestampColumn("updated_at")
		AttrIDColumn    = mysql.IntegerColumn("attr_id")
		ValueColumn     = mysql.StringColumn("value")
		allColumns      = mysql.ColumnList{RoleIDColumn, CreatedAtColumn, UpdatedAtColumn, AttrIDColumn, ValueColumn}
		mutableColumns  = mysql.ColumnList{CreatedAtColumn, UpdatedAtColumn, ValueColumn}
		defaultColumns  = mysql.ColumnList{CreatedAtColumn}
	)

	return fivenetRoleAttrsTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		RoleID:    RoleIDColumn,
		CreatedAt: CreatedAtColumn,
		UpdatedAt: UpdatedAtColumn,
		AttrID:    AttrIDColumn,
		Value:     ValueColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
