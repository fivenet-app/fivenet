package jobs

import (
	"context"
	"errors"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	jobs "github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	errorsjobs "github.com/galexrt/fivenet/gen/go/proto/services/jobs/errors"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tQuali = table.FivenetJobsQualifications
)

func (s *Server) ListQualifications(ctx context.Context, req *ListQualificationsRequest) (*ListQualificationsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	_ = userInfo

	condition := tQuali.DeletedAt.IS_NULL()

	// Get total count of values
	countStmt := tQuali.
		SELECT(
			jet.COUNT(tQuali.ID).AS("datacount.totalcount"),
		).
		FROM(
			tQuali,
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, 15)
	resp := &ListQualificationsResponse{
		Pagination:     pag,
		Qualifications: []*jobs.Qualification{},
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := tQuali.
		SELECT(
			tQuali.ID,
			tQuali.CreatedAt,
			tQuali.UpdatedAt,
			tQuali.DeletedAt,
			tQuali.Job,
			tQuali.Weight,
			tQuali.Closed,
			tQuali.Abbreviation,
			tQuali.Title,
			tQuali.Description,
			tQuali.CreatorID,
			tQuali.CreatorJob,
		).
		FROM(tQuali).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Qualifications); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
		}
	}

	resp.Pagination.Update(len(resp.Qualifications))

	return resp, nil
}

func (s *Server) CreateQualification(context.Context, *CreateQualificationRequest) (*CreateQualificationResponse, error) {

	// TODO

	return nil, nil
}

func (s *Server) UpdateQualification(context.Context, *UpdateQualificationRequest) (*UpdateQualificationResponse, error) {

	// TODO

	return nil, nil
}

func (s *Server) DeleteQualification(context.Context, *DeleteQualificationRequest) (*DeleteQualificationResponse, error) {

	// TODO

	return nil, nil
}

// TODO figure out database tables, functions, etc., needed
