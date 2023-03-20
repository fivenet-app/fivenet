package dmv

import (
	"context"
	"database/sql"
	"strings"

	"github.com/galexrt/arpanet/pkg/complhelper"
	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/galexrt/arpanet/proto/resources/common/database"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	ve = table.OwnedVehicles.AS("vehicle")
	us = table.Users.AS("usershortni")
)

type Server struct {
	DMVServiceServer

	db *sql.DB
	p  perms.Permissions
	c  *complhelper.Completor
}

func NewServer(db *sql.DB, p perms.Permissions, c *complhelper.Completor) *Server {
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
		condition = jet.AND(condition, jet.BoolExp(ve.Type.LIKE(jet.String(req.Type))))
	}
	userCondition := us.Identifier.EQ(ve.Owner)
	if req.UserId != 0 {
		condition = jet.AND(condition,
			jet.BoolExp(us.Identifier.EQ(ve.Owner)),
			us.ID.EQ(jet.Int32(req.UserId)),
		)
		userCondition = jet.AND(userCondition, us.ID.EQ(jet.Int32(req.UserId)))
	}

	countStmt := ve.
		SELECT(
			jet.COUNT(ve.Owner).AS("total_count"),
		).
		FROM(
			ve.
				LEFT_JOIN(us,
					userCondition,
				),
		).
		WHERE(condition)

	var count struct{ TotalCount int64 }
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, err
	}

	resp := &FindVehiclesResponse{
		Offset:     req.Offset,
		TotalCount: count.TotalCount,
		End:        0,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := ve.
		SELECT(
			ve.Plate,
			ve.Model,
			ve.Type,
			us.ID,
			us.Identifier,
			us.Job,
			us.JobGrade,
			us.Firstname,
			us.Lastname,
		).
		FROM(
			ve.
				LEFT_JOIN(us,
					userCondition,
				),
		).
		WHERE(condition).
		ORDER_BY(
			ve.Type.ASC(),
			ve.Plate.ASC(),
		).
		OFFSET(req.Offset).
		LIMIT(database.PaginationLimit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Vehicles); err != nil {
		return nil, err
	}

	resp.TotalCount = count.TotalCount
	if req.Offset >= resp.TotalCount {
		resp.Offset = 0
	} else {
		resp.Offset = req.Offset
	}
	resp.End = resp.Offset + int64(len(resp.Vehicles))

	for i := 0; i < len(resp.Vehicles); i++ {
		if resp.Vehicles[i].Owner == nil {
			continue
		}

		s.c.ResolveJob(resp.Vehicles[i].Owner)
	}

	return resp, nil
}
