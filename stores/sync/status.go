package syncstore

import (
	"context"
	"errors"

	database "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/common/database"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) CountJobs(ctx context.Context) (int64, error) {
	tJobs := table.FivenetJobs
	stmt := tJobs.
		SELECT(mysql.COUNT(tJobs.Name)).FROM(tJobs)

	var count database.DataCount
	if err := stmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) CountAccounts(ctx context.Context) (int64, error) {
	tUsers := table.FivenetUser
	stmt := tUsers.
		SELECT(mysql.COUNT(tUsers.ID)).FROM(tUsers)

	var count database.DataCount
	if err := stmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) CountUsers(ctx context.Context) (int64, error) {
	tUsers := table.FivenetUser
	stmt := tUsers.
		SELECT(mysql.COUNT(tUsers.ID)).FROM(tUsers)

	var count database.DataCount
	if err := stmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) CountVehicles(ctx context.Context) (int64, error) {
	tVehicles := table.FivenetOwnedVehicles
	stmt := tVehicles.
		SELECT(mysql.COUNT(tVehicles.Plate)).FROM(tVehicles)

	var count database.DataCount
	if err := stmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}

func (s *Store) CountLicenses(ctx context.Context) (int64, error) {
	tLicenses := table.FivenetLicenses
	stmt := tLicenses.
		SELECT(mysql.COUNT(tLicenses.Type)).FROM(tLicenses)

	var count database.DataCount
	if err := stmt.QueryContext(ctx, s.db, &count); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return 0, err
		}
	}

	return count.Total, nil
}
