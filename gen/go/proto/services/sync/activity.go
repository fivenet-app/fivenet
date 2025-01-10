package sync

import (
	"context"
	"fmt"
	"slices"

	jobs "github.com/fivenet-app/fivenet/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/pkg/config"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (s *Server) AddActivity(ctx context.Context, req *AddActivityRequest) (*AddActivityResponse, error) {
	resp := &AddActivityResponse{}

	switch d := req.Activity.(type) {
	case *AddActivityRequest_UserOauth2:
		if err := s.handleUserOauth2(ctx, d); err != nil {
			return nil, err
		}

	case *AddActivityRequest_UserActivity:
		if err := users.CreateUserActivities(ctx, s.db, d.UserActivity); err != nil {
			return nil, err
		}

	case *AddActivityRequest_UserProps:
		if err := s.handleUserProps(ctx, d); err != nil {
			return nil, err
		}

	case *AddActivityRequest_JobsUserActivity:
		if err := jobs.CreateJobsUserActivities(ctx, s.db, d.JobsUserActivity); err != nil {
			return nil, err
		}

	case *AddActivityRequest_JobsUserProps:
		if err := s.handleJobsUserProps(ctx, d); err != nil {
			return nil, err
		}

	case *AddActivityRequest_JobsTimeclock:
		if err := s.handleTimeclockEntry(ctx, d); err != nil {
			return nil, err
		}

	}

	return resp, nil
}

func (s *Server) handleUserOauth2(ctx context.Context, data *AddActivityRequest_UserOauth2) error {
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
		return err
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

func (s *Server) handleUserProps(ctx context.Context, data *AddActivityRequest_UserProps) error {
	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	// TODO retrieve current user props
	props := &users.UserProps{}

	activities, err := props.HandleChanges(ctx, tx, data.UserProps.Props, &data.UserProps.Props.UserId, data.UserProps.Reason)
	if err != nil {
		return err
	}

	if err := users.CreateUserActivities(ctx, tx, activities...); err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Server) handleJobsUserProps(ctx context.Context, data *AddActivityRequest_JobsUserProps) error {
	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	// TODO retrieve current jobs user props
	props := &jobs.JobsUserProps{}

	activities, err := props.HandleChanges(ctx, tx, data.JobsUserProps.Props, data.JobsUserProps.Props.Job, &data.JobsUserProps.Props.UserId, data.JobsUserProps.Reason)
	if err != nil {
		return err
	}

	if err := jobs.CreateJobsUserActivities(ctx, tx, activities...); err != nil {
		return err
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (s *Server) handleTimeclockEntry(ctx context.Context, data *AddActivityRequest_JobsTimeclock) error {
	// TODO

	return nil
}
