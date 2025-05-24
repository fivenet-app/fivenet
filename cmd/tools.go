package cmd

type ToolsCmd struct {
	DB DBCmd `cmd:""`

	SyncCmd SyncCmd `cmd:"" name:"sync"`
}
