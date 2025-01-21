package sync

import (
	"context"
	"errors"
	"fmt"
	"slices"

	jobs "github.com/fivenet-app/fivenet/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	pbsync "github.com/fivenet-app/fivenet/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils/tables"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Server) AddActivity(ctx context.Context, req *pbsync.AddActivityRequest) (*pbsync.AddActivityResponse, error) {
	resp := &pbsync.AddActivityResponse{}

	switch d := req.Activity.(type) {
	case *pbsync.AddActivityRequest_UserOauth2:
		if err := s.handleUserOauth2(ctx, d); err != nil {
			return nil, err
		}

	case *pbsync.AddActivityRequest_Dispatch:
		if _, err := s.centrum.CreateDispatch(ctx, d.Dispatch); err != nil {
			return nil, err
		}

	case *pbsync.AddActivityRequest_UserActivity:
		if err := users.CreateUserActivities(ctx, s.db, d.UserActivity); err != nil {
			return nil, err
		}

	case *pbsync.AddActivityRequest_UserProps:
		if err := s.handleUserProps(ctx, d); err != nil {
			return nil, err
		}

	case *pbsync.AddActivityRequest_JobsUserActivity:
		if err := jobs.CreateJobsUserActivities(ctx, s.db, d.JobsUserActivity); err != nil {
			return nil, err
		}

	case *pbsync.AddActivityRequest_JobsUserProps:
		if err := s.handleJobsUserProps(ctx, d); err != nil {
			return nil, err
		}

	case *pbsync.AddActivityRequest_JobsTimeclock:
		if err := s.handleTimeclockEntry(ctx, d); err != nil {
			return nil, err
		}

	case *pbsync.AddActivityRequest_UserUpdate:
		if err := s.handleUserUpdate(ctx, d); err != nil {
			return nil, err
		}
	}

	return resp, nil
}

func (s *Server) handleUserOauth2(ctx context.Context, data *pbsync.AddActivityRequest_UserOauth2) error {
	idx := slices.IndexFunc(s.cfg.OAuth2.Providers, func(in *config.OAuth2Provider) bool {
		return in.Name == data.UserOauth2.ProviderName
	})
	if idx == -1 {
		return fmt.Errorf("invalid provider name")
	}

	provider := s.cfg.OAuth2.Providers[idx]
	accountId := uint64(0)

	tAccounts := table.FivenetAccounts

	// Retrieve account via identifier
	stmt := tAccounts.
		SELECT(
			tAccounts.ID,
		).
		FROM(tAccounts).
		WHERE(tAccounts.License.EQ(jet.String(data.UserOauth2.Identifier))).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, &accountId); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	if accountId == 0 {
		return fmt.Errorf("no fivenet account found for identifier")
	}

	tOAuth2Accs := table.FivenetOauth2Accounts

	insertStmt := tOAuth2Accs.
		INSERT(
			tOAuth2Accs.AccountID,
			tOAuth2Accs.Provider,
			tOAuth2Accs.ExternalID,
			tOAuth2Accs.Username,
			tOAuth2Accs.Avatar,
		).
		VALUES(
			accountId,
			provider.Name,
			data.UserOauth2.ExternalId,
			data.UserOauth2.Username,
			provider.DefaultAvatar,
		)

	if _, err := insertStmt.ExecContext(ctx, s.db); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return nil
		}
	}

	return nil
}

