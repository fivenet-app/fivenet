package jobs

import (
	"context"
	"slices"
	"strings"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/audit"
	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	colleaguesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues/activity"
	jobslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/labels"
	notificationsclientview "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/clientview"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	pbjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	permsjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	grpc_audit "github.com/fivenet-app/fivenet/v2026/pkg/grpc/interceptors/audit"
	"github.com/fivenet-app/fivenet/v2026/pkg/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils/timeutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

const (
	defaultPageSize = 20
)

var (
	tColleagueProps    = table.FivenetJobColleagueProps.AS("colleague_props")
	tColleagueActivity = table.FivenetJobColleagueActivity
)

func (s *Server) ListColleagues(
	ctx context.Context,
	req *pbjobs.ListColleaguesRequest,
) (*pbjobs.ListColleaguesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Access Permission Check
	types, err := permsjobs.ColleaguesService.GetColleague.TypesTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	tColleague := table.FivenetUser.AS("colleague")
	q := jobsstore.ListColleaguesQuery{
		Job:        userInfo.GetJob(),
		Search:     req.GetSearch(),
		UserIDs:    req.GetUserIds(),
		UserOnly:   req.GetUserOnly(),
		Absent:     req.GetAbsent(),
		NamePrefix: req.GetNamePrefix(),
		NameSuffix: req.GetNameSuffix(),
		Sort:       req.GetSort(),
		Offset:     req.GetPagination().GetOffset(),
		Where:      s.customDB.Conditions.User.GetFilter(tColleague.Alias()),
	}

	if len(req.GetLabelIds()) > 0 &&
		(types.Contains(permsjobs.ColleaguesServiceGetColleagueTypesPermValueLabels) || userInfo.GetJobAdmin()) {
		q.LabelIDs = req.GetLabelIds()
	}

	count, err := s.store.CountColleagues(ctx, s.db, q)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count, defaultPageSize)
	resp := &pbjobs.ListColleaguesResponse{Pagination: pag}
	if count <= 0 {
		return resp, nil
	}

	q.Limit = limit
	resp.Colleagues, err = s.store.ListColleagues(ctx, s.db, q)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	if len(resp.GetColleagues()) > 0 &&
		(types.Contains(permsjobs.ColleaguesServiceGetColleagueTypesPermValueLabels) || userInfo.GetJobAdmin()) {
		userIds := make([]int32, 0, len(resp.GetColleagues()))
		for _, colleague := range resp.GetColleagues() {
			userIds = append(userIds, colleague.GetUserId())
		}

		labels, err := s.store.GetUsersLabels(
			ctx,
			s.db,
			userInfo.GetJob(),
			userIds,
			userInfo.GetJobAdmin(),
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}

		labelsByUser := map[int32]*jobslabels.Labels{}
		for _, props := range labels {
			labelsByUser[props.UserId] = props.Labels
		}

		for i := range resp.GetColleagues() {
			colleague := resp.GetColleagues()[i]
			if colleagueLabels, ok := labelsByUser[colleague.GetUserId()]; ok {
				if colleague.GetProps() == nil {
					colleague.Props = &jobscolleagues.ColleagueProps{
						UserId: colleague.GetUserId(),
						Job:    userInfo.GetJob(),
					}
				}

				colleague.Props.Labels = colleagueLabels
			}
		}
	}

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
) (*jobscolleagues.Colleague, error) {
	dest, err := s.store.GetColleague(
		ctx,
		s.db,
		job,
		userId,
		withColumns,
		userInfo != nil && userInfo.GetJobAdmin(),
	)
	if err != nil {
		return nil, err
	}
	if dest == nil {
		return nil, nil
	}

	s.enricher.EnrichJobInfo(dest)

	return dest, nil
}

