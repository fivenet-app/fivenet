package main

import (
	"os"

	"github.com/galexrt/arpanet/model"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
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
		gen.FieldType("grade", "int"),
		gen.FieldIgnore("name", "salary", "skin_male", "skim_female"))

	usersModel := g.GenerateModelAs("users", "Citizen",
		gen.FieldIgnore("license", "group", "skin", "loadout", "position", "last_property", "inventory", "tattoos", "levelData", "onDuty", "health", "armor"),
		gen.FieldType("sex", "Sex"),
		gen.FieldType("job_grade", "int"),
		gen.FieldType("accounts", "Accounts"), gen.FieldGORMTag("accounts", "serializer:json"))

	// Generate default DAO interface for those generated structs from database
	g.ApplyBasic(
		usersModel,
		jobsModel,
		jobGradesModel,
		model.Document{},
		model.DocumentAccess{},
	)

	// Generate the code
	g.Execute()
}
