package jobs

import (
	"context"
	"errors"
	"math"
	"time"

	database "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/database"
	jobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	pbjobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/jobs"
	permsjobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsjobs "github.com/fivenet-app/fivenet/v2025/services/jobs/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const TimeclockMaxDays = (365 / 2) * 24 * time.Hour // Half a year

var tTimeClock = table.FivenetJobTimeclock.AS("timeclock_entry")

func (s *Server) ListTimeclock(ctx context.Context, req *pbjobs.ListTimeclockRequest) (*pbjobs.ListTimeclockResponse, error) {
	trace.SpanFromContext(ctx).SetAttributes(attribute.IntSlice("fivenet.jobs.timeclock.user_ids", utils.SliceInt32ToInt(req.UserIds)))

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tColleague := tables.User().AS("colleague")

	condition := tColleague.Job.EQ(jet.String(userInfo.Job))
	statsCondition := tTimeClock.Job.EQ(jet.String(userInfo.Job))

	// Field Permission Check
	fields, err := s.ps.AttrStringList(userInfo, permsjobs.TimeclockServicePerm, permsjobs.TimeclockServiceListTimeclockPerm, permsjobs.TimeclockServiceListTimeclockAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	if !fields.Contains("All") {
		req.UserMode = jobs.TimeclockViewMode_TIMECLOCK_VIEW_MODE_SELF
	}

	if req.UserMode <= jobs.TimeclockViewMode_TIMECLOCK_VIEW_MODE_SELF {
		condition = condition.AND(tTimeClock.UserID.EQ(jet.Int32(userInfo.UserId)))
		statsCondition = statsCondition.AND(tTimeClock.UserID.EQ(jet.Int32(userInfo.UserId)))
	} else {
		if len(req.UserIds) > 0 {
			ids := make([]jet.Expression, len(req.UserIds))
			for i := range req.UserIds {
				ids[i] = jet.Int32(req.UserIds[i])
			}

			condition = condition.AND(
				tTimeClock.UserID.IN(ids...),
			)
			statsCondition = statsCondition.AND(
				tTimeClock.UserID.IN(ids...),
			)
		}
	}

	if req.Date != nil {
		if req.Mode <= jobs.TimeclockMode_TIMECLOCK_MODE_DAILY {
			if req.Date.End == nil {
				req.Date.End = timestamp.Now()
			}

			condition = condition.AND(tTimeClock.Date.EQ(
				jet.DateT(req.Date.End.AsTime()),
			))
		} else if req.Mode == jobs.TimeclockMode_TIMECLOCK_MODE_WEEKLY {
			if req.Date.End != nil {
				condition = condition.AND(jet.BoolExp(jet.Raw("YEARWEEK(`timeclock_entry`.`date`, 1) = YEARWEEK($date, 1)",
					jet.RawArgs{"$date": req.Date.End.AsTime()},
				)),
				)
			}
		} else {
			if req.Date.Start != nil {
				condition = condition.AND(tTimeClock.Date.GT_EQ(
					jet.DateT(req.Date.Start.AsTime()),
				))
			}
			if req.Date.End != nil {
				condition = condition.AND(tTimeClock.Date.LT_EQ(
					jet.DateT(req.Date.End.AsTime()),
				))
			}
		}

		// Make sure the provided dates are not "out of range"
		now := time.Now()
		if req.Date.Start != nil && now.Sub(req.Date.Start.AsTime()) >= TimeclockMaxDays {
			return nil, errorsjobs.ErrTimeclockOutOfRange
		}
		if req.Date.End != nil && now.Sub(req.Date.End.AsTime()) >= TimeclockMaxDays {
			return nil, errorsjobs.ErrTimeclockOutOfRange
		}
	}

	groupBys := []jet.GroupByClause{}
	if req.PerDay {
		groupBys = append(groupBys, tTimeClock.Date, tTimeClock.UserID)
	} else {
		groupBys = append(groupBys, tTimeClock.UserID)
	}

	var countStmt jet.SelectStatement
	if req.UserMode == jobs.TimeclockViewMode_TIMECLOCK_VIEW_MODE_ALL {
		if req.PerDay {
			countStmt = tTimeClock.
				SELECT(jet.RawString("COUNT(DISTINCT timeclock_entry.`date`, timeclock_entry.user_id)").AS("data_count.total")).
				FROM(
					tTimeClock.
						INNER_JOIN(tColleague,
							tColleague.ID.EQ(tTimeClock.UserID),
						),
				).
				WHERE(condition)
		} else {
			countStmt = tTimeClock.
				SELECT(jet.RawString("COUNT(DISTINCT timeclock_entry.`date`, timeclock_entry.user_id)").AS("data_count.total")).
				FROM(
					tTimeClock.
						INNER_JOIN(tColleague,
							tColleague.ID.EQ(tTimeClock.UserID),
						),
				).
				WHERE(condition)
		}
	} else {
		countStmt = tTimeClock.
			SELECT(jet.RawString("COUNT(DISTINCT timeclock_entry.`date`, timeclock_entry.user_id)").AS("data_count.total")).
			FROM(
				tTimeClock.
					INNER_JOIN(tColleague,
						tColleague.ID.EQ(tTimeClock.UserID),
					),
			).
			WHERE(condition)
	}

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.Total, 30)
	resp := &pbjobs.ListTimeclockResponse{
		Pagination: pag,
	}

	resp.Stats, err = s.getTimeclockStats(ctx, statsCondition)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	resp.StatsWeekly, err = s.getTimeclockWeeklyStats(ctx, statsCondition)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	if count.Total <= 0 {
		return resp, nil
	}

	spentTimeColumn := jet.StringColumn("timeclock_entry.spent_time")
	// Convert proto sort to db sorting
	orderBys := []jet.OrderByClause{}
	if req.Sort != nil {
		var staticColumns []jet.OrderByClause
		var columns []jet.Column
		switch req.Sort.Column {
		case "date":
			columns = append(columns, tTimeClock.Date)
		case "rank":
			staticColumns = append(staticColumns, tTimeClock.Date.DESC())
			columns = append(columns, tColleague.JobGrade)
		case "name":
			staticColumns = append(staticColumns, tTimeClock.Date.DESC())
			columns = append(columns, tColleague.Firstname)
		case "time":
			fallthrough
		default:
			columns = append(columns, spentTimeColumn)
		}

		for _, column := range columns {
			if req.Sort.Direction == database.AscSortDirection {
				if column == spentTimeColumn {
					orderBys = append(orderBys, column.ASC())
				} else {
					orderBys = append(orderBys, column.ASC(), spentTimeColumn.DESC())
				}
			} else {
				if column == spentTimeColumn {
					orderBys = append(orderBys, column.DESC())
				} else {
					orderBys = append(orderBys, column.DESC(), spentTimeColumn.DESC())
				}
			}
		}
		orderBys = append(staticColumns, orderBys...)
	} else {
		orderBys = append(orderBys,
			tTimeClock.Date.DESC(),
			spentTimeColumn.DESC(),
		)
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)

	if req.Mode <= jobs.TimeclockMode_TIMECLOCK_MODE_DAILY {
		resp.Entries = &pbjobs.ListTimeclockResponse_Daily{
			Daily: &pbjobs.TimeclockDay{},
		}

		stmt := tTimeClock.
			SELECT(
				tTimeClock.UserID,
				tTimeClock.StartTime,
				tTimeClock.EndTime,
				jet.SUM(tTimeClock.SpentTime).AS("timeclock_entry.spent_time"),
				tColleague.ID,
				tColleague.Job,
				tColleague.JobGrade,
				tColleague.Firstname,
				tColleague.Lastname,
				tColleague.Dateofbirth,
				tColleague.PhoneNumber,
				tUserProps.AvatarFileID.AS("colleague.avatar_file_id"),
				tAvatar.FilePath.AS("colleague.avatar"),
				tColleagueProps.UserID,
				tColleagueProps.Job,
				tColleagueProps.AbsenceBegin,
				tColleagueProps.AbsenceEnd,
				tColleagueProps.NamePrefix,
				tColleagueProps.NameSuffix,
			).
			FROM(
				tTimeClock.
					INNER_JOIN(tColleague,
						tColleague.ID.EQ(tTimeClock.UserID),
					).
					LEFT_JOIN(tUserProps,
						tUserProps.UserID.EQ(tColleague.ID),
					).
					LEFT_JOIN(tColleagueProps,
						tColleagueProps.UserID.EQ(tColleague.ID).
							AND(tColleague.Job.EQ(jet.String(userInfo.Job))),
					).
					LEFT_JOIN(tAvatar,
						tAvatar.ID.EQ(tUserProps.AvatarFileID),
					),
			).
			WHERE(condition).
			OFFSET(req.Pagination.Offset).
			GROUP_BY(groupBys...).
			ORDER_BY(orderBys...).
			LIMIT(limit)

		data := resp.GetDaily()
		if err := stmt.QueryContext(ctx, s.db, &data.Entries); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
			}
		}

		data.Date = req.Date.End
		for i := range data.Entries {
			if data.Entries[i].User != nil {
				if data.Entries[i].User.Job != userInfo.Job {
					jobInfoFn(data.Entries[i].User)
				} else {
					s.enricher.EnrichJobInfo(data.Entries[i].User)
				}
			}
			data.Sum += int64(math.Round(float64(data.Entries[i].SpentTime * 60 * 60)))
		}

		resp.Pagination.Update(len(data.Entries))
	} else if req.Mode == jobs.TimeclockMode_TIMECLOCK_MODE_WEEKLY {
		resp.Entries = &pbjobs.ListTimeclockResponse_Weekly{
			Weekly: &pbjobs.TimeclockWeekly{},
		}

		stmt := tTimeClock.
			SELECT(
				tTimeClock.UserID,
				tTimeClock.Date,
				tTimeClock.StartTime,
				tTimeClock.EndTime,
				jet.SUM(tTimeClock.SpentTime).AS("timeclock_entry.spent_time"),
				tColleague.ID,
				tColleague.Job,
				tColleague.JobGrade,
				tColleague.Firstname,
				tColleague.Lastname,
				tColleague.Dateofbirth,
				tColleague.PhoneNumber,
				tUserProps.AvatarFileID.AS("colleague.avatar_file_id"),
				tAvatar.FilePath.AS("colleague.avatar"),
				tColleagueProps.UserID,
				tColleagueProps.Job,
				tColleagueProps.AbsenceBegin,
				tColleagueProps.AbsenceEnd,
				tColleagueProps.NamePrefix,
				tColleagueProps.NameSuffix,
			).
			FROM(
				tTimeClock.
					INNER_JOIN(tColleague,
						tColleague.ID.EQ(tTimeClock.UserID),
					).
					LEFT_JOIN(tUserProps,
						tUserProps.UserID.EQ(tColleague.ID),
					).
					LEFT_JOIN(tColleagueProps,
						tColleagueProps.UserID.EQ(tColleague.ID).
							AND(tColleague.Job.EQ(jet.String(userInfo.Job))),
					).
					LEFT_JOIN(tAvatar,
						tAvatar.ID.EQ(tUserProps.AvatarFileID),
					),
			).
			WHERE(condition).
			OFFSET(req.Pagination.Offset).
			GROUP_BY(groupBys...).
			ORDER_BY(orderBys...).
			LIMIT(limit)

		data := resp.GetWeekly()
		if err := stmt.QueryContext(ctx, s.db, &data.Entries); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
			}
		}

		for i := range data.Entries {
			if data.Entries[i].User != nil {
				if data.Entries[i].User.Job != userInfo.Job {
					jobInfoFn(data.Entries[i].User)
				} else {
					s.enricher.EnrichJobInfo(data.Entries[i].User)
				}
			}
			data.Sum += int64(math.Round(float64(data.Entries[i].SpentTime * 60 * 60)))
		}

		resp.Pagination.Update(len(data.Entries))
	} else if req.Mode == jobs.TimeclockMode_TIMECLOCK_MODE_RANGE {
		resp.Entries = &pbjobs.ListTimeclockResponse_Range{
			Range: &pbjobs.TimeclockRange{},
		}

		stmt := tTimeClock.
			SELECT(
				tTimeClock.UserID,
				tTimeClock.Date,
				tTimeClock.StartTime,
				tTimeClock.EndTime,
				jet.SUM(tTimeClock.SpentTime).AS("timeclock_entry.spent_time"),
				tColleague.ID,
				tColleague.Job,
				tColleague.JobGrade,
				tColleague.Firstname,
				tColleague.Lastname,
				tColleague.Dateofbirth,
				tColleague.PhoneNumber,
				tUserProps.AvatarFileID.AS("colleague.avatar_file_id"),
				tAvatar.FilePath.AS("colleague.avatar"),
				tColleagueProps.UserID,
				tColleagueProps.Job,
				tColleagueProps.AbsenceBegin,
				tColleagueProps.AbsenceEnd,
				tColleagueProps.NamePrefix,
				tColleagueProps.NameSuffix,
			).
			FROM(
				tTimeClock.
					INNER_JOIN(tColleague,
						tColleague.ID.EQ(tTimeClock.UserID),
					).
					LEFT_JOIN(tUserProps,
						tUserProps.UserID.EQ(tColleague.ID),
					).
					LEFT_JOIN(tColleagueProps,
						tColleagueProps.UserID.EQ(tColleague.ID).
							AND(tColleague.Job.EQ(jet.String(userInfo.Job))),
					).
					LEFT_JOIN(tAvatar,
						tAvatar.ID.EQ(tUserProps.AvatarFileID),
					),
			).
			WHERE(condition).
			OFFSET(req.Pagination.Offset).
			GROUP_BY(groupBys...).
			ORDER_BY(orderBys...).
			LIMIT(limit)

		data := resp.GetRange()
		if err := stmt.QueryContext(ctx, s.db, &data.Entries); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
			}
		}

		data.Date = req.Date.End
		for i := range data.Entries {
			if data.Entries[i].User != nil {
				jobInfoFn(data.Entries[i].User)
			}

			data.Sum += int64(math.Round(float64(data.Entries[i].SpentTime * 60 * 60)))
		}

		resp.Pagination.Update(len(data.Entries))
	} else if req.Mode == jobs.TimeclockMode_TIMECLOCK_MODE_TIMELINE {
		resp.Entries = &pbjobs.ListTimeclockResponse_Range{
			Range: &pbjobs.TimeclockRange{},
		}

		stmt := tTimeClock.
			SELECT(
				tTimeClock.UserID,
				tTimeClock.Date,
				tTimeClock.StartTime,
				tTimeClock.EndTime,
				tTimeClock.SpentTime,
				tColleague.ID,
				tColleague.Job,
				tColleague.JobGrade,
				tColleague.Firstname,
				tColleague.Lastname,
				tColleague.Dateofbirth,
				tColleague.PhoneNumber,
				tUserProps.AvatarFileID.AS("colleague.avatar_file_id"),
				tAvatar.FilePath.AS("colleague.avatar"),
				tColleagueProps.UserID,
				tColleagueProps.Job,
				tColleagueProps.AbsenceBegin,
				tColleagueProps.AbsenceEnd,
				tColleagueProps.NamePrefix,
				tColleagueProps.NameSuffix,
			).
			FROM(
				tTimeClock.
					INNER_JOIN(tColleague,
						tColleague.ID.EQ(tTimeClock.UserID),
					).
					LEFT_JOIN(tUserProps,
						tUserProps.UserID.EQ(tColleague.ID),
					).
					LEFT_JOIN(tColleagueProps,
						tColleagueProps.UserID.EQ(tColleague.ID).
							AND(tColleague.Job.EQ(jet.String(userInfo.Job))),
					).
					LEFT_JOIN(tAvatar,
						tAvatar.ID.EQ(tUserProps.AvatarFileID),
					),
			).
			WHERE(condition).
			ORDER_BY(tTimeClock.Date.DESC(), tTimeClock.UserID.DESC())

		data := resp.GetRange()
		if err := stmt.QueryContext(ctx, s.db, &data.Entries); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
			}
		}

		data.Date = req.Date.End
		for i := range data.Entries {
			if data.Entries[i].User != nil {
				jobInfoFn(data.Entries[i].User)
			}

			data.Sum += int64(math.Round(float64(data.Entries[i].SpentTime * 60 * 60)))
		}

		resp.Pagination.Update(len(data.Entries))
	}

	return resp, nil
}

