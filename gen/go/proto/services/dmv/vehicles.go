package dmv

import (
	"context"
	"database/sql"
	"strings"

	"github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	permscitizenstore "github.com/galexrt/fivenet/gen/go/proto/services/citizenstore/perms"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	tVehicles = table.OwnedVehicles.AS("vehicle")
	tUsers    = table.Users.AS("usershort")
)

var (
	ErrFailedQuery = status.Error(codes.Internal, "errors.DMVService.ErrFailedQuery")
)

type Server struct {
	DMVServiceServer

	db *sql.DB
	ps perms.Permissions
	c  *mstlystcdata.Enricher
	a  audit.IAuditer
}

func NewServer(db *sql.DB, p perms.Permissions, c *mstlystcdata.Enricher, aud audit.IAuditer) *Server {
	return &Server{
		db: db,
		ps: p,
		c:  c,
		a:  aud,
	}
}

func (s *Server) ListVehicles(ctx context.Context, req *ListVehiclesRequest) (*ListVehiclesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := jet.Bool(true)
	userCondition := tUsers.Identifier.EQ(tVehicles.Owner)
	if req.Search != nil && *req.Search != "" {
		condition = jet.AND(condition, tVehicles.Plate.LIKE(jet.String(
			strings.ReplaceAll(*req.Search, "%", "")+"%",
		)))
	}
	if req.Model != nil && *req.Model != "" {
		condition = jet.AND(condition, tVehicles.Model.LIKE(jet.String(
			strings.ReplaceAll(*req.Model, "%", "")+"%",
		)))
	}
	if req.UserId != nil && *req.UserId != 0 {
		condition = jet.AND(condition,
			tUsers.Identifier.EQ(tVehicles.Owner),
			tUsers.ID.EQ(jet.Int32(*req.UserId)),
		)
		userCondition = jet.AND(userCondition, tUsers.ID.EQ(jet.Int32(*req.UserId)))
	}

	if req.Pagination.Offset <= 0 {
		s.a.Log(&model.FivenetAuditLog{
			Service: DMVService_ServiceDesc.ServiceName,
			Method:  "ListVehicles",
			UserID:  userInfo.UserId,
			UserJob: userInfo.Job,
			State:   int16(rector.EventType_EVENT_TYPE_VIEWED),
		}, req)
	}

	countStmt := tVehicles.
		SELECT(
			jet.COUNT(tVehicles.Owner).AS("datacount.totalcount"),
		).
		FROM(
			tVehicles.
				LEFT_JOIN(tUsers,
					userCondition,
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, ErrFailedQuery
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(15)
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
			case "model":
				column = tVehicles.Model
			case "plate":
				fallthrough
			default:
				column = tVehicles.Plate
			}

			if orderBy.Desc {
				orderBys = append(orderBys, column.DESC())
			} else {
				orderBys = append(orderBys, column.ASC())
			}
		}
	} else {
		orderBys = append(orderBys,
			tVehicles.Type.ASC(),
			tVehicles.Plate.ASC(),
		)
	}

	columns := []jet.Projection{
		tVehicles.Plate,
		tVehicles.Model,
		jet.REPLACE(tVehicles.Type, jet.String("_"), jet.String(" ")).AS("vehicle.type"),
		tUsers.ID,
		tUsers.Identifier,
		tUsers.Firstname,
		tUsers.Lastname,
	}

	// Field Permission Check
	fieldsAttr, err := s.ps.Attr(userInfo, permscitizenstore.CitizenStoreServicePerm, permscitizenstore.CitizenStoreServiceListCitizensPerm, permscitizenstore.CitizenStoreServiceListCitizensFieldsPermField)
	if err != nil {
		return nil, ErrFailedQuery
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}

	if utils.InSlice(fields, "PhoneNumber") {
		columns = append(columns, tUsers.PhoneNumber)
	}

	stmt := tVehicles.
		SELECT(
			columns[0],
			columns[1:]...,
		).
		FROM(
			tVehicles.
				LEFT_JOIN(tUsers,
					userCondition,
				),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(orderBys...).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Vehicles); err != nil {
		return nil, ErrFailedQuery
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Vehicles))

	return resp, nil
}
