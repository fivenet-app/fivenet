package qualifications

import (
	"context"
	"errors"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	"github.com/galexrt/fivenet/gen/go/proto/resources/qualifications"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	errorsqualifications "github.com/galexrt/fivenet/gen/go/proto/services/qualifications/errors"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tQualiResults = table.FivenetQualificationsResults.AS("qualificationresult")
)

func (s *Server) ListQualificationsResults(ctx context.Context, req *ListQualificationsResultsRequest) (*ListQualificationsResultsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tQuali := tQuali.AS("qualificationshort")

	condition := jet.Bool(true)

	if req.QualificationId != nil {
		check, err := s.checkIfUserHasAccessToQuali(ctx, *req.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_GRADE)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
		if !check {
			return nil, errorsqualifications.ErrFailedQuery
		}

		condition = condition.AND(tQualiResults.QualificationID.EQ(jet.Uint64(*req.QualificationId)))
	} else {
		condition = condition.AND(jet.AND(
			tQuali.DeletedAt.IS_NULL(),
			jet.OR(
				tQuali.CreatorID.EQ(jet.Int32(userInfo.UserId)),
				jet.AND(
					tQJobAccess.Access.IS_NOT_NULL(),
					tQJobAccess.Access.GT(jet.Int32(int32(qualifications.AccessLevel_ACCESS_LEVEL_BLOCKED))),
				),
			),
		))

		// TODO
		condition = condition.AND(tQualiResults.UserID.EQ(jet.Int32(userInfo.UserId)))
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
				LEFT_JOIN(tQJobAccess,
					tQJobAccess.QualificationID.EQ(tQuali.ID).
						AND(tQJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tQJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				),
		).
		GROUP_BY(tQualiResults.ID).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, errswrap.NewError(errorsqualifications.ErrFailedQuery, err)
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, QualificationsPageSize)
	resp := &ListQualificationsResultsResponse{
		Pagination: pag,
		Results:    []*qualifications.QualificationResult{},
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
			tQuali.Description,
			tQuali.Content,
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
				).
				LEFT_JOIN(tQJobAccess,
					tQJobAccess.QualificationID.EQ(tQuali.ID).
						AND(tQJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tQJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				),
		).
		GROUP_BY(tQualiResults.ID).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Results); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}
	}

	return resp, nil
}

func (s *Server) CreateOrUpdateQualificationResult(ctx context.Context, req *CreateOrUpdateQualificationResultRequest) (*CreateOrUpdateQualificationResultResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: QualificationsService_ServiceDesc.ServiceName,
		Method:  "CreateOrUpdateQualificationResult",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	ok, err := s.checkIfUserHasAccessToQuali(ctx, req.Result.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_GRADE)
	if err != nil {
		return nil, errswrap.NewError(errorsqualifications.ErrFailedQuery, err)
	}
	if !ok {
		return nil, errorsqualifications.ErrFailedQuery
	}

	tQualiResults := table.FivenetQualificationsResults
	if req.Result.Id <= 0 {
		stmt := tQualiResults.
			INSERT(
				tQualiResults.QualificationID,
				tQualiResults.UserID,
				tQualiResults.Status,
				tQualiResults.Score,
				tQualiResults.Summary,
				tQualiResults.CreatorID,
				tQualiResults.CreatorJob,
			).
			VALUES(
				req.Result.QualificationId,
				req.Result.UserId,
				req.Result.Status,
				req.Result.Score,
				req.Result.Summary,
				userInfo.UserId,
				userInfo.Job,
			)

		res, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		req.Result.Id = uint64(lastId)

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_CREATED)
	} else {
		stmt := tQualiResults.
			UPDATE(
				tQualiResults.QualificationID,
				tQualiResults.UserID,
				tQualiResults.Status,
				tQualiResults.Score,
				tQualiResults.Summary,
			).
			SET(
				req.Result.QualificationId,
				req.Result.UserId,
				req.Result.Status,
				req.Result.Score,
				req.Result.Summary,
			).
			WHERE(tQualiResults.ID.EQ(jet.Uint64(req.Result.Id)))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
		}

		auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)
	}

	result, err := s.getQualificationResult(ctx, req.Result.Id, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsqualifications.ErrFailedQuery)
	}

	return &CreateOrUpdateQualificationResultResponse{
		Result: result,
	}, nil
}

func (s *Server) getQualificationResult(ctx context.Context, resultId uint64, userInfo *userinfo.UserInfo) (*qualifications.QualificationResult, error) {
	stmt := tQualiResults.
		SELECT(
			tQualiResults.ID,
			tQualiResults.CreatedAt,
			tQualiResults.DeletedAt,
			tQualiResults.QualificationID,
			tQualiResults.UserID,
			tUser.ID,
			tUser.Identifier,
			tUser.Job,
			tUser.JobGrade,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tQualiResults.Status,
			tQualiResults.Score,
			tQualiResults.Summary,
			tQualiResults.CreatorID,
			tQualiResults.CreatorJob,
			tCreator.ID,
			tCreator.Identifier,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
		).
		FROM(tQualiResults.
			LEFT_JOIN(tUser,
				tUser.ID.EQ(tQualiResults.UserID),
			).
			LEFT_JOIN(tCreator,
				tCreator.ID.EQ(tQualiResults.CreatorID),
			),
		).
		GROUP_BY(tQualiResults.ID).
		WHERE(tQualiResults.ID.EQ(jet.Uint64(resultId))).
		LIMIT(1)

	var result qualifications.QualificationResult
	if err := stmt.QueryContext(ctx, s.db, &result); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if result.User != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, result.User)
	}

	if result.Creator != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, result.Creator)
	}

	return &result, nil
}

func (s *Server) DeleteQualificationResult(ctx context.Context, req *DeleteQualificationResultRequest) (*DeleteQualificationResultResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: QualificationsService_ServiceDesc.ServiceName,
		Method:  "DeleteQualificationResult",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	re, err := s.getQualificationResult(ctx, req.ResultId, userInfo)
	if err != nil {
		return nil, errswrap.NewError(errorsqualifications.ErrFailedQuery, err)
	}

	ok, err := s.checkIfUserHasAccessToQuali(ctx, re.QualificationId, userInfo, qualifications.AccessLevel_ACCESS_LEVEL_MANAGE)
	if err != nil {
		return nil, errswrap.NewError(errorsqualifications.ErrFailedQuery, err)
	}
	if !ok {
		return nil, errorsqualifications.ErrFailedQuery
	}

	stmt := tQualiResults.
		UPDATE(
			tQualiResults.DeletedAt,
		).
		SET(
			jet.CURRENT_TIMESTAMP(),
		).
		WHERE(
			tQualiResults.ID.EQ(jet.Uint64(re.Id)),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, err
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_DELETED)

	return &DeleteQualificationResultResponse{}, nil
}
