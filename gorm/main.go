package main

import (
	"os"

	"github.com/galexrt/arpanet/model"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

// Dynamic SQL
type Querier interface {
	// SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:           "./query",
		WithUnitTest:      true,
		FieldNullable:     true,
		FieldWithTypeTag:  true,
		FieldSignable:     true,
		FieldWithIndexTag: true,
		// generate mode
		Mode: gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})
	g.Revise()

	gormdb, err := gorm.Open(mysql.Open(os.Getenv("MYSQL_DSN")))
	if err != nil {
		panic(err)
	}
	g.UseDB(gormdb)

	// Generate Models of Tables
	jobsModel := g.GenerateModel("jobs")
	jobGradesModel := g.GenerateModel("job_grades",
		// Ignore certain fields
		gen.FieldIgnore("name", "salary", "skin_male", "skim_female"),

		// "Normalize" some data types
		gen.FieldType("grade", "int"),
	)

	vpcLSModel := g.GenerateModel("vpcLS",
		gen.FieldRename("NET", "net"),
		gen.FieldJSONTag("coordsx", "coords_x"),
		gen.FieldJSONTag("coordsy", "coords_y"),
	)

	userLicenses := g.GenerateModel("user_licenses",
		gen.FieldJSONTag("owner", "-"))

	usersModel := g.GenerateModel("users",
		// Ignore certain fields
		gen.FieldIgnore("license", "group", "skin", "loadout", "position", "is_dead", "last_property", "inventory", "tattoos", "disabled", "levelData", "onDuty", "health", "armor"),

		// Fixup some field types and column names
		gen.FieldType("sex", "Sex"),
		gen.FieldType("job_grade", "int"),

		gen.FieldRename("last_seen", "updated_at"),
		gen.FieldJSONTag("last_seen", "updated_at"),

		gen.FieldType("accounts", "MoneyAccounts"),
		gen.FieldGORMTag("accounts", "serializer:json"),
		gen.FieldJSONTag("accounts", "-"),

		// Add relations for lazy loading
		gen.FieldRelateModel(field.HasMany, "Documents", model.Document{},
			&field.RelateConfig{
				GORMTag: "foreignkey:Creator",
			}),
		gen.FieldRelate(field.HasMany, "UserLicenses", userLicenses,
			&field.RelateConfig{
				GORMTag: "foreignkey:Owner",
			}),
	)

	// Generate default DAO interface for those generated structs from database
	g.ApplyBasic(
		usersModel,
		jobsModel,
		jobGradesModel,
		vpcLSModel,
		model.Document{},
		model.DocumentJobAccess{},
		model.DocumentUserAccess{},
	)

	// Generate the code
	g.Execute()
}
