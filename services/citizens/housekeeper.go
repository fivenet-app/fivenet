package citizens

import (
	context "context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	"github.com/fivenet-app/fivenet/v2026/pkg/access"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	propertyregistry "github.com/fivenet-app/fivenet/v2026/pkg/properties"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	citizensstore "github.com/fivenet-app/fivenet/v2026/stores/citizens"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var HousekeeperModule = fx.Module(
	"citizens.housekeeper",
	fx.Provide(
		NewHousekeeper,
	),
)

const changedRowsAttributeKey = "changed_rows"

type Housekeeper struct {
	logger *zap.Logger
	tracer trace.Tracer

	db     *sql.DB
	appCfg appconfig.IConfig
	store  citizensstore.IStore
}

type HousekeeperParams struct {
	fx.In

	Logger    *zap.Logger
	DB        *sql.DB
	TP        *tracesdk.TracerProvider
	AppConfig appconfig.IConfig
}

type HousekeeperResult struct {
	fx.Out

	Housekeeper  *Housekeeper
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

func NewHousekeeper(p HousekeeperParams) HousekeeperResult {
	s := &Housekeeper{
		logger: p.Logger.Named("citizens.housekeeper"),
		tracer: p.TP.Tracer("citizens.housekeeper"),

		db:     p.DB,
		appCfg: p.AppConfig,
		store: citizensstore.New(citizensstore.Params{
			DB:           p.DB,
			CustomDB:     &config.CustomDB{},
			LabelsAccess: access.NewCitizenLabelsSubjectObjectAccess(p.DB),
		}),
	}

	return HousekeeperResult{
		Housekeeper:  s,
		CronRegister: s,
	}
}

func (s *Housekeeper) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "citizens.user_props.max_wanted_duration",
		Schedule: "*/2 * * * *", // Every two minutes
	}); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "citizens.user_props.expire_labels",
		Schedule: "*/2 * * * *", // Every five minutes
	}); err != nil {
		return err
	}

	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "citizens.user_props.expire_custom_properties",
		Schedule: "*/2 * * * *", // Every two minutes
	}); err != nil {
		return err
	}

	return nil
}

func (s *Housekeeper) RegisterCronjobHandlers(h *croner.Handlers) error {
	h.Add(
		"citizens.user_props.max_wanted_duration",
		func(ctx context.Context, data *cron.CronjobData) error {
			ctx, span := s.tracer.Start(ctx, "citizens.user_props.max_wanted_duration")
			defer span.End()

			dest := &cron.GenericCronData{
				Attributes: map[string]string{},
			}
			if err := data.Unmarshal(dest); err != nil {
				s.logger.Warn("failed to unmarshal user props cleanup cron data", zap.Error(err))
			}

			changedRows, err := s.maxWantedDurationHandling(ctx)
			if err != nil {
				s.logger.Error("error during user props cleanup", zap.Error(err))
				return err
			}

			dest.SetAttribute(changedRowsAttributeKey, strconv.Itoa(changedRows))

			// Marshal the updated cron data
			if err := data.MarshalFrom(dest); err != nil {
				return fmt.Errorf(
					"failed to marshal updated user props cleanup cron data. %w",
					err,
				)
			}

			return nil
		},
	)

	h.Add(
		"citizens.user_props.expire_labels",
		func(ctx context.Context, data *cron.CronjobData) error {
			ctx, span := s.tracer.Start(ctx, "citizens.user_props.expire_labels")
			defer span.End()

			dest := &cron.GenericCronData{
				Attributes: map[string]string{},
			}
			if err := data.Unmarshal(dest); err != nil {
				s.logger.Warn(
					"failed to unmarshal user props labels cleanup cron data",
					zap.Error(err),
				)
			}

			changedRows, err := s.expireLabelHandling(ctx)
			if err != nil {
				s.logger.Error("error during user props labels cleanup", zap.Error(err))
				return err
			}

			dest.SetAttribute(changedRowsAttributeKey, strconv.Itoa(changedRows))

			// Marshal the updated cron data
			if err := data.MarshalFrom(dest); err != nil {
				return fmt.Errorf(
					"failed to marshal updated user props labels cleanup cron data. %w",
					err,
				)
			}

			return nil
		},
	)

	h.Add(
		"citizens.user_props.expire_custom_properties",
		func(ctx context.Context, data *cron.CronjobData) error {
			ctx, span := s.tracer.Start(ctx, "citizens.user_props.expire_custom_properties")
			defer span.End()

			dest := &cron.GenericCronData{Attributes: map[string]string{}}
			if err := data.Unmarshal(dest); err != nil {
				s.logger.Warn(
					"failed to unmarshal expired custom properties cron data",
					zap.Error(err),
				)
			}

			startedAt := time.Now()
			changedRows, err := s.expireCustomPropertyHandling(ctx)
			propertyregistry.ObserveExpiryCleanup("user", time.Since(startedAt), changedRows, err)
			if err != nil {
				s.logger.Error("error during expired custom properties cleanup", zap.Error(err))
				return err
			}
			dest.SetAttribute(changedRowsAttributeKey, strconv.Itoa(changedRows))
			if err := data.MarshalFrom(dest); err != nil {
				return fmt.Errorf("failed to marshal expired custom properties cron data. %w", err)
			}
			return nil
		},
	)

	return nil
}

