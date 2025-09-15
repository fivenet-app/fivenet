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
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

const (
	defaultPageSize = 20

	nameColumn = "name"
	rankColumn = "rank"
)

var (
	tColleagueProps    = table.FivenetJobColleagueProps.AS("colleague_props")
	tColleagueActivity = table.FivenetJobColleagueActivity

	tAvatar = table.FivenetFiles.AS("profile_picture")
)

func (s *Server) ListColleagues(
	ctx context.Context,
	req *pbjobs.ListColleaguesRequest,
) (*pbjobs.ListColleaguesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Access Permission Check
	types, err := s.ps.AttrStringList(
		userInfo,
		permsjobs.JobsServicePerm,
		permsjobs.JobsServiceGetColleaguePerm,
		permsjobs.JobsServiceGetColleagueTypesPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	tColleague := tables.User().AS("colleague")

	condition := tColleague.Job.EQ(mysql.String(userInfo.GetJob())).
		AND(s.customDB.Conditions.User.GetFilter(tColleague.Alias()))

	userIds := []mysql.Expression{}
	for _, v := range req.GetUserIds() {
		userIds = append(userIds, mysql.Int32(v))
	}
	if len(userIds) > 0 && req.UserOnly != nil && req.GetUserOnly() {
		condition = condition.AND(tColleague.ID.IN(userIds...))
	} else {
		search := strings.TrimSpace(req.GetSearch())
		search = strings.ReplaceAll(search, "%", "")
		search = strings.ReplaceAll(search, " ", "%")
		if search != "" {
			search = "%" + search + "%"
			condition = condition.AND(
				mysql.CONCAT(tColleague.Firstname, mysql.String(" "), tColleague.Lastname).
					LIKE(mysql.String(search)),
			)
		}
	}

	if req.Absent != nil && req.GetAbsent() {
		condition = condition.AND(
			mysql.AND(
				tColleagueProps.AbsenceBegin.IS_NOT_NULL(),
				tColleagueProps.AbsenceEnd.IS_NOT_NULL(),
				tColleagueProps.AbsenceBegin.LT_EQ(mysql.CURRENT_DATE()),
				tColleagueProps.AbsenceEnd.GT_EQ(mysql.CURRENT_DATE()),
			))
	}

	if req.GetNamePrefix() != "" {
		namePrefix := strings.TrimSpace(req.GetNamePrefix())
		namePrefix = strings.ReplaceAll(namePrefix, "%", "")
		namePrefix = strings.ReplaceAll(namePrefix, " ", "%")
		if namePrefix != "" {
			namePrefix = "%" + namePrefix + "%"

			condition = condition.AND(mysql.AND(
				tColleagueProps.NamePrefix.IS_NOT_NULL(),
				tColleagueProps.NamePrefix.LIKE(mysql.String(namePrefix)),
			))
		}
	}

	if req.GetNameSuffix() != "" {
		nameSuffix := strings.TrimSpace(req.GetNameSuffix())
		nameSuffix = strings.ReplaceAll(nameSuffix, "%", "")
		nameSuffix = strings.ReplaceAll(nameSuffix, " ", "%")
		if nameSuffix != "" {
			nameSuffix = "%" + nameSuffix + "%"

			condition = condition.AND(mysql.AND(
				tColleagueProps.NameSuffix.IS_NOT_NULL(),
				tColleagueProps.NameSuffix.LIKE(mysql.String(nameSuffix)),
			))
		}
	}

	if len(req.GetLabelIds()) > 0 && (types.Contains("Labels") || userInfo.GetSuperuser()) {
		labelIDExprs := []mysql.Expression{}
		for _, labelId := range req.GetLabelIds() {
			labelIDExprs = append(labelIDExprs, mysql.Int64(labelId))
		}

		labelsExists := mysql.EXISTS(
			mysql.
				SELECT(mysql.Int(1)).
				FROM(tColleagueLabels).
				WHERE(
					tColleagueLabels.UserID.EQ(tColleague.ID).
						AND(tColleagueLabels.Job.EQ(mysql.String(userInfo.GetJob()))).
						AND(tColleagueLabels.LabelID.IN(labelIDExprs...)),
				),
		)

		condition = condition.AND(labelsExists)
	}

	// Get total count of values
	countStmt := tColleague.
		SELECT(
			mysql.COUNT(mysql.DISTINCT(tColleague.ID)).AS("data_count.total"),
		).
		OPTIMIZER_HINTS(mysql.OptimizerHint("idx_users_firstname_lastname_fulltext")).
		FROM(
			tColleague.
				LEFT_JOIN(tColleagueProps,
					tColleagueProps.UserID.EQ(tColleague.ID).
						AND(tColleagueProps.Job.EQ(mysql.String(userInfo.GetJob()))),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	pag, limit := req.GetPagination().
		GetResponseWithPageSize(count.Total, defaultPageSize)
	resp := &pbjobs.ListColleaguesResponse{
		Pagination: pag,
	}
	if count.Total <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []mysql.OrderByClause{}
	if len(userIds) > 0 {
		// Make sure to sort by the user ID if provided
		orderBys = append(orderBys, tColleague.ID.IN(userIds...).DESC())
	}

	if req.GetSort() != nil {
		var columns []mysql.Column
		switch req.GetSort().GetColumn() {
		case nameColumn:
			columns = append(columns, tColleague.Firstname, tColleague.Lastname)
		case rankColumn:
			fallthrough
		default:
			columns = append(columns, tColleague.JobGrade)
		}

		for _, column := range columns {
			if req.GetSort().GetDirection() == database.AscSortDirection {
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
			tUserProps.AvatarFileID.AS("colleague.profile_picture_file_id"),
			tAvatar.FilePath.AS("colleague.profile_picture"),
			tUserProps.Email.AS("colleague.email"),
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.AbsenceBegin,
			tColleagueProps.AbsenceEnd,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
		).
		OPTIMIZER_HINTS(mysql.OptimizerHint("idx_users_firstname_lastname_fulltext")).
		FROM(
			tColleague.
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tColleague.ID),
				).
				LEFT_JOIN(tColleagueProps,
					tColleagueProps.UserID.EQ(tColleague.ID).
						AND(tColleagueProps.Job.EQ(mysql.String(userInfo.GetJob()))),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		WHERE(condition).
		OFFSET(req.GetPagination().GetOffset()).
		ORDER_BY(orderBys...).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Colleagues); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	if len(resp.GetColleagues()) > 0 && (types.Contains("Labels") || userInfo.GetSuperuser()) {
		userIds := []mysql.Expression{}
		for _, colleague := range resp.GetColleagues() {
			userIds = append(userIds, mysql.Int32(colleague.GetUserId()))
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
			WHERE(mysql.AND(
				tJobLabels.Job.EQ(mysql.String(userInfo.GetJob())),
				tJobLabels.DeletedAt.IS_NULL(),
				tColleagueLabels.UserID.IN(userIds...),
			)).
			ORDER_BY(
				tJobLabels.Order.ASC(),
			)

		labels := []*struct {
			UserId int32 `alias:"userId" sql:"primary_key"`
			Labels *jobs.Labels
		}{}
		if err := labelsStmt.QueryContext(ctx, s.db, &labels); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
			}
		}

		for _, props := range labels {
			idx := slices.IndexFunc(resp.GetColleagues(), func(c *jobs.Colleague) bool {
				return c.GetUserId() == props.UserId
			})
			if idx == -1 {
				continue
			}

			colleague := resp.GetColleagues()[idx]
			if colleague.GetProps() == nil {
				colleague.Props = &jobs.ColleagueProps{
					UserId: colleague.GetUserId(),
					Job:    userInfo.GetJob(),
				}
			}

			colleague.Props.Labels = props.Labels
		}
	}

	resp.GetPagination().Update(len(resp.GetColleagues()))

	for i := range resp.GetColleagues() {
		s.enricher.EnrichJobInfo(resp.GetColleagues()[i])
	}

	return resp, nil
}

func (s *Server) getColleague(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	job string,
	userId int32,
	withColumns mysql.ProjectionList,
) (*jobs.Colleague, error) {
	tColleague := tables.User().AS("colleague")

	columns := mysql.ProjectionList{
		tColleague.Firstname,
		tColleague.Lastname,
		tColleague.Job,
		tColleague.JobGrade,
		tColleague.Dateofbirth,
		tColleague.PhoneNumber,
		tUserProps.AvatarFileID.AS("colleague.profile_picture_file_id"),
		tAvatar.FilePath.AS("colleague.profile_picture"),
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
						AND(tColleagueProps.Job.EQ(mysql.String(job))),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		WHERE(
			tColleague.ID.EQ(mysql.Int32(userId)),
		).
		LIMIT(1)

	dest := &jobs.Colleague{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.GetUserId() == 0 {
		return nil, nil
	}

	if dest.GetProps() == nil {
		dest.Props = &jobs.ColleagueProps{
			UserId: dest.GetUserId(),
			Job:    userInfo.GetJob(),
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

func (s *Server) GetColleague(
	ctx context.Context,
	req *pbjobs.GetColleagueRequest,
) (*pbjobs.GetColleagueResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.jobs.colleagues.user_id", req.GetUserId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbjobs.JobsService_ServiceDesc.ServiceName,
		Method:  "GetColleague",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	// Access Permission Check
	colleagueAccess, err := s.ps.AttrStringList(
		userInfo,
		permsjobs.JobsServicePerm,
		permsjobs.JobsServiceGetColleaguePerm,
		permsjobs.JobsServiceGetColleagueAccessPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	targetUser, err := s.getColleague(ctx, userInfo, userInfo.GetJob(), req.GetUserId(), nil)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if targetUser == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}

	infoOnly := req.GetInfoOnly()

	check := access.CheckIfHasOwnJobAccess(
		colleagueAccess,
		userInfo,
		targetUser.GetJob(),
		&users.UserShort{
			UserId:   targetUser.GetUserId(),
			Job:      targetUser.GetJob(),
			JobGrade: targetUser.GetJobGrade(),
		},
	)
	if !check {
		if userInfo.GetJob() != targetUser.GetJob() || !infoOnly {
			return nil, errorsjobs.ErrFailedQuery
		}
	}

	// Field Permission Check
	fields := &permissions.StringList{}
	if !infoOnly {
		fields, err = s.ps.AttrStringList(
			userInfo,
			permsjobs.JobsServicePerm,
			permsjobs.JobsServiceGetColleaguePerm,
			permsjobs.JobsServiceGetColleagueTypesPermField,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}
	if userInfo.GetSuperuser() {
		fields.Strings = []string{"Note"}
	}

	withColumns := mysql.ProjectionList{}
	for _, fType := range fields.GetStrings() {
		switch fType {
		case "Note":
			withColumns = append(withColumns, tColleagueProps.Note)
		}
	}

	colleague, err := s.getColleague(
		ctx,
		userInfo,
		userInfo.GetJob(),
		targetUser.GetUserId(),
		withColumns,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	auditEntry.State = audit.EventType_EVENT_TYPE_VIEWED

	return &pbjobs.GetColleagueResponse{
		Colleague: colleague,
	}, nil
}

func (s *Server) GetSelf(
	ctx context.Context,
	req *pbjobs.GetSelfRequest,
) (*pbjobs.GetSelfResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Field Permission Check
	types, err := s.ps.AttrStringList(
		userInfo,
		permsjobs.JobsServicePerm,
		permsjobs.JobsServiceGetColleaguePerm,
		permsjobs.JobsServiceGetColleagueTypesPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if userInfo.GetSuperuser() {
		types.Strings = append(types.Strings, "AbsenceDate", "Note")
	}

	withColumns := mysql.ProjectionList{}
	for _, fType := range types.GetStrings() {
		switch fType {
		case "Note":
			withColumns = append(withColumns, tColleagueProps.Note)
		}
	}

	colleague, err := s.getColleague(
		ctx,
		userInfo,
		userInfo.GetJob(),
		userInfo.GetUserId(),
		withColumns,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	return &pbjobs.GetSelfResponse{
		Colleague: colleague,
	}, nil
}

func (s *Server) SetColleagueProps(
	ctx context.Context,
	req *pbjobs.SetColleaguePropsRequest,
) (*pbjobs.SetColleaguePropsResponse, error) {
	logging.InjectFields(
		ctx,
		logging.Fields{"fivenet.jobs.colleagues.user_id", req.GetProps().GetUserId()},
	)

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	auditEntry := &audit.AuditEntry{
		Service: pbjobs.JobsService_ServiceDesc.ServiceName,
		Method:  "SetColleagueProps",
		UserId:  userInfo.GetUserId(),
		UserJob: userInfo.GetJob(),
		State:   audit.EventType_EVENT_TYPE_ERRORED,
	}
	defer s.aud.Log(auditEntry, req)

	if req.GetReason() == "" {
		return nil, errorsjobs.ErrReasonRequired
	}

	// Access Permission Check
	colleagueAccess, err := s.ps.AttrStringList(
		userInfo,
		permsjobs.JobsServicePerm,
		permsjobs.JobsServiceSetColleaguePropsPerm,
		permsjobs.JobsServiceSetColleaguePropsAccessPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	targetUser, err := s.getColleague(
		ctx,
		userInfo,
		userInfo.GetJob(),
		req.GetProps().GetUserId(),
		nil,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if targetUser == nil {
		return nil, errorsjobs.ErrNotFoundOrNoPerms
	}

	if !access.CheckIfHasOwnJobAccess(
		colleagueAccess,
		userInfo,
		targetUser.GetJob(),
		&users.UserShort{
			UserId:   targetUser.GetUserId(),
			Job:      targetUser.GetJob(),
			JobGrade: targetUser.GetJobGrade(),
		},
	) {
		return nil, errorsjobs.ErrFailedQuery
	}

	// Types Permission Check
	types, err := s.ps.AttrStringList(
		userInfo,
		permsjobs.JobsServicePerm,
		permsjobs.JobsServiceSetColleaguePropsPerm,
		permsjobs.JobsServiceSetColleaguePropsTypesPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if userInfo.GetSuperuser() {
		types.Strings = []string{"AbsenceDate", "Note", "Labels", "Name"}
	}

	props, err := jobs.GetColleagueProps(
		ctx,
		s.db,
		targetUser.GetJob(),
		targetUser.GetUserId(),
		types.GetStrings(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	reqProps := req.GetProps()

	if reqProps.GetAbsenceBegin() != nil && reqProps.GetAbsenceEnd() != nil {
		// Allow users to set their own absence date regardless of types perms check
		if userInfo.GetUserId() != targetUser.GetUserId() && !types.Contains("AbsenceDate") &&
			!userInfo.GetSuperuser() {
			return nil, errorsjobs.ErrPropsAbsenceDenied
		}

		jobProps, err := jobs.GetJobProps(ctx, s.db, userInfo.GetJob())
		if err != nil {
			return nil, err
		}

		if reqProps.GetAbsenceBegin().GetTimestamp() != nil &&
			reqProps.GetAbsenceEnd().GetTimestamp() != nil {
			// Check if absence begin and end are "valid"
			minStart := time.Now().
				Add(-(time.Duration(jobProps.GetSettings().GetAbsencePastDays()) * 24 * time.Hour))
			maxEnd := time.Now().
				Add(time.Duration(jobProps.GetSettings().GetAbsenceFutureDays()) * 24 * time.Hour)

			if !timeutils.InTimeSpan(minStart, maxEnd, reqProps.GetAbsenceBegin().AsTime()) {
				return nil, errorsjobs.ErrAbsenceBeginOutOfRange
			}
			if !timeutils.InTimeSpan(minStart, maxEnd, reqProps.GetAbsenceEnd().AsTime()) {
				return nil, errorsjobs.ErrAbsenceBeginOutOfRange
			}
		}
	}

	// Generate the update sets
	if reqProps.Note != nil {
		if !types.Contains("Note") && !userInfo.GetSuperuser() {
			return nil, errorsjobs.ErrPropsNoteDenied
		}
	}

	if reqProps.GetLabels() != nil {
		// Check if user is allowed to update labels
		if !types.Contains("Labels") && !userInfo.GetSuperuser() {
			return nil, errorsjobs.ErrPropsAbsenceDenied
		}

		added, _ := utils.SlicesDifferenceFunc(
			props.GetLabels().GetList(),
			reqProps.GetLabels().GetList(),
			func(in *jobs.Label) int64 {
				return in.GetId()
			},
		)

		valid, err := s.validateLabels(ctx, userInfo, added)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		if !valid {
			return nil, errorsjobs.ErrPropsLabelsDenied
		}
	}

	if reqProps.NamePrefix != nil || reqProps.NameSuffix != nil {
		if !types.Contains("Name") && !userInfo.GetSuperuser() {
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

	activities, err := props.HandleChanges(
		ctx,
		tx,
		reqProps,
		targetUser.GetJob(),
		&userInfo.UserId,
		req.GetReason(),
	)
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

	props, err = jobs.GetColleagueProps(
		ctx,
		s.db,
		targetUser.GetJob(),
		targetUser.GetUserId(),
		types.GetStrings(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	userId := int64(targetUser.GetUserId())
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

func (s *Server) getConditionForColleagueAccess(
	actTable *table.FivenetJobColleagueActivityTable,
	usersTable *tables.FivenetUserTable,
	levels []string,
	userInfo *userinfo.UserInfo,
) mysql.BoolExpression {
	condition := mysql.Bool(true)
	if userInfo.GetSuperuser() {
		return condition
	}

	// If no levels set, assume "Own" as a safe default
	if len(levels) == 0 {
		return actTable.TargetUserID.EQ(mysql.Int32(userInfo.GetUserId()))
	}

	if slices.Contains(levels, "Any") {
		return condition
	}
	if slices.Contains(levels, "Lower_Rank") {
		return usersTable.ID.LT(mysql.Int32(userInfo.GetJobGrade()))
	}
	if slices.Contains(levels, "Same_Rank") {
		return usersTable.ID.LT_EQ(mysql.Int32(userInfo.GetJobGrade()))
	}
	if slices.Contains(levels, "Own") {
		return usersTable.ID.EQ(mysql.Int32(userInfo.GetUserId()))
	}

	return mysql.Bool(false)
}

func (s *Server) ListColleagueActivity(
	ctx context.Context,
	req *pbjobs.ListColleagueActivityRequest,
) (*pbjobs.ListColleagueActivityResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.jobs.colleagues.user_ids", req.GetUserIds()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Access Field Permission Check
	colleagueAccess, err := s.ps.AttrStringList(
		userInfo,
		permsjobs.JobsServicePerm,
		permsjobs.JobsServiceGetColleaguePerm,
		permsjobs.JobsServiceGetColleagueAccessPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	tColleagueActivity := tColleagueActivity.AS("colleague_activity")
	tTargetColleague := tables.User().AS("target_user")
	tSourceUser := tTargetColleague.AS("source_user")

	condition := tColleagueActivity.Job.EQ(mysql.String(userInfo.GetJob()))

	resp := &pbjobs.ListColleagueActivityResponse{
		Pagination: &database.PaginationResponse{
			PageSize: defaultPageSize,
		},
	}

	if len(req.GetActivityTypes()) == 0 {
		return resp, nil
	}

	// If no user IDs given or more than 2, show all the user has access to
	reqUserIds := req.GetUserIds()
	if len(reqUserIds) == 0 || len(reqUserIds) >= 2 {
		condition = condition.AND(
			s.getConditionForColleagueAccess(
				tColleagueActivity,
				tTargetColleague,
				colleagueAccess.GetStrings(),
				userInfo,
			),
		)

		if len(reqUserIds) >= 2 {
			// More than 2 user ids
			userIds := make([]mysql.Expression, len(reqUserIds))
			for i := range userIds {
				userIds[i] = mysql.Int32(reqUserIds[i])
			}

			condition = condition.AND(tTargetColleague.ID.IN(userIds...))
		}
	} else {
		userId := reqUserIds[0]

		targetUser, err := s.getColleague(ctx, userInfo, userInfo.GetJob(), userId, nil)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		if targetUser == nil {
			return nil, errorsjobs.ErrNotFoundOrNoPerms
		}

		if !access.CheckIfHasOwnJobAccess(colleagueAccess, userInfo, targetUser.GetJob(), &users.UserShort{
			UserId:   targetUser.GetUserId(),
			Job:      targetUser.GetJob(),
			JobGrade: targetUser.GetJobGrade(),
		}) {
			return nil, errorsjobs.ErrFailedQuery
		}

		condition = condition.AND(tColleagueActivity.TargetUserID.EQ(mysql.Int32(userId)))
	}

	// Types Field Permission Check
	types, err := s.ps.AttrStringList(
		userInfo,
		permsjobs.JobsServicePerm,
		permsjobs.JobsServiceListColleagueActivityPerm,
		permsjobs.JobsServiceListColleagueActivityTypesPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if types.Len() == 0 {
		if !userInfo.GetSuperuser() {
			return resp, nil
		} else {
			types.Strings = append(types.Strings, "HIRED", "FIRED", "PROMOTED", "DEMOTED", "ABSENCE_DATE", "NOTE", "LABELS", "NAME")
		}
	}

	if len(req.GetActivityTypes()) > 0 {
		req.ActivityTypes = slices.DeleteFunc(
			req.GetActivityTypes(),
			func(t jobs.ColleagueActivityType) bool {
				return !slices.ContainsFunc(types.GetStrings(), func(s string) bool {
					return strings.Contains(t.String(), "COLLEAGUE_ACTIVITY_TYPE_"+s)
				})
			},
		)
	}

	condTypes := []mysql.Expression{}
	for _, aType := range req.GetActivityTypes() {
		condTypes = append(condTypes, mysql.Int32(int32(aType)))
	}

	if len(condTypes) == 0 {
		return resp, nil
	}

	condition = condition.AND(tColleagueActivity.ActivityType.IN(condTypes...))

	// Get total count of values
	countStmt := tColleagueActivity.
		SELECT(
			mysql.COUNT(tColleagueActivity.ID).AS("data_count.total"),
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

	pag, limit := req.GetPagination().
		GetResponseWithPageSize(count.Total, defaultPageSize)
	resp.Pagination = pag
	if count.Total <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []mysql.OrderByClause{}
	if req.GetSort() != nil {
		var column mysql.Column
		switch req.GetSort().GetColumn() {
		case "createdAt":
			fallthrough
		default:
			column = tColleagueActivity.CreatedAt
		}

		if req.GetSort().GetDirection() == database.AscSortDirection {
			orderBys = append(orderBys, column.ASC())
		} else {
			orderBys = append(orderBys, column.DESC())
		}
	} else {
		orderBys = append(orderBys, tColleagueActivity.CreatedAt.DESC())
	}

	tTargetUserProps := tUserProps.AS("target_user_props")
	tTargetUserAvatar := tAvatar.AS("target_user_profile_picture")
	tTargetColleagueProps := tColleagueProps.AS("fivenet_colleague_props")
	tSourceUserProps := tUserProps.AS("source_user_props")
	tSourceUserAvatar := tAvatar.AS("source_user_profile_picture")

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
			tTargetUserProps.AvatarFileID.AS("target_user.profile_picture_file_id"),
			tTargetUserAvatar.FilePath.AS("target_user.profile_picture"),
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
			tSourceUserProps.AvatarFileID.AS("source_user.profile_picture_file_id"),
			tSourceUserAvatar.FilePath.AS("source_user.profile_picture"),
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
						AND(tTargetColleague.Job.EQ(mysql.String(userInfo.GetJob()))),
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
		OFFSET(pag.GetOffset()).
		ORDER_BY(orderBys...).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Activity); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	pag.Update(len(resp.GetActivity()))

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetActivity() {
		if resp.GetActivity()[i].GetSourceUser() != nil {
			jobInfoFn(resp.GetActivity()[i].GetSourceUser())
		}

		if resp.GetActivity()[i].GetTargetUser() != nil {
			jobInfoFn(resp.GetActivity()[i].GetTargetUser())
		}
	}

	return resp, nil
}