func (s *Server) GetColleague(
	ctx context.Context,
	req *pbjobs.GetColleagueRequest,
) (*pbjobs.GetColleagueResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.jobs.colleagues.user_id", req.GetUserId()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Access Permission Check
	colleagueAccess, err := permsjobs.ColleaguesService.GetColleague.AccessTyped.Get(s.ps, userInfo)
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
		colleagueAccess.StringList(),
		userInfo,
		targetUser.GetJob(),
		&usershort.UserShort{
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
	fields := perms.NewTypedStringList[permsjobs.ColleaguesServiceGetColleagueTypesPermValue]()
	if !infoOnly {
		fields, err = permsjobs.ColleaguesService.GetColleague.TypesTyped.Get(s.ps, userInfo)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}
	if userInfo.GetJobAdmin() {
		fields.Set(permsjobs.ColleaguesServiceGetColleagueTypesPermValueNote)
	}

	withColumns := mysql.ProjectionList{}
	for _, fType := range fields.Values() {
		switch fType {
		case permsjobs.ColleaguesServiceGetColleagueTypesPermValueNote:
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

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_VIEWED)

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
	types, err := permsjobs.ColleaguesService.GetColleague.TypesTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if userInfo.GetJobAdmin() {
		types.Append(
			permsjobs.ColleaguesServiceGetColleagueTypesPermValueNote,
		)
	}

	withColumns := mysql.ProjectionList{}
	for _, fType := range types.Values() {
		switch fType {
		case permsjobs.ColleaguesServiceGetColleagueTypesPermValueNote:
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

	if req.GetReason() == "" {
		return nil, errorsjobs.ErrReasonRequired
	}

	// Access Permission Check
	colleagueAccess, err := permsjobs.ColleaguesService.SetColleagueProps.AccessTyped.Get(
		s.ps,
		userInfo,
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
		colleagueAccess.StringList(),
		userInfo,
		targetUser.GetJob(),
		&usershort.UserShort{
			UserId:   targetUser.GetUserId(),
			Job:      targetUser.GetJob(),
			JobGrade: targetUser.GetJobGrade(),
		},
	) {
		return nil, errorsjobs.ErrFailedQuery
	}

	// Types Permission Check
	types, err := permsjobs.ColleaguesService.SetColleagueProps.TypesTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if userInfo.GetJobAdmin() {
		types.Set(
			permsjobs.ColleaguesServiceSetColleaguePropsTypesPermValueAbsenceDate,
			permsjobs.ColleaguesServiceSetColleaguePropsTypesPermValueNote,
			permsjobs.ColleaguesServiceSetColleaguePropsTypesPermValueLabels,
			permsjobs.ColleaguesServiceSetColleaguePropsTypesPermValueName,
		)
	}

	props, err := s.store.GetColleagueProps(
		ctx,
		s.db,
		targetUser.GetJob(),
		targetUser.GetUserId(),
		types.Strings(),
		userInfo.GetJobAdmin(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	reqProps := req.GetProps()

	if reqProps.GetAbsenceBegin() != nil && reqProps.GetAbsenceEnd() != nil {
		// Allow users to set their own absence date regardless of types perms check
		if userInfo.GetUserId() != targetUser.GetUserId() &&
			!types.Contains(
				permsjobs.ColleaguesServiceSetColleaguePropsTypesPermValueAbsenceDate,
			) &&
			!userInfo.GetJobAdmin() {
			return nil, errorsjobs.ErrPropsAbsenceDenied
		}

		jobProps, err := s.store.GetJobProps(ctx, s.db, userInfo.GetJob())
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
		if !types.Contains(permsjobs.ColleaguesServiceSetColleaguePropsTypesPermValueNote) &&
			!userInfo.GetJobAdmin() {
			return nil, errorsjobs.ErrPropsNoteDenied
		}
	}

	if reqProps.GetLabels() != nil {
		// Check if user is allowed to update labels
		if !types.Contains(permsjobs.ColleaguesServiceSetColleaguePropsTypesPermValueLabels) &&
			!userInfo.GetJobAdmin() {
			return nil, errorsjobs.ErrPropsAbsenceDenied
		}

		added, _ := utils.SlicesDifferenceFunc(
			props.GetLabels().GetList(),
			reqProps.GetLabels().GetList(),
			func(in *jobslabels.Label) int64 {
				return in.GetId()
			},
		)

		valid, err := s.store.ValidateLabels(ctx, s.db, userInfo.GetJob(), added)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		if !valid {
			return nil, errorsjobs.ErrPropsLabelsDenied
		}
	}

	if reqProps.NamePrefix != nil || reqProps.NameSuffix != nil {
		if !types.Contains(permsjobs.ColleaguesServiceSetColleaguePropsTypesPermValueName) &&
			!userInfo.GetJobAdmin() {
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

	activities, err := s.store.HandleColleaguePropsChanges(
		ctx,
		tx,
		props,
		reqProps,
		targetUser.GetJob(),
		&userInfo.UserId,
		req.GetReason(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	if err := colleaguesactivity.CreateColleagueActivity(ctx, tx, activities...); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	grpc_audit.SetAction(ctx, audit.EventAction_EVENT_ACTION_UPDATED)

	props, err = s.store.GetColleagueProps(
		ctx,
		s.db,
		targetUser.GetJob(),
		targetUser.GetUserId(),
		types.Strings(),
		userInfo.GetJobAdmin(),
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	userId := int64(targetUser.GetUserId())
	s.notifi.SendObjectEvent(ctx, &notificationsclientview.ObjectEvent{
		Type:      notificationsclientview.ObjectType_OBJECT_TYPE_JOBS_COLLEAGUE,
		Id:        &userId,
		EventType: notificationsclientview.ObjectEventType_OBJECT_EVENT_TYPE_UPDATED,

		UserId: &userInfo.UserId,
		Job:    &userInfo.Job,
	})

	return &pbjobs.SetColleaguePropsResponse{
		Props: props,
	}, nil
}

func (s *Server) getConditionForColleagueAccess(
	actTable *table.FivenetJobColleagueActivityTable,
	usersTable *table.FivenetUserTable,
	levels []permsjobs.ColleaguesServiceGetColleagueAccessPermValue,
	userInfo *userinfo.UserInfo,
) mysql.BoolExpression {
	condition := mysql.Bool(true)
	if userInfo.GetJobAdmin() {
		return condition
	}

	// If no levels set, assume "Own" as a safe default
	if len(levels) == 0 {
		return actTable.TargetUserID.EQ(mysql.Int32(userInfo.GetUserId()))
	}

	if slices.Contains(levels, permsjobs.ColleaguesServiceGetColleagueAccessPermValueAny) {
		return condition
	}
	if slices.Contains(levels, permsjobs.ColleaguesServiceGetColleagueAccessPermValueLowerRank) {
		return usersTable.ID.LT(mysql.Int32(userInfo.GetJobGrade()))
	}
	if slices.Contains(levels, permsjobs.ColleaguesServiceGetColleagueAccessPermValueSameRank) {
		return usersTable.ID.LT_EQ(mysql.Int32(userInfo.GetJobGrade()))
	}
	if slices.Contains(levels, permsjobs.ColleaguesServiceGetColleagueAccessPermValueOwn) {
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
	colleagueAccess, err := permsjobs.ColleaguesService.GetColleague.AccessTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	tColleagueActivity := tColleagueActivity.AS("colleague_activity")
	tTargetColleague := table.FivenetUser.AS("target_user")

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
				colleagueAccess.Values(),
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

		if !access.CheckIfHasOwnJobAccess(
			colleagueAccess.StringList(),
			userInfo,
			targetUser.GetJob(),
			&usershort.UserShort{
				UserId:   targetUser.GetUserId(),
				Job:      targetUser.GetJob(),
				JobGrade: targetUser.GetJobGrade(),
			},
		) {
			return nil, errorsjobs.ErrFailedQuery
		}

		condition = condition.AND(tColleagueActivity.TargetUserID.EQ(mysql.Int32(userId)))
	}

	// Types Field Permission Check
	types, err := permsjobs.ColleaguesService.ListColleagueActivity.TypesTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if types.Len() == 0 {
		if !userInfo.GetJobAdmin() {
			return resp, nil
		} else {
			types.Set(
				permsjobs.ColleaguesServiceListColleagueActivityTypesPermValueHIRED,
				permsjobs.ColleaguesServiceListColleagueActivityTypesPermValueFIRED,
				permsjobs.ColleaguesServiceListColleagueActivityTypesPermValuePROMOTED,
				permsjobs.ColleaguesServiceListColleagueActivityTypesPermValueDEMOTED,
				permsjobs.ColleaguesServiceListColleagueActivityTypesPermValueABSENCEDATE,
				permsjobs.ColleaguesServiceListColleagueActivityTypesPermValueNOTE,
				permsjobs.ColleaguesServiceListColleagueActivityTypesPermValueLABELS,
				permsjobs.ColleaguesServiceListColleagueActivityTypesPermValueNAME,
			)
		}
	}

	if len(req.GetActivityTypes()) > 0 {
		req.ActivityTypes = slices.DeleteFunc(
			req.GetActivityTypes(),
			func(t colleaguesactivity.ColleagueActivityType) bool {
				return !slices.ContainsFunc(
					types.Values(),
					func(s permsjobs.ColleaguesServiceListColleagueActivityTypesPermValue) bool {
						return strings.Contains(t.String(), "COLLEAGUE_ACTIVITY_TYPE_"+string(s))
					},
				)
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

	query := jobsstore.ListQuery{
		Job:   userInfo.GetJob(),
		Where: condition,
		Sort:  req.GetSort(),
	}
	count, err := s.store.CountColleagueActivity(ctx, s.db, query)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count, defaultPageSize)
	resp.Pagination = pag
	if count <= 0 {
		return resp, nil
	}
	activityQuery := query
	activityQuery.Offset = pag.GetOffset()
	activityQuery.Limit = limit

	resp.Activity, err = s.store.ListColleagueActivity(ctx, s.db, activityQuery)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
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
