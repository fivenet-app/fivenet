package jobs

import (
	"context"
	"errors"
	"slices"

	database "github.com/fivenet-app/fivenet/gen/go/proto/resources/common/database"
	jobs "github.com/fivenet-app/fivenet/gen/go/proto/resources/jobs"
	errorsjobs "github.com/fivenet-app/fivenet/gen/go/proto/services/jobs/errors"
	permsjobs "github.com/fivenet-app/fivenet/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/pkg/perms"
	"github.com/fivenet-app/fivenet/pkg/utils"
	timeutils "github.com/fivenet-app/fivenet/pkg/utils/time"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

var tTimeClock = table.FivenetJobsTimeclock.AS("timeclock_entry")

func (s *Server) ListTimeclock(ctx context.Context, req *ListTimeclockRequest) (*ListTimeclockResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.IntSlice("fivenet.jobs.timeclock.user_ids", utils.SliceInt32ToInt(req.UserIds)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := jet.AND(tTimeClock.Job.EQ(jet.String(userInfo.Job)))
	statsCondition := jet.AND(tTimeClock.Job.EQ(jet.String(userInfo.Job)))

	// Field Permission Check
	fieldsAttr, err := s.ps.Attr(userInfo, permsjobs.JobsTimeclockServicePerm, permsjobs.JobsTimeclockServiceListTimeclockPerm, permsjobs.JobsTimeclockServiceListTimeclockAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}

	if len(fields) == 0 || !slices.Contains(fields, "All") {
		condition = condition.AND(tTimeClock.UserID.EQ(jet.Int32(userInfo.UserId)))
		statsCondition = statsCondition.AND(tTimeClock.UserID.EQ(jet.Int32(userInfo.UserId)))
	}

	if len(req.UserIds) > 0 {
		ids := make([]jet.Expression, len(req.UserIds))
		for i := 0; i < len(req.UserIds); i++ {
			ids[i] = jet.Int32(req.UserIds[i])
		}

		condition = condition.AND(
			tTimeClock.UserID.IN(ids...),
		)
		statsCondition = statsCondition.AND(
			tTimeClock.UserID.IN(ids...),
		)
	}

	if req.From != nil {
		condition = condition.AND(tTimeClock.Date.GT_EQ(
			jet.DateT(timeutils.TruncateToDay(req.From.AsTime())),
		))
	}
	if req.To != nil {
		condition = condition.AND(tTimeClock.Date.LT_EQ(
			jet.DateT(timeutils.TruncateToNight(req.To.AsTime())),
		))
	}

	countStmt := tTimeClock.
		SELECT(jet.COUNT(tTimeClock.Date).AS("datacount.totalcount")).
		FROM(tTimeClock).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.TotalCount, 30)
	resp := &ListTimeclockResponse{
		Pagination: pag,
	}

	resp.Stats, err = s.getTimeclockStats(ctx, statsCondition)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	resp.Weekly, err = s.getTimeclockWeeklyStats(ctx, statsCondition)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	if count.TotalCount <= 0 {
		return resp, nil
	}

	tUser := tUser.AS("colleague")
	// Convert proto sort to db sorting
	orderBys := []jet.OrderByClause{
		tTimeClock.Date.DESC(),
	}
	if req.Sort != nil {
		var column jet.Column
		switch req.Sort.Column {
		case "rank":
			column = tUser.JobGrade
		case "name":
			column = tUser.Firstname
		default:
			column = tTimeClock.SpentTime
		}

		if req.Sort.Direction == database.AscSortDirection {
			orderBys = append(orderBys, column.ASC(), tTimeClock.SpentTime.DESC())
		} else {
			orderBys = append(orderBys, column.DESC(), tTimeClock.SpentTime.DESC())
		}
	} else {
		orderBys = append(orderBys, tTimeClock.SpentTime.DESC())
	}

	stmt := tTimeClock.
		SELECT(
			tTimeClock.Job,
			tTimeClock.Date,
			tTimeClock.UserID,
			tTimeClock.StartTime,
			tTimeClock.EndTime,
			tTimeClock.SpentTime,
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
		FROM(
			tTimeClock.
				INNER_JOIN(tUser,
					tUser.ID.EQ(tTimeClock.UserID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tUser.ID),
				).
				LEFT_JOIN(tJobsUserProps,
					tJobsUserProps.UserID.EQ(tUser.ID).
						AND(tUser.Job.EQ(jet.String(userInfo.Job))),
				),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(orderBys...).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Entries); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Entries); i++ {
		if resp.Entries[i].User != nil {
			jobInfoFn(resp.Entries[i].User)
		}
	}

	resp.Pagination.Update(len(resp.Entries))

	return resp, nil
}

func (s *Server) GetTimeclockStats(ctx context.Context, req *GetTimeclockStatsRequest) (*GetTimeclockStatsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	userId := userInfo.UserId
	if req.UserId != nil && *req.UserId > 0 && *req.UserId != userInfo.UserId {
		// Field Permission Check
		fieldsAttr, err := s.ps.Attr(userInfo, permsjobs.JobsTimeclockServicePerm, permsjobs.JobsTimeclockServiceListTimeclockPerm, permsjobs.JobsTimeclockServiceListTimeclockAccessPermField)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		var fields perms.StringList
		if fieldsAttr != nil {
			fields = fieldsAttr.([]string)
		}
		if slices.Contains(fields, "All") {
			userId = *req.UserId
		}
	}

	condition := tTimeClock.Job.EQ(jet.String(userInfo.Job)).
		AND(tTimeClock.UserID.EQ(jet.Int32(userId)))

	stats, err := s.getTimeclockStats(ctx, condition)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	weekly, err := s.getTimeclockWeeklyStats(ctx, condition)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	return &GetTimeclockStatsResponse{
		Stats:  stats,
		Weekly: weekly,
	}, nil
}

func (s *Server) ListInactiveEmployees(ctx context.Context, req *ListInactiveEmployeesRequest) (*ListInactiveEmployeesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tUser := tUser.AS("colleague")

	condition := jet.AND(
		tTimeClock.Job.EQ(jet.String(userInfo.Job)),
		tUser.Job.EQ(jet.String(userInfo.Job)),
		jet.OR(
			jet.AND(
				tJobsUserProps.AbsenceBegin.IS_NULL(),
				tJobsUserProps.AbsenceEnd.IS_NULL(),
			),
			tJobsUserProps.AbsenceBegin.LT_EQ(
				jet.DateExp(jet.CURRENT_DATE().SUB(jet.INTERVAL(req.Days, jet.DAY))),
			),
			tJobsUserProps.AbsenceEnd.LT_EQ(
				jet.DateExp(jet.CURRENT_DATE().SUB(jet.INTERVAL(req.Days, jet.DAY))),
			),
		),
		tTimeClock.UserID.NOT_IN(
			tTimeClock.
				SELECT(
					tTimeClock.UserID,
				).
				FROM(tTimeClock).
				WHERE(jet.AND(
					tTimeClock.Job.EQ(jet.String(userInfo.Job)),
					tTimeClock.Date.GT_EQ(jet.DateExp(jet.CURRENT_DATE().SUB(jet.INTERVAL(req.Days, jet.DAY)))),
				)).
				GROUP_BY(tTimeClock.UserID),
		),
	)

	countStmt := tTimeClock.
		SELECT(
			jet.COUNT(jet.DISTINCT(tTimeClock.UserID)).AS("datacount.totalcount"),
		).
		FROM(
			tTimeClock.
				INNER_JOIN(tUser,
					tUser.ID.EQ(tTimeClock.UserID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tTimeClock.UserID),
				).
				LEFT_JOIN(tJobsUserProps,
					tJobsUserProps.UserID.EQ(tTimeClock.UserID).
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
	resp := &ListInactiveEmployeesResponse{
		Pagination: pag,
		Colleagues: []*jobs.Colleague{},
	}
	if count.TotalCount <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []jet.OrderByClause{}
	if req.Sort != nil {
		var column jet.Column
		switch req.Sort.Column {
		case "name":
			column = nil
		case "rank":
			fallthrough
		default:
			column = tUser.JobGrade
		}

		if column != nil && req.Sort.Direction == database.AscSortDirection {
			orderBys = append(orderBys, column.ASC())
		} else {
			orderBys = append(orderBys, column.DESC())
		}
	} else {
		orderBys = append(orderBys,
			tConduct.ID.DESC(),
		)
	}
	orderBys = append(orderBys,
		tUser.Firstname.ASC(),
		tUser.Lastname.ASC(),
	)

	stmt := tTimeClock.
		SELECT(
			tTimeClock.UserID,
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
		FROM(
			tTimeClock.
				INNER_JOIN(tUser,
					tUser.ID.EQ(tTimeClock.UserID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tTimeClock.UserID),
				).
				LEFT_JOIN(tJobsUserProps,
					tJobsUserProps.UserID.EQ(tTimeClock.UserID).
						AND(tUser.Job.EQ(jet.String(userInfo.Job))),
				),
		).
		WHERE(condition).
		ORDER_BY(orderBys...).
		GROUP_BY(tTimeClock.UserID).
		OFFSET(req.Pagination.Offset).
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
