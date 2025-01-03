package config

type DBSync struct {
	Enabled bool `default:"false" yaml:"enabled"`

	StateFile string `default:"dbsync.state.yaml" yaml:"stateFile"`

	Destination DBSyncDestination `yaml:"destination"`
	Source      DBSyncSource      `yaml:"source"`

	Tables DBSyncSourceTables `yaml:"tables"`
}

type DBSyncSource struct {
	// Refer to https://github.com/go-sql-driver/mysql#dsn-data-source-name for details
	DSN string `yaml:"dsn"`
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
	Enabled bool   `yaml:"enabled"`
	IDField string `yaml:"idField"`
	Query   string `yaml:"query"`
}
