package documentsstore

import (
	"database/sql"

	"github.com/fivenet-app/fivenet/v2026/pkg/access"
)

func testParams(db *sql.DB) Params {
	return Params{
		DB:             db,
		SubjectAccess:  access.NewDocumentsSubjectObjectAccess(db),
		StampAccess:    access.NewDocumentStampsSubjectObjectAccess(db),
		TemplateAccess: access.NewDocumentTemplatesSubjectObjectAccess(db),
	}
}
