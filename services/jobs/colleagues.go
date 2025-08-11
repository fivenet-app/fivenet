package jobs

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	jobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/notifications"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbjobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/jobs"
	permsjobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/access"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils/timeutils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsjobs "github.com/fivenet-app/fivenet/v2025/services/jobs/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

const (
	ColleaguesDefaultPageSize = 20
)

var (
	tColleagueProps    = table.FivenetJobColleagueProps.AS("colleague_props")
	tColleagueActivity = table.FivenetJobColleagueActivity

	tAvatar = table.FivenetFiles.AS("avatar")
)

func (s *Server) ListColleagues(ctx context.Context, req *pbjobs.ListColleaguesRequest) (*pbjobs.ListColleaguesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Access Permission Check
	types, err := s.ps.AttrStringList(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceGetColleaguePerm, permsjobs.JobsServiceGetColleagueTypesPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	tColleague := tables.User().AS("colleague")

	condition := tColleague.Job.EQ(jet.String(userInfo.Job)).
		AND(s.customDB.Conditions.User.GetFilter(tColleague.Alias()))

	userIds := []jet.Expression{}
	for _, v := range req.UserIds {
		userIds = append(userIds, jet.Int32(v))
	}
	if len(req.UserIds) > 0 && req.UserOnly != nil && *req.UserOnly {
		condition = condition.AND(tColleague.ID.IN(userIds...))
	} else {
		req.Search = strings.TrimSpace(req.Search)
		req.Search = strings.ReplaceAll(req.Search, "%", "")
		req.Search = strings.ReplaceAll(req.Search, " ", "%")
		if req.Search != "" {
			req.Search = "%" + req.Search + "%"
			condition = condition.AND(
				jet.CONCAT(tColleague.Firstname, jet.String(" "), tColleague.Lastname).
					LIKE(jet.String(req.Search)),
			)
		}
	}

	if req.Absent != nil && *req.Absent {
		condition = condition.AND(
			jet.AND(
				tColleagueProps.AbsenceBegin.IS_NOT_NULL(),
				tColleagueProps.AbsenceEnd.IS_NOT_NULL(),
				tColleagueProps.AbsenceBegin.LT_EQ(jet.CURRENT_DATE()),
				tColleagueProps.AbsenceEnd.GT_EQ(jet.CURRENT_DATE()),
			))
	}

	if len(req.LabelIds) > 0 && (types.Contains("Labels") || userInfo.Superuser) {
		labelIds := []jet.Expression{}
		for _, labelId := range req.LabelIds {
			labelIds = append(labelIds, jet.Uint64(labelId))
		}

		condition = condition.AND(
			tColleagueLabels.LabelID.IN(labelIds...),
		)
	}

	if req.NamePrefix != nil && *req.NamePrefix != "" {
		*req.NamePrefix = strings.TrimSpace(*req.NamePrefix)
		*req.NamePrefix = strings.ReplaceAll(*req.NamePrefix, "%", "")
		*req.NamePrefix = strings.ReplaceAll(*req.NamePrefix, " ", "%")
		if *req.NamePrefix != "" {
			*req.NamePrefix = "%" + *req.NamePrefix + "%"

			condition = condition.AND(jet.AND(
				tColleagueProps.NamePrefix.IS_NOT_NULL(),
				tColleagueProps.NamePrefix.LIKE(jet.String(*req.NamePrefix)),
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
				tColleagueProps.NameSuffix.IS_NOT_NULL(),
				tColleagueProps.NameSuffix.LIKE(jet.String(*req.NameSuffix)),
			))
		}
	}

	// Get total count of values
	countStmt := tColleague.
		SELECT(
			jet.COUNT(jet.DISTINCT(tColleague.ID)).AS("data_count.total"),
		).
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname_fulltext"))

	if len(req.LabelIds) > 0 && (types.Contains("Labels") || userInfo.Superuser) {
		countStmt = countStmt.
			FROM(
				tColleague.
					LEFT_JOIN(tColleagueProps,
						tColleagueProps.UserID.EQ(tColleague.ID).
							AND(tColleagueProps.Job.EQ(jet.String(userInfo.Job))),
					).
					INNER_JOIN(tColleagueLabels,
						tColleagueLabels.UserID.EQ(tColleague.ID).
							AND(tColleagueLabels.Job.EQ(jet.String(userInfo.Job))),
					).
					LEFT_JOIN(tJobLabels,
						tJobLabels.ID.EQ(tColleagueLabels.LabelID),
					),
			)
	} else {
		countStmt = countStmt.
			FROM(
				tColleague.
					LEFT_JOIN(tColleagueProps,
						tColleagueProps.UserID.EQ(tColleague.ID).
							AND(tColleagueProps.Job.EQ(jet.String(userInfo.Job))),
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

	pag, limit := req.Pagination.GetResponseWithPageSize(count.Total, ColleaguesDefaultPageSize)
	resp := &pbjobs.ListColleaguesResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []jet.OrderByClause{}
	if len(req.UserIds) > 0 {
		// Make sure to sort by the user ID if provided
		orderBys = append(orderBys, tColleague.ID.IN(userIds...).DESC())
	}

	if req.Sort != nil {
		var columns []jet.Column
		switch req.Sort.Column {
		case "name":
			columns = append(columns, tColleague.Firstname, tColleague.Lastname)
		case "rank":
			fallthrough
		default:
			columns = append(columns, tColleague.JobGrade)
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
			tColleague.JobGrade.ASC(),
			tColleague.Firstname.ASC(),
			tColleague.Lastname.ASC(),
		)
	}

	stmt := tColleague.
		SELECT(
			tColleague.ID,
			tColleague.Job,
			tColleague.JobGrade,
			tColleague.Firstname,
			tColleague.Lastname,
			tColleague.Dateofbirth,
			tColleague.PhoneNumber,
			tUserProps.AvatarFileID.AS("colleague.avatar_file_id"),
			tAvatar.FilePath.AS("colleague.avatar"),
			tUserProps.Email.AS("colleague.email"),
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.AbsenceBegin,
			tColleagueProps.AbsenceEnd,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
		).
		OPTIMIZER_HINTS(jet.OptimizerHint("idx_users_firstname_lastname_fulltext"))

	if len(req.LabelIds) > 0 && (types.Contains("Labels") || userInfo.Superuser) {
		stmt = stmt.
			FROM(
				tColleague.
					LEFT_JOIN(tUserProps,
						tUserProps.UserID.EQ(tColleague.ID),
					).
					LEFT_JOIN(tColleagueProps,
						tColleagueProps.UserID.EQ(tColleague.ID).
							AND(tColleagueProps.Job.EQ(jet.String(userInfo.Job))),
					).
					INNER_JOIN(tColleagueLabels,
						tColleagueLabels.UserID.EQ(tColleague.ID).
							AND(tColleagueLabels.Job.EQ(jet.String(userInfo.Job))),
					).
					LEFT_JOIN(tJobLabels,
						tJobLabels.ID.EQ(tColleagueLabels.LabelID),
					).
					LEFT_JOIN(tAvatar,
						tAvatar.ID.EQ(tUserProps.AvatarFileID),
					),
			)
	} else {
		stmt = stmt.
			FROM(
				tColleague.
					LEFT_JOIN(tUserProps,
						tUserProps.UserID.EQ(tColleague.ID),
					).
					LEFT_JOIN(tColleagueProps,
						tColleagueProps.UserID.EQ(tColleague.ID).
							AND(tColleagueProps.Job.EQ(jet.String(userInfo.Job))),
					).
					LEFT_JOIN(tAvatar,
						tAvatar.ID.EQ(tUserProps.AvatarFileID),
					),
			)
	}

	stmt = stmt.
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		GROUP_BY(tColleague.ID).
		ORDER_BY(orderBys...).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Colleagues); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	if len(resp.Colleagues) > 0 && (types.Contains("Labels") || userInfo.Superuser) {
		userIds := []jet.Expression{}
		for _, colleague := range resp.Colleagues {
			userIds = append(userIds, jet.Int32(colleague.UserId))
		}

		labelsStmt := tColleagueLabels.
			SELECT(
				tColleagueLabels.UserID.AS("user_id"),
				tJobLabels.ID,
				tJobLabels.Job,
				tJobLabels.Name,
				tJobLabels.Color,
			).
			FROM(
				tColleagueLabels.
					LEFT_JOIN(tJobLabels,
						tJobLabels.ID.EQ(tColleagueLabels.LabelID),
					),
			).
			WHERE(jet.AND(
				tJobLabels.Job.EQ(jet.String(userInfo.Job)),
				tJobLabels.DeletedAt.IS_NULL(),
				tColleagueLabels.UserID.IN(userIds...),
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
				colleague.Props = &jobs.ColleagueProps{
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
	tColleague := tables.User().AS("colleague")

	columns := []jet.Projection{
		tColleague.Firstname,
		tColleague.Lastname,
		tColleague.Job,
		tColleague.JobGrade,
		tColleague.Dateofbirth,
		tColleague.PhoneNumber,
		tUserProps.AvatarFileID.AS("colleague.avatar_file_id"),
		tAvatar.FilePath.AS("colleague.avatar"),
		tUserProps.Email.AS("colleague.email"),
		tColleagueProps.UserID,
		tColleagueProps.Job,
		tColleagueProps.AbsenceBegin,
		tColleagueProps.AbsenceEnd,
		tColleagueProps.NamePrefix,
		tColleagueProps.NameSuffix,
	}
	columns = append(columns, withColumns...)

	stmt := tColleague.
		SELECT(
			tColleague.ID,
			columns...,
		).
		FROM(
			tColleague.
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tColleague.ID),
				).
				LEFT_JOIN(tColleagueProps,
					tColleagueProps.UserID.EQ(tColleague.ID).
						AND(tColleagueProps.Job.EQ(jet.String(job))),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		WHERE(
			tColleague.ID.EQ(jet.Int32(userId)),
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
		dest.Props = &jobs.ColleagueProps{
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
	logging.InjectFields(ctx, logging.Fields{"fivenet.jobs.colleagues.user_id", req.UserId})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbjobs.JobsService_ServiceDesc.ServiceName,
		Method:  "GetColleague",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
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

	check := access.CheckIfHasOwnJobAccess(colleagueAccess, userInfo, targetUser.Job, &users.UserShort{
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
	fields := &permissions.StringList{}
	if !infoOnly {
		fields, err = s.ps.AttrStringList(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceGetColleaguePerm, permsjobs.JobsServiceGetColleagueTypesPermField)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}
	if userInfo.Superuser {
		fields.Strings = []string{"Note"}
	}

	withColumns := []jet.Projection{}
	for _, fType := range fields.Strings {
		switch fType {
		case "Note":
			withColumns = append(withColumns, tColleagueProps.Note)
		}
	}

	colleague, err := s.getColleague(ctx, userInfo, userInfo.Job, targetUser.UserId, withColumns)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_VIEWED

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
	if userInfo.Superuser {
		types.Strings = append(types.Strings, "AbsenceDate", "Note")
	}

	withColumns := []jet.Projection{}
	for _, fType := range types.Strings {
		switch fType {
		case "Note":
			withColumns = append(withColumns, tColleagueProps.Note)
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

func (s *Server) SetColleagueProps(ctx context.Context, req *pbjobs.SetColleaguePropsRequest) (*pbjobs.SetColleaguePropsResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.jobs.colleagues.user_id", req.Props.UserId})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbjobs.JobsService_ServiceDesc.ServiceName,
		Method:  "SetColleagueProps",
		UserId:  userInfo.UserId,
		UserJob: userInfo.Job,
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	if req.Reason == "" {
		return nil, errorsjobs.ErrReasonRequired
	}

	// Access Permission Check
	colleagueAccess, err := s.ps.AttrStringList(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceSetColleaguePropsPerm, permsjobs.JobsServiceSetColleaguePropsAccessPermField)
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

	if !access.CheckIfHasOwnJobAccess(colleagueAccess, userInfo, targetUser.Job, &users.UserShort{
		UserId:   targetUser.UserId,
		Job:      targetUser.Job,
		JobGrade: targetUser.JobGrade,
	}) {
		return nil, errorsjobs.ErrFailedQuery
	}

	// Types Permission Check
	types, err := s.ps.AttrStringList(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceSetColleaguePropsPerm, permsjobs.JobsServiceSetColleaguePropsTypesPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if userInfo.Superuser {
		types.Strings = []string{"AbsenceDate", "Note", "Labels", "Name"}
	}

	props, err := jobs.GetColleagueProps(ctx, s.db, targetUser.Job, targetUser.UserId, types.Strings)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	if req.Props.AbsenceBegin != nil && req.Props.AbsenceEnd != nil {
		// Allow users to set their own absence date regardless of types perms check
		if userInfo.UserId != targetUser.UserId && !types.Contains("AbsenceDate") && !userInfo.Superuser {
			return nil, errorsjobs.ErrPropsAbsenceDenied
		}

		jobProps, err := jobs.GetJobProps(ctx, s.db, userInfo.Job)
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
		if !types.Contains("Note") && !userInfo.Superuser {
			return nil, errorsjobs.ErrPropsNoteDenied
		}
	}

	if req.Props.Labels != nil {
		// Check if user is allowed to update labels
		if !types.Contains("Labels") && !userInfo.Superuser {
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
		if !types.Contains("Name") && !userInfo.Superuser {
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

	if err := jobs.CreateColleagueActivity(ctx, tx, activities...); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_UPDATED

	props, err = jobs.GetColleagueProps(ctx, s.db, targetUser.Job, targetUser.UserId, types.Strings)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	userId := uint64(targetUser.UserId)
	s.notifi.SendObjectEvent(ctx, &notifications.ObjectEvent{
		Type:      notifications.ObjectType_OBJECT_TYPE_JOBS_COLLEAGUE,
		Id:        &userId,
		EventType: notifications.ObjectEventType_OBJECT_EVENT_TYPE_UPDATED,

		UserId: &userInfo.UserId,
		Job:    &userInfo.Job,
	})

	return &pbjobs.SetColleaguePropsResponse{
		Props: props,
	}, nil
}

func (s *Server) getConditionForColleagueAccess(actTable *table.FivenetJobColleagueActivityTable, usersTable *tables.FivenetUserTable, levels []string, userInfo *userinfo.UserInfo) jet.BoolExpression {
	condition := jet.Bool(true)
	if userInfo.Superuser {
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
	logging.InjectFields(ctx, logging.Fields{"fivenet.jobs.colleagues.user_ids", req.UserIds})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Access Field Permission Check
	colleagueAccess, err := s.ps.AttrStringList(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceGetColleaguePerm, permsjobs.JobsServiceGetColleagueAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	tColleagueActivity := tColleagueActivity.AS("colleague_activity")
	tTargetColleague := tables.User().AS("target_user")
	tSourceUser := tTargetColleague.AS("source_user")

	condition := tColleagueActivity.Job.EQ(jet.String(userInfo.Job))

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
		condition = condition.AND(s.getConditionForColleagueAccess(tColleagueActivity, tTargetColleague, colleagueAccess.Strings, userInfo))

		if len(req.UserIds) >= 2 {
			// More than 2 user ids
			userIds := make([]jet.Expression, len(req.UserIds))
			for i := range req.UserIds {
				userIds[i] = jet.Int32(req.UserIds[i])
			}

			condition = condition.AND(tTargetColleague.ID.IN(userIds...))
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

		if !access.CheckIfHasOwnJobAccess(colleagueAccess, userInfo, targetUser.Job, &users.UserShort{
			UserId:   targetUser.UserId,
			Job:      targetUser.Job,
			JobGrade: targetUser.JobGrade,
		}) {
			return nil, errorsjobs.ErrFailedQuery
		}

		condition = condition.AND(tColleagueActivity.TargetUserID.EQ(jet.Int32(userId)))
	}

	// Types Field Permission Check
	types, err := s.ps.AttrStringList(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceListColleagueActivityPerm, permsjobs.JobsServiceListColleagueActivityTypesPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if types.Len() == 0 {
		if !userInfo.Superuser {
			return resp, nil
		} else {
			types.Strings = append(types.Strings, "HIRED", "FIRED", "PROMOTED", "DEMOTED", "ABSENCE_DATE", "NOTE", "LABELS", "NAME")
		}
	}

	if len(req.ActivityTypes) > 0 {
		req.ActivityTypes = slices.DeleteFunc(req.ActivityTypes, func(t jobs.ColleagueActivityType) bool {
			return !slices.ContainsFunc(types.Strings, func(s string) bool {
				return strings.Contains(t.String(), "COLLEAGUE_ACTIVITY_TYPE_"+s)
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

	condition = condition.AND(tColleagueActivity.ActivityType.IN(condTypes...))

	// Get total count of values
	countStmt := tColleagueActivity.
		SELECT(
			jet.COUNT(tColleagueActivity.ID).AS("data_count.total"),
		).
		FROM(
			tColleagueActivity.
				INNER_JOIN(tTargetColleague,
					tTargetColleague.ID.EQ(tColleagueActivity.TargetUserID),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.Total, ColleaguesDefaultPageSize)
	resp.Pagination = pag
	if count.Total <= 0 {
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
			column = tColleagueActivity.CreatedAt
		}

		if req.Sort.Direction == database.AscSortDirection {
			orderBys = append(orderBys, column.ASC())
		} else {
			orderBys = append(orderBys, column.DESC())
		}
	} else {
		orderBys = append(orderBys, tColleagueActivity.CreatedAt.DESC())
	}

	tTargetUserProps := tUserProps.AS("target_user_props")
	tTargetUserAvatar := tAvatar.AS("target_user_avatar")
	tTargetColleagueProps := tColleagueProps.AS("fivenet_colleague_props")
	tSourceUserProps := tUserProps.AS("source_user_props")
	tSourceUserAvatar := tAvatar.AS("source_user_avatar")

	stmt := tColleagueActivity.
		SELECT(
			tColleagueActivity.ID,
			tColleagueActivity.CreatedAt,
			tColleagueActivity.Job,
			tColleagueActivity.SourceUserID,
			tColleagueActivity.TargetUserID,
			tColleagueActivity.ActivityType,
			tColleagueActivity.Reason,
			tColleagueActivity.Data,
			tTargetColleague.ID,
			tTargetColleague.Job,
			tTargetColleague.JobGrade,
			tTargetColleague.Firstname,
			tTargetColleague.Lastname,
			tTargetColleague.Dateofbirth,
			tTargetColleague.PhoneNumber,
			tTargetUserProps.AvatarFileID.AS("target_user.avatar_file_id"),
			tTargetUserAvatar.FilePath.AS("target_user.avatar"),
			tTargetColleagueProps.UserID,
			tTargetColleagueProps.Job,
			tTargetColleagueProps.AbsenceBegin,
			tTargetColleagueProps.AbsenceEnd,
			tTargetColleagueProps.NamePrefix,
			tTargetColleagueProps.NameSuffix,
			tSourceUser.ID,
			tSourceUser.Job,
			tSourceUser.JobGrade,
			tSourceUser.Firstname,
			tSourceUser.Lastname,
			tSourceUser.Dateofbirth,
			tSourceUser.PhoneNumber,
			tSourceUserProps.AvatarFileID.AS("source_user.avatar_file_id"),
			tSourceUserAvatar.FilePath.AS("source_user.avatar"),
		).
		FROM(
			tColleagueActivity.
				INNER_JOIN(tTargetColleague,
					tTargetColleague.ID.EQ(tColleagueActivity.TargetUserID),
				).
				LEFT_JOIN(tTargetUserProps,
					tTargetUserProps.UserID.EQ(tTargetColleague.ID),
				).
				LEFT_JOIN(tTargetUserAvatar,
					tTargetUserAvatar.ID.EQ(tTargetUserProps.AvatarFileID),
				).
				LEFT_JOIN(tTargetColleagueProps,
					tTargetColleagueProps.UserID.EQ(tTargetColleague.ID).
						AND(tTargetColleague.Job.EQ(jet.String(userInfo.Job))),
				).
				LEFT_JOIN(tSourceUser,
					tSourceUser.ID.EQ(tColleagueActivity.SourceUserID),
				).
				LEFT_JOIN(tSourceUserProps,
					tSourceUserProps.UserID.EQ(tSourceUser.ID),
				).
				LEFT_JOIN(tSourceUserAvatar,
					tSourceUserAvatar.ID.EQ(tSourceUserProps.AvatarFileID),
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
