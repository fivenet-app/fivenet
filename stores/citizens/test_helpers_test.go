package citizensstore

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
)

func testParams(db *sql.DB) Params {
	return Params{
		DB:           db,
		CustomDB:     &config.CustomDB{},
		LabelsAccess: access.NewCitizenLabelsSubjectObjectAccess(db),
	}
}
