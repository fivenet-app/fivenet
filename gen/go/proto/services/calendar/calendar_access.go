package calendar

import (
	"context"
	"errors"
	"slices"

	calendar "github.com/fivenet-app/fivenet/gen/go/proto/resources/calendar"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) getAccess(ctx context.Context, calendarId uint64) (*calendar.CalendarAccess, error) {
	dest := &calendar.CalendarAccess{}

	jobStmt := tCJobAccess.
		SELECT(
			tCJobAccess.ID,
			tCJobAccess.CreatedAt,
			tCJobAccess.CalendarID,
			tCJobAccess.Job,
			tCJobAccess.MinimumGrade,
			tCJobAccess.Access,
		).
		FROM(tCJobAccess).
		WHERE(tCJobAccess.CalendarID.EQ(jet.Uint64(calendarId))).
		ORDER_BY(tCJobAccess.ID.ASC())

	if err := jobStmt.QueryContext(ctx, s.db, &dest.Jobs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	userStmt := tCUserAccess.
		SELECT(
			tCUserAccess.ID,
			tCUserAccess.CreatedAt,
			tCUserAccess.CalendarID,
			tCUserAccess.UserID,
			tCUserAccess.Access,
			tUsers.ID,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
		).
		FROM(
			tCUserAccess.
				LEFT_JOIN(tUsers,
					tUsers.ID.EQ(tCUserAccess.UserID),
				),
		).
		WHERE(tCUserAccess.CalendarID.EQ(jet.Uint64(calendarId))).
		ORDER_BY(tCUserAccess.ID.ASC())

	if err := userStmt.QueryContext(ctx, s.db, &dest.Jobs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

func (s *Server) checkIfUserHasAccessToCalendar(ctx context.Context, calendarId uint64, userInfo *userinfo.UserInfo, access calendar.AccessLevel, publicOk bool) (bool, error) {
	out, err := s.checkIfUserHasAccessToCalendarIDs(ctx, userInfo, access, publicOk, calendarId)
	return len(out) > 0, err
}

func (s *Server) checkIfUserHasAccessToCalendars(ctx context.Context, userInfo *userinfo.UserInfo, access calendar.AccessLevel, publicOk bool, calendarIds ...uint64) (bool, error) {
	out, err := s.checkIfUserHasAccessToCalendarIDs(ctx, userInfo, access, publicOk, calendarIds...)
	return len(out) == len(calendarIds), err
}

type calendarAccessEntry struct {
	ID     uint64 `alias:"calendar.id"`
	Public bool   `alias:"calendar.public"`
}

func (s *Server) checkIfUserHasAccessToCalendarIDs(ctx context.Context, userInfo *userinfo.UserInfo, access calendar.AccessLevel, publicOk bool, calendarIds ...uint64) ([]*calendarAccessEntry, error) {
	var dest []*calendarAccessEntry
	if len(calendarIds) == 0 {
		return dest, nil
	}

	// Allow superusers access to any docs
	if userInfo.SuperUser {
		for i := 0; i < len(calendarIds); i++ {
			dest = append(dest, &calendarAccessEntry{
				ID: calendarIds[i],
			})
		}
		return dest, nil
	}

	ids := make([]jet.Expression, len(calendarIds))
	for i := 0; i < len(calendarIds); i++ {
		ids[i] = jet.Uint64(calendarIds[i])
	}

	condition := jet.Bool(false)
	if publicOk {
		condition = tCalendar.Public.IS_TRUE()
	}

	stmt := tCalendar.
		SELECT(
			tCalendar.ID,
		).
		FROM(tCalendar.
			LEFT_JOIN(tCUserAccess,
				tCUserAccess.CalendarID.EQ(tCalendar.ID).
					AND(tCUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
			).
			LEFT_JOIN(tCJobAccess,
				tCJobAccess.CalendarID.EQ(tCalendar.ID).
					AND(tCJobAccess.Job.EQ(jet.String(userInfo.Job))).
					AND(tCJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
			).
			LEFT_JOIN(tCreator,
				tCalendar.CreatorID.EQ(tCreator.ID),
			),
		).
		GROUP_BY(tCalendar.ID).
		WHERE(jet.AND(
			tCalendar.ID.IN(ids...),
			tCalendar.DeletedAt.IS_NULL(),
			jet.OR(
				tCalendar.CreatorID.EQ(jet.Int32(userInfo.UserId)),
				tCalendar.CreatorJob.EQ(jet.String(userInfo.Job)),
				jet.AND(
					tCUserAccess.Access.IS_NOT_NULL(),
					tCUserAccess.Access.GT_EQ(jet.Int32(int32(access))),
				),
				jet.AND(
					tCUserAccess.Access.IS_NULL(),
					tCJobAccess.Access.IS_NOT_NULL(),
					tCJobAccess.Access.GT_EQ(jet.Int32(int32(access))),
				),
				condition,
			),
		)).
		ORDER_BY(tCalendar.ID.DESC(), tCJobAccess.MinimumGrade)

	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

func (s *Server) handleCalendarAccessChanges(ctx context.Context, tx qrm.DB, mode calendar.AccessLevelUpdateMode, calendarId uint64, access *calendar.CalendarAccess) error {
	// Get existing job and user accesses from database
	current, err := s.getAccess(ctx, calendarId)
	if err != nil {
		return err
	}

	switch mode {
	case calendar.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_UNSPECIFIED:
		fallthrough
	case calendar.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_UPDATE:
		toCreate, toUpdate, toDelete := s.compareCalendarAccess(current, access)

		if err := s.createCalendarAccess(ctx, tx, calendarId, toCreate); err != nil {
			return err
		}

		if err := s.updateCalendarAccess(ctx, tx, calendarId, toUpdate); err != nil {
			return err
		}

		if err := s.deleteCalendarAccess(ctx, tx, calendarId, toDelete); err != nil {
			return err
		}

	case calendar.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_DELETE:
		if err := s.deleteCalendarAccess(ctx, tx, calendarId, access); err != nil {
			return err
		}

	case calendar.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_CLEAR:
		if err := s.clearCalendarAccess(ctx, tx, calendarId); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) compareCalendarAccess(current, in *calendar.CalendarAccess) (toCreate *calendar.CalendarAccess, toUpdate *calendar.CalendarAccess, toDelete *calendar.CalendarAccess) {
	toCreate = &calendar.CalendarAccess{}
	toUpdate = &calendar.CalendarAccess{}
	toDelete = &calendar.CalendarAccess{}

	if current == nil || (len(current.Jobs) == 0 && len(current.Users) == 0) {
		return in, toUpdate, toDelete
	}

	slices.SortFunc(current.Jobs, func(a, b *calendar.CalendarJobAccess) int {
		return int(a.Id - b.Id)
	})

	if len(current.Jobs) == 0 {
		toCreate.Jobs = in.Jobs
	} else {
		foundTracker := []int{}
		for _, cj := range current.Jobs {
			var found *calendar.CalendarJobAccess
			var foundIdx int
			for i, uj := range in.Jobs {
				if cj.Job != uj.Job {
					continue
				}
				if cj.MinimumGrade != uj.MinimumGrade {
					continue
				}
				found = uj
				foundIdx = i
				break
			}
			// No match in incoming job access, needs to be deleted
			if found == nil {
				toDelete.Jobs = append(toDelete.Jobs, cj)
				continue
			}

			foundTracker = append(foundTracker, foundIdx)

			changed := false
			if cj.MinimumGrade != found.MinimumGrade {
				cj.MinimumGrade = found.MinimumGrade
				changed = true
			}
			if cj.Access != found.Access {
				cj.Access = found.Access
				changed = true
			}

			if changed {
				toUpdate.Jobs = append(toUpdate.Jobs, cj)
			}
		}

		for i, uj := range in.Jobs {
			idx := slices.Index(foundTracker, i)
			if idx == -1 {
				toCreate.Jobs = append(toCreate.Jobs, uj)
			}
		}
	}

	if len(current.Users) == 0 {
		toCreate.Users = in.Users
	} else {
		foundTracker := []int{}
		for _, cj := range current.Users {
			var found *calendar.CalendarUserAccess
			var foundIdx int
			for i, uj := range in.Users {
				if cj.UserId != uj.UserId {
					continue
				}
				found = uj
				foundIdx = i
				break
			}
			// No match in incoming job access, needs to be deleted
			if found == nil {
				toDelete.Users = append(toDelete.Users, cj)
				continue
			}

			foundTracker = append(foundTracker, foundIdx)

			changed := false
			if cj.Access != found.Access {
				cj.Access = found.Access
				changed = true
			}

			if changed {
				toUpdate.Users = append(toUpdate.Users, cj)
			}
		}

		for i, uj := range in.Users {
			idx := slices.Index(foundTracker, i)
			if idx == -1 {
				toCreate.Users = append(toCreate.Users, uj)
			}
		}
	}

	return
}

func (s *Server) createCalendarAccess(ctx context.Context, tx qrm.DB, calendarId uint64, access *calendar.CalendarAccess) error {
	if access == nil {
		return nil
	}

	if access.Jobs != nil {
		for k := 0; k < len(access.Jobs); k++ {
			// Create document job access
			tCJobAccess := table.FivenetCalendarJobAccess
			stmt := tCJobAccess.
				INSERT(
					tCJobAccess.CalendarID,
					tCJobAccess.Job,
					tCJobAccess.MinimumGrade,
					tCJobAccess.Access,
				).
				VALUES(
					calendarId,
					access.Jobs[k].Job,
					access.Jobs[k].MinimumGrade,
					access.Jobs[k].Access,
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return err
			}
		}
	}

	if access.Users != nil {
		for k := 0; k < len(access.Users); k++ {
			// Create document user access
			tCUserAccess := table.FivenetCalendarUserAccess
			stmt := tCUserAccess.
				INSERT(
					tCUserAccess.CalendarID,
					tCUserAccess.UserID,
					tCUserAccess.Access,
				).
				VALUES(
					calendarId,
					access.Users[k].UserId,
					access.Users[k].Access,
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Server) updateCalendarAccess(ctx context.Context, tx qrm.DB, calendarId uint64, access *calendar.CalendarAccess) error {
	if access == nil {
		return nil
	}

	if access.Jobs != nil {
		for k := 0; k < len(access.Jobs); k++ {
			// Create document job access
			tCJobAccess := table.FivenetCalendarJobAccess
			stmt := tCJobAccess.
				UPDATE(
					tCJobAccess.CalendarID,
					tCJobAccess.Job,
					tCJobAccess.MinimumGrade,
					tCJobAccess.Access,
				).
				SET(
					calendarId,
					access.Jobs[k].Job,
					access.Jobs[k].MinimumGrade,
					access.Jobs[k].Access,
				).
				WHERE(
					tCJobAccess.ID.EQ(jet.Uint64(access.Jobs[k].Id)),
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return err
			}
		}
	}

	if access.Users != nil {
		for k := 0; k < len(access.Users); k++ {
			// Create document user access
			tCUserAccess := table.FivenetCalendarUserAccess
			stmt := tCUserAccess.
				UPDATE(
					tCUserAccess.CalendarID,
					tCUserAccess.UserID,
					tCUserAccess.Access,
				).
				SET(
					calendarId,
					access.Users[k].UserId,
					access.Users[k].Access,
				).
				WHERE(
					tCUserAccess.ID.EQ(jet.Uint64(access.Users[k].Id)),
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Server) deleteCalendarAccess(ctx context.Context, tx qrm.DB, calendarId uint64, access *calendar.CalendarAccess) error {
	if access == nil {
		return nil
	}

	if access.Jobs != nil && len(access.Jobs) > 0 {
		jobIds := []jet.Expression{}
		for i := 0; i < len(access.Jobs); i++ {
			if access.Jobs[i].Id == 0 {
				continue
			}
			jobIds = append(jobIds, jet.Uint64(access.Jobs[i].Id))
		}

		tCJobAccess := table.FivenetCalendarJobAccess
		jobStmt := tCJobAccess.
			DELETE().
			WHERE(
				jet.AND(
					tCJobAccess.ID.IN(jobIds...),
					tCJobAccess.CalendarID.EQ(jet.Uint64(calendarId)),
				),
			).
			LIMIT(25)

		if _, err := jobStmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	if access.Users != nil && len(access.Users) > 0 {
		uaIds := []jet.Expression{}
		for i := 0; i < len(access.Users); i++ {
			if access.Users[i].Id == 0 {
				continue
			}
			uaIds = append(uaIds, jet.Uint64(access.Users[i].Id))
		}

		tCUserAccess := table.FivenetCalendarUserAccess
		userStmt := tCUserAccess.
			DELETE().
			WHERE(
				jet.AND(
					tCUserAccess.ID.IN(uaIds...),
					tCUserAccess.CalendarID.EQ(jet.Uint64(calendarId)),
				),
			).
			LIMIT(25)

		if _, err := userStmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) clearCalendarAccess(ctx context.Context, tx qrm.DB, calendarId uint64) error {
	jobStmt := tCJobAccess.
		DELETE().
		WHERE(jet.AND(
			tCJobAccess.CalendarID.EQ(jet.Uint64(calendarId)),
		))

	if _, err := jobStmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	userStmt := tCUserAccess.
		DELETE().
		WHERE(jet.AND(
			tCUserAccess.CalendarID.EQ(jet.Uint64(calendarId)),
		))

	if _, err := userStmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}
