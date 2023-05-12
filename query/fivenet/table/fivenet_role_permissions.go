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

var FivenetRolePermissions = newFivenetRolePermissionsTable("", "fivenet_role_permissions", "")

type fivenetRolePermissionsTable struct {
	mysql.Table

	// Columns
	RoleID       mysql.ColumnInteger
	PermissionID mysql.ColumnInteger
	Val          mysql.ColumnBool

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
}

type FivenetRolePermissionsTable struct {
	fivenetRolePermissionsTable

	NEW fivenetRolePermissionsTable
}

// AS creates new FivenetRolePermissionsTable with assigned alias
func (a FivenetRolePermissionsTable) AS(alias string) *FivenetRolePermissionsTable {
	return newFivenetRolePermissionsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new FivenetRolePermissionsTable with assigned schema name
func (a FivenetRolePermissionsTable) FromSchema(schemaName string) *FivenetRolePermissionsTable {
	return newFivenetRolePermissionsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new FivenetRolePermissionsTable with assigned table prefix
func (a FivenetRolePermissionsTable) WithPrefix(prefix string) *FivenetRolePermissionsTable {
	return newFivenetRolePermissionsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new FivenetRolePermissionsTable with assigned table suffix
func (a FivenetRolePermissionsTable) WithSuffix(suffix string) *FivenetRolePermissionsTable {
	return newFivenetRolePermissionsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetRolePermissionsTable(schemaName, tableName, alias string) *FivenetRolePermissionsTable {
	return &FivenetRolePermissionsTable{
		fivenetRolePermissionsTable: newFivenetRolePermissionsTableImpl(schemaName, tableName, alias),
		NEW:                         newFivenetRolePermissionsTableImpl("", "new", ""),
	}
}

func newFivenetRolePermissionsTableImpl(schemaName, tableName, alias string) fivenetRolePermissionsTable {
	var (
		RoleIDColumn       = mysql.IntegerColumn("role_id")
		PermissionIDColumn = mysql.IntegerColumn("permission_id")
		ValColumn          = mysql.BoolColumn("val")
		allColumns         = mysql.ColumnList{RoleIDColumn, PermissionIDColumn, ValColumn}
		mutableColumns     = mysql.ColumnList{ValColumn}
	)

	return fivenetRolePermissionsTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		RoleID:       RoleIDColumn,
		PermissionID: PermissionIDColumn,
		Val:          ValColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
	}
}
