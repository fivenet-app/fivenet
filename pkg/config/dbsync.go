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
	Jobs      DBSyncTable `yaml:"jobs"`
	JobGrades DBSyncTable `yaml:"jobGrades"`
	Licenses  DBSyncTable `yaml:"licenses"`

	Users        UsersDBSyncTable `yaml:"users"`
	UserLicenses DBSyncTable      `yaml:"userLicenses"`
	Vehicles     DBSyncTable      `yaml:"vehicles"`
}

type DBSyncTable struct {
	Enabled      bool    `yaml:"enabled"`
	IDField      string  `yaml:"idField"`
	Query        string  `yaml:"query"`
	UpdatedField *string `yaml:"updatedField"`
}

type UsersDBSyncTable struct {
	DBSyncTable `yaml:",inline" mapstructure:",squash"`

	SplitName bool `default:"false" yaml:"splitName"`
}
