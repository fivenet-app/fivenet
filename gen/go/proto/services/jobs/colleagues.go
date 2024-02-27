package jobs

import (
	"context"
	"errors"
	"slices"
	"strings"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	jobs "github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	"github.com/galexrt/fivenet/gen/go/proto/resources/users"
	errorsjobs "github.com/galexrt/fivenet/gen/go/proto/services/jobs/errors"
	permsjobs "github.com/galexrt/fivenet/gen/go/proto/services/jobs/perms"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/pkg/perms"
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

	stmt := tUser.
		SELECT(
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

func (s *Server) getColleague(ctx context.Context, userId int32) (*jobs.Colleague, error) {
	stmt := tUser.
		SELECT(
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
		).
		FROM(
			tUser.
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUser.ID),
				).
				LEFT_JOIN(tJobsUserProps,
					tJobsUserProps.UserID.EQ(tUser.ID),
				),
		).
		WHERE(
			tUser.ID.EQ(jet.Int32(userId)),
		).
		LIMIT(1)

	dest := &jobs.Colleague{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (s *Server) GetColleague(ctx context.Context, req *GetColleagueRequest) (*GetColleagueResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsService_ServiceDesc.ServiceName,
		Method:  "GetColleague",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	colleague, err := s.getColleague(ctx, req.UserId)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)

	return &GetColleagueResponse{
		Colleague: colleague,
	}, nil
}

func (s *Server) GetSelf(ctx context.Context, req *GetSelfRequest) (*GetSelfResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	colleague, err := s.getColleague(ctx, userInfo.UserId)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	return &GetSelfResponse{
		Colleague: colleague,
	}, nil
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

	// Field Permission Check
	fieldsAttr, err := s.ps.Attr(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceSetJobsUserPropsPerm, permsjobs.JobsServiceSetJobsUserPropsAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}

	targetUser, err := s.getColleague(ctx, req.Props.UserId)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	if !s.checkIfHasAccessToColleague(fields, userInfo, &users.UserShort{
		UserId:   targetUser.UserId,
		Job:      targetUser.Job,
		JobGrade: targetUser.JobGrade,
	}) {
		return nil, errorsjobs.ErrFailedQuery
	}

	props, err := s.getJobsUserProps(ctx, req.Props.UserId)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	absenceDate := jet.TimestampExp(jet.NULL)
	if req.Props.AbsenceDate != nil {
		if req.Props.AbsenceDate.Timestamp == nil {
			req.Props.AbsenceDate = nil
		} else {
			req.Props.AbsenceDate = timestamp.New(timeutils.TruncateToDay(req.Props.AbsenceDate.AsTime()))
			absenceDate = jet.DateTimeT(req.Props.AbsenceDate.AsTime())
		}
	} else {
		req.Props.AbsenceDate = props.AbsenceDate
	}

	// TODO add fivenet_jobs_user_activity entries

	stmt := tJobsUserProps.
		INSERT(
			tJobsUserProps.UserID,
			tJobsUserProps.AbsenceDate,
		).
		VALUES(
			req.Props.UserId,
			req.Props.AbsenceDate,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tJobsUserProps.AbsenceDate.SET(absenceDate),
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

func (s *Server) checkIfHasAccessToColleague(levels []string, userInfo *userinfo.UserInfo, target *users.UserShort) bool {
	if userInfo.SuperUser {
		return true
	}

	// If the creator is nil, treat it like a normal doc access check
	if target == nil {
		return true
	}

	// If no levels set, assume "Own" as default
	if len(levels) == 0 {
		return target.UserId == userInfo.UserId
	}

	if slices.Contains(levels, "Any") {
		return true
	}
	if slices.Contains(levels, "Lower_Rank") {
		if target.JobGrade < userInfo.JobGrade {
			return true
		}
	}
	if slices.Contains(levels, "Same_Rank") {
		if target.JobGrade <= userInfo.JobGrade {
			return true
		}
	}
	if slices.Contains(levels, "Own") {
		if target.UserId == userInfo.UserId {
			return true
		}
	}

	return false
}
