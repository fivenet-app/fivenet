package authstore

import (
	"context"

	"github.com/fivenet-app/fivenet/v2026/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

var tAccounts = table.FivenetAccounts

func (s *Store) getAccount(
	ctx context.Context,
	condition mysql.BoolExpression,
	withPass bool,
) (*model.FivenetAccounts, error) {
	tAccounts := table.FivenetAccounts.AS("account")
	columns := mysql.ProjectionList{
		tAccounts.ID,
		tAccounts.CreatedAt,
		tAccounts.UpdatedAt,
		tAccounts.DeletedAt,
		tAccounts.Enabled,
		tAccounts.Username,
		tAccounts.License,
		tAccounts.RegToken,
		tAccounts.Groups,
		tAccounts.LastChar,
	}
	if withPass {
		columns = append(columns, tAccounts.Password)
	}

	stmt := tAccounts.
		SELECT(
			columns[0],
			columns[1:]...,
		).
		FROM(tAccounts).
		WHERE(mysql.AND(
			tAccounts.Enabled.IS_TRUE(),
			tAccounts.DeletedAt.IS_NULL(),
			condition,
		)).
		LIMIT(1)

	acc := &model.FivenetAccounts{}
	if err := stmt.QueryContext(ctx, s.db, acc); err != nil {
		return nil, err
	}

	return acc, nil
}

func (s *Store) GetAccountByID(
	ctx context.Context,
	accountID int64,
	withPassword bool,
) (*model.FivenetAccounts, error) {
	return s.getAccount(ctx, tAccounts.ID.EQ(mysql.Int64(accountID)), withPassword)
}

func (s *Store) GetAccountByUsername(
	ctx context.Context,
	username string,
	withPassword bool,
) (*model.FivenetAccounts, error) {
	return s.getAccount(ctx, tAccounts.Username.EQ(mysql.String(username)), withPassword)
}

func (s *Store) GetLoginAccountByUsername(
	ctx context.Context,
	username string,
) (*model.FivenetAccounts, error) {
	return s.getAccount(ctx, mysql.AND(
		tAccounts.Username.EQ(mysql.String(username)),
		tAccounts.RegToken.IS_NULL(),
		tAccounts.Password.IS_NOT_NULL(),
	), true)
}

func (s *Store) GetAccountByIDAndUsername(
	ctx context.Context,
	accountID int64,
	username string,
	withPassword bool,
) (*model.FivenetAccounts, error) {
	return s.getAccount(ctx, mysql.AND(
		tAccounts.ID.EQ(mysql.Int64(accountID)),
		tAccounts.Username.EQ(mysql.String(username)),
	), withPassword)
}

func (s *Store) GetAccountByRegToken(
	ctx context.Context,
	regToken string,
	withPassword bool,
) (*model.FivenetAccounts, error) {
	return s.getAccount(ctx, tAccounts.RegToken.EQ(mysql.String(regToken)), withPassword)
}

func (s *Store) GetPasswordResetAccountByRegToken(
	ctx context.Context,
	regToken string,
) (*model.FivenetAccounts, error) {
	return s.getAccount(ctx, mysql.AND(
		tAccounts.RegToken.EQ(mysql.String(regToken)),
		tAccounts.Username.IS_NOT_NULL(),
		tAccounts.Password.IS_NULL(),
	), true)
}

func (s *Store) ActivateAccount(
	ctx context.Context,
	accountID int64,
	regToken, username, hashedPassword string,
) error {
	stmt := tAccounts.
		UPDATE(
			tAccounts.Username,
			tAccounts.Password,
			tAccounts.RegToken,
		).
		SET(
			tAccounts.Username.SET(mysql.String(username)),
			tAccounts.Password.SET(mysql.String(hashedPassword)),
			tAccounts.RegToken.SET(mysql.StringExp(mysql.NULL)),
		).
		WHERE(mysql.AND(
			tAccounts.ID.EQ(mysql.Int64(accountID)),
			tAccounts.RegToken.EQ(mysql.String(regToken)),
		)).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, s.db)
	return err
}

func (s *Store) UpdatePassword(ctx context.Context, accountID int64, hashedPassword string) error {
	stmt := tAccounts.
		UPDATE(tAccounts.Password).
		SET(tAccounts.Password.SET(mysql.String(hashedPassword))).
		WHERE(tAccounts.ID.EQ(mysql.Int64(accountID))).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, s.db)
	return err
}

func (s *Store) UpdateUsername(ctx context.Context, accountID int64, username string) error {
	stmt := tAccounts.
		UPDATE(tAccounts.Username).
		SET(tAccounts.Username.SET(mysql.String(username))).
		WHERE(tAccounts.ID.EQ(mysql.Int64(accountID))).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, s.db)
	return err
}

func (s *Store) ForgotPassword(ctx context.Context, accountID int64, hashedPassword string) error {
	stmt := tAccounts.
		UPDATE(
			tAccounts.Password,
			tAccounts.RegToken,
		).
		SET(
			tAccounts.Password.SET(mysql.String(hashedPassword)),
			tAccounts.RegToken.SET(mysql.StringExp(mysql.NULL)),
		).
		WHERE(tAccounts.ID.EQ(mysql.Int64(accountID))).
		LIMIT(1)

	_, err := stmt.ExecContext(ctx, s.db)
	return err
}
