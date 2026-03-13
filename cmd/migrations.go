package cmd

type MigrationsCmd struct {
	HTMLToJSON MigrationsHTMLToJSONCmd `cmd:"" help:"Migrate documents, comments, etc., from (raw) HTML format to the legacy custom JSON format." name:"htmltojson"`
	Filestore  MigrationsFilestoreCmd  `cmd:"" help:"Migrate files from the old database format to the new filestore format."                     name:"filestore"`

	StatsBackfill MigrationsStatsBackfillCmd `cmd:"" help:"Backfill stats for documents." name:"statsbackfill"`
}
