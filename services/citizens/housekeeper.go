package citizens

import (
	context "context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/cron"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	"github.com/fivenet-app/fivenet/v2026/pkg/config/appconfig"
	"github.com/fivenet-app/fivenet/v2026/pkg/croner"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
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

	return nil
}

func (s *Housekeeper) maxWantedDurationHandling(ctx context.Context) (int, error) {
	game := s.appCfg.Get().GetGame()

	// Skip if the max wanted duration feature is disabled or not properly configured
	if !game.GetMaxWantedDurationUserEnabled() || game.GetMaxWantedDurationUser() == nil {
		return 0, nil
	}

	maxDays := game.GetMaxWantedDurationUser().GetSeconds() / 24 / 3600

	stmt := tUserProps.
		SELECT(
			tUserProps.UserID,
		).
		FROM(tUserProps).
		WHERE(mysql.AND(
			tUserProps.Wanted.IS_TRUE(),
			mysql.OR(
				tUserProps.WantedAt.LT(
					mysql.CURRENT_TIMESTAMP().SUB(mysql.INTERVAL(maxDays, "DAY")),
				),
				tUserProps.WantedTill.LT(mysql.CURRENT_TIMESTAMP()),
			),
		)).
		LIMIT(100)

	var dest []struct {
		UserId int32
	}
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return 0, err
	}

	for _, row := range dest {
		err := s.unsetUserWantedState(ctx, s.db, row.UserId)
		if err != nil {
			s.logger.Error(
				"error updating user wanted state during cleanup",
				zap.Int32("user_id", row.UserId),
				zap.Error(err),
			)
			continue
		}
	}

	return len(dest), nil
}

func (s *Housekeeper) unsetUserWantedState(ctx context.Context, tx qrm.DB, userId int32) error {
	props, err := usersprops.GetUserProps(ctx, tx, userId, nil)
	if err != nil {
		return fmt.Errorf("error loading user %d props for cleanup. %w", userId, err)
	}

	in := proto.Clone(props).(*usersprops.UserProps)

	wanted := false
	in.Wanted = &wanted

	if _, err := props.HandleChanges(ctx, tx, in, nil, ""); err != nil {
		return fmt.Errorf("error handling user %d props changes during cleanup. %w", userId, err)
	}

	if err := usersactivity.CreateUserActivities(ctx, tx, &usersactivity.UserActivity{
		TargetUserId: userId,
		Type:         usersactivity.UserActivityType_USER_ACTIVITY_TYPE_WANTED,
		Reason:       "",
		Data: &usersactivity.UserActivityData{
			Data: &usersactivity.UserActivityData_WantedChange{
				WantedChange: &usersactivity.WantedChange{
					Wanted: wanted,
					Auto:   true,
				},
			},
		},
	}); err != nil {
		return fmt.Errorf(
			"error creating user activity for user %d during cleanup. %w",
			userId,
			err,
		)
	}

	return nil
}