func (s *Housekeeper) maxWantedDurationHandling(ctx context.Context) (int, error) {
	game := s.appCfg.Get().GetGame()

	// Skip if the max wanted duration feature is disabled or not properly configured
	if !game.GetMaxWantedDurationUserEnabled() || game.GetMaxWantedDurationUser() == nil {
		return 0, nil
	}
	wantedEnabled, err := s.store.IsSystemPropertyEnabled(ctx, "user.wanted")
	if err != nil {
		return 0, err
	}
	if !wantedEnabled {
		return 0, nil
	}

	maxSeconds := game.GetMaxWantedDurationUser().GetSeconds()

	dest, err := s.store.ListExpiredWantedUserProps(ctx, maxSeconds, 100)
	if err != nil {
		return 0, err
	}

	changed := 0
	for _, userId := range dest {
		cleared, err := s.clearExpiredWanted(ctx, userId, time.Duration(maxSeconds)*time.Second)
		if err != nil {
			s.logger.Error(
				"error updating user wanted state during cleanup",
				zap.Int32("user_id", userId),
				zap.Error(err),
			)
			continue
		}
		if cleared {
			changed++
		}
	}

	return changed, nil
}

// clearExpiredWanted locks the canonical wanted value and lifecycle state,
// then rechecks expiry inside the transaction. The initial candidate query is
// deliberately only a cheap work queue and must not decide whether to clear.
func (s *Housekeeper) clearExpiredWanted(
	ctx context.Context,
	userID int32,
	maxDuration time.Duration,
) (bool, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	defer tx.Rollback()

	tValues := table.FivenetUserPropertyValues
	value := struct {
		Wanted *bool `alias:"wanted"`
	}{}
	err = tValues.SELECT(tValues.ValueBool.AS("wanted")).FROM(tValues).WHERE(mysql.AND(
		tValues.UserID.EQ(mysql.Int32(userID)),
		tValues.PropertyKey.EQ(mysql.String("user.wanted")),
	)).LIMIT(1).FOR(mysql.UPDATE()).QueryContext(ctx, tx, &value)
	if errors.Is(err, qrm.ErrNoRows) {
		return false, tx.Commit()
	}
	if err != nil {
		return false, err
	}

	tState := table.FivenetUserPropertyState
	state := struct {
		ActiveAt  *time.Time `alias:"active_at"`
		ExpiresAt *time.Time `alias:"expires_at"`
	}{}
	err = tState.SELECT(tState.ActiveAt, tState.ExpiresAt).FROM(tState).WHERE(mysql.AND(
		tState.UserID.EQ(mysql.Int32(userID)),
		tState.PropertyKey.EQ(mysql.String("user.wanted")),
	)).LIMIT(1).FOR(mysql.UPDATE()).QueryContext(ctx, tx, &state)
	if errors.Is(err, qrm.ErrNoRows) {
		return false, tx.Commit()
	}
	if err != nil {
		return false, err
	}

	now := time.Now()
	if value.Wanted == nil || !*value.Wanted ||
		!wantedExpired(state.ActiveAt, state.ExpiresAt, maxDuration, now) {
		return false, tx.Commit()
	}
	if err := s.store.ClearWantedTx(ctx, tx, userID); err != nil {
		return false, fmt.Errorf("clearing user %d wanted state: %w", userID, err)
	}
	if err := tx.Commit(); err != nil {
		return false, err
	}
	return true, nil
}

