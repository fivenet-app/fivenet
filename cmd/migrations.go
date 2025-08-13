package cmd

type MigrationsCmd struct {
	Filestore FilestoreCmd `cmd:"" help:"Migrate files from the old database format to the new filestore format." name:"filestore"`
}
