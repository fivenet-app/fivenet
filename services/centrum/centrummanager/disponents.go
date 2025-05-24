package centrummanager

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	errorscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/errors"
	eventscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/events"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/protobuf/proto"
)

func (s *Manager) DisponentSignOn(ctx context.Context, job string, userId int32, signon bool) error {
	if signon {
		if um, ok := s.tracker.GetUserById(userId); !ok || um.Hidden {
			return errorscentrum.ErrNotOnDuty
		}

		stmt := tCentrumDisponents.
			INSERT(
				tCentrumDisponents.Job,
				tCentrumDisponents.UserID,
			).
			VALUES(
				job,
				userId,
			)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return err
			}
		}
	} else {
		stmt := tCentrumDisponents.
			DELETE().
			WHERE(jet.AND(
				tCentrumDisponents.Job.EQ(jet.String(job)),
				tCentrumDisponents.UserID.EQ(jet.Int32(userId)),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	// Load updated disponents into state
	if err := s.LoadDisponentsFromDB(ctx, job); err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	disponents, err := s.GetDisponents(ctx, job)
	if err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	change := &centrum.Disponents{
		Job:        job,
		Disponents: disponents,
	}
	data, err := proto.Marshal(change)
	if err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if _, err := s.js.Publish(ctx, eventscentrum.BuildSubject(eventscentrum.TopicGeneral, eventscentrum.TypeGeneralDisponents, job), data); err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	return nil
}
