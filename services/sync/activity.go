package sync

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	jobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/users"
	pbsync "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

func (s *Server) AddActivity(
	ctx context.Context,
	req *pbsync.AddActivityRequest,
) (*pbsync.AddActivityResponse, error) {
	resp := &pbsync.AddActivityResponse{}

	s.lastSyncedActivity.Store(time.Now().Unix())

	switch d := req.GetActivity().(type) {
	case *pbsync.AddActivityRequest_UserOauth2:
		if err := s.handleUserOauth2(ctx, d); err != nil {
			return nil, fmt.Errorf("failed to handle UserOauth2 activity. %w", err)
		}

	case *pbsync.AddActivityRequest_Dispatch:
		if _, err := s.dispatches.Create(ctx, d.Dispatch); err != nil {
			return nil, fmt.Errorf("failed to create dispatch. %w", err)
		}

	case *pbsync.AddActivityRequest_UserActivity:
		if err := users.CreateUserActivities(ctx, s.db, d.UserActivity); err != nil {
			return nil, fmt.Errorf("failed to create user activities. %w", err)
		}

	case *pbsync.AddActivityRequest_UserProps:
		if err := s.handleUserProps(ctx, d); err != nil {
			return nil, fmt.Errorf("failed to handle UserProps activity. %w", err)
		}

	case *pbsync.AddActivityRequest_ColleagueActivity:
		if err := jobs.CreateColleagueActivity(ctx, s.db, d.ColleagueActivity); err != nil {
			return nil, fmt.Errorf("failed to create jobs user activities. %w", err)
		}

	case *pbsync.AddActivityRequest_ColleagueProps:
		if err := s.handleColleagueProps(ctx, d); err != nil {
			return nil, fmt.Errorf("failed to handle ColleagueProps activity. %w", err)
		}

	case *pbsync.AddActivityRequest_JobTimeclock:
		if err := s.handleTimeclockEntry(ctx, d); err != nil {
			return nil, fmt.Errorf("failed to handle JobTimeclock activity. %w", err)
		}

	case *pbsync.AddActivityRequest_UserUpdate:
		if s.esxCompat {
			return nil, fmt.Errorf("ESX compatibility mode is enabled, cannot send data. %w", ErrSendDataDisabled)
		}

		if err := s.handleUserUpdate(ctx, d); err != nil {
			return nil, fmt.Errorf("failed to handle UserUpdate activity. %w", err)
		}
	}

	return resp, nil
}

func (s *Server) handleUserOauth2(
	ctx context.Context,
	data *pbsync.AddActivityRequest_UserOauth2,
) error {
	idx := slices.IndexFunc(s.cfg.OAuth2.Providers, func(in *config.OAuth2Provider) bool {
		return in.Name == data.UserOauth2.GetProviderName()
	})
	if idx == -1 {
		return fmt.Errorf("invalid provider name. %s", data.UserOauth2.GetProviderName())
	}

	provider := s.cfg.OAuth2.Providers[idx]
	tAccounts := table.FivenetAccounts

	// Struct to hold the query result
	type Account struct {
		ID int64
	}
	var account Account

	// Retrieve account via identifier
	stmt := tAccounts.
		SELECT(
			tAccounts.ID,
		).
		FROM(tAccounts).
		WHERE(tAccounts.License.EQ(jet.String(data.UserOauth2.GetIdentifier()))).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, &account); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query account by identifier. %w", err)
		}
	}

	if account.ID == 0 {
		s.logger.Warn(
			"no fivenet account found for identifier in user oauth2 sync connect",
			zap.String("provider", data.UserOauth2.GetProviderName()),
			zap.String(
				"identifier",
				data.UserOauth2.GetIdentifier(),
			),
			zap.String("external_id", data.UserOauth2.GetExternalId()),
		)
		return nil
	}

	tOAuth2Accs := table.FivenetAccountsOauth2

	insertStmt := tOAuth2Accs.
		INSERT(
			tOAuth2Accs.AccountID,
			tOAuth2Accs.Provider,
			tOAuth2Accs.ExternalID,
			tOAuth2Accs.Username,
			tOAuth2Accs.Avatar,
		).
		VALUES(
			account.ID,
			provider.Name,
			data.UserOauth2.GetExternalId(),
			data.UserOauth2.GetUsername(),
			provider.DefaultAvatar,
		)

	if _, err := insertStmt.ExecContext(ctx, s.db); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return fmt.Errorf("failed to insert OAuth2 account. %w", err)
		}
	}

	return nil
}

