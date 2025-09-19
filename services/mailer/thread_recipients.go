package mailer

import (
	"context"
	"errors"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/mailer"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/userinfo"
	"github.com/fivenet-app/fivenet/v2025/pkg/grpc/errswrap"
	errorsmailer "github.com/fivenet-app/fivenet/v2025/services/mailer/errors"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

type canAccessIdsHelper struct {
	IDs []int64 `alias:"id"`
}

func (s *Server) checkIfEmailPartOfThread(
	ctx context.Context,
	userInfo *userinfo.UserInfo,
	threadId int64,
	emailId int64,
	accessLevel mailer.AccessLevel,
) error {
	check, err := s.access.CanUserAccessTarget(ctx, emailId, userInfo, accessLevel)
	if err != nil {
		return errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return errorsmailer.ErrThreadAccessDenied
	}
	check, err = s.checkIfEmailIdPartOfThread(ctx, threadId, emailId)
	if err != nil {
		return errswrap.NewError(err, errorsmailer.ErrFailedQuery)
	}
	if !check && !userInfo.GetSuperuser() {
		return errorsmailer.ErrThreadAccessDenied
	}

	return nil
}

func (s *Server) checkIfEmailIdPartOfThread(
	ctx context.Context,
	threadId int64,
	emailId int64,
) (bool, error) {
	stmt := tThreadsRecipients.
		SELECT(
			tThreadsRecipients.ID.AS("id"),
		).
		FROM(tThreadsRecipients).
		WHERE(mysql.AND(
			tThreadsRecipients.ThreadID.EQ(mysql.Int64(threadId)),
			tThreadsRecipients.EmailID.EQ(mysql.Int64(emailId)),
		))

	dest := &canAccessIdsHelper{}
	if err := stmt.QueryContext(ctx, s.db, &dest.IDs); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return false, err
		}
	}

	return len(dest.IDs) == 1, nil
}
