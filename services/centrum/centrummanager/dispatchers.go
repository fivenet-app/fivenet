package centrummanager

import (
	"context"

	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	errorscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/errors"
	eventscentrum "github.com/fivenet-app/fivenet/v2025/services/centrum/events"
	jet "github.com/go-jet/jet/v2/mysql"
	"google.golang.org/protobuf/proto"
)

func (s *Manager) DispatcherSignOn(ctx context.Context, job string, userId int32, signon bool) error {
	if signon {
		if um, ok := s.tracker.GetUserById(userId); !ok || um.Hidden {
			return errorscentrum.ErrNotOnDuty
		}

		stmt := tCentrumDispatchers.
			INSERT(
				tCentrumDispatchers.Job,
				tCentrumDispatchers.UserID,
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
		stmt := tCentrumDispatchers.
			DELETE().
			WHERE(jet.AND(
				tCentrumDispatchers.Job.EQ(jet.String(job)),
				tCentrumDispatchers.UserID.EQ(jet.Int32(userId)),
			)).
			LIMIT(1)

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
		}
	}

	// Load updated dispatchers into state
	if err := s.LoadDispatchersFromDB(ctx, job); err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	dispatchers, err := s.GetDispatchers(ctx, job)
	if err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	data, err := proto.Marshal(dispatchers)
	if err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	if _, err := s.js.Publish(ctx, eventscentrum.BuildSubject(eventscentrum.TopicGeneral, eventscentrum.TypeGeneralDispatchers, job), data); err != nil {
		return errswrap.NewError(err, errorscentrum.ErrFailedQuery)
	}

	return nil
}
