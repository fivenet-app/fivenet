package mailerstore

import (
	"context"
	"errors"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	maileraccess "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/access"
	maileremails "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/emails"
	mailersettings "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/settings"
	mailerthreads "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/mailer/threads"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	usershort "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/short"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

var tEmails = table.FivenetMailerEmails.AS("email")

func (s *Store) dbOr(q qrm.DB) qrm.DB {
	if q != nil {
		return q
	}
	return s.db
}

func (s *Store) CountEmails(
	ctx context.Context,
	q qrm.DB,
	condition mysql.BoolExpression,
) (int64, error) {
	stmt := tEmails.
		SELECT(
			mysql.COUNT(tEmails.ID).AS("data_count.total"),
		).
		FROM(tEmails).
		WHERE(condition)

	var count database.DataCount
	if err := stmt.QueryContext(ctx, s.dbOr(q), &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) ListEmails(
	ctx context.Context,
	db qrm.DB,
	userInfo *userinfo.UserInfo,
	pag *database.PaginationRequest,
	all bool,
) (*database.PaginationResponse, []*maileremails.Email, error) {
	if pag == nil {
		pag = &database.PaginationRequest{}
	}

	if userInfo != nil && userInfo.GetSuperuser() && all {
		return s.listAllEmails(ctx, db, pag)
	}

	includeDeleted := userInfo != nil && userInfo.GetSuperuser()

	visibleIDs := s.subjectAccess.VisibleIDsByConditionQuery(
		userInfo,
		int32(maileraccess.AccessLevel_ACCESS_LEVEL_READ),
		includeDeleted,
		mysql.Bool(true),
	)
	ctes := visibleIDs.CTEs
	visibleEmailID := mysql.IntegerColumn("id").From(visibleIDs.Table)

	var countStmt mysql.Statement = tEmails.
		SELECT(mysql.COUNT(visibleEmailID).AS("data_count.total")).
		FROM(visibleIDs.Table)

	if len(ctes) > 0 {
		countStmt = mysql.WITH(ctes...)(countStmt)
	}

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.dbOr(db), &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, nil, err
		}
	}

	pagination, limit := pag.GetResponseWithPageSize(count.Total, 20)
	if count.Total <= 0 {
		return pagination, []*maileremails.Email{}, nil
	}

	var stmt mysql.Statement = tEmails.
		SELECT(
			tEmails.ID,
			tEmails.CreatedAt,
			tEmails.UpdatedAt,
			tEmails.DeletedAt,
			tEmails.Deactivated,
			tEmails.Job,
			tEmails.UserID,
			tEmails.Email,
			tEmails.EmailChanged,
			tEmails.Label,
		).
		FROM(
			visibleIDs.Table.
				INNER_JOIN(tEmails,
					tEmails.ID.EQ(visibleEmailID),
				),
		).
		ORDER_BY(
			tEmails.Job.ASC(),
			tEmails.Label.ASC(),
		).
		OFFSET(pagination.GetOffset()).
		LIMIT(limit)

	if len(ctes) > 0 {
		stmt = mysql.WITH(ctes...)(stmt)
	}

	emails := []*maileremails.Email{}
	if err := stmt.QueryContext(ctx, s.dbOr(db), &emails); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, nil, err
		}
	}

	return pagination, emails, nil
}

func (s *Store) listAllEmails(
	ctx context.Context,
	db qrm.DB,
	pag *database.PaginationRequest,
) (*database.PaginationResponse, []*maileremails.Email, error) {
	var countStmt mysql.Statement = tEmails.
		SELECT(mysql.COUNT(tEmails.ID).AS("data_count.total")).
		FROM(tEmails)

	var count database.DataCount
	if err := countStmt.QueryContext(ctx, s.dbOr(db), &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, nil, err
		}
	}

	pagination, limit := pag.GetResponseWithPageSize(count.Total, 20)
	if count.Total <= 0 {
		return pagination, []*maileremails.Email{}, nil
	}

	var stmt mysql.Statement = tEmails.
		SELECT(
			tEmails.ID,
			tEmails.CreatedAt,
			tEmails.UpdatedAt,
			tEmails.DeletedAt,
			tEmails.Deactivated,
			tEmails.Job,
			tEmails.UserID,
			tEmails.Email,
			tEmails.EmailChanged,
			tEmails.Label,
		).
		FROM(tEmails).
		ORDER_BY(
			tEmails.Job.ASC(),
			tEmails.Label.ASC(),
		).
		OFFSET(pagination.GetOffset()).
		LIMIT(limit)

	emails := []*maileremails.Email{}
	if err := stmt.QueryContext(ctx, s.dbOr(db), &emails); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, nil, err
		}
	}

	return pagination, emails, nil
}