func (s *Server) GetTimeclockStats(ctx context.Context, req *pbjobs.GetTimeclockStatsRequest) (*pbjobs.GetTimeclockStatsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	userId := userInfo.UserId
	if req.UserId != nil && *req.UserId > 0 && *req.UserId != userInfo.UserId {
		// Field Permission Check
		fields, err := s.ps.AttrStringList(userInfo, permsjobs.TimeclockServicePerm, permsjobs.TimeclockServiceListTimeclockPerm, permsjobs.TimeclockServiceListTimeclockAccessPermField)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		if fields.Contains("All") {
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

	return &pbjobs.GetTimeclockStatsResponse{
		Stats:  stats,
		Weekly: weekly,
	}, nil
}

func (s *Server) ListInactiveEmployees(ctx context.Context, req *pbjobs.ListInactiveEmployeesRequest) (*pbjobs.ListInactiveEmployeesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tColleague := tables.User().AS("colleague")

	condition := jet.AND(
		tTimeClock.Job.EQ(jet.String(userInfo.Job)),
		tColleague.Job.EQ(jet.String(userInfo.Job)),
		jet.OR(
			jet.AND(
				tColleagueProps.AbsenceBegin.IS_NULL(),
				tColleagueProps.AbsenceEnd.IS_NULL(),
			),
			tColleagueProps.AbsenceBegin.LT_EQ(
				jet.DateExp(jet.CURRENT_DATE().SUB(jet.INTERVAL(req.Days, jet.DAY))),
			),
			tColleagueProps.AbsenceEnd.LT_EQ(
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
			jet.COUNT(jet.DISTINCT(tTimeClock.UserID)).AS("data_count.total"),
		).
		FROM(
			tTimeClock.
				INNER_JOIN(tColleague,
					tColleague.ID.EQ(tTimeClock.UserID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tTimeClock.UserID),
				).
				LEFT_JOIN(tColleagueProps,
					tColleagueProps.UserID.EQ(tTimeClock.UserID).
						AND(tColleague.Job.EQ(jet.String(userInfo.Job))),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(count.Total, 20)
	resp := &pbjobs.ListInactiveEmployeesResponse{
		Pagination: pag,
		Colleagues: []*jobs.Colleague{},
	}
	if count.Total <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []jet.OrderByClause{}
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
		orderBys = append(orderBys, tColleague.JobGrade.ASC())
	}

	stmt := tTimeClock.
		SELECT(
			tTimeClock.UserID,
			tColleague.ID,
			tColleague.Job,
			tColleague.JobGrade,
			tColleague.Firstname,
			tColleague.Lastname,
			tColleague.Dateofbirth,
			tColleague.PhoneNumber,
			tUserProps.AvatarFileID.AS("colleague.avatar_file_id"),
			tAvatar.FilePath.AS("colleague.avatar"),
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.AbsenceBegin,
			tColleagueProps.AbsenceEnd,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
		).
		FROM(
			tTimeClock.
				INNER_JOIN(tColleague,
					tColleague.ID.EQ(tTimeClock.UserID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tTimeClock.UserID),
				).
				LEFT_JOIN(tColleagueProps,
					tColleagueProps.UserID.EQ(tTimeClock.UserID).
						AND(tColleague.Job.EQ(jet.String(userInfo.Job))),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
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

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.Colleagues {
		jobInfoFn(resp.Colleagues[i])
	}

	resp.Pagination.Update(len(resp.Colleagues))

	return resp, nil
}
