package jobs

import (
	"context"
	"errors"
	"strings"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	jobs "github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	errorsjobs "github.com/galexrt/fivenet/gen/go/proto/services/jobs/errors"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	timeutils "github.com/galexrt/fivenet/pkg/utils/time"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tJobsUserProps = table.FivenetJobsUserProps
)

func (s *Server) ListColleagues(ctx context.Context, req *ListColleaguesRequest) (*ListColleaguesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tUser := tUser.AS("colleague")
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
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname_fulltext")).
		FROM(tUser).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(15)
	resp := &ListColleaguesResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	selectors := jet.ProjectionList{
		tUser.ID,
		tUser.Identifier,
		tUser.Firstname,
		tUser.Lastname,
		tUser.Job,
		tUser.JobGrade,
		tUser.Dateofbirth,
		tUser.PhoneNumber,
		tUserProps.Avatar.AS("colleague.avatar"),
		tJobsUserProps.UserID,
		tJobsUserProps.AbsenceDate,
	}

	stmt := tUser.
		SELECT(
			selectors[0], selectors[1:]...,
		).
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname_fulltext")).
		FROM(
			tUser.
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUser.ID),
				).
				LEFT_JOIN(tJobsUserProps,
					tJobsUserProps.UserID.EQ(tUser.ID),
				),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(
			tUser.JobGrade.ASC(),
			tUser.Firstname.ASC(),
			tUser.Lastname.ASC(),
		).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Colleagues); err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Colleagues))

	for i := 0; i < len(resp.Colleagues); i++ {
		s.enricher.EnrichJobInfo(resp.Colleagues[i])
	}

	return resp, nil
}

func (s *Server) GetColleague(ctx context.Context, req *GetColleagueRequest) (*GetColleagueResponse, error) {
	resp := &GetColleagueResponse{}

	// TODO

	return resp, nil
}

func (s *Server) getJobsUserProps(ctx context.Context, userId int32) (*jobs.JobsUserProps, error) {
	stmt := tJobsUserProps.
		SELECT(
			tJobsUserProps.UserID,
			tJobsUserProps.AbsenceDate,
		).
		FROM(tJobsUserProps).
		WHERE(tJobsUserProps.UserID.EQ(jet.Int32(userId))).
		LIMIT(1)

	dest := &jobs.JobsUserProps{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
		}
	}

	return dest, nil
}

func (s *Server) SetJobsUserProps(ctx context.Context, req *SetJobsUserPropsRequest) (*SetJobsUserPropsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsService_ServiceDesc.ServiceName,
		Method:  "SetJobsUserProps",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	// TODO add option to set the user props for others
	req.Props.UserId = userInfo.UserId

	props, err := s.getJobsUserProps(ctx, req.Props.UserId)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	if req.Props.AbsenceDate != nil {
		req.Props.AbsenceDate = timestamp.New(timeutils.TruncateToDay(req.Props.AbsenceDate.AsTime()))
	} else {
		req.Props.AbsenceDate = props.AbsenceDate
	}

	// TODO add fivenet_jobs_user_activity entries

	stmt := tJobProps.
		INSERT(
			tJobsUserProps.UserID,
			tJobsUserProps.AbsenceDate,
		).
		VALUES(
			req.Props.UserId,
			req.Props.AbsenceDate,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tJobsUserProps.AbsenceDate.SET(jet.DateTimeT(req.Props.AbsenceDate.AsTime())),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	props, err = s.getJobsUserProps(ctx, req.Props.UserId)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	return &SetJobsUserPropsResponse{
		Props: props,
	}, nil
}
