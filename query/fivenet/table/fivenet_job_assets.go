package table

import "github.com/go-jet/jet/v2/mysql"

var FivenetJobAssets = newFivenetJobAssetsTable("", "fivenet_job_assets", "")

type fivenetJobAssetsTable struct {
	mysql.Table

	Job             mysql.ColumnString
	FileID          mysql.ColumnInteger
	CreatedByUserID mysql.ColumnInteger
	DisplayName     mysql.ColumnString
	CreatedAt       mysql.ColumnTimestamp
	AllColumns      mysql.ColumnList
	MutableColumns  mysql.ColumnList
	DefaultColumns  mysql.ColumnList
}

type FivenetJobAssetsTable struct {
	fivenetJobAssetsTable

	NEW fivenetJobAssetsTable
}

func (a FivenetJobAssetsTable) AS(alias string) *FivenetJobAssetsTable {
	return newFivenetJobAssetsTable(a.SchemaName(), a.TableName(), alias)
}

func (a FivenetJobAssetsTable) FromSchema(schemaName string) *FivenetJobAssetsTable {
	return newFivenetJobAssetsTable(schemaName, a.TableName(), a.Alias())
}

func (a FivenetJobAssetsTable) WithPrefix(prefix string) *FivenetJobAssetsTable {
	return newFivenetJobAssetsTable(a.SchemaName(), prefix+a.TableName(), a.TableName())
}

func (a FivenetJobAssetsTable) WithSuffix(suffix string) *FivenetJobAssetsTable {
	return newFivenetJobAssetsTable(a.SchemaName(), a.TableName()+suffix, a.TableName())
}

func newFivenetJobAssetsTable(schema, name, alias string) *FivenetJobAssetsTable {
	return &FivenetJobAssetsTable{
		fivenetJobAssetsTable: newFivenetJobAssetsTableImpl(schema, name, alias),
		NEW:                   newFivenetJobAssetsTableImpl("", "new", ""),
	}
}

func newFivenetJobAssetsTableImpl(schema, name, alias string) fivenetJobAssetsTable {
	job := mysql.StringColumn("job")
	fileID := mysql.IntegerColumn("file_id")
	createdByUserID := mysql.IntegerColumn("created_by_user_id")
	displayName := mysql.StringColumn("display_name")
	createdAt := mysql.TimestampColumn("created_at")
	all := mysql.ColumnList{job, fileID, createdByUserID, displayName, createdAt}
	return fivenetJobAssetsTable{
		Table: mysql.NewTable(
			schema,
			name,
			alias,
			all...),
		Job:             job,
		FileID:          fileID,
		CreatedByUserID: createdByUserID,
		DisplayName:     displayName,
		CreatedAt:       createdAt,
		AllColumns:      all,
		MutableColumns: mysql.ColumnList{
			createdByUserID,
			displayName,
			createdAt,
		},
		DefaultColumns: mysql.ColumnList{createdAt},
	}
}
