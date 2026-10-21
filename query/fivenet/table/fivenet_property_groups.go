package table

import "github.com/go-jet/jet/v2/mysql"

var FivenetPropertyGroups = newFivenetPropertyGroupsTable("", "fivenet_property_groups", "")

type fivenetPropertyGroupsTable struct {
	mysql.Table

	Scope          mysql.ColumnString
	GroupKey       mysql.ColumnString
	Label          mysql.ColumnString
	SortOrder      mysql.ColumnInteger
	DeletedAt      mysql.ColumnTimestamp
	CreatedAt      mysql.ColumnTimestamp
	UpdatedAt      mysql.ColumnTimestamp
	AllColumns     mysql.ColumnList
	MutableColumns mysql.ColumnList
	DefaultColumns mysql.ColumnList
}

type FivenetPropertyGroupsTable struct {
	fivenetPropertyGroupsTable

	NEW fivenetPropertyGroupsTable
}

func (a FivenetPropertyGroupsTable) AS(alias string) *FivenetPropertyGroupsTable {
	return newFivenetPropertyGroupsTable(a.SchemaName(), a.TableName(), alias)
}

func (a FivenetPropertyGroupsTable) FromSchema(schemaName string) *FivenetPropertyGroupsTable {
	return newFivenetPropertyGroupsTable(schemaName, a.TableName(), a.Alias())
}

func (a FivenetPropertyGroupsTable) WithPrefix(prefix string) *FivenetPropertyGroupsTable {
	return newFivenetPropertyGroupsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

func (a FivenetPropertyGroupsTable) WithSuffix(suffix string) *FivenetPropertyGroupsTable {
	return newFivenetPropertyGroupsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetPropertyGroupsTable(
	schemaName, tableName, alias string,
) *FivenetPropertyGroupsTable {
	return &FivenetPropertyGroupsTable{
		fivenetPropertyGroupsTable: newFivenetPropertyGroupsTableImpl(schemaName, tableName, alias),
		NEW:                        newFivenetPropertyGroupsTableImpl("", "new", ""),
	}
}

func newFivenetPropertyGroupsTableImpl(
	schemaName, tableName, alias string,
) fivenetPropertyGroupsTable {
	scope := mysql.StringColumn("scope")
	groupKey := mysql.StringColumn("group_key")
	label := mysql.StringColumn("label")
	sortOrder := mysql.IntegerColumn("sort_order")
	deletedAt := mysql.TimestampColumn("deleted_at")
	createdAt := mysql.TimestampColumn("created_at")
	updatedAt := mysql.TimestampColumn("updated_at")
	all := mysql.ColumnList{scope, groupKey, label, sortOrder, deletedAt, createdAt, updatedAt}
	return fivenetPropertyGroupsTable{
		Table:          mysql.NewTable(schemaName, tableName, alias, all...),
		Scope:          scope,
		GroupKey:       groupKey,
		Label:          label,
		SortOrder:      sortOrder,
		DeletedAt:      deletedAt,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
		AllColumns:     all,
		MutableColumns: mysql.ColumnList{label, sortOrder, deletedAt, updatedAt},
		DefaultColumns: mysql.ColumnList{label, sortOrder, createdAt, updatedAt},
	}
}
