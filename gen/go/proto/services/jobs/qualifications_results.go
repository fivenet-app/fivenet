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
	tQualiResults = table.FivenetJobsQualificationsResults.AS("qualificationresult")
)

func (s *Server) ListQualificationsResults(ctx context.Context, req *ListQualificationsResultsRequest) (*ListQualificationsResultsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tQualiResults.UserID.EQ(jet.Int32(userInfo.UserId))

	if req.QualificationId != nil {
		check, err := s.checkIfUserHasAccessToQuali(ctx, *req.QualificationId, userInfo, jobs.AccessLevel_ACCESS_LEVEL_GRADE)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		if !check {
			return nil, errorsjobs.ErrFailedQuery
		}

		condition = condition.AND(tQualiResults.QualificationID.EQ(jet.Uint64(*req.QualificationId)))
	}
	if len(req.Status) > 0 {
		statuses := []jet.Expression{}
		for i := 0; i < len(req.Status); i++ {
			statuses = append(statuses, jet.Int16(int16(req.Status[i])))
		}

		condition = condition.AND(tQualiResults.Status.IN(statuses...))
	}

	countStmt := tQualiResults.
		SELECT(
			jet.COUNT(tQualiResults.ID).AS("datacount.totalcount"),
		).
		FROM(
			tQualiResults.
				INNER_JOIN(tQuali,
					tQuali.ID.EQ(tQualiResults.QualificationID),
				).
				LEFT_JOIN(tCreator,
					tQuali.CreatorID.EQ(tCreator.ID),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, 15)
	resp := &ListQualificationsResultsResponse{
		Pagination: pag,
		Results:    []*jobs.Qualification{},
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := tQualiResults.
		SELECT(
			tQualiResults.ID,
			tQualiResults.CreatedAt,
			tQualiResults.QualificationID,
			tQualiResults.UserID,
			tQualiResults.Status,
			tQualiResults.Score,
			tQualiResults.Summary,
			tQuali.ID,
			tQuali.CreatedAt,
			tQuali.UpdatedAt,
			tQuali.Job,
			tQuali.Closed,
			tQuali.Abbreviation,
			tQuali.Title,
			tQuali.Summary,
			tQuali.Description,
			tQuali.CreatorJob,
			tQuali.CreatorID,
			tCreator.ID,
			tCreator.Identifier,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
		).
		FROM(
			tQualiResults.
				INNER_JOIN(tQuali,
					tQuali.ID.EQ(tQualiResults.QualificationID),
				).
				LEFT_JOIN(tCreator,
					tQuali.CreatorID.EQ(tCreator.ID),
				),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Results); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	return resp, nil
}
