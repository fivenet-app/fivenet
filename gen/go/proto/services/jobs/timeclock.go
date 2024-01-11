package jobs

import (
	"context"
	"errors"
	"slices"
	"time"

	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	errorsjobs "github.com/galexrt/fivenet/gen/go/proto/services/jobs/errors"
	permsjobs "github.com/galexrt/fivenet/gen/go/proto/services/jobs/perms"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	"github.com/galexrt/fivenet/pkg/perms"
	timeutils "github.com/galexrt/fivenet/pkg/utils/time"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var (
	tTimeClock = table.FivenetJobsTimeclock.AS("timeclock_entry")
)

func (s *Server) TimeclockListEntries(ctx context.Context, req *TimeclockListEntriesRequest) (*TimeclockListEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := jet.Bool(true)
	condition = condition.AND(tTimeClock.Job.EQ(jet.String(userInfo.Job)))

	// Field Permission Check
	fieldsAttr, err := s.p.Attr(userInfo, permsjobs.JobsServicePerm, permsjobs.JobsServiceTimeclockListEntriesPerm, permsjobs.JobsServiceTimeclockListEntriesAccessPermField)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}
	var fields perms.StringList
	if fieldsAttr != nil {
		fields = fieldsAttr.([]string)
	}

	if len(fields) == 0 || !slices.Contains(fields, "All") {
		condition = condition.AND(tTimeClock.UserID.EQ(jet.Int32(userInfo.UserId)))
	}

	if len(req.UserIds) > 0 {
		ids := make([]jet.Expression, len(req.UserIds))
		for i := 0; i < len(req.UserIds); i++ {
			ids[i] = jet.Int32(req.UserIds[i])
		}

		condition = condition.AND(
			tTimeClock.UserID.IN(ids...),
		)
	}

	if req.From != nil {
		condition = condition.AND(tTimeClock.Date.LT_EQ(
			jet.TimestampT(timeutils.TruncateToDay(req.From.AsTime())),
		))
	}
	if req.To != nil {
		condition = condition.AND(tTimeClock.Date.GT_EQ(
			jet.TimestampT(timeutils.TruncateToDay(req.To.AsTime())),
		))
	}

	countStmt := tTimeClock.
		SELECT(jet.COUNT(tTimeClock.Date).AS("datacount.totalcount")).
		FROM(tTimeClock).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	pag, limit := req.Pagination.GetResponseWithPageSize(25)
	resp := &TimeclockListEntriesResponse{
		Pagination: pag,
	}

	resp.Stats, err = s.getTimeclockstats(ctx, condition)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	if count.TotalCount <= 0 {
		return resp, nil
	}

	tUser := tUser.AS("user_short")
	stmt := tTimeClock.
		SELECT(
			tTimeClock.Job,
			tTimeClock.Date,
			tTimeClock.UserID,
			tTimeClock.StartTime,
			tTimeClock.EndTime,
			tTimeClock.SpentTime,
			tUser.ID,
			tUser.Identifier,
			tUser.Job,
			tUser.JobGrade,
			tUser.Firstname,
			tUser.Lastname,
			tUser.PhoneNumber,
		).
		FROM(
			tTimeClock.
				INNER_JOIN(tUser,
					tUser.ID.EQ(tTimeClock.UserID),
				),
		).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		ORDER_BY(
			tTimeClock.Date.DESC(),
			tTimeClock.SpentTime.DESC(),
		).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Entries); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
		}
	}

	for i := 0; i < len(resp.Entries); i++ {
		if resp.Entries[i].User != nil {
			s.enricher.EnrichJobInfo(resp.Entries[i].User)
		}
	}

	resp.Pagination.Update(count.TotalCount, len(resp.Entries))

	return resp, nil
}

const TimeclockStatsSpan = 7 * 24 * time.Hour

func (s *Server) TimeclockStats(ctx context.Context, req *TimeclockStatsRequest) (*TimeclockStatsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := tTimeClock.Job.EQ(jet.String(userInfo.Job)).
		AND(tTimeClock.UserID.EQ(jet.Int32(userInfo.UserId)))

	stats, err := s.getTimeclockstats(ctx, condition)
	if err != nil {
		return nil, errswrap.NewError(errorsjobs.ErrFailedQuery, err)
	}

	return &TimeclockStatsResponse{
		Stats: stats,
	}, nil
}