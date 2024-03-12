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
	tQualiRequests = table.FivenetJobsQualificationsRequests.AS("qualificationrequest")
)

func (s *Server) ListQualificationRequests(ctx context.Context, req *ListQualificationRequestsRequest) (*ListQualificationRequestsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tQualiRequests.UserID.EQ(jet.Int32(userInfo.UserId))

	if req.QualificationId != nil {
		condition = condition.AND(tQualiRequests.QualificationID.EQ(jet.Uint64(*req.QualificationId)))
	}

	countStmt := tQualiRequests.
		SELECT(
			jet.COUNT(tQualiRequests.ID).AS("datacount.totalcount"),
		).
		FROM(tQualiRequests).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, 15)
	resp := &ListQualificationRequestsResponse{
		Pagination: pag,
		Requests:   []*jobs.QualificationRequest{},
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := tQualiRequests.
		SELECT(
			tQualiRequests.ID,
			tQualiRequests.CreatedAt,
			tQualiRequests.QualificationID,
			tQualiRequests.UserID,
			tQualiRequests.UserComment,
			tQualiRequests.Approved,
			tQualiRequests.ApprovedAt,
			tQualiRequests.ApproverComment,
			tQualiRequests.ApproverID,
			tQualiRequests.ApproverJob,
			tCreator.ID,
			tCreator.Identifier,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
		).
		FROM(
			tQualiRequests.
				LEFT_JOIN(tCreator,
					tQuali.CreatorID.EQ(tCreator.ID),
				),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Requests); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	return resp, nil
}
