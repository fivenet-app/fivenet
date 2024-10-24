package messenger

import (
	"context"
	"errors"
	"slices"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/messenger"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) checkIfUserHasAccessToThread(ctx context.Context, threadId uint64, userInfo *userinfo.UserInfo, access messenger.AccessLevel) (bool, error) {
	out, err := s.checkIfUserHasAccessToThreadIDs(ctx, userInfo, access, threadId)
	return len(out) > 0, err
}

func (s *Server) checkIfUserHasAccessToThreads(ctx context.Context, userInfo *userinfo.UserInfo, access messenger.AccessLevel, threadIds ...uint64) (bool, error) {
	out, err := s.checkIfUserHasAccessToThreadIDs(ctx, userInfo, access, threadIds...)
	return len(out) == len(threadIds), err
}

func (s *Server) checkIfUserHasAccessToThreadIDs(ctx context.Context, userInfo *userinfo.UserInfo, access messenger.AccessLevel, threadIds ...uint64) ([]uint64, error) {
	if len(threadIds) == 0 {
		return threadIds, nil
	}

	// Allow superusers access to any docs
	if userInfo.SuperUser {
		return threadIds, nil
	}

	ids := make([]jet.Expression, len(threadIds))
	for i := 0; i < len(threadIds); i++ {
		ids[i] = jet.Uint64(threadIds[i])
	}

	stmt := tThreads.
		SELECT(
			tThreads.ID,
		).
		FROM(
			tThreads.
				LEFT_JOIN(tThreadsUserAccess,
					tThreadsUserAccess.ThreadID.EQ(tThreads.ID).
						AND(tThreadsUserAccess.UserID.EQ(jet.Int32(userInfo.UserId))),
				).
				LEFT_JOIN(tThreadsJobAccess,
					tThreadsJobAccess.ThreadID.EQ(tThreads.ID).
						AND(tThreadsJobAccess.Job.EQ(jet.String(userInfo.Job))).
						AND(tThreadsJobAccess.MinimumGrade.LT_EQ(jet.Int32(userInfo.JobGrade))),
				),
		).
		WHERE(jet.AND(
			tThreads.ID.IN(ids...),
			tThreads.DeletedAt.IS_NULL(),
			jet.OR(
				jet.AND(
					tThreads.CreatorID.EQ(jet.Int32(userInfo.UserId)),
					tThreads.CreatorJob.EQ(jet.String(userInfo.Job)),
				),
				jet.AND(
					tThreadsUserAccess.Access.IS_NOT_NULL(),
					tThreadsUserAccess.Access.GT_EQ(jet.Int32(int32(access))),
				),
				jet.AND(
					tThreadsUserAccess.Access.IS_NULL(),
					tThreadsJobAccess.Access.IS_NOT_NULL(),
					tThreadsJobAccess.Access.GT_EQ(jet.Int32(int32(access))),
				),
			),
		)).
		GROUP_BY(tThreads.ID).
		ORDER_BY(tThreads.ID.DESC(), tThreadsJobAccess.MinimumGrade)

	var dest struct {
		IDs []uint64 `alias:"document.id"`
	}
	if err := stmt.QueryContext(ctx, s.db, &dest.IDs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest.IDs, nil
}

func (s *Server) checkIfHasAccess(levels []string, userInfo *userinfo.UserInfo, creatorJob string, creator *users.UserShort) bool {
	if userInfo.SuperUser {
		return true
	}

	// If the document creator job is not equal to the creator's current job, normal access checks need to be applied
	// and not the rank attributes checks
	if creatorJob != userInfo.Job {
		return true
	}

	// If the creator is nil, treat it like a normal doc access check
	if creator == nil {
		return true
	}

	// If no levels set, assume "Own" as a safe default
	if len(levels) == 0 {
		return creator.UserId == userInfo.UserId
	}

	if slices.Contains(levels, "Any") {
		return true
	}
	if slices.Contains(levels, "Lower_Rank") {
		if creator.JobGrade < userInfo.JobGrade {
			return true
		}
	}
	if slices.Contains(levels, "Same_Rank") {
		if creator.JobGrade <= userInfo.JobGrade {
			return true
		}
	}
	if slices.Contains(levels, "Own") {
		if creator.UserId == userInfo.UserId {
			return true
		}
	}

	return false
}
