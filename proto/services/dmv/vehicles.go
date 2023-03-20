package dmv

import (
	"context"
	"database/sql"

	"github.com/galexrt/arpanet/pkg/complhelper"
	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/galexrt/arpanet/proto/resources/common/database"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

var (
	ve = table.OwnedVehicles.AS("vehicle")
	us = table.Users.AS("usershort")
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
	var condition jet.BoolExpression
	if req.Search != "" {
		condition = jet.BoolExp(jet.Raw(
			"MATCH(plate) AGAINST ($search IN NATURAL LANGUAGE MODE)",
			jet.RawArgs{"$search": req.Search},
		))
	} else {
		condition = jet.Bool(true)
	}

	countStmt := ve.
		SELECT(
			jet.COUNT(ve.Owner).AS("total_count"),
		).
		FROM(
			ve,
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
					us.Identifier.EQ(ve.Owner),
				),
		).
		WHERE(condition).
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
		s.c.ResolveJob(resp.Vehicles[i].Owner)
	}

	return resp, nil
}
