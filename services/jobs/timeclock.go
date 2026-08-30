package jobs

import (
	"context"
	"math"
	"time"

	jobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs"
	jobscolleagues "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues"
	jobstimeclock "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/timeclock"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	pbjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs"
	permsjobs "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/auth"
	"github.com/fivenet-app/fivenet/v2026/pkg/grpc/errswrap"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	errorsjobs "github.com/fivenet-app/fivenet/v2026/services/jobs/errors"
	jobsstore "github.com/fivenet-app/fivenet/v2026/stores/jobs"
	colleaguehydrator "github.com/fivenet-app/fivenet/v2026/stores/jobs/colleagues/hydrator"
	"github.com/fivenet-app/fivenet/v2026/stores/jobs/usersel"
)

const TimeclockMaxDays = (365 / 2) * 24 * time.Hour // Half a year

var tTimeClock = table.FivenetJobTimeclock.AS("timeclock_entry")

func appendTimeclockColleagueTarget(
	targets []colleaguehydrator.Target,
	userID int32,
	set func(*jobscolleagues.Colleague),
) []colleaguehydrator.Target {
	if userID <= 0 || set == nil {
		return targets
	}

	return append(targets, colleaguehydrator.Target{
		UserID: userID,
		Set:    set,
	})
}

func appendTimeclockEntryColleagueTargets(
	targets []colleaguehydrator.Target,
	entries []*jobstimeclock.TimeclockEntry,
) []colleaguehydrator.Target {
	for _, entry := range entries {
		targets = appendTimeclockColleagueTarget(targets, entry.GetUserId(), entry.SetUser)
	}

	return targets
}

func (s *Server) hydrateTimeclockEntryColleagues(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	entries ...*jobstimeclock.TimeclockEntry,
) error {
	if len(entries) == 0 {
		return nil
	}

	targets := appendTimeclockEntryColleagueTargets(nil, entries)
	if len(targets) == 0 {
		return nil
	}

	if err := s.colleagueHydrator.HydrateTargets(
		ctx,
		s.db,
		userInfo,
		targets,
		colleaguehydrator.ResolveOpts{},
	); err != nil {
		return errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	return nil
}

func (s *Server) ListTimeclock(
	ctx context.Context,
	req *pbjobs.ListTimeclockRequest,
) (*pbjobs.ListTimeclockResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Field Permission Check
	fields, err := permsjobs.TimeclockService.ListTimeclock.AccessTyped.Get(s.perms, userInfo)
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

	selector := req.GetUsers()
	if !fields.Contains(permsjobs.TimeclockServiceListTimeclockAccessPermValueAll) {
		selector = &jobs.UserSelector{
			UserIds: []int32{userInfo.GetUserId()},
		}
	}
	resolvedUserIDs, err := s.userSel.Resolve(
		ctx,
		s.db,
		userInfo,
		selector,
		usersel.ResolveOpts{},
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if usersel.HasSelection(selector) && len(resolvedUserIDs) == 0 {
		pag, _ := req.GetPagination().GetResponseWithPageSize(0, 30)
		return &pbjobs.ListTimeclockResponse{Pagination: pag}, nil
	}

	countQuery := jobsstore.TimeclockQuery{
		Job:      userInfo.GetJob(),
		UserMode: req.GetUserMode(),
		Mode:     req.GetMode(),
		Date:     req.GetDate(),
		PerDay:   req.GetPerDay(),
		UserIDs:  resolvedUserIDs,
		UserID:   userInfo.GetUserId(),
		Sort:     req.GetSort(),
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
		if err := s.hydrateTimeclockEntryColleagues(
			ctx,
			userInfo,
			data.GetEntries()...); err != nil {
			return nil, err
		}
		data.Date = req.GetDate().GetEnd()
		for i := range data.GetEntries() {
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
		if err := s.hydrateTimeclockEntryColleagues(
			ctx,
			userInfo,
			data.GetEntries()...); err != nil {
			return nil, err
		}
		for i := range data.GetEntries() {
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
		if err := s.hydrateTimeclockEntryColleagues(
			ctx,
			userInfo,
			data.GetEntries()...); err != nil {
			return nil, err
		}
		data.Date = req.GetDate().GetEnd()
		for i := range data.GetEntries() {
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
		if err := s.hydrateTimeclockEntryColleagues(
			ctx,
			userInfo,
			data.GetEntries()...); err != nil {
			return nil, err
		}
		data.Date = req.GetDate().GetEnd()
		for i := range data.GetEntries() {
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

	fields, err := permsjobs.TimeclockService.ListTimeclock.AccessTyped.Get(s.perms, userInfo)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	allowAll := fields.Contains(permsjobs.TimeclockServiceListTimeclockAccessPermValueAll)
	requestedSelection := usersel.HasSelection(req.GetUsers())
	selector := req.GetUsers()
	if !allowAll || !requestedSelection {
		selector = &jobs.UserSelector{
			UserIds: []int32{userInfo.GetUserId()},
		}
	}

	resolvedUserIDs, err := s.userSel.Resolve(
		ctx,
		s.db,
		userInfo,
		selector,
		usersel.ResolveOpts{},
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if allowAll && requestedSelection && len(resolvedUserIDs) == 0 {
		return &pbjobs.GetTimeclockStatsResponse{
			Stats:  &jobstimeclock.TimeclockStats{},
			Weekly: []*jobstimeclock.TimeclockWeeklyStats{},
		}, nil
	}

	statsQuery := jobsstore.TimeclockQuery{
		Job:     userInfo.GetJob(),
		UserIDs: resolvedUserIDs,
	}
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

	resolvedUserIDs, err := s.userSel.Resolve(
		ctx,
		s.db,
		userInfo,
		req.GetUsers(),
		usersel.ResolveOpts{},
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}
	if usersel.HasSelection(req.GetUsers()) && len(resolvedUserIDs) == 0 {
		pag, _ := req.GetPagination().GetResponseWithPageSize(0, 20)
		return &pbjobs.ListInactiveEmployeesResponse{
			Pagination: pag,
			Colleagues: []*jobscolleagues.Colleague{},
		}, nil
	}

	count, err := s.store.CountInactiveEmployees(ctx, s.db, jobsstore.InactiveEmployeesQuery{
		Days:    req.GetDays(),
		UserIDs: resolvedUserIDs,
		Job:     userInfo.GetJob(),
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
			Days:    req.GetDays(),
			UserIDs: resolvedUserIDs,
			Sort:    req.GetSort(),
			Offset:  req.GetPagination().GetOffset(),
			Limit:   limit,
			Job:     userInfo.GetJob(),
		},
	)
	if err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	targets := make([]colleaguehydrator.Target, 0, len(resp.GetColleagues()))
	for i := range resp.GetColleagues() {
		colleague := resp.GetColleagues()[i]
		targets = append(targets, colleaguehydrator.Target{
			UserID: colleague.GetUserId(),
			Set: func(hydrated *jobscolleagues.Colleague) {
				resp.Colleagues[i] = hydrated
			},
		})
	}

	if err := s.colleagueHydrator.HydrateTargets(
		ctx,
		s.db,
		userInfo,
		targets,
		colleaguehydrator.ResolveOpts{},
	); err != nil {
		return nil, errswrap.NewError(err, errorsjobs.ErrFailedQuery)
	}

	return resp, nil
}
