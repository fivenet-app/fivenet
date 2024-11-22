package notificator

import (
	"context"
	"database/sql"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/pkg/events"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	"go.uber.org/zap"
	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var tNotifications = table.FivenetNotifications

var (
	ErrFailedRequest = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.NotificatorService.ErrFailedRequest"}, nil)
	ErrFailedStream  = common.I18nErr(codes.InvalidArgument, &common.TranslateItem{Key: "errors.NotificatorService.ErrFailedStream"}, nil)
)

type Server struct {
	NotificatorServiceServer

	logger   *zap.Logger
	db       *sql.DB
	p        perms.Permissions
	tm       *auth.TokenMgr
	ui       userinfo.UserInfoRetriever
	js       *events.JSWrapper
	enricher *mstlystcdata.Enricher
	appCfg   appconfig.IConfig
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger    *zap.Logger
	DB        *sql.DB
	Perms     perms.Permissions
	TM        *auth.TokenMgr
	UI        userinfo.UserInfoRetriever
	JS        *events.JSWrapper
	Enricher  *mstlystcdata.Enricher
	AppConfig appconfig.IConfig
}

func NewServer(p Params) *Server {
	s := &Server{
		logger:   p.Logger,
		db:       p.DB,
		p:        p.Perms,
		tm:       p.TM,
		ui:       p.UI,
		js:       p.JS,
		enricher: p.Enricher,
		appCfg:   p.AppConfig,
	}

	return s
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterNotificatorServiceServer(srv, s)
}

func (s *Server) GetNotifications(ctx context.Context, req *GetNotificationsRequest) (*GetNotificationsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tNotifications := tNotifications.AS("notification")
	condition := tNotifications.UserID.EQ(jet.Int32(userInfo.UserId))
	if req.IncludeRead != nil && !*req.IncludeRead {
		condition = condition.AND(tNotifications.ReadAt.IS_NULL())
	}

	if len(req.Categories) > 0 {
		categoryIds := make([]jet.Expression, len(req.Categories))
		for i := 0; i < len(req.Categories); i++ {
			categoryIds[i] = jet.Int16(int16(req.Categories[i]))
		}

		condition = condition.AND(tNotifications.Category.IN(categoryIds...))
	}

	countStmt := tNotifications.
		SELECT(
			jet.COUNT(tNotifications.ID).AS("datacount.totalcount"),
		).
		FROM(tNotifications).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, ErrFailedRequest)
		}
	}

	pag, limit := req.Pagination.GetResponse(count.TotalCount)
	resp := &GetNotificationsResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := tNotifications.
		SELECT(
			tNotifications.ID,
			tNotifications.CreatedAt,
			tNotifications.ReadAt,
			tNotifications.UserID,
			tNotifications.Title,
			tNotifications.Type,
			tNotifications.Content,
			tNotifications.Category,
			tNotifications.Data,
		).
		FROM(tNotifications).
		WHERE(
			condition,
		).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(tNotifications.ID.DESC()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Notifications); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, ErrFailedRequest)
		}
	}

	resp.Pagination.Update(len(resp.Notifications))

	return resp, nil
}

func (s *Server) MarkNotifications(ctx context.Context, req *MarkNotificationsRequest) (*MarkNotificationsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tNotifications.UserID.EQ(
		jet.Int32(userInfo.UserId)).AND(
		tNotifications.ReadAt.IS_NULL(),
	)
	// If not all
	if len(req.Ids) > 0 {
		ids := make([]jet.Expression, len(req.Ids))
		for i := 0; i < len(req.Ids); i++ {
			ids[i] = jet.Uint64(req.Ids[i])
		}
		condition = condition.AND(tNotifications.ID.IN(ids...))
	} else if req.All == nil || !*req.All {
		return &MarkNotificationsResponse{}, nil
	}

	stmt := tNotifications.
		UPDATE(
			tNotifications.ReadAt,
		).
		SET(
			tNotifications.ReadAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(condition)

	res, err := stmt.ExecContext(ctx, s.db)
	if err != nil {
		return nil, errswrap.NewError(err, ErrFailedRequest)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return nil, errswrap.NewError(err, ErrFailedRequest)
	}

	return &MarkNotificationsResponse{
		Updated: uint64(rows),
	}, nil
}

func (s *Server) getNotificationCount(ctx context.Context, userId int32) (int32, error) {
	stmt := tNotifications.
		SELECT(
			jet.COUNT(tNotifications.ID).AS("count"),
		).
		FROM(tNotifications).
		WHERE(jet.AND(
			tNotifications.UserID.EQ(jet.Int32(userId)),
			tNotifications.ReadAt.IS_NULL(),
		)).
		ORDER_BY(tNotifications.ID.DESC())

	var dest struct {
		Count int32
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, errswrap.NewError(err, ErrFailedStream)
		}
	}

	return dest.Count, nil
}
