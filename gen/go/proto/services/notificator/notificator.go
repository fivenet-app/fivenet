package notificator

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	nots = table.FivenetNotifications
)

type Server struct {
	NotificatorServiceServer

	logger *zap.Logger
	db     *sql.DB
	p      perms.Permissions
}

func NewServer(logger *zap.Logger, db *sql.DB, p perms.Permissions) *Server {
	return &Server{
		logger: logger,
		db:     db,
		p:      p,
	}
}

func (s *Server) PermissionStreamFuncOverride(ctx context.Context, srv interface{}, info *grpc.StreamServerInfo) (context.Context, error) {
	// Skip permission check for the notificator services
	return ctx, nil
}

func (s *Server) GetNotifications(ctx context.Context, req *GetNotificationsRequest) (*GetNotificationsResponse, error) {
	userInfo := auth.GetUserInfoFromContext(ctx)

	condition := nots.UserID.EQ(jet.Int32(userInfo.UserId))
	if req.IncludeRead {
		condition = jet.AND(
			condition,
			nots.ReadAt.IS_NOT_NULL(),
		)
	}

	countStmt := nots.
		SELECT(
			jet.COUNT(nots.ID).AS("datacount.totalcount"),
		).
		FROM(nots).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, err
	}

	pag, limit := req.Pagination.GetResponse()
	resp := &GetNotificationsResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := nots.
		SELECT(
			nots.AllColumns,
		).
		FROM(nots).
		WHERE(
			condition,
		).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, resp.Notifications); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, err
		}
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Notifications))

	return resp, nil
}

func (s *Server) ReadNotifications(ctx context.Context, req *ReadNotificationsRequest) (*ReadNotificationsResponse, error) {
	userInfo := auth.GetUserInfoFromContext(ctx)

	ids := make([]jet.Expression, len(req.Ids))
	for i := 0; i < len(req.Ids); i++ {
		ids[i] = jet.Uint64(req.Ids[i])
	}

	stmt := nots.
		UPDATE(
			nots.ReadAt,
		).
		SET(
			nots.ReadAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(
			jet.AND(
				nots.UserID.EQ(jet.Int32(userInfo.UserId)),
				nots.ID.IN(ids...),
			),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	return &ReadNotificationsResponse{}, nil
}

func (s *Server) Stream(req *StreamRequest, srv NotificatorService_StreamServer) error {
	userInfo := auth.GetUserInfoFromContext(srv.Context())

	nots := nots.AS("notification")
	stmt := nots.
		SELECT(
			nots.ID,
			nots.Title,
			nots.Type,
			nots.Content,
			nots.Data,
		).
		FROM(nots).
		ORDER_BY(nots.ID.DESC()).
		LIMIT(10)

	resp := &StreamResponse{
		LastId: req.LastId,
	}
	for {
		q := stmt.
			WHERE(
				jet.AND(
					nots.UserID.EQ(jet.Int32(userInfo.UserId)),
					nots.ID.GT(jet.Uint64(req.LastId)),
				),
			)

		if err := q.QueryContext(srv.Context(), s.db, &resp.Notifications); err != nil {
			return err
		}

		// Update last id for user
		if len(resp.Notifications) > 0 {
			req.LastId = resp.Notifications[0].Id
			resp.LastId = resp.Notifications[0].Id
		}

		if err := srv.Send(resp); err != nil {
			return err
		}
		resp.Notifications = nil

		select {
		case <-srv.Context().Done():
			return nil
		case <-time.After(30 * time.Second):
		}
	}
}
