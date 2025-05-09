package jobs

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	jobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/rector"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbjobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/jobs"
	permsjobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/timeutils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsjobs "github.com/fivenet-app/fivenet/v2025/services/jobs/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const (
	ColleaguesDefaultPageSize = 20
)

var (
	tJobsUserProps    = table.FivenetJobsUserProps.AS("jobs_user_props")
	tJobsUserActivity = table.FivenetJobsUserActivity
)

func (s *Server) ListColleagues(ctx context.Context, req *pbjobs.ListColleaguesRequest) (*pbjobs.ListColleaguesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Access Permission Check
	types, err := s.ps.AttrStringList(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceGetColleaguePerm, permsjobs.JobsServiceGetColleagueTypesPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	tUser := tables.Users().AS("colleague")

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

	if len(req.LabelIds) > 0 && (types.Contains("Labels") || userInfo.SuperUser) {
		labelIds := []jet.Expression{}
		for _, labelId := range req.LabelIds {
			labelIds = append(labelIds, jet.Uint64(labelId))
		}

		condition = condition.AND(
			tUserLabels.LabelID.IN(labelIds...),
		)
	}

	if req.NamePrefix != nil && *req.NamePrefix != "" {
		*req.NamePrefix = strings.TrimSpace(*req.NamePrefix)
		*req.NamePrefix = strings.ReplaceAll(*req.NamePrefix, "%", "")
		*req.NamePrefix = strings.ReplaceAll(*req.NamePrefix, " ", "%")
		if *req.NamePrefix != "" {
			*req.NamePrefix = "%" + *req.NamePrefix + "%"

			condition = condition.AND(jet.AND(
				tJobsUserProps.NamePrefix.IS_NOT_NULL(),
				tJobsUserProps.NamePrefix.LIKE(jet.String(*req.NamePrefix)),
			))
		}
	}

	if req.NameSuffix != nil {
		*req.NameSuffix = strings.TrimSpace(*req.NameSuffix)
		*req.NameSuffix = strings.ReplaceAll(*req.NameSuffix, "%", "")
		*req.NameSuffix = strings.ReplaceAll(*req.NameSuffix, " ", "%")
		if *req.NameSuffix != "" {
			*req.NameSuffix = "%" + *req.NameSuffix + "%"

			condition = condition.AND(jet.AND(
				tJobsUserProps.NameSuffix.IS_NOT_NULL(),
				tJobsUserProps.NameSuffix.LIKE(jet.String(*req.NameSuffix)),
			))
		}
	}

	// Get total count of values
	countStmt := tUser.
		SELECT(
			jet.COUNT(jet.DISTINCT(tUser.ID)).AS("datacount.totalcount"),
		).
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname_fulltext"))

	if len(req.LabelIds) > 0 && (types.Contains("Labels") || userInfo.SuperUser) {
		countStmt = countStmt.
			FROM(
				tUser.
					LEFT_JOIN(tJobsUserProps,
						tJobsUserProps.UserID.EQ(tUser.ID).
							AND(tJobsUserProps.Job.EQ(jet.String(userInfo.Job))),
					).
					INNER_JOIN(tUserLabels,
						tUserLabels.UserID.EQ(tUser.ID).
							AND(tUserLabels.Job.EQ(jet.String(userInfo.Job))),
					).
					LEFT_JOIN(tJobLabels,
						tJobLabels.ID.EQ(tUserLabels.LabelID),
					),
			)
	} else {
		countStmt = countStmt.
			FROM(
				tUser.
					LEFT_JOIN(tJobsUserProps,
						tJobsUserProps.UserID.EQ(tUser.ID).
							AND(tJobsUserProps.Job.EQ(jet.String(userInfo.Job))),
					),
			)
	}

	countStmt = countStmt.
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, ColleaguesDefaultPageSize)
	resp := &pbjobs.ListColleaguesResponse{
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
			tUserProps.Email.AS("colleague.email"),
			tJobsUserProps.UserID,
			tJobsUserProps.Job,
			tJobsUserProps.AbsenceBegin,
			tJobsUserProps.AbsenceEnd,
			tJobsUserProps.NamePrefix,
			tJobsUserProps.NameSuffix,
		).
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname_fulltext"))

	if len(req.LabelIds) > 0 && (types.Contains("Labels") || userInfo.SuperUser) {
		stmt = stmt.
			FROM(
				tUser.
					LEFT_JOIN(tUserProps,
						tUserProps.UserID.EQ(tUser.ID),
					).
					LEFT_JOIN(tJobsUserProps,
						tJobsUserProps.UserID.EQ(tUser.ID).
							AND(tJobsUserProps.Job.EQ(jet.String(userInfo.Job))),
					).
					INNER_JOIN(tUserLabels,
						tUserLabels.UserID.EQ(tUser.ID).
							AND(tUserLabels.Job.EQ(jet.String(userInfo.Job))),
					).
					LEFT_JOIN(tJobLabels,
						tJobLabels.ID.EQ(tUserLabels.LabelID),
					),
			)
	} else {
		stmt = stmt.
			FROM(
				tUser.
					LEFT_JOIN(tUserProps,
						tUserProps.UserID.EQ(tUser.ID),
					).
					LEFT_JOIN(tJobsUserProps,
						tJobsUserProps.UserID.EQ(tUser.ID).
							AND(tJobsUserProps.Job.EQ(jet.String(userInfo.Job))),
					),
			)
	}

	stmt = stmt.
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		GROUP_BY(tUser.ID).
		ORDER_BY(orderBys...).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Colleagues); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	if len(resp.Colleagues) > 0 && (types.Contains("Labels") || userInfo.SuperUser) {
		userIds := []jet.Expression{}
		for _, colleague := range resp.Colleagues {
			userIds = append(userIds, jet.Int32(colleague.UserId))
		}

		labelsStmt := tUserLabels.
			SELECT(
				tUserLabels.UserID.AS("user_id"),
				tJobLabels.ID,
				tJobLabels.Job,
				tJobLabels.Name,
				tJobLabels.Color,
			).
			FROM(
				tUserLabels.
					LEFT_JOIN(tJobLabels,
						tJobLabels.ID.EQ(tUserLabels.LabelID),
					),
			).
			WHERE(jet.AND(
				tJobLabels.Job.EQ(jet.String(userInfo.Job)),
				tUserLabels.UserID.IN(userIds...),
			)).
			ORDER_BY(
				tJobLabels.Order.ASC(),
			)

		labels := []*struct {
			UserId int32 `sql:"primary_key" alias:"userId"`
			Labels *jobs.Labels
		}{}
		if err := labelsStmt.QueryContext(ctx, s.db, &labels); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
			}
		}

		for _, props := range labels {
			idx := slices.IndexFunc(resp.Colleagues, func(c *jobs.Colleague) bool {
				return c.UserId == props.UserId
			})
			if idx == -1 {
				continue
			}

			colleague := resp.Colleagues[idx]
			if colleague.Props == nil {
				colleague.Props = &jobs.JobsUserProps{
					UserId: colleague.UserId,
					Job:    userInfo.Job,
				}
			}

			colleague.Props.Labels = props.Labels
		}
	}

	resp.Pagination.Update(len(resp.Colleagues))

	for i := range resp.Colleagues {
		s.enricher.EnrichJobInfo(resp.Colleagues[i])
	}

	return resp, nil
}

