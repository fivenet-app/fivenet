package completor

import (
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStoreCompleteDocumentCategoriesAppliesSearchAndOrdering(t *testing.T) {
	t.Parallel()

	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	store := New(db, &config.CustomDB{})

	expectedQuery := regexp.QuoteMeta(`FROM fivenet_documents_categories AS category`) +
		`(?s).*` + regexp.QuoteMeta(`category.job IN (?, ?)`) +
		`(?s).*` + regexp.QuoteMeta(`category.name LIKE ?`) +
		`(?s).*` + regexp.QuoteMeta(`ORDER BY category.id IN (?, ?) DESC, category.job = ? DESC, category.sort_key ASC LIMIT ?;`)

	mock.ExpectQuery(expectedQuery).
		WithArgs("medic", "police", "%Permit%", int64(8), int64(9), "police", int64(15)).
		WillReturnRows(sqlmock.NewRows([]string{
			"category.id",
			"category.name",
			"category.description",
			"category.job",
			"category.color",
			"category.icon",
		}).AddRow(int64(8), "Permit", nil, "police", nil, nil))

	categories, err := store.CompleteDocumentCategories(t.Context(), DocumentCategoriesQuery{
		Search:      "Permit",
		CategoryIDs: []int64{8, 9},
		Jobs:        []string{"medic", "police"},
		CurrentJob:  "police",
	})
	require.NoError(t, err)
	require.Len(t, categories, 1)
	assert.Equal(t, int64(8), categories[0].GetId())
	assert.Equal(t, "Permit", categories[0].GetName())
	require.NoError(t, mock.ExpectationsWereMet())
}
