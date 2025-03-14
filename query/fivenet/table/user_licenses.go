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

var UserLicenses = newUserLicensesTable("", "user_licenses", "")

type userLicensesTable struct {
	mysql.Table

	// Columns
	Type  mysql.ColumnString
	Owner mysql.ColumnString

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type UserLicensesTable struct {
	userLicensesTable

	NEW userLicensesTable
}

// AS creates new UserLicensesTable with assigned alias
func (a UserLicensesTable) AS(alias string) *UserLicensesTable {
	return newUserLicensesTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new UserLicensesTable with assigned schema name
func (a UserLicensesTable) FromSchema(schemaName string) *UserLicensesTable {
	return newUserLicensesTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new UserLicensesTable with assigned table prefix
func (a UserLicensesTable) WithPrefix(prefix string) *UserLicensesTable {
	return newUserLicensesTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new UserLicensesTable with assigned table suffix
func (a UserLicensesTable) WithSuffix(suffix string) *UserLicensesTable {
	return newUserLicensesTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newUserLicensesTable(schemaName, tableName, alias string) *UserLicensesTable {
	return &UserLicensesTable{
		userLicensesTable: newUserLicensesTableImpl(schemaName, tableName, alias),
		NEW:               newUserLicensesTableImpl("", "new", ""),
	}
}

func newUserLicensesTableImpl(schemaName, tableName, alias string) userLicensesTable {
	var (
		TypeColumn     = mysql.StringColumn("type")
		OwnerColumn    = mysql.StringColumn("owner")
		allColumns     = mysql.ColumnList{TypeColumn, OwnerColumn}
		mutableColumns = mysql.ColumnList{}
		defaultColumns = mysql.ColumnList{}
	)

	return userLicensesTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		Type:  TypeColumn,
		Owner: OwnerColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
