package model

func Register() {
	DB.AutoMigrate(
		&Citizen{},
		&Document{},
		&CitizenAccess{},
		&JobAccess{},
		&Job{},
		&JobGrade{},
	)
}
