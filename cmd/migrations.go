package cmd

type MigrationsCmd struct {
	Filestore FilestoreCmd `cmd:"" name:"filestore"`
}
