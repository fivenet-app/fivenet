package jobs

import (
	"context"
	"database/sql"
	"strings"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/server/audit"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	tUser = table.Users.AS("user")
)

var (
	ErrFailedQuery = status.Error(codes.Internal, "errors.JobsService.ErrFailedQuery")
)

type Server struct {
	JobsServiceServer

	db       *sql.DB
	p        perms.Permissions
	enricher *mstlystcdata.Enricher
	aud      audit.IAuditer
}

func NewServer(db *sql.DB, p perms.Permissions, enricher *mstlystcdata.Enricher, aud audit.IAuditer) *Server {
	return &Server{
		db:       db,
		p:        p,
		enricher: enricher,
		aud:      aud,
	}
}

// TODO employee warns/notes management
// TODO time clock

func (s *Server) ColleaguesList(ctx context.Context, req *ColleaguesListRequest) (*ColleaguesListResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	selectors := jet.ProjectionList{
		tUser.ID,
		tUser.Identifier,
		tUser.Firstname,
		tUser.Lastname,
		tUser.Job,
		tUser.JobGrade,
		tUser.Dateofbirth,
		tUser.PhoneNumber,
	}

	condition := tUser.Job.EQ(jet.String(userInfo.Job))

	req.SearchName = strings.TrimSpace(req.SearchName)
	req.SearchName = strings.ReplaceAll(req.SearchName, "%", "")
	req.SearchName = strings.ReplaceAll(req.SearchName, " ", "%")
	if req.SearchName != "" {
		req.SearchName = "%" + req.SearchName + "%"
		condition = condition.AND(
			jet.CONCAT(tUser.Firstname, jet.String(" "), tUser.Lastname).
				LIKE(jet.String(req.SearchName)),
		)
	}

	// Get total count of values
	countStmt := tUser.
		SELECT(
			jet.COUNT(tUser.ID).AS("datacount.totalcount"),
		).
		FROM(tUser).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, ErrFailedQuery
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(15)
	resp := &ColleaguesListResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := tUser.
		SELECT(
			selectors[0], selectors[1:]...,
		).
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname_fulltext")).
		FROM(tUser).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(
			tUser.JobGrade.ASC(),
			tUser.Firstname.ASC(),
			tUser.Lastname.ASC(),
		).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Users); err != nil {
		return nil, ErrFailedQuery
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Users))

	for i := 0; i < len(resp.Users); i++ {
		s.enricher.EnrichJobInfo(resp.Users[i])
	}

	return resp, nil
}
