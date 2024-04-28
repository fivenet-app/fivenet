package calendar

import (
	"context"
	"errors"

	calendar "github.com/galexrt/fivenet/gen/go/proto/resources/calendar"
	database "github.com/galexrt/fivenet/gen/go/proto/resources/common/database"
	errorscalendar "github.com/galexrt/fivenet/gen/go/proto/services/calendar/errors"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/errswrap"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) ListCalendars(ctx context.Context, req *ListCalendarsRequest) (*ListCalendarsResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	condition := jet.AND(
		tCalendar.DeletedAt.IS_NULL(),
		jet.OR(
			jet.OR(
				tCalendar.Public.IS_TRUE(),
				tCalendar.CreatorID.EQ(jet.Int32(userInfo.UserId)),
			),
			jet.OR(
				jet.AND(
					tCUserAccess.Access.IS_NOT_NULL(),
					tCUserAccess.Access.NOT_EQ(jet.Int32(int32(calendar.AccessLevel_ACCESS_LEVEL_BLOCKED))),
				),
				jet.AND(
					tCUserAccess.Access.IS_NULL(),
					tCJobAccess.Access.IS_NOT_NULL(),
					tCJobAccess.Access.NOT_EQ(jet.Int32(int32(calendar.AccessLevel_ACCESS_LEVEL_BLOCKED))),
				),
			),
		),
	)

	countStmt := tCalendar.
		SELECT(
			jet.COUNT(jet.DISTINCT(tCalendar.ID)).AS("datacount.totalcount"),
		).
		FROM(tCalendar.
			LEFT_JOIN(tCUserAccess,
				tCUserAccess.CalendarID.EQ(tCalendar.ID).
					AND(tCUserAccess.EntryID.IS_NULL()).
					AND(tCUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
			).
			LEFT_JOIN(tCJobAccess,
				tCJobAccess.CalendarID.EQ(tCalendar.ID).
					AND(tCJobAccess.EntryID.IS_NULL()).
					AND(tCJobAccess.Job.EQ(jet.String(userInfo.Job))).
					AND(tCJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
			).
			LEFT_JOIN(tCreator,
				tCalendar.CreatorID.EQ(tCreator.ID),
			),
		).
		GROUP_BY(tCalendar.ID).
		WHERE(condition)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.db, &count); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	pag, limit := req.Pagination.GetResponse(count.TotalCount)
	resp := &ListCalendarsResponse{
		Pagination: pag,
	}

	if count.TotalCount <= 0 {
		return resp, nil
	}

	stmt := tCalendar.
		SELECT(
			tCalendar.ID,
			tCalendar.CreatedAt,
			tCalendar.UpdatedAt,
			tCalendar.DeletedAt,
			tCalendar.Job,
			tCalendar.Name,
			tCalendar.Description,
			tCalendar.Public,
			tCalendar.Closed,
			tCalendar.CreatorID,
			tCalendar.CreatorJob,
		).
		FROM(tCalendar.
			LEFT_JOIN(tCUserAccess,
				tCUserAccess.CalendarID.EQ(tCalendar.ID).
					AND(tCUserAccess.EntryID.IS_NULL()).
					AND(tCUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
			).
			LEFT_JOIN(tCJobAccess,
				tCJobAccess.CalendarID.EQ(tCalendar.ID).
					AND(tCJobAccess.EntryID.IS_NULL()).
					AND(tCJobAccess.Job.EQ(jet.String(userInfo.Job))).
					AND(tCJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
			).
			LEFT_JOIN(tCreator,
				tCalendar.CreatorID.EQ(tCreator.ID),
			),
		).
		GROUP_BY(tCalendar.ID).
		WHERE(condition).
		OFFSET(req.Pagination.Offset).
		LIMIT(limit)

	if err := stmt.QueryContext(ctx, s.db, &resp.Calendars); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	}

	jobInfoFn := s.enricher.EnrichJobInfoSafeFunc(userInfo)
	for i := 0; i < len(resp.Calendars); i++ {
		if resp.Calendars[i].Creator != nil {
			jobInfoFn(resp.Calendars[i].Creator)
		}
	}

	resp.Pagination.Update(len(resp.Calendars))

	return resp, nil
}

func (s *Server) CreateOrUpdateCalendar(ctx context.Context, req *CreateOrUpdateCalendarRequest) (*CreateOrUpdateCalendarResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	_ = userInfo

	// TODO

	return nil, nil
}
func (s *Server) DeleteCalendar(ctx context.Context, req *DeleteCalendarRequest) (*DeleteCalendarResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	_ = userInfo

	// TODO

	return nil, nil
}
