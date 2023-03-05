package main

import (
	"os"

	"github.com/galexrt/arpanet/model"
	"github.com/galexrt/arpanet/pkg/permify/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		panic("Error loading .env file")
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:           "./query",
		WithUnitTest:      true,
		FieldNullable:     false,
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
		gen.FieldType("coordsx", "float32"),
		gen.FieldJSONTag("coordsy", "coords_y"),
		gen.FieldType("coordsy", "float32"),
	)

	userLicenses := g.GenerateModel("user_licenses",
		gen.FieldType("type", "LicenseType"),
		gen.FieldJSONTag("owner", "-"),
	)

	usersModel := g.GenerateModel("users",
		// Ignore certain fields
		gen.FieldIgnore("license", "group", "skin", "loadout", "position", "is_dead", "last_property", "inventory", "tattoos", "levelData", "onDuty", "health", "armor"),

		// Fixup some field types and column names
		gen.FieldType("sex", "Sex"),
		gen.FieldType("job_grade", "int"),

		gen.FieldRename("last_seen", "updated_at"),
		gen.FieldJSONTag("last_seen", "updated_at"),

		gen.FieldType("accounts", "MoneyAccounts"),
		gen.FieldGORMTag("accounts", "serializer:json"),
		gen.FieldJSONTag("accounts", "-"),

		// Add relations for lazy loading
		// gen.FieldRelateModel(field.HasOne, "UserProps", model.UserProps{},
		// 	&field.RelateConfig{
		// 		GORMTag:       "foreignKey:ID;references:UserID",
		// 		RelatePointer: true,
		// 	}),
		gen.FieldRelate(field.HasMany, "UserLicenses", userLicenses,
			&field.RelateConfig{
				GORMTag: "foreignKey:Owner;references:Identifier",
			}),
		gen.FieldRelateModel(field.HasMany, "Documents", model.Document{},
			&field.RelateConfig{
				GORMTag: "foreignKey:CreatorID;references:ID",
			}),

		// Activity
		/*
			gen.FieldRelateModel(field.HasMany, "TargetActivity", model.UserActivity{},
				&field.RelateConfig{
					GORMTag: "foreignKey:TargetUserID;references:ID",
				}),
			gen.FieldRelateModel(field.HasMany, "CauseActivity", model.UserActivity{},
				&field.RelateConfig{
					GORMTag: "foreignKey:CauseUserID;references:ID",
				}),
		*/

		// User Roles + Permissions for Permify
		gen.FieldRelateModel(field.Many2Many, "Roles", models.Role{},
			&field.RelateConfig{
				GORMTag: "many2many:arpanet_user_roles;constraint:OnUpdate:CASCADE,OnDelete:CASCADE",
			}),
		gen.FieldRelateModel(field.Many2Many, "Permissions", models.Permission{},
			&field.RelateConfig{
				GORMTag: "many2many:arpanet_user_permissions;constraint:OnUpdate:CASCADE,OnDelete:CASCADE",
			}),
	)

	// Generate default DAO interface for those generated structs from database
	g.ApplyBasic(
		// User related
		model.Account{},
		usersModel,
		userLicenses,
		model.UserActivity{},
		model.UserProps{},
		jobsModel,
		jobGradesModel,
		// User location
		model.UserLocation{},
		vpcLSModel,
		// Document related
		model.Document{},
		model.DocumentJobAccess{},
		model.DocumentMentions{},
		model.DocumentUserAccess{},
	)

	// Generate the code
	g.Execute()
}
