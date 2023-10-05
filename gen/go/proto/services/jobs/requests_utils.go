package jobs

import (
	"context"
	"errors"

	jobs "github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) getRequest(ctx context.Context, job string, id uint64) (*jobs.Request, error) {
	tCreator := tUser.AS("creator")
	tApprover := tUser.AS("approver")
	stmt := tRequests.
		SELECT(
			tRequests.ID,
			tRequests.CreatedAt,
			tRequests.UpdatedAt,
			tRequests.DeletedAt,
			tRequests.Job,
			tRequests.TypeID,
			tReqTypes.ID,
			tReqTypes.Name,
			tReqTypes.Description,
			tRequests.Title,
			tRequests.Message,
			tRequests.Status,
			tRequests.CreatorID,
			tRequests.Approved,
			tRequests.ApproverID,
			tRequests.Closed,
			tRequests.BeginsAt,
			tRequests.EndsAt,
			tCreator.ID,
			tCreator.Identifier,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tApprover.ID,
			tApprover.Identifier,
			tApprover.Job,
			tApprover.JobGrade,
			tApprover.Firstname,
			tApprover.Lastname,
		).
		FROM(
			tRequests.
				LEFT_JOIN(tReqTypes,
					tReqTypes.ID.EQ(tRequests.TypeID),
				).
				INNER_JOIN(tCreator,
					tCreator.ID.EQ(tRequests.CreatorID),
				).
				LEFT_JOIN(tApprover,
					tApprover.ID.EQ(tRequests.ApproverID),
				),
		).
		WHERE(jet.AND(
			tRequests.Job.EQ(jet.String(job)),
			tRequests.ID.EQ(jet.Uint64(id)),
			tRequests.DeletedAt.IS_NULL(),
		)).
		LIMIT(1)

	var dest jobs.Request
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, ErrFailedQuery
		}
	}

	return nil, nil
}

func (s *Server) getRequestType(ctx context.Context, job string, id uint64) (*jobs.RequestType, error) {
	stmt := tReqTypes.
		SELECT(
			tReqTypes.ID,
			tReqTypes.CreatedAt,
			tReqTypes.UpdatedAt,
			tReqTypes.DeletedAt,
			tReqTypes.Job,
			tReqTypes.Name,
			tReqTypes.Description,
			tReqTypes.Weight,
		).
		FROM(
			tReqTypes,
		).
		WHERE(jet.AND(
			tReqTypes.Job.EQ(jet.String(job)),
			tReqTypes.ID.EQ(jet.Uint64(id)),
			tReqTypes.DeletedAt.IS_NULL(),
		)).
		LIMIT(1)

	var dest jobs.RequestType
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, ErrFailedQuery
		}
	}

	return &dest, nil
}

func (s *Server) getRequestComment(ctx context.Context, job string, id uint64) (*jobs.RequestComment, error) {
	tCreator := tUser.AS("creator")
	stmt := tReqComments.
		SELECT(
			tReqComments.ID,
			tReqComments.CreatedAt,
			tReqComments.UpdatedAt,
			tReqComments.DeletedAt,
			tReqComments.RequestID,
			tReqComments.Comment,
			tReqComments.CreatorID,
			tCreator.ID,
			tCreator.Identifier,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
		).
		FROM(
			tReqComments.
				INNER_JOIN(tRequests,
					tRequests.ID.EQ(tReqComments.RequestID),
				),
		).
		WHERE(jet.AND(
			tReqComments.ID.EQ(jet.Uint64(id)),
			tReqComments.RequestID.EQ(tRequests.ID),
			tRequests.Job.EQ(jet.String(job)),
			tReqComments.DeletedAt.IS_NULL(),
		)).
		LIMIT(1)

	var dest jobs.RequestType
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(qrm.ErrNoRows, err) {
			return nil, ErrFailedQuery
		}
	}

	return nil, nil
}
