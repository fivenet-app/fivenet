package syncstore

import (
	"context"
	"errors"
	"fmt"

	colleaguesactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/colleagues/activity"
	jobstimeclock "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/timeclock"
	activity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/activity"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) AddActivity(
	ctx context.Context,
	req *pbsync.AddActivityRequest,
) (*pbsync.AddActivityResponse, error) {
	resp := &pbsync.AddActivityResponse{}
	if req == nil {
		return resp, nil
	}

	switch d := req.GetActivity().(type) {
	case *pbsync.AddActivityRequest_UserOauth2:
		if err := s.handleUserOauth2(ctx, d.UserOauth2); err != nil {
			return nil, fmt.Errorf("failed to handle UserOauth2 activity. %w", err)
		}
	case *pbsync.AddActivityRequest_Dispatch:
		if s.dispatches != nil {
			if _, err := s.dispatches.Create(ctx, d.Dispatch); err != nil {
				return nil, fmt.Errorf("failed to create dispatch. %w", err)
			}
		}
	case *pbsync.AddActivityRequest_UserActivity:
		if err := s.createUserActivities(ctx, s.db, d.UserActivity); err != nil {
			return nil, fmt.Errorf("failed to create user activities. %w", err)
		}
	case *pbsync.AddActivityRequest_UserProps:
		if err := s.handleUserProps(ctx, d.UserProps); err != nil {
			return nil, fmt.Errorf("failed to handle UserProps activity. %w", err)
		}
	case *pbsync.AddActivityRequest_ColleagueActivity:
		if err := s.createColleagueActivity(ctx, s.db, d.ColleagueActivity); err != nil {
			return nil, fmt.Errorf("failed to create jobs user activities. %w", err)
		}
	case *pbsync.AddActivityRequest_ColleagueProps:
		if err := s.handleColleagueProps(ctx, d.ColleagueProps); err != nil {
			return nil, fmt.Errorf("failed to handle ColleagueProps activity. %w", err)
		}
	case *pbsync.AddActivityRequest_JobTimeclock:
		if err := s.handleTimeclockEntry(ctx, d.JobTimeclock); err != nil {
			return nil, fmt.Errorf("failed to handle JobTimeclock activity. %w", err)
		}
	case *pbsync.AddActivityRequest_UserUpdate:
		if err := s.handleUserUpdate(ctx, d.UserUpdate); err != nil {
			return nil, fmt.Errorf("failed to handle UserUpdate activity. %w", err)
		}
	case *pbsync.AddActivityRequest_AccountUpdate:
		if err := s.handleAccountUpdate(ctx, d.AccountUpdate); err != nil {
			return nil, fmt.Errorf("failed to handle AccountUpdate activity. %w", err)
		}
	}

	return resp, nil
}

func (s *Store) AddUserActivity(
	ctx context.Context,
	req *pbsync.AddUserActivityRequest,
) (*pbsync.AddActivityResponse, error) {
	if err := s.createUserActivities(ctx, s.db, req.GetUserActivity()); err != nil {
		return nil, fmt.Errorf("failed to create user activities. %w", err)
	}
	return &pbsync.AddActivityResponse{}, nil
}

func (s *Store) AddUserUpdate(
	ctx context.Context,
	req *pbsync.AddUserUpdateRequest,
) (*pbsync.AddActivityResponse, error) {
	if err := s.handleUserUpdate(ctx, req.GetUserUpdate()); err != nil {
		return nil, fmt.Errorf("failed to handle user update data. %w", err)
	}
	return &pbsync.AddActivityResponse{}, nil
}

func (s *Store) AddUserProps(
	ctx context.Context,
	req *pbsync.AddUserPropsRequest,
) (*pbsync.AddActivityResponse, error) {
	if err := s.handleUserProps(ctx, req.GetUserProps()); err != nil {
		return nil, fmt.Errorf("failed to handle UserProps activity. %w", err)
	}
	return &pbsync.AddActivityResponse{}, nil
}

func (s *Store) AddColleagueActivity(
	ctx context.Context,
	req *pbsync.AddColleagueActivityRequest,
) (*pbsync.AddActivityResponse, error) {
	if err := s.createColleagueActivity(ctx, s.db, req.GetColleagueActivity()); err != nil {
		return nil, fmt.Errorf("failed to create jobs user activities. %w", err)
	}
	return &pbsync.AddActivityResponse{}, nil
}

