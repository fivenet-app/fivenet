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
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tJobsUserProps    = table.FivenetJobsUserProps
	tJobsUserActivity = table.FivenetJobsUserActivity
)

func (s *Server) ListColleagues(ctx context.Context, req *ListColleaguesRequest) (*ListColleaguesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tUser := tUser.AS("colleague")
	condition := tUser.Job.EQ(jet.String(userInfo.Job)).
		AND(s.customDB.Conditions.User.GetFilter(tUser.Alias()))

	if req.UserId != nil && *req.UserId > 0 {
		condition = condition.AND(tUser.ID.EQ(jet.Int32(*req.UserId)))
	} else {
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
	}

	if req.Absent != nil && *req.Absent {
		condition = condition.AND(jet.AND(
			tJobsUserProps.AbsenceDate.IS_NOT_NULL(),
			tJobsUserProps.AbsenceDate.GT_EQ(jet.CURRENT_DATE()),
		))
	}

	// Get total count of values
	countStmt := tUser.
		SELECT(
			jet.COUNT(tUser.ID).AS("datacount.totalcount"),
		).
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname_fulltext")).
		FROM(
			tUser.
				LEFT_JOIN(tJobsUserProps,
					tJobsUserProps.UserID.EQ(tUser.ID),
				),
		).
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
	tUser := tUser.AS("colleague")
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

	s.enricher.EnrichJobInfo(dest)

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

	// Field Permission Check
	fieldsAttr, err := s.ps.Attr(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceGetColleaguePerm, permsjobs.JobsServiceGetColleagueAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}

	targetUser, err := s.getColleague(ctx, req.UserId)
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

	colleague, err := s.getColleague(ctx, targetUser.UserId)
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
	tJobsUserProps := tJobsUserProps.AS("jobsuserprops")
	stmt := tJobsUserProps.
		SELECT(
			tJobsUserProps.UserID,
			tJobsUserProps.AbsenceDate,
		).
		FROM(tJobsUserProps).
		WHERE(tJobsUserProps.UserID.EQ(jet.Int32(userId))).
		LIMIT(1)

	dest := &jobs.JobsUserProps{
		UserId: userId,
	}
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

	absenceDate := jet.DateExp(jet.NULL)
	if req.Props.AbsenceDate != nil {
		if req.Props.AbsenceDate.Timestamp == nil {
			req.Props.AbsenceDate = nil
		} else {
			req.Props.AbsenceDate = timestamp.New(req.Props.AbsenceDate.AsTime())
			absenceDate = jet.DateT(req.Props.AbsenceDate.AsTime())
		}
	} else {
		req.Props.AbsenceDate = props.AbsenceDate
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

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

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	// Compare absence date if any was set
	if req.Props.AbsenceDate != nil && (props.AbsenceDate == nil || req.Props.AbsenceDate.AsTime().Compare(props.AbsenceDate.AsTime()) != 0) {
		if err := s.addJobsUserActivity(ctx, tx, &jobs.JobsUserActivity{
			Job:          userInfo.Job,
			SourceUserId: userInfo.UserId,
			TargetUserId: req.Props.UserId,
			ActivityType: jobs.JobsUserActivityType_JOBS_USER_ACTIVITY_TYPE_ABSENCE_DATE,
			Reason:       req.Reason,
			Data: &jobs.JobsUserActivityData{
				Data: &jobs.JobsUserActivityData_AbsenceDate{
					AbsenceDate: &jobs.ColleagueAbsenceDate{
						AbsenceDate: timestamp.New(req.Props.AbsenceDate.AsTime()),
					},
				},
			},
		}); err != nil {
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
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

func (s *Server) addJobsUserActivity(ctx context.Context, tx qrm.DB, activity *jobs.JobsUserActivity) error {
	stmt := tJobsUserActivity.
		INSERT(
			tJobsUserActivity.Job,
			tJobsUserActivity.SourceUserID,
			tJobsUserActivity.TargetUserID,
			tJobsUserActivity.ActivityType,
			tJobsUserActivity.Reason,
			tJobsUserActivity.Data,
		).
		VALUES(
			activity.Job,
			activity.SourceUserId,
			activity.TargetUserId,
			activity.ActivityType,
			activity.Reason,
			activity.Data,
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}

func (s *Server) checkIfHasAccessToColleague(levels []string, userInfo *userinfo.UserInfo, target *users.UserShort) bool {
	if userInfo.SuperUser {
		return true
	}

	// Not the same job, can't do anything
	if userInfo.Job != target.Job {
		return false
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

func (s *Server) ListColleagueActivity(ctx context.Context, req *ListColleagueActivityRequest) (*ListColleagueActivityResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Field Permission Check
	fieldsAttr, err := s.ps.Attr(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceSetJobsUserPropsPerm, permsjobs.JobsServiceSetJobsUserPropsAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}

	targetUser, err := s.getColleague(ctx, req.UserId)
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

	tJobsUserActivity := tJobsUserActivity.AS("jobsuseractivity")
	condition := tJobsUserActivity.Job.EQ(jet.String(userInfo.Job)).
		AND(tJobsUserActivity.TargetUserID.EQ(jet.Int32(req.UserId)))

	// Get total count of values
	countStmt := tJobsUserActivity.
		SELECT(
			jet.COUNT(tJobsUserActivity.ID).AS("datacount.totalcount"),
		).
		FROM(tJobsUserActivity).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(15)
	resp := &ListColleagueActivityResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	tUTarget := tUser.AS("target_user")
	tUSource := tUser.AS("source_user")
	stmt := tJobsUserActivity.
		SELECT(
			tJobsUserActivity.ID,
			tJobsUserActivity.CreatedAt,
			tJobsUserActivity.Job,
			tJobsUserActivity.SourceUserID,
			tJobsUserActivity.TargetUserID,
			tJobsUserActivity.ActivityType,
			tJobsUserActivity.Reason,
			tJobsUserActivity.Data,
			tUTarget.ID,
			tUTarget.Identifier,
			tUTarget.Job,
			tUTarget.JobGrade,
			tUTarget.Firstname,
			tUTarget.Lastname,
			tUSource.ID,
			tUSource.Identifier,
			tUSource.Job,
			tUSource.JobGrade,
			tUSource.Firstname,
			tUSource.Lastname,
		).
		FROM(
			tJobsUserActivity.
				LEFT_JOIN(tUTarget,
					tUTarget.ID.EQ(tJobsUserActivity.TargetUserID),
				).
				LEFT_JOIN(tUSource,
					tUSource.ID.EQ(tJobsUserActivity.SourceUserID),
				),
		).
		WHERE(condition).
		OFFSET(pag.Offset).
		ORDER_BY(tJobsUserActivity.CreatedAt.DESC(), tJobsUserActivity.ID.DESC()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Activity); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	pag.Update(count.TotalCount, len(resp.Activity))

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Activity); i++ {
		if resp.Activity[i].SourceUser != nil {
			jobInfoFn(resp.Activity[i].SourceUser)
		}
		if resp.Activity[i].TargetUser != nil {
			jobInfoFn(resp.Activity[i].TargetUser)
		}
	}

	return resp, nil
}
