package dmv

import (
	"context"
	"database/sql"
	"strings"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/pkg/audit"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	vehicle = table.OwnedVehicles.AS("vehicle")
	user    = table.Users.AS("usershort")
)

var (
	FailedQueryErr = status.Error(codes.Internal, "Failed to retrieve vehicles data!")
)

type Server struct {
	DMVServiceServer

	db *sql.DB
	p  perms.Permissions
	c  *mstlystcdata.Enricher
	a  audit.IAuditer
}

func NewServer(db *sql.DB, p perms.Permissions, c *mstlystcdata.Enricher, aud audit.IAuditer) *Server {
	return &Server{
		db: db,
		p:  p,
		c:  c,
		a:  aud,
	}
}

func (s *Server) ListVehicles(ctx context.Context, req *ListVehiclesRequest) (*ListVehiclesResponse, error) {
	condition := jet.Bool(true)
	userCondition := user.Identifier.EQ(vehicle.Owner)
	if req.Search != "" {
		req.Search = strings.ReplaceAll(req.Search, "%", "") + "%"
		condition = jet.AND(condition, vehicle.Plate.LIKE(jet.String(req.Search)))
	}
	if req.Model != "" {
		req.Model = strings.ReplaceAll(req.Model, "%", "") + "%"
		condition = jet.AND(condition, vehicle.Model.LIKE(jet.String(req.Model)))
	}
	if req.UserId != 0 {
		condition = jet.AND(condition,
			user.Identifier.EQ(vehicle.Owner),
			user.ID.EQ(jet.Int32(req.UserId)),
		)
		userCondition = jet.AND(userCondition, user.ID.EQ(jet.Int32(req.UserId)))
	}

	if req.Pagination.Offset <= 0 {
		userInfo := auth.MustGetUserInfoFromContext(ctx)

		s.a.AddEntryWithData(&model.FivenetAuditLog{
			Service: DMVService_ServiceDesc.ServiceName,
			Method:  "ListVehicles",
			UserID:  userInfo.UserId,
			UserJob: userInfo.Job,
			State:   int16(rector.EVENT_TYPE_VIEWED),
		}, req)
	}

	countStmt := vehicle.
		SELECT(
			jet.COUNT(vehicle.Owner).AS("datacount.totalcount"),
		).
		FROM(
			vehicle.
				LEFT_JOIN(user,
					userCondition,
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, FailedQueryErr
	}

	pag, limit := req.Pagination.GetResponse()
	resp := &ListVehiclesResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	// Convert our proto abstracted `common.OrderBy` to actual gorm order by instructions
	orderBys := []jet.OrderByClause{}
	if len(req.OrderBy) > 0 {
		for _, orderBy := range req.OrderBy {
			var column jet.Column
			switch orderBy.Column {
			case "plate":
				column = vehicle.Plate
			case "model":
			default:
				column = vehicle.Model
			}

			if orderBy.Desc {
				orderBys = append(orderBys, column.DESC())
			} else {
				orderBys = append(orderBys, column.ASC())
			}
		}
	} else {
		orderBys = append(orderBys,
			vehicle.Type.ASC(),
			vehicle.Plate.ASC(),
		)
	}

	stmt := vehicle.
		SELECT(
			vehicle.Plate,
			vehicle.Model,
			vehicle.Type,
			user.ID,
			user.Identifier,
			user.Firstname,
			user.Lastname,
		).
		FROM(
			vehicle.
				LEFT_JOIN(user,
					userCondition,
				),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(orderBys...).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Vehicles); err != nil {
		return nil, FailedQueryErr
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Vehicles))

	return resp, nil
}
