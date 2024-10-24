package jobs

import (
	"context"
	"errors"
	"slices"
	"strings"

	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	jobs "github.com/fivenet-app/fivenet/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/rector"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	errorsjobs "github.com/fivenet-app/fivenet/gen/go/proto/services/jobs/errors"
	permsjobs "github.com/fivenet-app/fivenet/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/utils"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
		req.Search = strings.TrimSpace(req.Search)
		req.Search = strings.ReplaceAll(req.Search, "%", "")
		req.Search = strings.ReplaceAll(req.Search, " ", "%")
		if req.Search != "" {
			req.Search = "%" + req.Search + "%"
			condition = condition.AND(
				jet.CONCAT(tUser.Firstname, jet.String(" "), tUser.Lastname).
					LIKE(jet.String(req.Search)),
			)
		}
	}

	if req.Absent != nil && *req.Absent {
		condition = condition.AND(
			jet.AND(
				tJobsUserProps.AbsenceBegin.IS_NOT_NULL(),
				tJobsUserProps.AbsenceEnd.IS_NOT_NULL(),
				tJobsUserProps.AbsenceBegin.LT_EQ(jet.CURRENT_DATE()),
				tJobsUserProps.AbsenceEnd.GT_EQ(jet.CURRENT_DATE()),
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
					tJobsUserProps.UserID.EQ(tUser.ID).
						AND(tUser.Job.EQ(jet.String(userInfo.Job))),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, 16)
	resp := &ListColleaguesResponse{
		Pagination: pag,
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []jet.OrderByClause{}
	if req.Sort != nil {
		var columns []jet.Column
		switch req.Sort.Column {
		case "name":
			columns = append(columns, tUser.Firstname, tUser.Lastname)
		case "rank":
			fallthrough
		default:
			columns = append(columns, tUser.JobGrade)
		}

		for _, column := range columns {
			if req.Sort.Direction == database.AscSortDirection {
				orderBys = append(orderBys, column.ASC())
			} else {
				orderBys = append(orderBys, column.DESC())
			}
		}
	} else {
		orderBys = append(orderBys,
			tUser.JobGrade.ASC(),
			tUser.Firstname.ASC(),
			tUser.Lastname.ASC(),
		)
	}

	stmt := tUser.
		SELECT(
			tUser.ID,
			tUser.Job,
			tUser.JobGrade,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.PhoneNumber,
			tUserProps.Avatar.AS("colleague.avatar"),
			tJobsUserProps.UserID,
			tJobsUserProps.Job,
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
					tJobsUserProps.UserID.EQ(tUser.ID).
						AND(tJobsUserProps.Job.EQ(jet.String(userInfo.Job))),
				),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(orderBys...).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Colleagues); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	resp.Pagination.Update(len(resp.Colleagues))

	for i := 0; i < len(resp.Colleagues); i++ {
		s.enricher.EnrichJobInfo(resp.Colleagues[i])
	}

	return resp, nil
}

func (s *Server) getColleague(ctx context.Context, job string, userId int32, withColumns []jet.Projection) (*jobs.Colleague, error) {
	tUser := tUser.AS("colleague")
	columns := []jet.Projection{
		tUser.Firstname,
		tUser.Lastname,
		tUser.Job,
		tUser.JobGrade,
		tUser.Dateofbirth,
		tUser.PhoneNumber,
		tUserProps.Avatar.AS("colleague.avatar"),
		tJobsUserProps.UserID,
		tJobsUserProps.Job,
		tJobsUserProps.AbsenceBegin,
		tJobsUserProps.AbsenceEnd,
	}
	columns = append(columns, withColumns...)

	stmt := tUser.
		SELECT(
			tUser.ID,
			columns...,
		).
		FROM(
			tUser.
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUser.ID),
				).
				LEFT_JOIN(tJobsUserProps,
					tJobsUserProps.UserID.EQ(tUser.ID).
						AND(tJobsUserProps.Job.EQ(jet.String(job))),
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
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.jobs.colleague.id", int64(req.UserId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsService_ServiceDesc.ServiceName,
		Method:  "GetColleague",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	// Access Permission Check
	accessAttr, err := s.ps.Attr(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceGetColleaguePerm, permsjobs.JobsServiceGetColleagueAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	var access perms.StringList
	if accessAttr != nil {
		access = accessAttr.([]string)
	}

	targetUser, err := s.getColleague(ctx, userInfo.Job, req.UserId, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	infoOnly := req.InfoOnly != nil && *req.InfoOnly

	check := s.checkIfHasAccessToColleague(access, userInfo, &users.UserShort{
		UserId:   targetUser.UserId,
		Job:      targetUser.Job,
		JobGrade: targetUser.JobGrade,
	})
	if !check {
		if userInfo.Job != targetUser.Job || !infoOnly {
			return nil, errorsjobs.ErrFailedQuery
		}
	}

	// Field Permission Check
	var types perms.StringList
	if !infoOnly {
		typesAttr, err := s.ps.Attr(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceGetColleaguePerm, permsjobs.JobsServiceGetColleagueTypesPermField)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		if typesAttr != nil {
			types = typesAttr.([]string)
		}
	}
	if userInfo.SuperUser {
		types = []string{"Note"}
	}

	withColumns := []jet.Projection{}
	for _, fType := range types {
		switch fType {
		case "Note":
			withColumns = append(withColumns, tJobsUserProps.Note)
		}
	}

	colleague, err := s.getColleague(ctx, userInfo.Job, targetUser.UserId, withColumns)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)

	return &GetColleagueResponse{
		Colleague: colleague,
	}, nil
}

func (s *Server) GetSelf(ctx context.Context, req *GetSelfRequest) (*GetSelfResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Field Permission Check
	typesAttr, err := s.ps.Attr(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceGetColleaguePerm, permsjobs.JobsServiceGetColleagueTypesPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	var types perms.StringList
	if typesAttr != nil {
		types = typesAttr.([]string)
	}
	if len(types) == 0 && userInfo.SuperUser {
		types = append(types, "AbsenceDate", "Note")
	}

	withColumns := []jet.Projection{}
	for _, fType := range types {
		switch fType {
		case "Note":
			withColumns = append(withColumns, tJobsUserProps.Note)
		}
	}

	colleague, err := s.getColleague(ctx, userInfo.Job, userInfo.UserId, withColumns)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	return &GetSelfResponse{
		Colleague: colleague,
	}, nil
}

func (s *Server) getJobsUserProps(ctx context.Context, userId int32, job string, fields []string) (*jobs.JobsUserProps, error) {
	tJobsUserProps := table.FivenetJobsUserProps.AS("jobsuserprops")
	columns := []jet.Projection{
		tJobsUserProps.Job,
		tJobsUserProps.AbsenceBegin,
		tJobsUserProps.AbsenceEnd,
	}

	for _, field := range fields {
		switch field {
		case "Note":
			columns = append(columns, tJobsUserProps.Note)
		}
	}

	stmt := tJobsUserProps.
		SELECT(
			tJobsUserProps.UserID,
			columns...,
		).
		FROM(tJobsUserProps).
		WHERE(jet.AND(
			tJobsUserProps.UserID.EQ(jet.Int32(userId)),
			tJobsUserProps.Job.EQ(jet.String(job)),
		)).
		LIMIT(1)

	dest := &jobs.JobsUserProps{
		UserId: userId,
	}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	return dest, nil
}

func (s *Server) SetJobsUserProps(ctx context.Context, req *SetJobsUserPropsRequest) (*SetJobsUserPropsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.jobs.colleague.id", int64(req.Props.UserId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: JobsService_ServiceDesc.ServiceName,
		Method:  "SetJobsUserProps",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	if req.Reason == "" {
		return nil, errorsjobs.ErrReasonRequired
	}

	// Access Permission Check
	accessAttr, err := s.ps.Attr(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceSetJobsUserPropsPerm, permsjobs.JobsServiceSetJobsUserPropsAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	var access perms.StringList
	if accessAttr != nil {
		access = accessAttr.([]string)
	}

	targetUser, err := s.getColleague(ctx, userInfo.Job, req.Props.UserId, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	if !s.checkIfHasAccessToColleague(access, userInfo, &users.UserShort{
		UserId:   targetUser.UserId,
		Job:      targetUser.Job,
		JobGrade: targetUser.JobGrade,
	}) {
		return nil, errorsjobs.ErrFailedQuery
	}

	// Types Permission Check
	typesAttr, err := s.ps.Attr(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceSetJobsUserPropsPerm, permsjobs.JobsServiceSetJobsUserPropsTypesPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	var types perms.StringList
	if typesAttr != nil {
		types = typesAttr.([]string)
	}
	if userInfo.SuperUser {
		types = []string{"AbsenceDate", "Note"}
	}

	props, err := s.getJobsUserProps(ctx, targetUser.UserId, targetUser.Job, types)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	absenceBegin := jet.DateExp(jet.NULL)
	absenceEnd := jet.DateExp(jet.NULL)
	if req.Props.AbsenceBegin != nil && req.Props.AbsenceEnd != nil {
		// Allow users to set their own absence date regardless of types perms check
		if userInfo.UserId != targetUser.UserId && !slices.Contains(types, "AbsenceDate") && !userInfo.SuperUser {
			return nil, errorsjobs.ErrPropsAbsenceDenied
		}

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

	tJobsUserProps := table.FivenetJobsUserProps
	updateSets := []jet.ColumnAssigment{
		tJobsUserProps.AbsenceBegin.SET(jet.DateExp(jet.Raw("VALUES(`absence_begin`)"))),
		tJobsUserProps.AbsenceEnd.SET(jet.DateExp(jet.Raw("VALUES(`absence_end`)"))),
	}
	// Generate the update sets
	if req.Props.Note != nil {
		if !slices.Contains(types, "Note") && !userInfo.SuperUser {
			return nil, errorsjobs.ErrPropsNoteDenied
		}

		updateSets = append(updateSets, tJobsUserProps.Note.SET(jet.String(*req.Props.Note)))
	} else {
		req.Props.Note = props.Note
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	stmt := tJobsUserProps.
		INSERT(
			tJobsUserProps.UserID,
			tJobsUserProps.Job,
			tJobsUserProps.AbsenceBegin,
			tJobsUserProps.AbsenceEnd,
			tJobsUserProps.Note,
		).
		VALUES(
			targetUser.UserId,
			targetUser.Job,
			absenceBegin,
			absenceEnd,
			req.Props.Note,
		).
		ON_DUPLICATE_KEY_UPDATE(
			updateSets...,
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	// Compare absence dates if any were set
	if req.Props.AbsenceBegin != nil && (props.AbsenceBegin == nil || req.Props.AbsenceBegin.AsTime().Compare(props.AbsenceBegin.AsTime()) != 0) ||
		req.Props.AbsenceEnd != nil && (props.AbsenceEnd == nil || req.Props.AbsenceEnd.AsTime().Compare(props.AbsenceEnd.AsTime()) != 0) {
		if err := s.addJobsUserActivity(ctx, tx, &jobs.JobsUserActivity{
			Job:          userInfo.Job,
			SourceUserId: userInfo.UserId,
			TargetUserId: targetUser.UserId,
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
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}
	if ((req.Props.Note == nil && props.Note == nil) ||
		(req.Props.Note == nil && props.Note != nil) || (req.Props.Note != nil && props.Note == nil)) ||
		*req.Props.Note != *props.Note {
		if err := s.addJobsUserActivity(ctx, tx, &jobs.JobsUserActivity{
			Job:          userInfo.Job,
			SourceUserId: userInfo.UserId,
			TargetUserId: targetUser.UserId,
			ActivityType: jobs.JobsUserActivityType_JOBS_USER_ACTIVITY_TYPE_NOTE,
			Reason:       req.Reason,
		}); err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	props, err = s.getJobsUserProps(ctx, targetUser.UserId, targetUser.Job, types)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
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

	// If no levels set, assume "Own" as a safe default
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

	// If no levels set, assume "Own" as a safe default
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
	trace.SpanFromContext(ctx).SetAttributes(attribute.IntSlice("fivenet.jobs.colleagues.user_ids", utils.SliceInt32ToInt(req.UserIds)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Access Field Permission Check
	accessAttr, err := s.ps.Attr(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceGetColleaguePerm, permsjobs.JobsServiceGetColleagueAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	var access perms.StringList
	if accessAttr != nil {
		access = accessAttr.([]string)
	}

	tJobsUserActivity := tJobsUserActivity.AS("jobsuseractivity")
	tTargetUser := tUser.AS("target_user")
	tSourceUser := tUser.AS("source_user")

	condition := tJobsUserActivity.Job.EQ(jet.String(userInfo.Job))

	resp := &ListColleagueActivityResponse{
		Pagination: &database.PaginationResponse{
			TotalCount: 0,
			Offset:     0,
			End:        0,
			PageSize:   15,
		},
	}

	if len(req.ActivityTypes) == 0 {
		return resp, nil
	}

	// If no user IDs given or more than 2, show all the user has access to
	if len(req.UserIds) == 0 || len(req.UserIds) >= 2 {
		condition = condition.AND(s.getConditionForColleagueAccess(tJobsUserActivity, tTargetUser, access, userInfo))

		if len(req.UserIds) >= 2 {
			// More than 2 user ids
			userIds := make([]jet.Expression, len(req.UserIds))
			for i := 0; i < len(req.UserIds); i++ {
				userIds[i] = jet.Int32(req.UserIds[i])
			}

			condition = condition.AND(tTargetUser.ID.IN(userIds...))
		}
	} else {
		userId := req.UserIds[0]

		targetUser, err := s.getColleague(ctx, userInfo.Job, userId, nil)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
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
	typesAttr, err := s.ps.Attr(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceListColleagueActivityPerm, permsjobs.JobsServiceListColleagueActivityTypesPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	var types perms.StringList
	if typesAttr != nil {
		types = typesAttr.([]string)
	}
	if len(types) == 0 {
		if !userInfo.SuperUser {
			return resp, nil
		} else {
			types = append(types, "HIRED", "FIRED", "PROMOTED", "DEMOTED", "ABSENCE_DATE", "NOTE")
		}
	}

	if len(req.ActivityTypes) > 0 {
		req.ActivityTypes = slices.DeleteFunc(req.ActivityTypes, func(t jobs.JobsUserActivityType) bool {
			return !slices.ContainsFunc(types, func(s string) bool {
				return strings.Contains(t.String(), "JOBS_USER_ACTIVITY_TYPE_"+s)
			})
		})
	}

	condTypes := []jet.Expression{}
	for _, aType := range req.ActivityTypes {
		condTypes = append(condTypes, jet.Int16(int16(aType)))
	}

	if len(condTypes) == 0 {
		return resp, nil
	}

	condition = condition.AND(tJobsUserActivity.ActivityType.IN(condTypes...))

	// Get total count of values
	countStmt := tJobsUserActivity.
		SELECT(
			jet.COUNT(tJobsUserActivity.ID).AS("datacount.totalcount"),
		).
		FROM(
			tJobsUserActivity.
				INNER_JOIN(tTargetUser,
					tTargetUser.ID.EQ(tJobsUserActivity.TargetUserID),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, 16)
	resp.Pagination = pag
	if count.TotalCount <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []jet.OrderByClause{}
	if req.Sort != nil {
		var column jet.Column
		switch req.Sort.Column {
		case "createdAt":
			fallthrough
		default:
			column = tJobsUserActivity.CreatedAt
		}

		if req.Sort.Direction == database.AscSortDirection {
			orderBys = append(orderBys, column.ASC())
		} else {
			orderBys = append(orderBys, column.DESC())
		}
	} else {
		orderBys = append(orderBys, tJobsUserActivity.CreatedAt.DESC())
	}

	tTargetUserProps := tUserProps.AS("target_user_props")
	tTargetJobsUserProps := tJobsUserProps.AS("fivenet_jobs_user_props")
	tSourceUserProps := tUserProps.AS("source_user_props")

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
			tTargetUser.ID,
			tTargetUser.Job,
			tTargetUser.JobGrade,
			tTargetUser.Firstname,
			tTargetUser.Lastname,
			tTargetUser.Dateofbirth,
			tTargetUser.PhoneNumber,
			tTargetUserProps.Avatar.AS("target_user.avatar"),
			tTargetJobsUserProps.UserID,
			tTargetJobsUserProps.Job,
			tTargetJobsUserProps.AbsenceBegin,
			tTargetJobsUserProps.AbsenceEnd,
			tSourceUser.ID,
			tSourceUser.Job,
			tSourceUser.JobGrade,
			tSourceUser.Firstname,
			tSourceUser.Lastname,
			tSourceUser.Dateofbirth,
			tSourceUser.PhoneNumber,
			tSourceUserProps.Avatar.AS("source_user.avatar"),
		).
		FROM(
			tJobsUserActivity.
				INNER_JOIN(tTargetUser,
					tTargetUser.ID.EQ(tJobsUserActivity.TargetUserID),
				).
				LEFT_JOIN(tTargetUserProps,
					tTargetUserProps.UserID.EQ(tTargetUser.ID),
				).
				LEFT_JOIN(tTargetJobsUserProps,
					tTargetJobsUserProps.UserID.EQ(tTargetUser.ID).
						AND(tTargetUser.Job.EQ(jet.String(userInfo.Job))),
				).
				LEFT_JOIN(tSourceUser,
					tSourceUser.ID.EQ(tJobsUserActivity.SourceUserID),
				).
				LEFT_JOIN(tSourceUserProps,
					tSourceUserProps.UserID.EQ(tSourceUser.ID),
				),
		).
		WHERE(condition).
		OFFSET(pag.Offset).
		ORDER_BY(orderBys...).
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
