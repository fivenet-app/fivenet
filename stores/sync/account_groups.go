package syncstore

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"slices"

	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	notificationsevents "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications/events"
	activity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/activity"
	pkguserinfo "github.com/fivenet-app/fivenet/v2026/pkg/userinfo"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

type accountGroupState struct {
	ID      int64
	License string
	Groups  *accounts.AccountGroups
}

type accountGroupChange struct {
	accountID int64
	license   string
	groups    *accounts.AccountGroups
}

func accountGroupsFromSyncUpdate(update *activity.AccountUpdate) *accounts.AccountGroups {
	if update == nil {
		return nil
	}

	if update.GetGroups() != nil && len(update.GetGroups().GetGroups()) > 0 {
		return &accounts.AccountGroups{Groups: slices.Clone(update.GetGroups().GetGroups())}
	}

	if update.GetGroup() != "" {
		return &accounts.AccountGroups{Groups: []string{update.GetGroup()}}
	}

	return nil
}

func accountGroupsEqual(a, b *accounts.AccountGroups) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	return a.Equal(b)
}

func accountGroupsFromState(state *accountGroupState) *accounts.AccountGroups {
	if state == nil || state.Groups == nil || len(state.Groups.GetGroups()) == 0 {
		return nil
	}

	return &accounts.AccountGroups{Groups: slices.Clone(state.Groups.GetGroups())}
}

func (s *Store) loadAccountGroupState(
	ctx context.Context,
	tx *sql.Tx,
	license string,
) (*accountGroupState, error) {
	if license == "" {
		return nil, nil
	}

	tAccounts := table.FivenetAccounts
	stmt := tAccounts.
		SELECT(
			tAccounts.ID.AS("id"),
			tAccounts.License.AS("license"),
			tAccounts.Groups.AS("groups"),
		).
		FROM(tAccounts).
		WHERE(tAccounts.License.EQ(mysql.String(license))).
		LIMIT(1)

	dest := &accountGroupState{}
	sqlStmt, args := stmt.Sql()
	row := tx.QueryRowContext(ctx, sqlStmt, args...)
	var groups sql.NullString
	if err := row.Scan(&dest.ID, &dest.License, &groups); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to query account groups for license %s. %w", license, err)
	}

	if dest.ID == 0 {
		return nil, nil
	}
	if groups.Valid {
		dest.Groups = &accounts.AccountGroups{}
		if err := dest.Groups.Scan(groups.String); err != nil {
			return nil, fmt.Errorf("failed to scan account groups for license %s. %w", license, err)
		}
	}

	return dest, nil
}

func (s *Store) publishAccountGroupsChanged(
	ctx context.Context,
	accountID int64,
	license string,
	groups *accounts.AccountGroups,
) {
	if s.notifi == nil {
		return
	}

	superuserGroups, superuserUsers := []string(nil), []string(nil)
	if s.cfg != nil {
		superuserGroups = s.cfg.Auth.SuperuserGroups
		superuserUsers = s.cfg.Auth.SuperuserUsers
	}

	event := pkguserinfo.BuildAccountGroupsChangedEvent(
		accountID,
		nil,
		groups,
		pkguserinfo.CanBeSuperuser(groups, license, superuserGroups, superuserUsers),
	)
	if err := s.notifi.SendAccountEvent(ctx, accountID, &notificationsevents.UserEvent{
		Data: &notificationsevents.UserEvent_AccountGroupsChanged{
			AccountGroupsChanged: event,
		},
	}); err != nil {
		groupNames := []string(nil)
		if groups != nil {
			groupNames = groups.GetGroups()
		}
		s.logger.Warn(
			"failed to publish account group change event",
			zap.Int64("account_id", accountID),
			zap.Strings("groups", groupNames),
			zap.Error(err),
		)
	}
}
