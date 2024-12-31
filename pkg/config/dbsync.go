package config

type DBSync struct {
	Enabled bool         `default:"false" yaml:"enabled"`
	Source  DBSyncSource `yaml:"source"`
}

type DBSyncSource struct {
	// Refer to https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	DSN string `yaml:"dsn"`

	Tables DBSyncSourceTables `yaml:"tables"`
}

type DBSyncSourceTables struct {
	Users         DBSyncTable `yaml:"users"`
	Jobs          DBSyncTable `yaml:"jobs"`
	JobGrades     DBSyncTable `yaml:"jobGrades"`
	Licenses      DBSyncTable `yaml:"licenses"`
	UserLicenses  DBSyncTable `yaml:"userLicenses"`
	OwnedVehicles DBSyncTable `yaml:"ownedVehicles"`
}

type DBSyncTable struct {
	Enabled bool     `yaml:"enabled"`
	Query   string   `yaml:"query"`
	Queries []string `yaml:"queries"`
}
