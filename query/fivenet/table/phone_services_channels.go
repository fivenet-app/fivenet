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

var PhoneServicesChannels = newPhoneServicesChannelsTable("", "phone_services_channels", "")

type phoneServicesChannelsTable struct {
	mysql.Table

	// Columns
	ID          mysql.ColumnInteger
	PhoneNumber mysql.ColumnString
	Company     mysql.ColumnString
	Timestamp   mysql.ColumnTimestamp

	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type PhoneServicesChannelsTable struct {
	phoneServicesChannelsTable

	NEW phoneServicesChannelsTable
}

// AS creates new PhoneServicesChannelsTable with assigned alias
func (a PhoneServicesChannelsTable) AS(alias string) *PhoneServicesChannelsTable {
	return newPhoneServicesChannelsTable(a.SchemaName(), a.TableName(), alias)
}

// Schema creates new PhoneServicesChannelsTable with assigned schema name
func (a PhoneServicesChannelsTable) FromSchema(schemaName string) *PhoneServicesChannelsTable {
	return newPhoneServicesChannelsTable(schemaName, a.TableName(), a.Alias())
}

// WithPrefix creates new PhoneServicesChannelsTable with assigned table prefix
func (a PhoneServicesChannelsTable) WithPrefix(prefix string) *PhoneServicesChannelsTable {
	return newPhoneServicesChannelsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

// WithSuffix creates new PhoneServicesChannelsTable with assigned table suffix
func (a PhoneServicesChannelsTable) WithSuffix(suffix string) *PhoneServicesChannelsTable {
	return newPhoneServicesChannelsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newPhoneServicesChannelsTable(schemaName, tableName, alias string) *PhoneServicesChannelsTable {
	return &PhoneServicesChannelsTable{
		phoneServicesChannelsTable: newPhoneServicesChannelsTableImpl(schemaName, tableName, alias),
		NEW:                        newPhoneServicesChannelsTableImpl("", "new", ""),
	}
}

func newPhoneServicesChannelsTableImpl(schemaName, tableName, alias string) phoneServicesChannelsTable {
	var (
		IDColumn          = mysql.IntegerColumn("id")
		PhoneNumberColumn = mysql.StringColumn("phone_number")
		CompanyColumn     = mysql.StringColumn("company")
		TimestampColumn   = mysql.TimestampColumn("timestamp")
		allColumns        = mysql.ColumnList{IDColumn, PhoneNumberColumn, CompanyColumn, TimestampColumn}
		mutableColumns    = mysql.ColumnList{PhoneNumberColumn, CompanyColumn, TimestampColumn}
		defaultColumns    = mysql.ColumnList{TimestampColumn}
	)

	return phoneServicesChannelsTable{
		Table: mysql.NewTable(schemaName, tableName, alias, allColumns...),

		//Columns
		ID:          IDColumn,
		PhoneNumber: PhoneNumberColumn,
		Company:     CompanyColumn,
		Timestamp:   TimestampColumn,

		AllColumns:     allColumns,
		MutableColumns: mutableColumns,
		DefaultColumns: defaultColumns,
	}
}
