package dmv

import (
	"context"
	"database/sql"
	"strings"

	"github.com/galexrt/arpanet/pkg/mstlystcdata"
	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/galexrt/arpanet/proto/resources/common/database"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	vehicle = table.OwnedVehicles.AS("vehicle")
	user    = table.Users.AS("usershortni")
)

type Server struct {
	DMVServiceServer

	db *sql.DB
	p  perms.Permissions
	c  *mstlystcdata.Enricher
}

func NewServer(db *sql.DB, p perms.Permissions, c *mstlystcdata.Enricher) *Server {
	return &Server{
		db: db,
		p:  p,
		c:  c,
	}
}

func (s *Server) FindVehicles(ctx context.Context, req *FindVehiclesRequest) (*FindVehiclesResponse, error) {
	condition := jet.Bool(true)
	if req.Search != "" {
		condition = jet.AND(condition, jet.BoolExp(jet.Raw(
			"MATCH(plate) AGAINST ($search IN NATURAL LANGUAGE MODE)",
			jet.RawArgs{"$search": req.Search},
		)))
	}
	if req.Type != "" {
		req.Type = strings.ReplaceAll(req.Type, "%", "") + "%"
		condition = jet.AND(condition, jet.BoolExp(vehicle.Type.LIKE(jet.String(req.Type))))
	}
	userCondition := user.Identifier.EQ(vehicle.Owner)
	if req.UserId != 0 {
		condition = jet.AND(condition,
			jet.BoolExp(user.Identifier.EQ(vehicle.Owner)),
			user.ID.EQ(jet.Int32(req.UserId)),
		)
		userCondition = jet.AND(userCondition, user.ID.EQ(jet.Int32(req.UserId)))
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
		return nil, err
	}

	resp := &FindVehiclesResponse{
		Pagination: database.EmptyPaginationResponse(req.Pagination.Offset),
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := vehicle.
		SELECT(
			vehicle.Plate,
			vehicle.Model,
			vehicle.Type,
			user.ID,
			user.Identifier,
			user.Job,
			user.JobGrade,
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
		ORDER_BY(
			vehicle.Type.ASC(),
			vehicle.Plate.ASC(),
		).
		OFFSET(req.Pagination.Offset).
		LIMIT(database.DefaultPageLimit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Vehicles); err != nil {
		return nil, err
	}

	database.PaginationHelper(resp.Pagination,
		count.TotalCount,
		req.Pagination.Offset,
		len(resp.Vehicles))

	for i := 0; i < len(resp.Vehicles); i++ {
		if resp.Vehicles[i].Owner == nil {
			continue
		}

		s.c.EnrichJobInfo(resp.Vehicles[i].Owner)
	}

	return resp, nil
}