func (s *Server) handleUserProps(ctx context.Context, data *pbsync.AddActivityRequest_UserProps) error {
	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	props, err := users.GetUserProps(ctx, tx, data.UserProps.Props.UserId, nil)
	if err != nil {
		return err
	}

	reason := ""
	if data.UserProps.Reason != nil {
		reason = *data.UserProps.Reason
	}

	activities, err := props.HandleChanges(ctx, tx, data.UserProps.Props, &data.UserProps.Props.UserId, reason)
	if err != nil {
		return err
	}

	if data.UserProps.Reason != nil && *data.UserProps.Reason != "" {
		if err := users.CreateUserActivities(ctx, tx, activities...); err != nil {
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Server) handleJobsUserProps(ctx context.Context, data *pbsync.AddActivityRequest_JobsUserProps) error {
	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	props, err := jobs.GetJobsUserProps(ctx, tx, data.JobsUserProps.Props.Job, data.JobsUserProps.Props.UserId, nil)
	if err != nil {
		return err
	}

	reason := ""
	if data.JobsUserProps.Reason != nil {
		reason = *data.JobsUserProps.Reason
	}

	activities, err := props.HandleChanges(ctx, tx, data.JobsUserProps.Props, data.JobsUserProps.Props.Job, &data.JobsUserProps.Props.UserId, reason)
	if err != nil {
		return err
	}

	if data.JobsUserProps.Reason != nil && *data.JobsUserProps.Reason != "" {
		if err := jobs.CreateJobsUserActivities(ctx, tx, activities...); err != nil {
			return err
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Server) handleTimeclockEntry(ctx context.Context, data *pbsync.AddActivityRequest_JobsTimeclock) error {
	tTimeClock := table.FivenetJobsTimeclock

	stmt := tTimeClock.
		INSERT(
			tTimeClock.Job,
			tTimeClock.UserID,
			tTimeClock.Date,
			tTimeClock.StartTime,
			tTimeClock.EndTime,
			tTimeClock.SpentTime,
		).
		VALUES(
			data.JobsTimeclock.Job,
			data.JobsTimeclock.UserId,
			data.JobsTimeclock.Date,
			data.JobsTimeclock.StartTime,
			data.JobsTimeclock.EndTime,
			data.JobsTimeclock.SpentTime,
		).
		ON_DUPLICATE_KEY_UPDATE(
			tTimeClock.StartTime.SET(jet.TimestampExp(jet.Raw("VALUES(`start_time`)"))),
			tTimeClock.EndTime.SET(jet.TimestampExp(jet.Raw("VALUES(`end_time`)"))),
			tTimeClock.SpentTime.SET(jet.FloatExp(jet.Raw("VALUES(`spent_time`)"))),
		)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return err
	}

	return nil
}

func (s *Server) handleUserUpdate(ctx context.Context, data *pbsync.AddActivityRequest_UserUpdate) error {
	tUser := tables.Users()

	selectStmt := tUser.
		SELECT(
			tUser.ID,
		).
		FROM(tUser).
		WHERE(tUser.ID.EQ(jet.Int32(data.UserUpdate.UserId))).
		LIMIT(1)

	user := &users.User{}
	if err := selectStmt.QueryContext(ctx, s.db, user); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}

		return nil
	}

	updateSets := []jet.ColumnAssigment{}
	if data.UserUpdate.Group != nil {
		updateSets = append(updateSets, tUser.Group.SET(jet.String(*data.UserUpdate.Group)))
	}
	if data.UserUpdate.Job != nil {
		updateSets = append(updateSets, tUser.Job.SET(jet.String(*data.UserUpdate.Job)))
	}
	if data.UserUpdate.JobGrade != nil {
		updateSets = append(updateSets, tUser.JobGrade.SET(jet.Int32(*data.UserUpdate.JobGrade)))
	}
	if data.UserUpdate.Firstname != nil {
		updateSets = append(updateSets, tUser.Firstname.SET(jet.String(*data.UserUpdate.Firstname)))
	}
	if data.UserUpdate.Lastname != nil {
		updateSets = append(updateSets, tUser.Lastname.SET(jet.String(*data.UserUpdate.Lastname)))
	}

	if len(updateSets) > 0 {
		stmt := tUser.
			UPDATE().
			SET(updateSets[0]).
			WHERE(tUser.ID.EQ(jet.Int32(data.UserUpdate.UserId)))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return err
		}
	}

	return nil
}
