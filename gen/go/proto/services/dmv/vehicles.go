package dmv

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	permscitizenstore "github.com/fivenet-app/fivenet/gen/go/proto/services/citizenstore/perms"
	errorsdmv "github.com/fivenet-app/fivenet/gen/go/proto/services/dmv/errors"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/mstlystcdata"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/server/audit"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	grpc "google.golang.org/grpc"
)

var (
	tVehicles = table.OwnedVehicles.AS("vehicle")
	tUsers    = table.Users.AS("usershort")
)

type Server struct {
	DMVServiceServer

	db       *sql.DB
	ps       perms.Permissions
	enricher *mstlystcdata.Enricher
	aud      audit.IAuditer
	customDB config.CustomDB
}

type Params struct {
	fx.In

	DB       *sql.DB
	Ps       perms.Permissions
	Enricher *mstlystcdata.Enricher
	Aud      audit.IAuditer
	Config   *config.Config
}

func NewServer(p Params) *Server {
	return &Server{
		db:       p.DB,
		ps:       p.Ps,
		enricher: p.Enricher,
		aud:      p.Aud,
		customDB: p.Config.Database.Custom,
	}
}

func (s *Server) RegisterServer(srv *grpc.Server) {
	RegisterDMVServiceServer(srv, s)
}

func (s *Server) ListVehicles(ctx context.Context, req *ListVehiclesRequest) (*ListVehiclesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := jet.Bool(true)
	userCondition := tUsers.Identifier.EQ(tVehicles.Owner)
	if req.LicensePlate != nil && *req.LicensePlate != "" {
		condition = jet.AND(condition, tVehicles.Plate.LIKE(jet.String(
			strings.ReplaceAll(*req.LicensePlate, "%", "")+"%",
		)))
	}

	// Make sure the model column is available
	modelColumn := s.customDB.Columns.Vehicle.GetModel(tVehicles.Alias())
	if modelColumn != nil && req.Model != nil && *req.Model != "" {
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
		s.aud.Log(&model.FivenetAuditLog{
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
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsdmv.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, 16)
	resp := &ListVehiclesResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []jet.OrderByClause{
		tVehicles.Type.ASC(),
	}
	if req.Sort != nil {
		var column jet.Column
		switch req.Sort.Column {
		case "model":
			column = tVehicles.Model
		case "plate":
			fallthrough
		default:
			column = tVehicles.Plate
		}

		if req.Sort.Direction == database.AscSortDirection {
			orderBys = append(orderBys, column.ASC())
		} else {
			orderBys = append(orderBys, column.DESC())
		}
	} else {
		orderBys = append(orderBys, tVehicles.Plate.ASC())
	}

	columns := dbutils.Columns{
		modelColumn,
		jet.REPLACE(tVehicles.Type, jet.String("_"), jet.String(" ")).AS("vehicle.type"),
		tUsers.ID,
		tUsers.Firstname,
		tUsers.Lastname,
		tUsers.Dateofbirth,
	}

	// Field Permission Check
	fieldsAttr, err := s.ps.Attr(userInfo, permscitizenstore.CitizenStoreServicePerm, permscitizenstore.CitizenStoreServiceListCitizensPerm, permscitizenstore.CitizenStoreServiceListCitizensFieldsPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsdmv.ErrFailedQuery)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}

	if slices.Contains(fields, "PhoneNumber") {
		columns = append(columns, tUsers.PhoneNumber)
	}

	stmt := tVehicles.
		SELECT(
			tVehicles.Plate,
			columns.Get()...,
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
		return nil, errswrap.NewError(err, errorsdmv.ErrFailedQuery)
	}

	resp.Pagination.Update(len(resp.Vehicles))

	return resp, nil
}
