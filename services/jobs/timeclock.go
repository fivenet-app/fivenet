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
	"github.com/go-jet/jet/v2/mysql"
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

	condition := tColleague.Job.EQ(mysql.String(userInfo.GetJob()))
	statsCondition := tTimeClock.Job.EQ(mysql.String(userInfo.GetJob()))

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
		condition = condition.AND(tTimeClock.UserID.EQ(mysql.Int32(userInfo.GetUserId())))
		statsCondition = statsCondition.AND(tTimeClock.UserID.EQ(mysql.Int32(userInfo.GetUserId())))
	} else if len(req.GetUserIds()) > 0 {
		ids := make([]mysql.Expression, len(req.GetUserIds()))
		for i := range req.GetUserIds() {
			ids[i] = mysql.Int32(req.GetUserIds()[i])
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
				mysql.DateT(req.GetDate().GetEnd().AsTime()),
			))
		} else if req.GetMode() == jobs.TimeclockMode_TIMECLOCK_MODE_WEEKLY {
			if req.GetDate().GetEnd() != nil {
				condition = condition.AND(mysql.BoolExp(mysql.Raw("YEARWEEK(`timeclock_entry`.`date`, 1) = YEARWEEK($date, 1)",
					mysql.RawArgs{"$date": req.GetDate().GetEnd().AsTime()},
				)),
				)
			}
		} else {
			if req.GetDate().GetStart() != nil {
				condition = condition.AND(tTimeClock.Date.GT_EQ(
					mysql.DateT(req.GetDate().GetStart().AsTime()),
				))
			}
			if req.GetDate().GetEnd() != nil {
				condition = condition.AND(tTimeClock.Date.LT_EQ(
					mysql.DateT(req.GetDate().GetEnd().AsTime()),
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

	groupBys := []mysql.GroupByClause{}
	if req.GetPerDay() {
		groupBys = append(groupBys, tTimeClock.Date, tTimeClock.UserID)
	} else {
		groupBys = append(groupBys, tTimeClock.UserID)
	}

	// User mode doesn't change the count query
	countStmt := tTimeClock.
		SELECT(
			mysql.COUNT(mysql.RawString("DISTINCT `timeclock_entry`.`date`, `timeclock_entry`.`user_id`")).
				AS("data_count.total"),
		).
		FROM(
			tTimeClock.
				INNER_JOIN(tColleague,
					tColleague.ID.EQ(tTimeClock.UserID),
				),
		).
		WHERE(condition)

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

	spentTimeColumn := mysql.StringColumn("timeclock_entry.spent_time")
	// Convert proto sort to db sorting
	orderBys := []mysql.OrderByClause{}
	if req.GetSort() != nil {
		var staticColumns []mysql.OrderByClause
		var columns []mysql.Column
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

	agg := tTimeClock.
		SELECT(
			tTimeClock.UserID.AS("agg.user_id"),
			tTimeClock.Date.AS("agg.date"),
			mysql.MIN(tTimeClock.StartTime).AS("agg.start_time"),
			mysql.MAX(tTimeClock.EndTime).AS("agg.end_time"),
			mysql.SUM(tTimeClock.SpentTime).AS("agg.spent_time"),
		).
		FROM(
			tTimeClock.INNER_JOIN(
				tColleague,
				tColleague.ID.EQ(tTimeClock.UserID),
			),
		).
		WHERE(condition).
		GROUP_BY(groupBys...).
		AsTable("agg")

	stmt := agg.
		SELECT(
			mysql.IntegerColumn("agg.user_id").AS("timeclock_entry.user_id"),
			mysql.DateColumn("agg.date").AS("timeclock_entry.date"),
			mysql.DateTimeColumn("agg.start_time").AS("timeclock_entry.start_time"),
			mysql.DateTimeColumn("agg.end_time").AS("timeclock_entry.end_time"),
			mysql.FloatColumn("agg.spent_time").AS("timeclock_entry.spent_time"),

			tColleague.ID,
			tColleague.Job,
			tColleague.JobGrade,
			tColleague.Firstname,
			tColleague.Lastname,
			tColleague.Dateofbirth,
			tColleague.PhoneNumber,
			tUserProps.AvatarFileID.AS("colleague.profile_picture_file_id"),
			tAvatar.FilePath.AS("colleague.profile_picture"),
			tColleagueProps.UserID,
			tColleagueProps.Job,
			tColleagueProps.AbsenceBegin,
			tColleagueProps.AbsenceEnd,
			tColleagueProps.NamePrefix,
			tColleagueProps.NameSuffix,
		).
		FROM(
			agg.
				INNER_JOIN(tColleague,
					tColleague.ID.EQ(mysql.IntegerColumn("agg.user_id")),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tColleague.ID),
				).
				LEFT_JOIN(tColleagueProps,
					mysql.AND(
						tColleagueProps.UserID.EQ(tColleague.ID),
						tColleague.Job.EQ(mysql.String(userInfo.GetJob())),
					),
				).
				LEFT_JOIN(tAvatar,
					tAvatar.ID.EQ(tUserProps.AvatarFileID),
				),
		).
		ORDER_BY(orderBys...).
		OFFSET(req.GetPagination().GetOffset()).
		LIMIT(limit)

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)

	switch req.GetMode() {
	case jobs.TimeclockMode_TIMECLOCK_MODE_UNSPECIFIED:
		fallthrough

	case jobs.TimeclockMode_TIMECLOCK_MODE_DAILY:
		resp.Entries = &pbjobs.ListTimeclockResponse_Daily{
			Daily: &pbjobs.TimeclockDay{},
		}

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

	case jobs.TimeclockMode_TIMECLOCK_MODE_WEEKLY:
		resp.Entries = &pbjobs.ListTimeclockResponse_Weekly{
			Weekly: &pbjobs.TimeclockWeekly{},
		}

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

	case jobs.TimeclockMode_TIMECLOCK_MODE_RANGE:
		resp.Entries = &pbjobs.ListTimeclockResponse_Range{
			Range: &pbjobs.TimeclockRange{},
		}

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

	case jobs.TimeclockMode_TIMECLOCK_MODE_TIMELINE:
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
				tUserProps.AvatarFileID.AS("colleague.profile_picture_file_id"),
				tAvatar.FilePath.AS("colleague.profile_picture"),
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
						mysql.AND(
							tColleagueProps.UserID.EQ(tColleague.ID),
							tColleague.Job.EQ(mysql.String(userInfo.GetJob())),
						),
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

	condition := mysql.AND(
		tTimeClock.Job.EQ(mysql.String(userInfo.GetJob())),
		tTimeClock.UserID.EQ(mysql.Int32(userId)),
	)

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

	condition := mysql.AND(
		tTimeClock.Job.EQ(mysql.String(userInfo.GetJob())),
		tColleague.Job.EQ(mysql.String(userInfo.GetJob())),
		mysql.OR(
			mysql.AND(
				tColleagueProps.AbsenceBegin.IS_NULL(),
				tColleagueProps.AbsenceEnd.IS_NULL(),
			),
			tColleagueProps.AbsenceBegin.LT_EQ(
				mysql.DateExp(mysql.CURRENT_DATE().SUB(mysql.INTERVAL(req.GetDays(), mysql.DAY))),
			),
			tColleagueProps.AbsenceEnd.LT_EQ(
				mysql.DateExp(mysql.CURRENT_DATE().SUB(mysql.INTERVAL(req.GetDays(), mysql.DAY))),
			),
		),
		tTimeClock.UserID.NOT_IN(
			tTimeClock.
				SELECT(
					tTimeClock.UserID,
				).
				FROM(tTimeClock).
				WHERE(mysql.AND(
					tTimeClock.Job.EQ(mysql.String(userInfo.GetJob())),
					tTimeClock.Date.GT_EQ(
						mysql.DateExp(
							mysql.CURRENT_DATE().SUB(mysql.INTERVAL(req.GetDays(), mysql.DAY)),
						),
					),
				)).
				GROUP_BY(tTimeClock.UserID),
		),
	)

	countStmt := tTimeClock.
		SELECT(
			mysql.COUNT(mysql.DISTINCT(tTimeClock.UserID)).AS("data_count.total"),
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
					mysql.AND(
						tColleagueProps.UserID.EQ(tTimeClock.UserID),
						tColleague.Job.EQ(mysql.String(userInfo.GetJob())),
					),
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
	orderBys := []mysql.OrderByClause{}
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
			tUserProps.AvatarFileID.AS("colleague.profile_picture_file_id"),
			tAvatar.FilePath.AS("colleague.profile_picture"),
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
					mysql.AND(
						tColleagueProps.UserID.EQ(tTimeClock.UserID),
						tColleague.Job.EQ(mysql.String(userInfo.GetJob())),
					),
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

	return resp, nil
}
