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
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	errorsjobs "github.com/fivenet-app/fivenet/v2025/services/jobs/errors"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
)

const TimeclockMaxDays = (365 / 2) * 24 * time.Hour // Half a year

var tTimeClock = table.FivenetJobTimeclock.AS("timeclock_entry")

func (s *Server) ListTimeclock(
	ctx context.Context,
	req *pbjobs.ListTimeclockRequest,
) (*pbjobs.ListTimeclockResponse, error) {
	logging.InjectFields(ctx, logging.Fields{"fivenet.jobs.timeclock.user_ids", req.GetUserIds()})

	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tColleague := tables.User().AS("colleague")

	condition := tColleague.Job.EQ(jet.String(userInfo.GetJob()))
	statsCondition := tTimeClock.Job.EQ(jet.String(userInfo.GetJob()))

	// Field Permission Check
	fields, err := s.ps.AttrStringList(
		userInfo,
		permsjobs.TimeclockServicePerm,
		permsjobs.TimeclockServiceListTimeclockPerm,
		permsjobs.TimeclockServiceListTimeclockAccessPermField,
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	if !fields.Contains("All") {
		req.UserMode = jobs.TimeclockViewMode_TIMECLOCK_VIEW_MODE_SELF
	}

	if req.GetUserMode() <= jobs.TimeclockViewMode_TIMECLOCK_VIEW_MODE_SELF {
		condition = condition.AND(tTimeClock.UserID.EQ(jet.Int32(userInfo.GetUserId())))
		statsCondition = statsCondition.AND(tTimeClock.UserID.EQ(jet.Int32(userInfo.GetUserId())))
	} else if len(req.GetUserIds()) > 0 {
		ids := make([]jet.Expression, len(req.GetUserIds()))
		for i := range req.GetUserIds() {
			ids[i] = jet.Int32(req.GetUserIds()[i])
		}

		condition = condition.AND(
			tTimeClock.UserID.IN(ids...),
		)
		statsCondition = statsCondition.AND(
			tTimeClock.UserID.IN(ids...),
		)
	}

	if req.GetDate() != nil {
		if req.GetMode() <= jobs.TimeclockMode_TIMECLOCK_MODE_DAILY {
			if req.GetDate().GetEnd() == nil {
				req.Date.End = timestamp.Now()
			}

			condition = condition.AND(tTimeClock.Date.EQ(
				jet.DateT(req.GetDate().GetEnd().AsTime()),
			))
		} else if req.GetMode() == jobs.TimeclockMode_TIMECLOCK_MODE_WEEKLY {
			if req.GetDate().GetEnd() != nil {
				condition = condition.AND(jet.BoolExp(jet.Raw("YEARWEEK(`timeclock_entry`.`date`, 1) = YEARWEEK($date, 1)",
					jet.RawArgs{"$date": req.GetDate().GetEnd().AsTime()},
				)),
				)
			}
		} else {
			if req.GetDate().GetStart() != nil {
				condition = condition.AND(tTimeClock.Date.GT_EQ(
					jet.DateT(req.GetDate().GetStart().AsTime()),
				))
			}
			if req.GetDate().GetEnd() != nil {
				condition = condition.AND(tTimeClock.Date.LT_EQ(
					jet.DateT(req.GetDate().GetEnd().AsTime()),
				))
			}
		}

		// Make sure the provided dates are not "out of range"
		now := time.Now()
		if req.GetDate().GetStart() != nil &&
			now.Sub(req.GetDate().GetStart().AsTime()) >= TimeclockMaxDays {
			return nil, errorsjobs.ErrTimeclockOutOfRange
		}
		if req.GetDate().GetEnd() != nil &&
			now.Sub(req.GetDate().GetEnd().AsTime()) >= TimeclockMaxDays {
			return nil, errorsjobs.ErrTimeclockOutOfRange
		}
	}

	groupBys := []jet.GroupByClause{}
	if req.GetPerDay() {
		groupBys = append(groupBys, tTimeClock.Date, tTimeClock.UserID)
	} else {
		groupBys = append(groupBys, tTimeClock.UserID)
	}

	var countStmt jet.SelectStatement
	if req.GetUserMode() == jobs.TimeclockViewMode_TIMECLOCK_VIEW_MODE_ALL {
		var countCol jet.Projection
		if req.GetPerDay() {
			countCol = jet.RawString("COUNT(DISTINCT timeclock_entry.`date`, timeclock_entry.user_id)").
				AS("data_count.total")
		} else {
			countCol = jet.RawString("COUNT(DISTINCT timeclock_entry.`date`, timeclock_entry.user_id)").AS("data_count.total")
		}
		countStmt = tTimeClock.
			SELECT(countCol).
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

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, 30)
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
	if req.GetSort() != nil {
		var staticColumns []jet.OrderByClause
		var columns []jet.Column
		switch req.GetSort().GetColumn() {
		case "date":
			columns = append(columns, tTimeClock.Date)
		case rankColumn:
			staticColumns = append(staticColumns, tTimeClock.Date.DESC())
			columns = append(columns, tColleague.JobGrade)
		case nameColumn:
			staticColumns = append(staticColumns, tTimeClock.Date.DESC())
			columns = append(columns, tColleague.Firstname)
		case "time":
			fallthrough
		default:
			columns = append(columns, spentTimeColumn)
		}

		for _, column := range columns {
			if req.GetSort().GetDirection() == database.AscSortDirection {
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

	if req.GetMode() <= jobs.TimeclockMode_TIMECLOCK_MODE_DAILY {
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
							AND(tColleague.Job.EQ(jet.String(userInfo.GetJob()))),
					).
					LEFT_JOIN(tAvatar,
						tAvatar.ID.EQ(tUserProps.AvatarFileID),
					),
			).
			WHERE(condition).
			OFFSET(req.GetPagination().GetOffset()).
			GROUP_BY(groupBys...).
			ORDER_BY(orderBys...).
			LIMIT(limit)

		data := resp.GetDaily()
		if err := stmt.QueryContext(ctx, s.db, &data.Entries); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
			}
		}

		data.Date = req.GetDate().GetEnd()
		for i := range data.GetEntries() {
			if data.GetEntries()[i].GetUser() != nil {
				if data.GetEntries()[i].GetUser().GetJob() != userInfo.GetJob() {
					jobInfoFn(data.GetEntries()[i].GetUser())
				} else {
					s.enricher.EnrichJobInfo(data.GetEntries()[i].GetUser())
				}
			}
			data.Sum += int64(math.Round(float64(data.GetEntries()[i].GetSpentTime() * 60 * 60)))
		}

		resp.GetPagination().Update(len(data.GetEntries()))
	} else if req.GetMode() == jobs.TimeclockMode_TIMECLOCK_MODE_WEEKLY {
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
							AND(tColleague.Job.EQ(jet.String(userInfo.GetJob()))),
					).
					LEFT_JOIN(tAvatar,
						tAvatar.ID.EQ(tUserProps.AvatarFileID),
					),
			).
			WHERE(condition).
			OFFSET(req.GetPagination().GetOffset()).
			GROUP_BY(groupBys...).
			ORDER_BY(orderBys...).
			LIMIT(limit)

		data := resp.GetWeekly()
		if err := stmt.QueryContext(ctx, s.db, &data.Entries); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
			}
		}

		for i := range data.GetEntries() {
			if data.GetEntries()[i].GetUser() != nil {
				if data.GetEntries()[i].GetUser().GetJob() != userInfo.GetJob() {
					jobInfoFn(data.GetEntries()[i].GetUser())
				} else {
					s.enricher.EnrichJobInfo(data.GetEntries()[i].GetUser())
				}
			}
			data.Sum += int64(math.Round(float64(data.GetEntries()[i].GetSpentTime() * 60 * 60)))
		}

		resp.GetPagination().Update(len(data.GetEntries()))
	} else if req.GetMode() == jobs.TimeclockMode_TIMECLOCK_MODE_RANGE {
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
							AND(tColleague.Job.EQ(jet.String(userInfo.GetJob()))),
					).
					LEFT_JOIN(tAvatar,
						tAvatar.ID.EQ(tUserProps.AvatarFileID),
					),
			).
			WHERE(condition).
			OFFSET(req.GetPagination().GetOffset()).
			GROUP_BY(groupBys...).
			ORDER_BY(orderBys...).
			LIMIT(limit)

		data := resp.GetRange()
		if err := stmt.QueryContext(ctx, s.db, &data.Entries); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
			}
		}

		data.Date = req.GetDate().GetEnd()
		for i := range data.GetEntries() {
			if data.GetEntries()[i].GetUser() != nil {
				jobInfoFn(data.GetEntries()[i].GetUser())
			}

			data.Sum += int64(math.Round(float64(data.GetEntries()[i].GetSpentTime() * 60 * 60)))
		}

		resp.GetPagination().Update(len(data.GetEntries()))
	} else if req.GetMode() == jobs.TimeclockMode_TIMECLOCK_MODE_TIMELINE {
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
							AND(tColleague.Job.EQ(jet.String(userInfo.GetJob()))),
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

		data.Date = req.GetDate().GetEnd()
		for i := range data.GetEntries() {
			if data.GetEntries()[i].GetUser() != nil {
				jobInfoFn(data.GetEntries()[i].GetUser())
			}

			data.Sum += int64(math.Round(float64(data.GetEntries()[i].GetSpentTime() * 60 * 60)))
		}

		resp.GetPagination().Update(len(data.GetEntries()))
	}

	return resp, nil
}

