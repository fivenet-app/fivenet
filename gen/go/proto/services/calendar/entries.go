package calendar

import (
	"context"
	"errors"
	"time"

	calendar "github.com/galexrt/fivenet/gen/go/proto/resources/calendar"
	errorscalendar "github.com/galexrt/fivenet/gen/go/proto/services/calendar/errors"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) ListCalendarEntries(ctx context.Context, req *ListCalendarEntriesRequest) (*ListCalendarEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := jet.AND(
		tCalendarEntry.DeletedAt.IS_NULL(),
		jet.OR(
			jet.OR(
				tCalendarEntry.Public.IS_TRUE(),
				tCalendarEntry.CreatorID.EQ(jet.Int32(userInfo.UserId)),
			),
			jet.OR(
				jet.AND(
					tCUserAccess.Access.IS_NOT_NULL(),
					tCUserAccess.Access.GT(jet.Int32(int32(calendar.AccessLevel_ACCESS_LEVEL_BLOCKED))),
				),
				jet.AND(
					tCUserAccess.Access.IS_NULL(),
					tCJobAccess.Access.IS_NOT_NULL(),
					tCJobAccess.Access.GT(jet.Int32(int32(calendar.AccessLevel_ACCESS_LEVEL_BLOCKED))),
				),
			),
		),
	)

	condition = condition.AND(tCalendarEntry.StartTime.GT_EQ(jet.DateTime(int(req.Year), time.Month(req.Month), 1, 0, 0, 0)))

	resp := &ListCalendarEntriesResponse{}

	stmt := tCalendarEntry.
		SELECT(
			tCalendarEntry.ID,
			tCalendarEntry.CreatedAt,
			tCalendarEntry.UpdatedAt,
			tCalendarEntry.DeletedAt,
			tCalendarEntry.CalendarID,
			tCalendarEntry.Job,
			tCalendarEntry.StartTime,
			tCalendarEntry.EndTime,
			tCalendarEntry.Title,
			tCalendarEntry.Content,
			tCalendarEntry.Public,
			tCalendarEntry.CreatorID,
			tCalendarEntry.CreatorJob,
		).
		FROM(tCalendarEntry.
			LEFT_JOIN(tCUserAccess,
				jet.OR(
					tCUserAccess.CalendarID.EQ(tCalendarEntry.CalendarID).
						AND(tCUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
					tCUserAccess.EntryID.EQ(tCalendarEntry.ID).
						AND(tCUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
				),
			).
			LEFT_JOIN(tCJobAccess,
				jet.OR(
					tCJobAccess.CalendarID.EQ(tCalendarEntry.CalendarID).
						AND(tCJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tCJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
					tCJobAccess.EntryID.EQ(tCalendarEntry.ID).
						AND(tCJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tCJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				),
			).
			LEFT_JOIN(tCreator,
				tCalendarEntry.CreatorID.EQ(tCreator.ID),
			),
		).
		GROUP_BY(tCalendarEntry.ID).
		WHERE(condition).
		LIMIT(200)

	if err := stmt.QueryContext(ctx, s.db, &resp.Entries); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Entries); i++ {
		if resp.Entries[i].Creator != nil {
			jobInfoFn(resp.Entries[i].Creator)
		}
	}

	return resp, nil
}

func (s *Server) CreateOrUpdateCalendarEntries(ctx context.Context, req *CreateOrUpdateCalendarEntriesRequest) (*CreateOrUpdateCalendarEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	_ = userInfo

	// TODO

	return nil, nil
}
func (s *Server) DeleteCalendarEntries(ctx context.Context, req *DeleteCalendarEntriesRequest) (*DeleteCalendarEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	_ = userInfo

	// TODO

	return nil, nil
}
func (s *Server) ShareCalendarEntry(ctx context.Context, req *ShareCalendarEntryRequest) (*ShareCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	_ = userInfo

	// TODO

	return nil, nil
}
