package jobs

import (
	"context"
	"errors"
	"slices"
	"strings"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	jobs "github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/gen/go/proto/resources/rector"
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
			tJobsUserProps.AbsenceBegin.IS_NOT_NULL(),
			tJobsUserProps.AbsenceBegin.GT_EQ(jet.CURRENT_DATE()),
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

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, 15)
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
			tJobsUserProps.AbsenceBegin,
			tJobsUserProps.AbsenceEnd,
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
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
		}
	}

	resp.Pagination.Update(len(resp.Colleagues))

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
			tJobsUserProps.AbsenceBegin,
			tJobsUserProps.AbsenceEnd,
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
			tJobsUserProps.AbsenceBegin,
			tJobsUserProps.AbsenceEnd,
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

	absenceBegin := jet.DateExp(jet.NULL)
	absenceEnd := jet.DateExp(jet.NULL)
	if req.Props.AbsenceBegin != nil && req.Props.AbsenceEnd != nil {
		if req.Props.AbsenceBegin.Timestamp == nil {
			req.Props.AbsenceBegin = nil
		} else {
			absenceBegin = jet.DateT(req.Props.AbsenceBegin.AsTime())
		}

		if req.Props.AbsenceEnd.Timestamp == nil {
			req.Props.AbsenceEnd = nil
		} else {
			absenceEnd = jet.DateT(req.Props.AbsenceEnd.AsTime())
		}
	} else {
		req.Props.AbsenceBegin = props.AbsenceBegin
		req.Props.AbsenceEnd = props.AbsenceEnd
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
			tJobsUserProps.AbsenceBegin,
			tJobsUserProps.AbsenceEnd,
		).
		VALUES(
			req.Props.UserId,
			absenceBegin,
			absenceEnd,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tJobsUserProps.AbsenceBegin.SET(jet.DateExp(jet.Raw("VALUES(`absence_begin`)"))),
			tJobsUserProps.AbsenceEnd.SET(jet.DateExp(jet.Raw("VALUES(`absence_end`)"))),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	// Compare absence dates if any were set
	if req.Props.AbsenceBegin != nil && req.Props.AbsenceEnd != nil &&
		(props.AbsenceBegin == nil ||
			props.AbsenceEnd == nil ||
			req.Props.AbsenceBegin.AsTime().Compare(props.AbsenceBegin.AsTime()) != 0 ||
			req.Props.AbsenceEnd.AsTime().Compare(props.AbsenceEnd.AsTime()) != 0) {
		if err := s.addJobsUserActivity(ctx, tx, &jobs.JobsUserActivity{
			Job:          userInfo.Job,
			SourceUserId: userInfo.UserId,
			TargetUserId: req.Props.UserId,
			ActivityType: jobs.JobsUserActivityType_JOBS_USER_ACTIVITY_TYPE_ABSENCE_DATE,
			Reason:       req.Reason,
			Data: &jobs.JobsUserActivityData{
				Data: &jobs.JobsUserActivityData_AbsenceDate{
					AbsenceDate: &jobs.ColleagueAbsenceDate{
						AbsenceBegin: req.Props.AbsenceBegin,
						AbsenceEnd:   req.Props.AbsenceEnd,
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

func (s *Server) getConditionForColleagueAccess(actTable *table.FivenetJobsUserActivityTable, usersTable *table.UsersTable, levels []string, userInfo *userinfo.UserInfo) jet.BoolExpression {
	condition := jet.Bool(true)
	if userInfo.SuperUser {
		return condition
	}

	// If no levels set, assume "Own" as default
	if len(levels) == 0 {
		return actTable.TargetUserID.EQ(jet.Int32(userInfo.UserId))
	}

	if slices.Contains(levels, "Any") {
		return condition
	}
	if slices.Contains(levels, "Lower_Rank") {
		return usersTable.ID.LT(jet.Int32(userInfo.JobGrade))
	}
	if slices.Contains(levels, "Same_Rank") {
		return usersTable.ID.LT_EQ(jet.Int32(userInfo.JobGrade))
	}
	if slices.Contains(levels, "Own") {
		return usersTable.ID.EQ(jet.Int32(userInfo.UserId))
	}

	return jet.Bool(false)
}

func (s *Server) ListColleagueActivity(ctx context.Context, req *ListColleagueActivityRequest) (*ListColleagueActivityResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Access Field Permission Check
	accessAttr, err := s.ps.Attr(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceGetColleaguePerm, permsjobs.JobsServiceGetColleagueAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	var access perms.StringList
	if accessAttr != nil {
		access = accessAttr.([]string)
	}

	tJobsUserActivity := tJobsUserActivity.AS("jobsuseractivity")
	tUTarget := tUser.AS("target_user")
	tUSource := tUser.AS("source_user")

	condition := tJobsUserActivity.Job.EQ(jet.String(userInfo.Job))

	// If no user IDs given or more than 2, show all the user has access to
	if len(req.UserIds) == 0 || len(req.UserIds) >= 2 {
		condition = condition.AND(s.getConditionForColleagueAccess(tJobsUserActivity, tUTarget, access, userInfo))

		if len(req.UserIds) >= 2 {
			// More than 2 user ids
			userIds := make([]jet.Expression, len(req.UserIds))
			for i := 0; i < len(req.UserIds); i++ {
				userIds[i] = jet.Int32(req.UserIds[i])
			}

			condition = condition.AND(tUTarget.ID.IN(userIds...))
		}
	} else {
		userId := req.UserIds[0]

		targetUser, err := s.getColleague(ctx, userId)
		if err != nil {
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
		}

		if !s.checkIfHasAccessToColleague(access, userInfo, &users.UserShort{
			UserId:   targetUser.UserId,
			Job:      targetUser.Job,
			JobGrade: targetUser.JobGrade,
		}) {
			return nil, errorsjobs.ErrFailedQuery
		}

		condition = condition.AND(tJobsUserActivity.TargetUserID.EQ(jet.Int32(userId)))
	}

	// Types Field Permission Check
	typesAttr, err := s.ps.Attr(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceListColleaguesPerm, permsjobs.JobsServiceListColleagueActivityTypesPermField)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	var types perms.StringList
	if typesAttr != nil {
		types = typesAttr.([]string)
	}
	if len(types) == 0 {
		if !userInfo.SuperUser {
			return &ListColleagueActivityResponse{}, nil
		}
		types = append(types, "HIRED", "FIRED", "PROMOTED", "DEMOTED", "ABSENCE_DATE")
	}

	condTypes := []jet.Expression{}
	for _, aType := range types {
		switch strings.ToUpper(aType) {
		case "HIRED":
			condTypes = append(condTypes, jet.Int32(int32(jobs.JobsUserActivityType_JOBS_USER_ACTIVITY_TYPE_HIRED)))
		case "FIRED":
			condTypes = append(condTypes, jet.Int32(int32(jobs.JobsUserActivityType_JOBS_USER_ACTIVITY_TYPE_FIRED)))
		case "PROMOTED":
			condTypes = append(condTypes, jet.Int32(int32(jobs.JobsUserActivityType_JOBS_USER_ACTIVITY_TYPE_PROMOTED)))
		case "DEMOTED":
			condTypes = append(condTypes, jet.Int32(int32(jobs.JobsUserActivityType_JOBS_USER_ACTIVITY_TYPE_DEMOTED)))
		case "ABSENCE_DATE":
			condTypes = append(condTypes, jet.Int32(int32(jobs.JobsUserActivityType_JOBS_USER_ACTIVITY_TYPE_ABSENCE_DATE)))
		}
	}

	if len(condTypes) == 0 {
		return &ListColleagueActivityResponse{}, nil
	}

	condition = condition.AND(tJobsUserActivity.ActivityType.IN(condTypes...))

	// Get total count of values
	countStmt := tJobsUserActivity.
		SELECT(
			jet.COUNT(tJobsUserActivity.ID).AS("datacount.totalcount"),
		).
		FROM(
			tJobsUserActivity.
				INNER_JOIN(tUTarget,
					tUTarget.ID.EQ(tJobsUserActivity.TargetUserID),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, 15)
	resp := &ListColleagueActivityResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

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
				INNER_JOIN(tUTarget,
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

	pag.Update(len(resp.Activity))

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
