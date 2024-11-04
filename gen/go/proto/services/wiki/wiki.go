package wiki

import (
	"context"
	"database/sql"
	"errors"

	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/wiki"
	errorswiki "github.com/fivenet-app/fivenet/gen/go/proto/services/wiki/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	tPage       = table.FivenetWikiPages.AS("page")
	tPageShort  = table.FivenetWikiPages.AS("pageshort")
	tPJobAccess = table.FivenetWikiPageJobAccess.AS("job_access")
	tUsers      = table.Users
	tCreator    = tUsers.AS("creator")
)

type Server struct {
	WikiServiceServer

	logger *zap.Logger
	db     *sql.DB
	aud    audit.IAuditer
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	DB     *sql.DB
	Audit  audit.IAuditer
}

func NewServer(p Params) *Server {
	s := &Server{
		logger: p.Logger.Named("wiki"),
		db:     p.DB,
		aud:    p.Audit,
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterWikiServiceServer(srv, s)
}

func (s *Server) ListPages(ctx context.Context, req *ListPagesRequest) (*ListPagesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := jet.AND(
		tPage.Job.EQ(jet.String(userInfo.Job)),
		tPageShort.DeletedAt.IS_NULL(),
		jet.OR(
			tPageShort.CreatorID.EQ(jet.Int32(userInfo.UserId)),
			jet.AND(
				tPJobAccess.Access.IS_NOT_NULL(),
				tPJobAccess.Access.GT(jet.Int32(int32(wiki.AccessLevel_ACCESS_LEVEL_BLOCKED))),
			),
		),
	)

	countStmt := tPageShort.
		SELECT(
			jet.COUNT(jet.DISTINCT(tPageShort.ID)),
		).
		FROM(
			tPageShort.
				LEFT_JOIN(tPJobAccess,
					tPJobAccess.PageID.EQ(tPageShort.ID).
						AND(tPJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tPJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorswiki.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponse(count.TotalCount)
	resp := &ListPagesResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := tPageShort.
		SELECT(
			tPageShort.ID,
			tPageShort.Job,
			tPageShort.Path,
			tPageShort.CreatedAt.AS("page_meta.created_at"),
			tPageShort.UpdatedAt.AS("page_meta.updated_at"),
			tPageShort.DeletedAt.AS("page_meta.deleted_at"),
			tPageShort.Title.AS("page_meta.title"),
			tPageShort.Description.AS("page_meta.description"),
			tPageShort.CreatorID.AS("page_meta.creator_id"),
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tPageShort.ContentType.AS("page_meta.content_Type"),
		).
		FROM(
			tPageShort.
				LEFT_JOIN(tPJobAccess,
					tPJobAccess.PageID.EQ(tPageShort.ID).
						AND(tPJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tPJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				).
				LEFT_JOIN(tCreator,
					tPageShort.CreatorID.EQ(tCreator.ID),
				),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(tPageShort.Path.ASC()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, resp.Pages); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errorswiki.ErrFailedQuery
		}
	}

	return resp, nil
}

func (s *Server) GetPage(ctx context.Context, req *GetPageRequest) (*GetPageResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := jet.AND(
		tPage.Job.EQ(jet.String(userInfo.Job)),
		tPage.DeletedAt.IS_NULL(),
		jet.OR(
			tPage.CreatorID.EQ(jet.Int32(userInfo.UserId)),
			jet.AND(
				tPJobAccess.Access.IS_NOT_NULL(),
				tPJobAccess.Access.GT(jet.Int32(int32(wiki.AccessLevel_ACCESS_LEVEL_BLOCKED))),
			),
		),
	)

	if req.Id != nil && *req.Id > 0 {
		condition = tPage.ID.EQ(jet.Uint64(*req.Id))
	} else if req.Path != nil {
		if *req.Path == "" {
			*req.Path = "/"
		}
		condition = tPage.Path.EQ(jet.String(*req.Path))
	} else {
		return nil, errorswiki.ErrFailedQuery
	}

	stmt := tPage.
		SELECT(
			tPage.ID,
			tPage.Job,
			tPage.Path,
			tPage.CreatedAt.AS("page_meta.created_at"),
			tPage.UpdatedAt.AS("page_meta.updated_at"),
			tPage.DeletedAt.AS("page_meta.deleted_at"),
			tPage.Title.AS("page_meta.title"),
			tPage.Description.AS("page_meta.description"),
			tPage.CreatorID.AS("page_meta.creator_id"),
			tCreator.ID,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tPage.ContentType.AS("page_meta.content_Type"),
			tPage.Content,
			tPage.Data,
		).
		FROM(
			tPage.
				LEFT_JOIN(tPJobAccess,
					tPJobAccess.PageID.EQ(tPage.ID).
						AND(tPJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tPJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				).
				LEFT_JOIN(tCreator,
					tPage.CreatorID.EQ(tCreator.ID),
				),
		).
		WHERE(condition).
		LIMIT(1)

	resp := &GetPageResponse{
		Page: &wiki.Page{},
	}

	if err := stmt.QueryContext(ctx, s.db, resp.Page); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errorswiki.ErrFailedQuery
		}
	}

	if resp.Page.Id <= 0 {
		resp.Page = nil
	}

	return resp, nil
}

func (s *Server) CreateOrUpdatePage(ctx context.Context, req *CreateOrUpdatePageRequest) (*CreateOrUpdatePageResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	_ = userInfo

	// TODO

	return nil, nil
}

func (s *Server) DeletePage(ctx context.Context, req *DeletePageRequest) (*DeletePageResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	_ = userInfo

	// TODO

	return nil, nil
}
