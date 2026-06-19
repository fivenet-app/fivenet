package completorstore

import (
	"context"
	"database/sql"

	documentscategory "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/documents/category"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
)

type IStore interface {
	CompleteCitizens(ctx context.Context, q CitizensQuery) ([]*usershort.UserShort, error)
	CompleteDocumentCategories(
		ctx context.Context,
		q DocumentCategoriesQuery,
	) ([]*documentscategory.Category, error)
}

type Store struct {
	db       *sql.DB
	customDB *config.CustomDB
}

func New(db *sql.DB, customDB *config.CustomDB) IStore {
	return &Store{db: db, customDB: customDB}
}