func wantedExpired(activeAt, expiresAt *time.Time, maxDuration time.Duration, now time.Time) bool {
	return (activeAt != nil && !activeAt.Add(maxDuration).After(now)) ||
		(expiresAt != nil && !expiresAt.After(now))
}

// expireCustomPropertyHandling is definition driven: it scans definitions
// first, then performs a small, per-key lookup of expired state/value pairs.
// This prevents generic cleanup from ever touching managed properties.
func (s *Housekeeper) expireCustomPropertyHandling(ctx context.Context) (int, error) {
	keys, err := s.store.ListExpiringCustomPropertyKeys(ctx)
	if err != nil {
		return 0, err
	}

	changed := 0
	for _, key := range keys {
		userIDs, err := s.store.ListExpiredCustomPropertyUserIDs(ctx, key, 100)
		if err != nil {
			return changed, err
		}
		for _, userID := range userIDs {
			cleared, err := s.store.ClearExpiredCustomProperty(ctx, userID, key)
			if err != nil {
				propertyregistry.RecordExpiryCleanupItemFailure("user")
				s.logger.Error(
					"error clearing expired custom user property",
					zap.Int32("user_id", userID),
					zap.String("property_key", key),
					zap.Error(err),
				)
				continue
			}
			if cleared {
				changed++
			}
		}
	}
	return changed, nil
}

func (s *Housekeeper) expireLabelHandling(ctx context.Context) (int, error) {
	expiredLabel := tCitizenLabels.AS("expired_label")
	expiredLabelJob := tCitizensLabelsJob.AS("expired_label_job")

	expiredLabels := expiredLabel.
		SELECT(
			expiredLabel.UserID.AS("user_id"),
			mysql.MIN(expiredLabel.LabelID).AS("label_id"),
		).
		FROM(expiredLabel).
		WHERE(mysql.AND(
			expiredLabel.ExpiresAt.IS_NOT_NULL(),
			expiredLabel.ExpiresAt.LT_EQ(mysql.CURRENT_TIMESTAMP()),
		)).
		GROUP_BY(expiredLabel.UserID).
		LIMIT(100).
		AsTable("expired_labels")

	labelID := mysql.IntegerColumn("label_id").From(expiredLabels)
	userID := mysql.IntegerColumn("user_id").From(expiredLabels)

	stmt := mysql.
		SELECT(
			labelID,
			userID,
		).
		FROM(
			expiredLabels.
				INNER_JOIN(expiredLabelJob,
					expiredLabelJob.ID.EQ(labelID),
				),
		)

	var dest []struct {
		LabelID int64
		UserId  int32
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return 0, err
	}

	for _, row := range dest {
		if err := s.removeLabelFromUser(
			ctx,
			s.db,
			row.LabelID,
			row.UserId,
		); err != nil {
			s.logger.Error(
				"error updating user labels during cleanup",
				zap.Int32("user_id", row.UserId),
				zap.Error(err),
			)
			continue
		}
	}

	return len(dest), nil
}

func (s *Housekeeper) removeLabelFromUser(
	ctx context.Context,
	tx qrm.DB,
	labelId int64,
	targetUserId int32,
) error {
	stmt := tCitizenLabels.
		DELETE().
		WHERE(mysql.AND(
			tCitizenLabels.UserID.EQ(mysql.Int32(targetUserId)),
			tCitizenLabels.LabelID.EQ(mysql.Int64(labelId)),
		)).
		LIMIT(1)

	if _, err := stmt.ExecContext(ctx, s.db); err != nil {
		return fmt.Errorf(
			"failed to delete citizen label user %d assignment. %w",
			targetUserId,
			err,
		)
	}

	if err := usersactivity.CreateUserActivities(ctx, tx, &usersactivity.UserActivity{
		TargetUserId: targetUserId,
		Type:         usersactivity.UserActivityType_USER_ACTIVITY_TYPE_LABELS,
		Reason:       "",
		Data: &usersactivity.UserActivityData{
			Data: &usersactivity.UserActivityData_LabelsChange{
				LabelsChange: &usersactivity.LabelsChange{
					RemovedIds: []int64{labelId},
					Expired:    true,
				},
			},
		},
	}); err != nil {
		return fmt.Errorf(
			"error creating user activity for user %d during cleanup. %w",
			targetUserId,
			err,
		)
	}

	return nil
}
