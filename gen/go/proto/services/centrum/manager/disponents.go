package manager

import (
	"context"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/centrum"
	errorscentrum "github.com/fivenet-app/fivenet/gen/go/proto/services/centrum/errors"
	eventscentrum "github.com/fivenet-app/fivenet/gen/go/proto/services/centrum/events"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/protobuf/proto"
)

func (s *Manager) DisponentSignOn(ctx context.Context, job string, userId int32, signon bool) error {
	if signon {
		if _, ok := s.tracker.GetUserById(userId); !ok {
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
				return err
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
			return err
		}
	}

	// Load updated disponents into state
	if err := s.LoadDisponentsFromDB(ctx, job); err != nil {
		return err
	}

	disponents, err := s.GetDisponents(ctx, job)
	if err != nil {
		return errorscentrum.ErrFailedQuery
	}

	change := &centrum.Disponents{
		Job:        job,
		Disponents: disponents,
	}
	data, err := proto.Marshal(change)
	if err != nil {
		return err
	}

	if _, err := s.js.Publish(ctx, eventscentrum.BuildSubject(eventscentrum.TopicGeneral, eventscentrum.TypeGeneralDisponents, job), data); err != nil {
		return err
	}

	return nil
}
