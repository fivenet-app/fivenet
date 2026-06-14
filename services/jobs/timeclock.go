package jobs

import (
	"context"
	"math"
	"time"

	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	jobstimeclock "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/timeclock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	pbjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	permsjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
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

	// Field Permission Check
	fields, err := permsjobs.TimeclockService.ListTimeclock.AccessTyped.Get(s.ps, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if !fields.Contains(permsjobs.TimeclockServiceListTimeclockAccessPermValueAll) {
		req.UserMode = jobstimeclock.TimeclockViewMode_TIMECLOCK_VIEW_MODE_SELF
	}

	if req.GetDate() != nil && req.GetMode() <= jobstimeclock.TimeclockMode_TIMECLOCK_MODE_DAILY &&
		req.GetDate().GetEnd() == nil {
		req.Date.End = timestamp.Now()
	}
	if req.GetDate() != nil {
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

	countQuery := jobsstore.TimeclockQuery{
		Job:      userInfo.GetJob(),
		UserMode: req.GetUserMode(),
		Mode:     req.GetMode(),
		Date:     req.GetDate(),
		PerDay:   req.GetPerDay(),
		UserIDs:  req.GetUserIds(),
		UserID:   userInfo.GetUserId(),
	}

	count, err := s.store.CountTimeclock(ctx, s.db, countQuery)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count, 30)
	resp := &pbjobs.ListTimeclockResponse{
		Pagination: pag,
	}

	resp.Stats, err = s.store.GetTimeclockStats(ctx, s.db, countQuery)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	resp.StatsWeekly, err = s.store.GetTimeclockWeeklyStats(ctx, s.db, countQuery)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	if count <= 0 {
		return resp, nil
	}

	listQuery := countQuery
	listQuery.Offset = req.GetPagination().GetOffset()
	listQuery.Limit = limit

	switch req.GetMode() {
	case jobstimeclock.TimeclockMode_TIMECLOCK_MODE_UNSPECIFIED:
		fallthrough

	case jobstimeclock.TimeclockMode_TIMECLOCK_MODE_DAILY:
		resp.Entries = &pbjobs.ListTimeclockResponse_Daily{Daily: &pbjobs.TimeclockDay{}}
		data := resp.GetDaily()
		data.Entries, err = s.store.ListTimeclock(ctx, s.db, listQuery)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		data.Date = req.GetDate().GetEnd()
		for i := range data.GetEntries() {
			if data.GetEntries()[i].GetUser() != nil {
				s.enricher.EnrichJobInfo(data.GetEntries()[i].GetUser())
			}
			data.Sum += int64(math.Round(float64(data.GetEntries()[i].GetSpentTime() * 60 * 60)))
		}
		resp.GetPagination().Update(len(data.GetEntries()))

	case jobstimeclock.TimeclockMode_TIMECLOCK_MODE_WEEKLY:
		resp.Entries = &pbjobs.ListTimeclockResponse_Weekly{Weekly: &pbjobs.TimeclockWeekly{}}
		data := resp.GetWeekly()
		data.Entries, err = s.store.ListTimeclock(ctx, s.db, listQuery)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		for i := range data.GetEntries() {
			if data.GetEntries()[i].GetUser() != nil {
				s.enricher.EnrichJobInfo(data.GetEntries()[i].GetUser())
			}
			data.Sum += int64(math.Round(float64(data.GetEntries()[i].GetSpentTime() * 60 * 60)))
		}
		resp.GetPagination().Update(len(data.GetEntries()))

	case jobstimeclock.TimeclockMode_TIMECLOCK_MODE_RANGE:
		resp.Entries = &pbjobs.ListTimeclockResponse_Range{Range: &pbjobs.TimeclockRange{}}
		data := resp.GetRange()
		data.Entries, err = s.store.ListTimeclock(ctx, s.db, listQuery)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		data.Date = req.GetDate().GetEnd()
		for i := range data.GetEntries() {
			if data.GetEntries()[i].GetUser() != nil {
				s.enricher.EnrichJobInfo(data.GetEntries()[i].GetUser())
			}
			data.Sum += int64(math.Round(float64(data.GetEntries()[i].GetSpentTime() * 60 * 60)))
		}
		resp.GetPagination().Update(len(data.GetEntries()))

	case jobstimeclock.TimeclockMode_TIMECLOCK_MODE_TIMELINE:
		resp.Entries = &pbjobs.ListTimeclockResponse_Range{Range: &pbjobs.TimeclockRange{}}
		data := resp.GetRange()
		data.Entries, err = s.store.ListTimeclockTimeline(ctx, s.db, listQuery)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		data.Date = req.GetDate().GetEnd()
		for i := range data.GetEntries() {
			if data.GetEntries()[i].GetUser() != nil {
				s.enricher.EnrichJobInfo(data.GetEntries()[i].GetUser())
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
		fields, err := permsjobs.TimeclockService.ListTimeclock.AccessTyped.Get(s.ps, userInfo)
		if err != nil {
			return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
		}
		if fields.Contains(permsjobs.TimeclockServiceListTimeclockAccessPermValueAll) {
			userId = req.GetUserId()
		}
	}

	statsQuery := jobsstore.TimeclockQuery{Job: userInfo.GetJob(), UserID: userId}
	stats, err := s.store.GetTimeclockStats(ctx, s.db, statsQuery)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	weekly, err := s.store.GetTimeclockWeeklyStats(ctx, s.db, statsQuery)
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

	count, err := s.store.CountInactiveEmployees(ctx, s.db, jobsstore.InactiveEmployeesQuery{
		Days: req.GetDays(),
		Job:  userInfo.GetJob(),
	})
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	pag, limit := req.GetPagination().GetResponseWithPageSize(count, 20)
	resp := &pbjobs.ListInactiveEmployeesResponse{
		Pagination: pag,
		Colleagues: []*jobscolleagues.Colleague{},
	}
	if count <= 0 {
		return resp, nil
	}

	resp.Colleagues, err = s.store.ListInactiveEmployees(
		ctx,
		s.db,
		jobsstore.InactiveEmployeesQuery{
			Days:   req.GetDays(),
			Sort:   req.GetSort(),
			Offset: req.GetPagination().GetOffset(),
			Limit:  limit,
			Job:    userInfo.GetJob(),
		},
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := range resp.GetColleagues() {
		jobInfoFn(resp.GetColleagues()[i])
	}

	return resp, nil
}