func (s *Server) handleUserProps(
	ctx context.Context,
	data *pbsync.AddActivityRequest_UserProps,
) error {
	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction. %w", err)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	reqP := data.UserProps.GetProps()
	props, err := users.GetUserProps(ctx, tx, reqP.GetUserId(), nil)
	if err != nil {
		return fmt.Errorf("failed to get user props. %w", err)
	}

	reason := ""
	if data.UserProps.Reason != nil {
		reason = data.UserProps.GetReason()
	}

	activities, err := props.HandleChanges(
		ctx,
		tx,
		reqP,
		&reqP.UserId,
		reason,
	)
	if err != nil {
		return fmt.Errorf("failed to handle user props changes. %w", err)
	}

	if data.UserProps.Reason != nil && data.UserProps.GetReason() != "" {
		if err := users.CreateUserActivities(ctx, tx, activities...); err != nil {
			return fmt.Errorf("failed to create user activities. %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction. %w", err)
	}

	return nil
}

func (s *Server) handleColleagueProps(
	ctx context.Context,
	data *pbsync.AddActivityRequest_ColleagueProps,
) error {
	// Begin transaction
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction. %w", err)
	}
	// Defer a rollback in case anything fails
	defer tx.Rollback()

	current := data.ColleagueProps.GetProps()
	props, err := jobs.GetColleagueProps(
		ctx,
		tx,
		current.GetJob(),
		current.GetUserId(),
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to get jobs user props. %w", err)
	}

	reason := ""
	if data.ColleagueProps.Reason != nil {
		reason = data.ColleagueProps.GetReason()
	}

	activities, err := props.HandleChanges(
		ctx,
		tx,
		current,
		current.GetJob(),
		&current.UserId,
		reason,
	)
	if err != nil {
		return fmt.Errorf("failed to handle jobs user props changes. %w", err)
	}

	if data.ColleagueProps.Reason != nil && data.ColleagueProps.GetReason() != "" {
		if err := jobs.CreateColleagueActivity(ctx, tx, activities...); err != nil {
			return fmt.Errorf("failed to create jobs user activities. %w", err)
		}
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction. %w", err)
	}

	return nil
}

func (s *Server) handleTimeclockEntry(
	ctx context.Context,
	data *pbsync.AddActivityRequest_JobTimeclock,
) error {
	tTimeClock := table.FivenetJobTimeclock

	d := data.JobTimeclock
	if d.GetStart() {
		// Run select query to see if a timeclock entry needs to be created
		stmt := tTimeClock.
			SELECT(
				tTimeClock.UserID,
				tTimeClock.Date,
				tTimeClock.EndTime,
			).
			FROM(tTimeClock).
			WHERE(jet.AND(
				tTimeClock.Job.EQ(jet.String(d.GetJob())),
				tTimeClock.UserID.EQ(jet.Int32(d.GetUserId())),
			)).
			ORDER_BY(tTimeClock.Date.DESC()).
			LIMIT(1)

		dest := &jobs.TimeclockEntry{}
		if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return fmt.Errorf("failed to query timeclock entry. %w", err)
			}
		}

		// Found an entry, no need to create a new one
		if dest.GetUserId() > 0 {
			return nil
		}

		updateStmt := tTimeClock.
			INSERT(
				tTimeClock.Job,
				tTimeClock.UserID,
				tTimeClock.Date,
			).
			VALUES(
				d.GetJob(),
				d.GetUserId(),
				jet.CURRENT_DATE(),
			)

		if _, err := updateStmt.ExecContext(ctx, s.db); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return fmt.Errorf("failed to insert timeclock entry. %w", err)
			}
		}
	} else {
		stmt := tTimeClock.
			UPDATE().
			SET(
				tTimeClock.SpentTime.SET(jet.FloatExp(jet.Raw("`spent_time` + CAST((TIMESTAMPDIFF(SECOND, `start_time`, `end_time`) / 3600) AS DECIMAL(10,2))"))),
				tTimeClock.EndTime.SET(jet.CURRENT_TIMESTAMP()),
			).
			WHERE(jet.AND(
				tTimeClock.Job.EQ(jet.String(d.GetJob())),
				tTimeClock.UserID.EQ(jet.Int32(d.GetUserId())),
				tTimeClock.StartTime.IS_NOT_NULL(),
				tTimeClock.EndTime.IS_NULL(),
			))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return fmt.Errorf("failed to update timeclock entry. %w", err)
		}
	}

	return nil
}

func (s *Server) handleUserUpdate(
	ctx context.Context,
	data *pbsync.AddActivityRequest_UserUpdate,
) error {
	d := data.UserUpdate

	tUser := tables.User()

	selectStmt := tUser.
		SELECT(
			tUser.ID,
		).
		FROM(tUser).
		WHERE(tUser.ID.EQ(jet.Int32(d.GetUserId()))).
		LIMIT(1)

	user := &users.User{}
	if err := selectStmt.QueryContext(ctx, s.db, user); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query user by ID. %w", err)
		}

		return nil
	}

	updateSets := []jet.ColumnAssigment{}
	if d.Group != nil {
		updateSets = append(updateSets, tUser.Group.SET(jet.String(d.GetGroup())))
	}
	if d.Job != nil {
		updateSets = append(updateSets, tUser.Job.SET(jet.String(d.GetJob())))
	}
	if d.JobGrade != nil {
		updateSets = append(
			updateSets,
			tUser.JobGrade.SET(jet.Int32(d.GetJobGrade())),
		)
	}
	if d.Firstname != nil {
		updateSets = append(
			updateSets,
			tUser.Firstname.SET(jet.String(d.GetFirstname())),
		)
	}
	if d.Lastname != nil {
		updateSets = append(
			updateSets,
			tUser.Lastname.SET(jet.String(d.GetLastname())),
		)
	}

	if len(updateSets) > 0 {
		stmt := tUser.
			UPDATE().
			SET(updateSets[0]).
			WHERE(tUser.ID.EQ(jet.Int32(d.GetUserId())))

		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return fmt.Errorf("failed to update user. %w", err)
		}
	}

	return nil
}
