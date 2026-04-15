package sync

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"time"

	activity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/sync/activity"
	pbsync "github.com/fivenet-app/fivenet/v2026/gen/go/proto/services/sync"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/zap"
)

func (s *Server) AddUserOAuth2Conn(
	ctx context.Context,
	req *pbsync.AddUserOAuth2ConnRequest,
) (*pbsync.AddActivityResponse, error) {
	s.lastSyncedActivity.Store(time.Now().Unix())

	if err := s.handleUserOauth2(ctx, req.UserOauth2); err != nil {
		return nil, fmt.Errorf("failed to handle UserOauth2 activity. %w", err)
	}

	return &pbsync.AddActivityResponse{}, nil
}

func (s *Server) handleUserOauth2(
	ctx context.Context,
	data *activity.UserOAuth2Conn,
) error {
	idx := slices.IndexFunc(s.cfg.OAuth2.Providers, func(in *config.OAuth2Provider) bool {
		return in.Name == data.GetProviderName()
	})
	if idx == -1 {
		return fmt.Errorf("invalid provider name. %s", data.GetProviderName())
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
		WHERE(tAccounts.License.EQ(mysql.String(data.GetIdentifier()))).
		LIMIT(1)

	if err := stmt.QueryContext(ctx, s.db, &account); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query account by identifier. %w", err)
		}
	}

	if account.ID == 0 {
		s.logger.Warn(
			"no fivenet account found for identifier in user oauth2 sync connect",
			zap.String("provider", data.GetProviderName()),
			zap.String(
				"identifier",
				data.GetIdentifier(),
			),
			zap.String("external_id", data.GetExternalId()),
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
			data.GetExternalId(),
			data.GetUsername(),
			provider.DefaultAvatar,
		)

	if _, err := insertStmt.ExecContext(ctx, s.db); err != nil {
		if !dbutils.IsDuplicateError(err) {
			return fmt.Errorf("failed to insert OAuth2 account. %w", err)
		}
	}

	return nil
}
