package manager

import (
	"context"

	"github.com/galexrt/fivenet/gen/go/proto/resources/dispatch"
	errorscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/errors"
	eventscentrum "github.com/galexrt/fivenet/gen/go/proto/services/centrum/events"
	"github.com/galexrt/fivenet/pkg/utils/dbutils"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/protobuf/proto"
)

func (s *Manager) DisponentSignOn(ctx context.Context, job string, userId int32, signon bool) error {
	if signon {
		if _, ok := s.tracker.GetUserByJobAndID(job, userId); !ok {
			return errorscentrum.ErrNotOnDuty
		}

		stmt := tCentrumUsers.
			INSERT(
				tCentrumUsers.Job,
				tCentrumUsers.UserID,
				tCentrumUsers.Identifier,
			).
			VALUES(
				job,
				userId,
				tUsers.
					SELECT(
						tUsers.Identifier.AS("identifier"),
					).
					FROM(tUsers).
					WHERE(
						tUsers.ID.EQ(jet.Int32(userId)),
					).
					LIMIT(1),
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return errorscentrum.ErrFailedQuery
			}
		}
	} else {
		stmt := tCentrumUsers.
			DELETE().
			WHERE(jet.AND(
				tCentrumUsers.Job.EQ(jet.String(job)),
				tCentrumUsers.UserID.EQ(jet.Int32(userId)),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return errorscentrum.ErrFailedQuery
		}
	}

	// Load updated disponents into state
	if err := s.LoadDisponents(ctx, job); err != nil {
		return errorscentrum.ErrFailedQuery
	}

	disponents := s.GetDisponents(job)
	change := &dispatch.DisponentsChange{
		Job:        job,
		Disponents: disponents,
	}
	data, err := proto.Marshal(change)
	if err != nil {
		return errorscentrum.ErrFailedQuery
	}
	s.events.JS.PublishAsync(eventscentrum.BuildSubject(eventscentrum.TopicGeneral, eventscentrum.TypeGeneralDisponents, job, 0), data)

	return nil
}