func (s *Store) AddColleagueProps(
	ctx context.Context,
	req *pbsync.AddColleaguePropsRequest,
) (*pbsync.AddActivityResponse, error) {
	if err := s.handleColleagueProps(ctx, req.GetColleagueProps()); err != nil {
		return nil, fmt.Errorf("failed to handle ColleagueProps activity. %w", err)
	}
	return &pbsync.AddActivityResponse{}, nil
}

func (s *Store) AddJobTimeclock(
	ctx context.Context,
	req *pbsync.AddJobTimeclockRequest,
) (*pbsync.AddActivityResponse, error) {
	if err := s.handleTimeclockEntry(ctx, req.GetJobTimeclock()); err != nil {
		return nil, fmt.Errorf("failed to handle JobTimeclock activity. %w", err)
	}
	return &pbsync.AddActivityResponse{}, nil
}

func (s *Store) AddDispatch(
	ctx context.Context,
	req *pbsync.AddDispatchRequest,
) (*pbsync.AddActivityResponse, error) {
	if s.dispatches != nil {
		if _, err := s.dispatches.Create(ctx, req.GetDispatch()); err != nil {
			return nil, fmt.Errorf("failed to create dispatch. %w", err)
		}
	}
	return &pbsync.AddActivityResponse{}, nil
}

func (s *Store) createUserActivities(
	ctx context.Context,
	tx qrm.DB,
	activities ...*usersactivity.UserActivity,
) error {
	if len(activities) == 0 {
		return nil
	}

	tUserActivity := table.FivenetUserActivity
	stmt := tUserActivity.
		INSERT(
			tUserActivity.SourceUserID,
			tUserActivity.TargetUserID,
			tUserActivity.Type,
			tUserActivity.Reason,
			tUserActivity.Data,
		)
	for _, activity := range activities {
		stmt = stmt.VALUES(
			activity.SourceUserId,
			activity.GetTargetUserId(),
			activity.GetType(),
			activity.GetReason(),
			activity.Data,
		)
	}

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) createColleagueActivity(
	ctx context.Context,
	tx qrm.DB,
	activities ...*colleaguesactivity.ColleagueActivity,
) error {
	if len(activities) == 0 {
		return nil
	}

	tJobColleagueActivity := table.FivenetJobColleagueActivity
	stmt := tJobColleagueActivity.
		INSERT(
			tJobColleagueActivity.Job,
			tJobColleagueActivity.SourceUserID,
			tJobColleagueActivity.TargetUserID,
			tJobColleagueActivity.ActivityType,
			tJobColleagueActivity.Reason,
			tJobColleagueActivity.Data,
		)
	for _, activity := range activities {
		stmt = stmt.VALUES(
			activity.GetJob(),
			activity.SourceUserId,
			activity.TargetUserId,
			activity.GetActivityType(),
			activity.GetReason(),
			activity.Data,
		)
	}

	_, err := stmt.ExecContext(ctx, tx)
	return err
}

func (s *Store) handleUserProps(ctx context.Context, data *activity.UserProps) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction. %w", err)
	}
	defer tx.Rollback()

	reqP := data.GetProps()
	props, err := s.citizensStore.GetUserProps(ctx, tx, reqP.GetUserId())
	if err != nil {
		return fmt.Errorf("failed to get user props. %w", err)
	}

	reason := ""
	if data.Reason != nil {
		reason = data.GetReason()
	}

	activities, err := s.citizensStore.HandleUserPropsChanges(
		ctx,
		tx,
		props,
		reqP,
		&reqP.UserId,
		reason,
	)
	if err != nil {
		return fmt.Errorf("failed to handle user props changes. %w", err)
	}

	if data.Reason != nil && data.GetReason() != "" {
		if err := s.createUserActivities(ctx, tx, activities...); err != nil {
			return fmt.Errorf("failed to create user activities. %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction. %w", err)
	}

	return nil
}