func (s *Server) GetTimeclockStats(
	ctx context.Context,
	req *pbjobs.GetTimeclockStatsRequest,
) (*pbjobs.GetTimeclockStatsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	userId := userInfo.GetUserId()
	if req.UserId != nil && req.GetUserId() > 0 && req.GetUserId() != userInfo.GetUserId() {
		// Field Permission Check
		fields, err := s.ps.AttrStringList(
			userInfo,
			permsjobs.TimeclockServicePerm,
			permsjobs.TimeclockServiceListTimeclockPerm,
			permsjobs.TimeclockServiceListTimeclockAccessPermField,
		)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		if fields.Contains("All") {
			userId = req.GetUserId()
		}
	}

	condition := tTimeClock.Job.EQ(jet.String(userInfo.GetJob())).
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

func (s *Server) ListInactiveEmployees(
	ctx context.Context,
	req *pbjobs.ListInactiveEmployeesRequest,
) (*pbjobs.ListInactiveEmployeesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	tColleague := tables.User().AS("colleague")

	condition := jet.AND(
		tTimeClock.Job.EQ(jet.String(userInfo.GetJob())),
		tColleague.Job.EQ(jet.String(userInfo.GetJob())),
		jet.OR(
			jet.AND(
				tColleagueProps.AbsenceBegin.IS_NULL(),
				tColleagueProps.AbsenceEnd.IS_NULL(),
			),
			tColleagueProps.AbsenceBegin.LT_EQ(
				jet.DateExp(jet.CURRENT_DATE().SUB(jet.INTERVAL(req.GetDays(), jet.DAY))),
			),
			tColleagueProps.AbsenceEnd.LT_EQ(
				jet.DateExp(jet.CURRENT_DATE().SUB(jet.INTERVAL(req.GetDays(), jet.DAY))),
			),
		),
		tTimeClock.UserID.NOT_IN(
			tTimeClock.
				SELECT(
					tTimeClock.UserID,
				).
				FROM(tTimeClock).
				WHERE(jet.AND(
					tTimeClock.Job.EQ(jet.String(userInfo.GetJob())),
					tTimeClock.Date.GT_EQ(
						jet.DateExp(jet.CURRENT_DATE().SUB(jet.INTERVAL(req.GetDays(), jet.DAY))),
					),
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
						AND(tColleague.Job.EQ(jet.String(userInfo.GetJob()))),
				),
		).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count.Total, 20)
	resp := &pbjobs.ListInactiveEmployeesResponse{
		Pagination: pag,
		Colleagues: []*jobs.Colleague{},
	}
	if count.Total <= 0 {
		return resp, nil
	}

	// Convert proto sort to db sorting
	orderBys := []jet.OrderByClause{}
	if req.GetSort() != nil {
		var columns []jet.Column
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
						AND(tColleague.Job.EQ(jet.String(userInfo.GetJob()))),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		WHERE(condition).
		ORDER_BY(orderBys...).
		GROUP_BY(tTimeClock.UserID).
		OFFSET(req.GetPagination().GetOffset()).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Colleagues); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetColleagues() {
		jobInfoFn(resp.GetColleagues()[i])
	}

	resp.GetPagination().Update(len(resp.GetColleagues()))

	return resp, nil
}
