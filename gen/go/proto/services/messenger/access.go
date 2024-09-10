package messenger

import (
	"context"
	"errors"
	"slices"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/messenger"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) getThreadAccess(ctx context.Context, threadId uint64) (*messenger.ThreadAccess, error) {
	tThreadsJobAccess := tThreadsJobAccess.AS("threadjobaccess")
	jobStmt := tThreadsJobAccess.
		SELECT(
			tThreadsJobAccess.ID,
			tThreadsJobAccess.ThreadID,
			tThreadsJobAccess.Job,
			tThreadsJobAccess.MinimumGrade,
			tThreadsJobAccess.Access,
		).
		FROM(
			tThreadsJobAccess,
		).
		WHERE(
			tThreadsJobAccess.ThreadID.EQ(jet.Uint64(threadId)),
		).
		ORDER_BY(
			tThreadsJobAccess.Job.ASC(),
			tThreadsJobAccess.MinimumGrade.ASC(),
		)

	var jobAccess []*messenger.ThreadJobAccess
	if err := jobStmt.QueryContext(ctx, s.db, &jobAccess); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	tThreadsUserAccess := tThreadsUserAccess.AS("threaduseraccess")
	userStmt := tThreadsUserAccess.
		SELECT(
			tThreadsUserAccess.ID,
			tThreadsUserAccess.ThreadID,
			tThreadsUserAccess.UserID,
			tThreadsUserAccess.Access,
			tUsers.ID,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Dateofbirth,
			tUsers.PhoneNumber,
			tUserProps.Avatar.AS("usershort.avatar"),
		).
		FROM(
			tThreadsUserAccess.
				LEFT_JOIN(tUsers,
					tUsers.ID.EQ(tThreadsUserAccess.UserID),
				).
				LEFT_JOIN(tUserProps,
					tUserProps.UserID.EQ(tThreadsUserAccess.UserID),
				),
		).
		WHERE(
			tThreadsUserAccess.ThreadID.EQ(jet.Uint64(threadId)),
		).
		ORDER_BY(
			tThreadsUserAccess.ID.ASC(),
		)

	var userAccess []*messenger.ThreadUserAccess
	if err := userStmt.QueryContext(ctx, s.db, &userAccess); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return &messenger.ThreadAccess{
		Jobs:  jobAccess,
		Users: userAccess,
	}, nil
}

func (s *Server) handleThreadAccessChanges(ctx context.Context, tx qrm.DB, mode messenger.AccessLevelUpdateMode, threadId uint64, access *messenger.ThreadAccess) (*messenger.ThreadAccess, error) {
	switch mode {
	case messenger.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_UNSPECIFIED:
		fallthrough
	case messenger.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_UPDATE:
		// Get existing job and user accesses from database
		current, err := s.getThreadAccess(ctx, threadId)
		if err != nil {
			return nil, err
		}

		toCreate, toUpdate, toDelete := s.compareThreadAccess(current, access)

		if err := s.createThreadAccess(ctx, tx, threadId, toCreate); err != nil {
			return nil, err
		}

		if err := s.updateThreadAccess(ctx, tx, threadId, toUpdate); err != nil {
			return nil, err
		}

		if err := s.deleteThreadAccess(ctx, tx, threadId, toDelete); err != nil {
			return nil, err
		}

		return toDelete, nil

	case messenger.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_DELETE:
		if err := s.deleteThreadAccess(ctx, tx, threadId, access); err != nil {
			return nil, err
		}

	case messenger.AccessLevelUpdateMode_ACCESS_LEVEL_UPDATE_MODE_CLEAR:
		if err := s.clearThreadAccess(ctx, tx, threadId); err != nil {
			return nil, err
		}
	}

	return nil, nil
}