func (s *Store) handleColleagueProps(ctx context.Context, data *activity.ColleagueProps) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction. %w", err)
	}
	defer tx.Rollback()

	input := data.GetProps()
	props, err := s.jobsStore.GetColleagueProps(ctx, tx, input.GetJob(), input.GetUserId(), nil)
	if err != nil {
		return fmt.Errorf("failed to get jobs user props. %w", err)
	}

	reason := ""
	if data.Reason != nil {
		reason = data.GetReason()
	}

	activities, err := s.jobsStore.HandleColleaguePropsChanges(
		ctx,
		tx,
		props,
		input,
		input.GetJob(),
		&input.UserId,
		reason,
	)
	if err != nil {
		return fmt.Errorf("failed to handle jobs user props changes. %w", err)
	}

	if data.Reason != nil && data.GetReason() != "" {
		if err := s.createColleagueActivity(ctx, tx, activities...); err != nil {
			return fmt.Errorf("failed to create jobs user activities. %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction. %w", err)
	}

	return nil
}

func (s *Store) handleTimeclockEntry(ctx context.Context, data *activity.TimeclockUpdate) error {
	tTimeClock := table.FivenetJobTimeclock

	if data.GetStart() {
		stmt := tTimeClock.
			SELECT(tTimeClock.UserID, tTimeClock.Date, tTimeClock.EndTime).
			FROM(tTimeClock).
			WHERE(mysql.AND(tTimeClock.Job.EQ(mysql.String(data.GetJob())), tTimeClock.UserID.EQ(mysql.Int32(data.GetUserId())))).
			ORDER_BY(tTimeClock.Date.DESC()).
			LIMIT(1)

		dest := &jobstimeclock.TimeclockEntry{}
		if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return fmt.Errorf("failed to query timeclock entry. %w", err)
			}
		}

		if dest.GetUserId() > 0 {
			return nil
		}

		updateStmt := tTimeClock.
			INSERT(tTimeClock.Job, tTimeClock.UserID, tTimeClock.Date).
			VALUES(data.GetJob(), data.GetUserId(), mysql.CURRENT_DATE())
		if _, err := updateStmt.ExecContext(ctx, s.db); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return fmt.Errorf("failed to insert timeclock entry. %w", err)
			}
		}
	} else {
		stmt := tTimeClock.
			UPDATE().
			SET(
				tTimeClock.SpentTime.SET(mysql.RawFloat("`spent_time` + CAST((TIMESTAMPDIFF(SECOND, `start_time`, `end_time`) / 3600) AS DECIMAL(10,2))")),
				tTimeClock.EndTime.SET(mysql.CURRENT_TIMESTAMP()),
			).
			WHERE(mysql.AND(tTimeClock.Job.EQ(mysql.String(data.GetJob())), tTimeClock.UserID.EQ(mysql.Int32(data.GetUserId())), tTimeClock.StartTime.IS_NOT_NULL(), tTimeClock.EndTime.IS_NULL())).
			LIMIT(1)
		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return fmt.Errorf("failed to update timeclock entry. %w", err)
		}
	}

	return nil
}

func (s *Store) handleUserUpdate(ctx context.Context, data *activity.UserUpdate) error {
	tUser := table.FivenetUser
	selectStmt := tUser.SELECT(tUser.ID).
		FROM(tUser).
		WHERE(tUser.ID.EQ(mysql.Int32(data.GetUserId()))).
		LIMIT(1)

	user := &users.User{}
	if err := selectStmt.QueryContext(ctx, s.db, user); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query user by ID. %w", err)
		}
		return nil
	}

	updateSets := []mysql.ColumnAssigment{}
	if data.Group != nil {
		updateSets = append(updateSets, tUser.Group.SET(mysql.String(data.GetGroup())))
	}
	if data.Job != nil {
		updateSets = append(updateSets, tUser.Job.SET(mysql.String(data.GetJob())))
	}
	if data.JobGrade != nil {
		updateSets = append(updateSets, tUser.JobGrade.SET(mysql.Int32(data.GetJobGrade())))
	}
	if data.Firstname != nil {
		updateSets = append(updateSets, tUser.Firstname.SET(mysql.String(data.GetFirstname())))
	}
	if data.Lastname != nil {
		updateSets = append(updateSets, tUser.Lastname.SET(mysql.String(data.GetLastname())))
	}

	if len(updateSets) > 0 {
		stmt := tUser.UPDATE().
			SET(updateSets[0]).
			WHERE(tUser.ID.EQ(mysql.Int32(data.GetUserId()))).
			LIMIT(1)
		if _, err := stmt.ExecContext(ctx, s.db); err != nil {
			return fmt.Errorf("failed to update user. %w", err)
		}
	}

	return nil
}
