package calendar

import (
	"context"
	"errors"
	"time"

	calendar "github.com/galexrt/fivenet/gen/go/proto/resources/calendar"
	errorscalendar "github.com/galexrt/fivenet/gen/go/proto/services/calendar/errors"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
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
			tCreator.ID,
			tCreator.Identifier,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
			tUserProps.Avatar.AS("creator.avatar"),
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
			).
			LEFT_JOIN(tUserProps,
				tUserProps.UserID.EQ(tCreator.ID),
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

func (s *Server) GetCalendarEntry(ctx context.Context, req *GetCalendarEntryRequest) (*GetCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	// Check if user has access to existing calendar
	check, err := s.checkIfUserHasAccessToCalendarEntry(ctx, req.CalendarId, req.EntryId, userInfo, calendar.AccessLevel_ACCESS_LEVEL_VIEW)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
	}

	entry, err := s.getEntry(ctx, userInfo, tCalendarEntry.ID.EQ(jet.Uint64(req.EntryId)))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if entry == nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
	}

	return &GetCalendarEntryResponse{
		Entry: entry,
	}, nil
}

func (s *Server) CreateOrUpdateCalendarEntries(ctx context.Context, req *CreateOrUpdateCalendarEntriesRequest) (*CreateOrUpdateCalendarEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	if req.Entry.Id > 0 {
		check, err := s.checkIfUserHasAccessToCalendarEntry(ctx, req.Entry.CalendarId, req.Entry.Id, userInfo, calendar.AccessLevel_ACCESS_LEVEL_EDIT)
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
		if !check {
			return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
		}

		startTime := jet.TimestampExp(jet.NULL)
		if req.Entry.StartTime != nil {
			startTime = jet.TimestampT(req.Entry.StartTime.AsTime())
		}
		endTime := jet.TimestampExp(jet.NULL)
		if req.Entry.EndTime != nil {
			endTime = jet.TimestampT(req.Entry.EndTime.AsTime())
		}

		stmt := tCalendarEntry.
			UPDATE(
				tCalendarEntry.Title,
				tCalendarEntry.Content,
				tCalendarEntry.StartTime,
				tCalendarEntry.EndTime,
				tCalendarEntry.Public,
			).
			SET(
				tCalendarEntry.Title.SET(jet.String(req.Entry.Title)),
				tCalendarEntry.Content.SET(jet.String(req.Entry.Content)),
				tCalendarEntry.StartTime.SET(startTime),
				tCalendarEntry.EndTime.SET(endTime),
				tCalendarEntry.Public.SET(jet.Bool(req.Entry.Public)),
			).
			WHERE(jet.AND(
				tCalendarEntry.ID.EQ(jet.Uint64(req.Entry.Id)),
				tCalendarEntry.CalendarID.EQ(jet.Uint64(req.Entry.CalendarId)),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
	} else {
		check, err := s.checkIfUserHasAccessToCalendar(ctx, req.Entry.CalendarId, userInfo, calendar.AccessLevel_ACCESS_LEVEL_EDIT)
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}
		if !check {
			return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
		}

		stmt := tCalendarEntry.
			INSERT(
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
			VALUES(
				req.Entry.CalendarId,
				userInfo.Job,
				req.Entry.StartTime,
				req.Entry.EndTime,
				req.Entry.Title,
				req.Entry.Content,
				req.Entry.Public,
				userInfo.UserId,
				userInfo.Job,
			)

		res, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}

		lastId, err := res.LastInsertId()
		if err != nil {
			return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
		}

		req.Entry.Id = uint64(lastId)
	}

	entry, err := s.getEntry(ctx, userInfo, tCalendarEntry.ID.EQ(jet.Uint64(req.Entry.Id)))
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	return &CreateOrUpdateCalendarEntriesResponse{
		Entry: entry,
	}, nil
}

func (s *Server) DeleteCalendarEntries(ctx context.Context, req *DeleteCalendarEntriesRequest) (*DeleteCalendarEntriesResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)

	check, err := s.checkIfUserHasAccessToCalendarEntry(ctx, req.CalendarId, req.EntryId, userInfo, calendar.AccessLevel_ACCESS_LEVEL_MANAGE)
	if err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}
	if !check {
		return nil, errswrap.NewError(err, errorscalendar.ErrNoPerms)
	}

	stmt := tCalendarEntry.
		UPDATE(
			tCalendarEntry.DeletedAt,
		).
		SET(
			tCalendarEntry.DeletedAt.SET(jet.CURRENT_TIMESTAMP()),
		).
		WHERE(jet.AND(
			tCalendarEntry.CalendarID.EQ(jet.Uint64(req.CalendarId)),
			tCalendarEntry.ID.EQ(jet.Uint64(req.EntryId)),
		))

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return nil, errswrap.NewError(err, errorscalendar.ErrFailedQuery)
	}

	return &DeleteCalendarEntriesResponse{}, nil
}

func (s *Server) ShareCalendarEntry(ctx context.Context, req *ShareCalendarEntryRequest) (*ShareCalendarEntryResponse, error) {
	userInfo := auth.MustGetUserInfoFromContext(ctx)
	_ = userInfo

	// TODO

	return nil, nil
}

func (s *Server) getEntry(ctx context.Context, userInfo *userinfo.UserInfo, condition jet.BoolExpression) (*calendar.CalendarEntry, error) {
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
			tCreator.ID,
			tCreator.Identifier,
			tCreator.Job,
			tCreator.JobGrade,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
		).
		FROM(tCalendarEntry.
			LEFT_JOIN(tCreator,
				tCalendar.CreatorID.EQ(tCreator.ID),
			),
		).
		GROUP_BY(tCalendarEntry.ID)

	if condition != nil {
		stmt = stmt.WHERE(condition)

	}

	dest := &calendar.CalendarEntry{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.Id == 0 {
		return nil, nil
	}

	if dest.Creator != nil {
		s.enricher.EnrichJobInfoSafe(userInfo, dest.Creator)
	}

	return dest, nil
}
