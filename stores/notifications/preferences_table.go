package notificationsstore

import "github.com/go-jet/jet/v2/mysql"

// preferenceTable is kept local until the migration has run against the query
// generator database and a generated binding can replace it.
type preferenceTable struct {
	mysql.Table

	UserID       mysql.ColumnInteger
	Category     mysql.ColumnInteger
	Kind         mysql.ColumnInteger
	InboxEnabled mysql.ColumnBool
	ToastEnabled mysql.ColumnBool
	SoundEnabled mysql.ColumnBool
}

func (a preferenceTable) AS(alias string) *preferenceTable {
	return newPreferenceTable(a.SchemaName(), a.TableName(), alias)
}

func newPreferenceTable(schemaName, tableName, alias string) *preferenceTable {
	userID := mysql.IntegerColumn("user_id")
	category := mysql.IntegerColumn("category")
	kind := mysql.IntegerColumn("kind")
	inboxEnabled := mysql.BoolColumn("inbox_enabled")
	toastEnabled := mysql.BoolColumn("toast_enabled")
	soundEnabled := mysql.BoolColumn("sound_enabled")

	return &preferenceTable{
		Table: mysql.NewTable(
			schemaName,
			tableName,
			alias,
			userID,
			category,
			kind,
			inboxEnabled,
			toastEnabled,
			soundEnabled,
		),
		UserID:       userID,
		Category:     category,
		Kind:         kind,
		InboxEnabled: inboxEnabled,
		ToastEnabled: toastEnabled,
		SoundEnabled: soundEnabled,
	}
}

var tPreferences = newPreferenceTable("", "fivenet_user_notification_preferences", "")
