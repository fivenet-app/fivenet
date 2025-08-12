//nolint:forbidigo // This is a CLI tool that uses `fmt.Println` for output
package main

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/go-jet/jet/v2/generator/metadata"
	genmysql "github.com/go-jet/jet/v2/generator/mysql"
	"github.com/go-jet/jet/v2/generator/template"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigFile("./query/gen/gen.yaml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	destDir := "./query"

	if err := genmysql.GenerateDSN(viper.GetString("dsn"), destDir, genTemplate()); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// https://github.com/go-jet/jet/blob/f2e4b8551c48b97d0cd2d3deff47dc1b2aa2f04e/tests/postgres/generator_template_test.go#L274-L314
func genTemplate() template.Template {
	return template.Default(mysql.Dialect).
		UseSchema(func(schemaMetaData metadata.Schema) template.Schema {
			return template.DefaultSchema(schemaMetaData).
				UseModel(template.DefaultModel().
					UseTable(func(table metadata.Table) template.TableModel {
						if shouldSkipTable(table.Name) {
							return template.TableModel{Skip: true}
						}

						return template.DefaultTableModel(table).
							UseField(func(columnMetaData metadata.Column) template.TableModelField {
								defaultColumn := template.DefaultTableModelField(columnMetaData)
								if shouldSkipField(table.Name, columnMetaData.Name) {
									defaultColumn.Skip = true
								}

								return defaultColumn.UseTags(
									fmt.Sprintf(`json:"%s"`, columnMetaData.Name),
								)
							})
					}),
				).
				UseSQLBuilder(template.DefaultSQLBuilder().
					UseTable(func(table metadata.Table) template.TableSQLBuilder {
						if shouldSkipTable(table.Name) {
							return template.TableSQLBuilder{Skip: true}
						}

						return template.DefaultTableSQLBuilder(table).
							UseColumn(func(columnMetaData metadata.Column) template.TableSQLBuilderColumn {
								defaultColumn := template.DefaultTableSQLBuilderColumn(
									columnMetaData,
								)
								if shouldSkipField(table.Name, columnMetaData.Name) {
									defaultColumn.Skip = true
								}

								return defaultColumn
							})
					}),
				)
		})
}

func shouldSkipTable(table string) bool {
	table = strings.ToLower(table)

	excludeTables := viper.GetStringSlice("excludeTables")
	includeTables := viper.GetStringMapStringSlice("includeTables")

	excludeIdx := slices.IndexFunc(excludeTables, func(c string) bool {
		c = strings.ToLower(c)
		if strings.HasSuffix(c, "*") {
			return strings.HasPrefix(table, c[:len(c)-1])
		}
		return table == c
	})
	if excludeIdx > -1 {
		return true
	}

	includeKeys := []string{}
	for key := range includeTables {
		includeKeys = append(includeKeys, key)
	}

	includeIdx := slices.IndexFunc(includeKeys, func(c string) bool {
		c = strings.ToLower(c)
		if strings.HasSuffix(c, "*") {
			return strings.HasPrefix(table, c[:len(c)-1])
		}
		return table == c
	})

	return includeIdx == -1
}

func shouldSkipField(table string, field string) bool {
	includeTables := viper.GetStringMapStringSlice("includeTables")

	fields, ok := includeTables[table]
	if !ok {
		return false
	}

	if len(fields) == 0 || fields[0] == "*" {
		return false
	}

	return !slices.Contains(fields, field)
}
