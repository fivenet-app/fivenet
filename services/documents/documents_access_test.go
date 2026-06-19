package documents

import (
	"testing"

	resourcesaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/access"
	documentsaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/access"
	"github.com/stretchr/testify/assert"
)

func TestDocumentAccessHasDuplicates(t *testing.T) {
	t.Parallel()

	access := &documentsaccess.DocumentAccess{
		Jobs: []*resourcesaccess.JobAccess{
			{
				Job:          "police",
				MinimumGrade: 3,
				Access:       2,
			},
			{
				Job:          "police",
				MinimumGrade: 3,
				Access:       7,
			},
		},
	}

	assert.True(t, documentsaccess.DocumentAccessHasDuplicates(access))
	assert.False(t, documentsaccess.DocumentAccessHasDuplicates(&documentsaccess.DocumentAccess{
		Jobs: []*resourcesaccess.JobAccess{
			{
				Job:          "police",
				MinimumGrade: 3,
				Access:       2,
			},
			{
				Job:          "police",
				MinimumGrade: 4,
				Access:       7,
			},
		},
	}))
}