func (s *Store) GetEmailByCondition(
	ctx context.Context,
	q qrm.DB,
	condition mysql.BoolExpression,
) (*maileremails.Email, error) {
	stmt := tEmails.
		SELECT(
			tEmails.ID,
			tEmails.CreatedAt,
			tEmails.UpdatedAt,
			tEmails.DeletedAt,
			tEmails.Deactivated,
			tEmails.Job,
			tEmails.UserID,
			tEmails.Email,
			tEmails.EmailChanged,
			tEmails.Label,
		).
		FROM(tEmails).
		WHERE(condition).
		LIMIT(1)

	dest := &maileremails.Email{}
	if err := stmt.QueryContext(ctx, s.dbOr(q), dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	if dest.GetId() == 0 {
		return nil, nil
	}

	return dest, nil
}

func (s *Store) GetEmail(
	ctx context.Context,
	q qrm.DB,
	emailID int64,
	includeDeleted bool,
) (*maileremails.Email, error) {
	condition := tEmails.ID.EQ(mysql.Int64(emailID))
	if !includeDeleted {
		condition = mysql.AND(condition, tEmails.DeletedAt.IS_NULL())
	}
	return s.GetEmailByCondition(ctx, q, condition)
}

func (s *Store) GetUserShort(
	ctx context.Context,
	q qrm.DB,
	userID int32,
) (*usershort.UserShort, error) {
	tUsers := table.FivenetUser.AS("user_short")

	stmt := tUsers.
		SELECT(
			tUsers.Firstname,
			tUsers.Lastname,
			tUsers.Dateofbirth,
		).
		FROM(tUsers).
		WHERE(tUsers.ID.EQ(mysql.Int32(userID))).
		LIMIT(1)

	dest := &usershort.UserShort{}
	if err := stmt.QueryContext(ctx, s.dbOr(q), dest); err != nil {
		return nil, err
	}

	return dest, nil
}

func (s *Store) ListRecipientsByEmails(
	ctx context.Context,
	q qrm.DB,
	recipients []string,
) ([]*mailerthreads.ThreadRecipientEmail, error) {
	if len(recipients) == 0 {
		return []*mailerthreads.ThreadRecipientEmail{}, nil
	}

	expr := make([]mysql.Expression, len(recipients))
	for idx := range recipients {
		expr[idx] = mysql.String(recipients[idx])
	}

	stmt := tEmails.
		SELECT(
			tEmails.ID.AS("thread_recipient_email.email_id"),
			tEmails.Email,
			tEmails.Deactivated,
		).
		FROM(tEmails).
		WHERE(mysql.AND(
			tEmails.Email.IN(expr...),
			tEmails.DeletedAt.IS_NULL(),
		)).
		LIMIT(int64(len(recipients)))

	dest := []*mailerthreads.ThreadRecipientEmail{}
	if err := stmt.QueryContext(ctx, s.dbOr(q), &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}

func (s *Store) GetEmailSettings(
	ctx context.Context,
	q qrm.DB,
	emailID int64,
) (*mailersettings.EmailSettings, error) {
	tSettings := table.FivenetMailerSettings.AS("email_settings")
	stmt := tSettings.
		SELECT(
			tSettings.EmailID,
			tSettings.Signature,
			table.FivenetMailerSettingsBlocked.TargetEmail.AS("email_settings.blocked_emails"),
		).
		FROM(
			tSettings.
				LEFT_JOIN(table.FivenetMailerSettingsBlocked,
					table.FivenetMailerSettingsBlocked.EmailID.EQ(tSettings.EmailID),
				),
		).
		WHERE(tSettings.EmailID.EQ(mysql.Int64(emailID))).
		LIMIT(25)

	dest := &mailersettings.EmailSettings{EmailId: emailID}
	if err := stmt.QueryContext(ctx, s.dbOr(q), dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return dest, nil
}
