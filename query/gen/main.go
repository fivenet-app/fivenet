package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-jet/jet/v2/generator/metadata"
	genmysql "github.com/go-jet/jet/v2/generator/mysql"
	"github.com/go-jet/jet/v2/generator/template"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
)

var rootCmd = &cobra.Command{
	Use: "gen",
	RunE: func(cmd *cobra.Command, args []string) error {
		destDir := "./query"

		return genmysql.GenerateDSN(viper.GetString("dsn"), destDir, genTemplate())
	},
}

func main() {
	viper.SetConfigFile("./query/gen/gen.yaml")
	viper.ReadInConfig()

	// Run Command
	if err := rootCmd.Execute(); err != nil {
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
								defaultTableModelField := template.DefaultTableModelField(columnMetaData)
								return defaultTableModelField.UseTags(
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
						return template.DefaultTableSQLBuilder(table)
					}),
				)
		})
}

func shouldSkipTable(table string) bool {
	table = strings.ToLower(table)

	excludeTables := viper.GetStringSlice("excludeTables")
	includeTables := viper.GetStringSlice("includeTables")

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

	includeIdx := slices.IndexFunc(includeTables, func(c string) bool {
		c = strings.ToLower(c)
		if strings.HasSuffix(c, "*") {
			return strings.HasPrefix(table, c[:len(c)-1])
		}
		return table == c
	})

	return includeIdx == -1
}