func (s *Server) getColleague(ctx context.Context, userInfo *userinfo.UserInfo, job string, userId int32, withColumns []jet.Projection) (*jobs.Colleague, error) {
	tUser := tables.Users().AS("colleague")

	columns := []jet.Projection{
		tUser.Firstname,
		tUser.Lastname,
		tUser.Job,
		tUser.JobGrade,
		tUser.Dateofbirth,
		tUser.PhoneNumber,
		tUserProps.Avatar.AS("colleague.avatar"),
		tUserProps.Email.AS("colleague.email"),
		tJobsUserProps.UserID,
		tJobsUserProps.Job,
		tJobsUserProps.AbsenceBegin,
		tJobsUserProps.AbsenceEnd,
		tJobsUserProps.NamePrefix,
		tJobsUserProps.NameSuffix,
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
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.UserId == 0 {
		return nil, nil
	}

	if dest.Props == nil {
		dest.Props = &jobs.JobsUserProps{
			UserId: dest.UserId,
			Job:    userInfo.Job,
		}
	}

	labels, err := s.getUserLabels(ctx, userInfo, userId)
	if err != nil {
		return nil, err
	}
	dest.Props.Labels = labels

	s.enricher.EnrichJobInfo(dest)

	return dest, nil
}

func (s *Server) GetColleague(ctx context.Context, req *pbjobs.GetColleagueRequest) (*pbjobs.GetColleagueResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.jobs.colleague.id", int64(req.UserId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbjobs.JobsService_ServiceDesc.ServiceName,
		Method:  "GetColleague",
		UserID:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   int16(rector.EventType_EVENT_TYPE_ERRORED),
	}
	defer s.aud.Log(auditEntry, req)

	// Access Permission Check
	colleagueAccess, err := s.ps.AttrStringList(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceGetColleaguePerm, permsjobs.JobsServiceGetColleagueAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	targetUser, err := s.getColleague(ctx, userInfo, userInfo.Job, req.UserId, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if targetUser == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}

	infoOnly := req.InfoOnly != nil && *req.InfoOnly

	check := access.CheckIfHasAccess(colleagueAccess, userInfo, targetUser.Job, &users.UserShort{
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
	var fields *permissions.StringList
	if !infoOnly {
		fields, err = s.ps.AttrStringList(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceGetColleaguePerm, permsjobs.JobsServiceGetColleagueTypesPermField)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}
	if userInfo.SuperUser {
		fields.Strings = []string{"Note"}
	}

	withColumns := []jet.Projection{}
	for _, fType := range fields.Strings {
		switch fType {
		case "Note":
			withColumns = append(withColumns, tJobsUserProps.Note)
		}
	}

	colleague, err := s.getColleague(ctx, userInfo, userInfo.Job, targetUser.UserId, withColumns)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_VIEWED)

	return &pbjobs.GetColleagueResponse{
		Colleague: colleague,
	}, nil
}

func (s *Server) GetSelf(ctx context.Context, req *pbjobs.GetSelfRequest) (*pbjobs.GetSelfResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Field Permission Check
	types, err := s.ps.AttrStringList(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceGetColleaguePerm, permsjobs.JobsServiceGetColleagueTypesPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if userInfo.SuperUser {
		types.Strings = append(types.Strings, "AbsenceDate", "Note")
	}

	withColumns := []jet.Projection{}
	for _, fType := range types.Strings {
		switch fType {
		case "Note":
			withColumns = append(withColumns, tJobsUserProps.Note)
		}
	}

	colleague, err := s.getColleague(ctx, userInfo, userInfo.Job, userInfo.UserId, withColumns)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	return &pbjobs.GetSelfResponse{
		Colleague: colleague,
	}, nil
}

func (s *Server) SetJobsUserProps(ctx context.Context, req *pbjobs.SetJobsUserPropsRequest) (*pbjobs.SetJobsUserPropsResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.Int64("fivenet.jobs.colleague.id", int64(req.Props.UserId)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &model.FivenetAuditLog{
		Service: pbjobs.JobsService_ServiceDesc.ServiceName,
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
	colleagueAccess, err := s.ps.AttrStringList(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceSetJobsUserPropsPerm, permsjobs.JobsServiceSetJobsUserPropsAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	targetUser, err := s.getColleague(ctx, userInfo, userInfo.Job, req.Props.UserId, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if targetUser == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}

	if !access.CheckIfHasAccess(colleagueAccess, userInfo, targetUser.Job, &users.UserShort{
		UserId:   targetUser.UserId,
		Job:      targetUser.Job,
		JobGrade: targetUser.JobGrade,
	}) {
		return nil, errorsjobs.ErrFailedQuery
	}

	// Types Permission Check
	types, err := s.ps.AttrStringList(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceSetJobsUserPropsPerm, permsjobs.JobsServiceSetJobsUserPropsTypesPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if userInfo.SuperUser {
		types.Strings = []string{"AbsenceDate", "Note", "Labels", "Name"}
	}

	props, err := jobs.GetJobsUserProps(ctx, s.db, targetUser.Job, targetUser.UserId, types.Strings)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	if req.Props.AbsenceBegin != nil && req.Props.AbsenceEnd != nil {
		// Allow users to set their own absence date regardless of types perms check
		if userInfo.UserId != targetUser.UserId && !types.Contains("AbsenceDate") && !userInfo.SuperUser {
			return nil, errorsjobs.ErrPropsAbsenceDenied
		}

		jobProps, err := users.GetJobProps(ctx, s.db, userInfo.Job)
		if err != nil {
			return nil, err
		}

		if req.Props.AbsenceBegin.Timestamp != nil && req.Props.AbsenceEnd.Timestamp != nil {
			// Check if absence begin and end are "valid"
			minStart := time.Now().Add(-(time.Duration(jobProps.Settings.AbsencePastDays) * 24 * time.Hour))
			maxEnd := time.Now().Add(time.Duration(jobProps.Settings.AbsenceFutureDays) * 24 * time.Hour)

			if !timeutils.InTimeSpan(minStart, maxEnd, req.Props.AbsenceBegin.AsTime()) {
				return nil, errorsjobs.ErrAbsenceBeginOutOfRange
			}
			if !timeutils.InTimeSpan(minStart, maxEnd, req.Props.AbsenceEnd.AsTime()) {
				return nil, errorsjobs.ErrAbsenceBeginOutOfRange
			}
		}
	}

	// Generate the update sets
	if req.Props.Note != nil {
		if !types.Contains("Note") && !userInfo.SuperUser {
			return nil, errorsjobs.ErrPropsNoteDenied
		}
	}

	if req.Props.Labels != nil {
		// Check if user is allowed to update labels
		if !types.Contains("Labels") && !userInfo.SuperUser {
			return nil, errorsjobs.ErrPropsAbsenceDenied
		}

		added, _ := utils.SlicesDifferenceFunc(props.Labels.List, req.Props.Labels.List,
			func(in *jobs.Label) uint64 {
				return in.Id
			})

		valid, err := s.validateLabels(ctx, userInfo, added)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		if !valid {
			return nil, errorsjobs.ErrPropsLabelsDenied
		}
	}

	if req.Props.NamePrefix != nil || req.Props.NameSuffix != nil {
		if !types.Contains("Name") && !userInfo.SuperUser {
			return nil, errorsjobs.ErrPropsNameDenied
		}
	}

	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	activities, err := props.HandleChanges(ctx, tx, req.Props, targetUser.Job, &userInfo.UserId, req.Reason)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	if err := jobs.CreateJobsUserActivities(ctx, tx, activities...); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	auditEntry.State = int16(rector.EventType_EVENT_TYPE_UPDATED)

	props, err = jobs.GetJobsUserProps(ctx, s.db, targetUser.Job, targetUser.UserId, types.Strings)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	return &pbjobs.SetJobsUserPropsResponse{
		Props: props,
	}, nil
}

func (s *Server) getConditionForColleagueAccess(actTable *table.FivenetJobsUserActivityTable, usersTable *tables.FivenetUsersTable, levels []string, userInfo *userinfo.UserInfo) jet.BoolExpression {
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

func (s *Server) ListColleagueActivity(ctx context.Context, req *pbjobs.ListColleagueActivityRequest) (*pbjobs.ListColleagueActivityResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.IntSlice("fivenet.jobs.colleagues.user_ids", utils.SliceInt32ToInt(req.UserIds)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Access Field Permission Check
	colleagueAccess, err := s.ps.AttrStringList(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceGetColleaguePerm, permsjobs.JobsServiceGetColleagueAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	tJobsUserActivity := tJobsUserActivity.AS("jobsuseractivity")
	tTargetUser := tables.Users().AS("target_user")
	tSourceUser := tTargetUser.AS("source_user")

	condition := tJobsUserActivity.Job.EQ(jet.String(userInfo.Job))

	resp := &pbjobs.ListColleagueActivityResponse{
		Pagination: &database.PaginationResponse{
			PageSize: ColleaguesDefaultPageSize,
		},
	}

	if len(req.ActivityTypes) == 0 {
		return resp, nil
	}

	// If no user IDs given or more than 2, show all the user has access to
	if len(req.UserIds) == 0 || len(req.UserIds) >= 2 {
		condition = condition.AND(s.getConditionForColleagueAccess(tJobsUserActivity, tTargetUser, colleagueAccess.Strings, userInfo))

		if len(req.UserIds) >= 2 {
			// More than 2 user ids
			userIds := make([]jet.Expression, len(req.UserIds))
			for i := range req.UserIds {
				userIds[i] = jet.Int32(req.UserIds[i])
			}

			condition = condition.AND(tTargetUser.ID.IN(userIds...))
		}
	} else {
		userId := req.UserIds[0]

		targetUser, err := s.getColleague(ctx, userInfo, userInfo.Job, userId, nil)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		if targetUser == nil {
			return nil, errorsjobs.ErrNotFoundOrNoPerms
		}

		if !access.CheckIfHasAccess(colleagueAccess, userInfo, targetUser.Job, &users.UserShort{
			UserId:   targetUser.UserId,
			Job:      targetUser.Job,
			JobGrade: targetUser.JobGrade,
		}) {
			return nil, errorsjobs.ErrFailedQuery
		}

		condition = condition.AND(tJobsUserActivity.TargetUserID.EQ(jet.Int32(userId)))
	}

	// Types Field Permission Check
	types, err := s.ps.AttrStringList(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceListColleagueActivityPerm, permsjobs.JobsServiceListColleagueActivityTypesPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if types.Len() == 0 {
		if !userInfo.SuperUser {
			return resp, nil
		} else {
			types.Strings = append(types.Strings, "HIRED", "FIRED", "PROMOTED", "DEMOTED", "ABSENCE_DATE", "NOTE", "LABELS", "NAME")
		}
	}

	if len(req.ActivityTypes) > 0 {
		req.ActivityTypes = slices.DeleteFunc(req.ActivityTypes, func(t jobs.JobsUserActivityType) bool {
			return !slices.ContainsFunc(types.Strings, func(s string) bool {
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

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, ColleaguesDefaultPageSize)
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
			tTargetJobsUserProps.NamePrefix,
			tTargetJobsUserProps.NameSuffix,
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
	for i := range resp.Activity {
		if resp.Activity[i].SourceUser != nil {
			jobInfoFn(resp.Activity[i].SourceUser)
		}

		if resp.Activity[i].TargetUser != nil {
			jobInfoFn(resp.Activity[i].TargetUser)
		}
	}

	return resp, nil
}
