package wikistore

import (
	"context"
	"database/sql"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/content"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	reswiki "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki"
	wikiaccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki/access"
	wikiactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/wiki/activity"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/go-jet/jet/v2/qrm"
)

type IStore interface {
	ListPages(ctx context.Context, q ListPagesQuery) (*ListPagesResult, error)
	GetPage(ctx context.Context, pageID int64, withContent bool) (*reswiki.Page, error)
	CreatePage(
		ctx context.Context,
		tx qrm.DB,
		userInfo *userinfo.UserInfo,
		parentID *int64,
		contentType content.ContentType,
		pageAccess *wikiaccess.PageAccess,
	) (int64, *wikiaccess.PageAccess, error)
	UpdatePage(
		ctx context.Context,
		tx qrm.DB,
		userInfo *userinfo.UserInfo,
		page *reswiki.Page,
		sortRank string,
	) (*wikiaccess.PageAccess, error)
	DeletePage(
		ctx context.Context,
		tx qrm.DB,
		pageId int64,
		deletedAtTime *timestamp.Timestamp,
		parentId int64,
	) error
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
	db             *sql.DB
	access         *access.SubjectObjectAccess
	accessResolver *access.SubjectResolver
}

func New(db *sql.DB) IStore {
	return &Store{
		db:             db,
		access:         access.NewWikiPageSubjectObjectAccess(db),
		accessResolver: access.NewSubjectResolver(db),
	}
}
