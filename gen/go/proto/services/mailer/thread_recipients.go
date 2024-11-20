package mailer

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/mailer"
	errorsmailer "github.com/fivenet-app/fivenet/gen/go/proto/services/mailer/errors"
	"github.com/fivenet-app/fivenet/pkg/grpc/auth/userinfo"
	"github.com/fivenet-app/fivenet/pkg/grpc/errswrap"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type canAccessIdsHelper struct {
	IDs []uint64 `alias:"id"`
}

func (s *Server) checkIfEmailPartOfThread(ctx context.Context, userInfo *userinfo.UserInfo, threadId uint64, emailId uint64, accessLevel mailer.AccessLevel) error {
	check, err := s.access.CanUserAccessTarget(ctx, emailId, userInfo, accessLevel)
	if err != nil {
		return errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return errorsmailer.ErrThreadAccessDenied
	}
	check, err = s.checkIfEmailIdPartOfThread(ctx, threadId, emailId)
	if err != nil {
		return errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check && !userInfo.SuperUser {
		return errorsmailer.ErrThreadAccessDenied
	}

	return nil
}

func (s *Server) checkIfEmailIdPartOfThread(ctx context.Context, threadId uint64, emailId uint64) (bool, error) {
	stmt := tThreadsRecipients.
		SELECT(
			tThreadsRecipients.ID.AS("id"),
		).
		FROM(tThreadsRecipients).
		WHERE(jet.AND(
			tThreadsRecipients.ThreadID.EQ(jet.Uint64(threadId)),
			tThreadsRecipients.EmailID.EQ(jet.Uint64(emailId)),
		))

	dest := &canAccessIdsHelper{}
	if err := stmt.QueryContext(ctx, s.db, &dest.IDs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return false, err
		}
	}

	return len(dest.IDs) == 1, nil
}
