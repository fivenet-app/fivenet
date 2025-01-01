package config

type DBSync struct {
	Enabled bool         `default:"false" yaml:"enabled"`
	Source  DBSyncSource `yaml:"source"`
}

type DBSyncSource struct {
	// Refer to https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	DSN string `yaml:"dsn"`

	Destination DBSyncDestination `yaml:"destination"`

	Tables DBSyncSourceTables `yaml:"tables"`
}

type DBSyncDestination struct {
	URL   string `yaml:"url"`
	Token string `yaml:"token"`
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
	Queries []string `yaml:"queries"`
}
