package mailer

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/content"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) UpsertEmailSettingsSignature(
	ctx context.Context,
	q qrm.DB,
	emailID int64,
	signature *content.Content,
) error {
	stmt := table.FivenetMailerSettings.
		INSERT(
			table.FivenetMailerSettings.EmailID,
			table.FivenetMailerSettings.Signature,
		).
		VALUES(
			emailID,
			signature,
		).
		ON_DUPLICATE_KEY_UPDATE(
			table.FivenetMailerSettings.Signature.SET(mysql.String("VALUES(`signature`)")),
		)

	if _, err := stmt.ExecContext(ctx, s.dbOr(q)); err != nil {
		return err
	}

	return nil
}

func (s *Store) AddBlockedEmails(
	ctx context.Context,
	q qrm.DB,
	emailID int64,
	blockedEmails []string,
) error {
	if len(blockedEmails) == 0 {
		return nil
	}

	stmt := table.FivenetMailerSettingsBlocked.
		INSERT(
			table.FivenetMailerSettingsBlocked.EmailID,
			table.FivenetMailerSettingsBlocked.TargetEmail,
		)

	for _, be := range blockedEmails {
		stmt = stmt.VALUES(emailID, be)
	}

	if _, err := stmt.ExecContext(ctx, s.dbOr(q)); err != nil {
		return err
	}

	return nil
}

func (s *Store) DeleteBlockedEmails(
	ctx context.Context,
	q qrm.DB,
	emailID int64,
	blockedEmails []string,
) error {
	if len(blockedEmails) == 0 {
		return nil
	}

	targets := make([]mysql.Expression, 0, len(blockedEmails))
	for _, be := range blockedEmails {
		targets = append(targets, mysql.String(be))
	}

	stmt := table.FivenetMailerSettingsBlocked.
		DELETE().
		WHERE(mysql.AND(
			table.FivenetMailerSettingsBlocked.EmailID.EQ(mysql.Int64(emailID)),
			table.FivenetMailerSettingsBlocked.TargetEmail.IN(targets...),
		)).
		LIMIT(int64(len(blockedEmails)))

	if _, err := stmt.ExecContext(ctx, s.dbOr(q)); err != nil {
		return err
	}

	return nil
}