func (s *Server) compareThreadAccess(current, in *messenger.ThreadAccess) (toCreate *messenger.ThreadAccess, toUpdate *messenger.ThreadAccess, toDelete *messenger.ThreadAccess) {
	toCreate = &messenger.ThreadAccess{}
	toUpdate = &messenger.ThreadAccess{}
	toDelete = &messenger.ThreadAccess{}

	if current == nil || (len(current.Jobs) == 0 && len(current.Users) == 0) {
		return in, toUpdate, toDelete
	}

	slices.SortFunc(current.Jobs, func(a, b *messenger.ThreadJobAccess) int {
		return int(a.Id - b.Id)
	})

	if len(current.Jobs) == 0 {
		toCreate.Jobs = in.Jobs
	} else {
		foundTracker := []int{}
		for _, cj := range current.Jobs {
			var found *messenger.ThreadJobAccess
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
			var found *messenger.ThreadUserAccess
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

func (s *Server) createThreadAccess(ctx context.Context, tx qrm.DB, threadId uint64, access *messenger.ThreadAccess) error {
	if access == nil {
		return nil
	}

	if access.Jobs != nil {
		for k := 0; k < len(access.Jobs); k++ {
			// Create thread job access
			stmt := tThreadsJobAccess.
				INSERT(
					tThreadsJobAccess.ThreadID,
					tThreadsJobAccess.Job,
					tThreadsJobAccess.MinimumGrade,
					tThreadsJobAccess.Access,
				).
				VALUES(
					threadId,
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
			stmt := tThreadsUserAccess.
				INSERT(
					tThreadsUserAccess.ThreadID,
					tThreadsUserAccess.UserID,
					tThreadsUserAccess.Access,
				).
				VALUES(
					threadId,
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

func (s *Server) updateThreadAccess(ctx context.Context, tx qrm.DB, threadId uint64, access *messenger.ThreadAccess) error {
	if access == nil {
		return nil
	}

	if access.Jobs != nil {
		for k := 0; k < len(access.Jobs); k++ {
			// Create document job access
			stmt := tThreadsJobAccess.
				UPDATE(
					tThreadsJobAccess.ThreadID,
					tThreadsJobAccess.Job,
					tThreadsJobAccess.MinimumGrade,
					tThreadsJobAccess.Access,
				).
				SET(
					threadId,
					access.Jobs[k].Job,
					access.Jobs[k].MinimumGrade,
					access.Jobs[k].Access,
				).
				WHERE(
					tThreadsJobAccess.ID.EQ(jet.Uint64(access.Jobs[k].Id)),
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return err
			}
		}
	}

	if access.Users != nil {
		for k := 0; k < len(access.Users); k++ {
			// Create document user access
			stmt := tThreadsUserAccess.
				UPDATE(
					tThreadsUserAccess.ThreadID,
					tThreadsUserAccess.UserID,
					tThreadsUserAccess.Access,
				).
				SET(
					threadId,
					access.Users[k].UserId,
					access.Users[k].Access,
				).
				WHERE(
					tThreadsUserAccess.ID.EQ(jet.Uint64(access.Users[k].Id)),
				)

			if _, err := stmt.ExecContext(ctx, tx); err != nil {
				return err
			}
		}
	}

	return nil
}

func (s *Server) deleteThreadAccess(ctx context.Context, tx qrm.DB, threadId uint64, access *messenger.ThreadAccess) error {
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

		stmt := tThreadsJobAccess.
			DELETE().
			WHERE(jet.AND(
				tThreadsJobAccess.ID.IN(jobIds...),
				tThreadsJobAccess.ThreadID.EQ(jet.Uint64(threadId)),
			)).
			LIMIT(25)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
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

		userStmt := tThreadsUserAccess.
			DELETE().
			WHERE(
				jet.AND(
					tThreadsUserAccess.ID.IN(uaIds...),
					tThreadsUserAccess.ThreadID.EQ(jet.Uint64(threadId)),
				),
			).
			LIMIT(25)

		if _, err := userStmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) clearThreadAccess(ctx context.Context, tx qrm.DB, threadId uint64) error {
	jobStmt := tThreadsJobAccess.
		DELETE().
		WHERE(tThreadsJobAccess.ThreadID.EQ(jet.Uint64(threadId)))

	if _, err := jobStmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	userStmt := tThreadsUserAccess.
		DELETE().
		WHERE(tThreadsUserAccess.ThreadID.EQ(jet.Uint64(threadId)))

	if _, err := userStmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}
