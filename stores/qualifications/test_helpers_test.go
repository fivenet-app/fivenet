package qualificationsstore

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/v2026/pkg/access"
)

func testParams(db *sql.DB) Params {
	return Params{
		DB:     db,
		Access: access.NewQualificationsSubjectObjectAccess(db),
	}
}
