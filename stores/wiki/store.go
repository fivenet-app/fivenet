package wikistore

import (
	"context"
	"database/sql"

	reswiki "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki"
	wikiactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki/activity"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/go-jet/jet/v2/qrm"
)

type IStore interface {
	ListPages(ctx context.Context, q ListPagesQuery) (*ListPagesResult, error)
	GetPage(ctx context.Context, pageID int64, withContent bool) (*reswiki.Page, error)
	GetPageOrderInfo(ctx context.Context, q qrm.DB, pageID int64) (*PageOrderInfo, error)
	NextPageGroupRank(
		ctx context.Context,
		q qrm.DB,
		job string,
		parentID *int64,
		startpage bool,
		excludeID int64,
	) (string, error)
	InsertPageGroupRank(
		ctx context.Context,
		q qrm.DB,
		job string,
		parentID *int64,
		startpage bool,
		excludeID int64,
		beforeID, afterID *int64,
	) (string, error)
	CountPageActivity(ctx context.Context, q PageActivityQuery) (int64, error)
	ListPageActivity(ctx context.Context, q PageActivityQuery) ([]*wikiactivity.PageActivity, error)
	AddPageActivity(
		ctx context.Context,
		tx qrm.DB,
		activity *wikiactivity.PageActivity,
	) (int64, error)
	CountPageChildren(ctx context.Context, pageID int64) (int64, error)
}

type Store struct {
	db     *sql.DB
	access *access.SubjectObjectAccess
}

func New(db *sql.DB) IStore {
	return &Store{db: db, access: access.NewWikiPageSubjectObjectAccess(db)}
}
