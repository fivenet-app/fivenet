package cmd

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/kong"
	"github.com/fivenet-app/fivenet/cmd/envs"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/users"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ToolsCmd struct {
	DB DBCmd `cmd:""`

	UserActivityMigrate UserActivityMigrateCmd `cmd:""`
}

type UserActivityMigrateCmd struct {
	Limit int64 `default:"250" help:"Limit to set on the select query."`
}

func (c *UserActivityMigrateCmd) Run(ctx *kong.Context) error {
	fxOpts := getFxBaseOpts(Cli.StartTimeout, false)

	if err := os.Setenv(envs.SkipDBMigrationsEnv, "true"); err != nil {
		return err
	}

	fxOpts = append(fxOpts,
		fx.Invoke(func(lifecycle fx.Lifecycle, logger *zap.Logger, db *sql.DB, shutdowner fx.Shutdowner) {
			lifecycle.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					go func() {
						exitCode := 0
						if err := c.run(context.Background(), logger, db); err != nil {
							logger.Error("error during migration", zap.Error(err))
							// handle error, set non-zero exit code so caller knows the job failed
							exitCode = 1
						}
						_ = shutdowner.Shutdown(fx.ExitCode(exitCode))
					}()
					return nil
				},
			})
		}),
	)

	app := fx.New(fxOpts...)
	app.Run()

	return nil
}

func (c *UserActivityMigrateCmd) run(ctx context.Context, logger *zap.Logger, db *sql.DB) error {
	tUserActivity := table.FivenetUserActivity

	limit := c.Limit
	logger.Info("starting user activity migration", zap.Int64("limit", limit))

	for {
		stmt := tUserActivity.
			SELECT(
				tUserActivity.ID,
				jet.String("fivenet_user_activity.key"),
				jet.String("fivenet_user_activity.old_value"),
				jet.String("fivenet_user_activity.new_value"),
				tUserActivity.Reason,
			).
			FROM(tUserActivity).
			WHERE(jet.AND(
				jet.String("fivenet_user_activity.key").IS_NOT_NULL(),
			)).
			ORDER_BY(tUserActivity.ID.ASC()).
			LIMIT(limit)

		dest := []*users.UserActivity{}
		if err := stmt.QueryContext(ctx, db, &dest); err != nil {
			if !errors.Is(err, qrm.ErrNoRows) {
				return err
			}
		}

		if len(dest) == 0 {
			logger.Info("no more user activity found in database that needs to be migrated")
			return nil
		}

		for _, activity := range dest {
			switch activity.Key {
			case "DocStore.Relation":
				activity.Type = users.UserActivityType_USER_ACTIVITY_TYPE_DOCUMENT

				val := activity.NewValue
				if activity.OldValue != "" {
					val = activity.OldValue
				}

				docId, err := strconv.Atoi(val)
				if err != nil {
					return err
				}

				activity.Reason = strings.ReplaceAll(activity.Reason, "DOC_RELATION_", "")
				relation := 3 // CAUSED
				if activity.Reason == "MENTIONED" {
					relation = 1
				} else if activity.Reason == "TARGETS" {
					relation = 2
				}
				activity.Reason = ""

				activity.Data = &users.UserActivityData{
					Data: &users.UserActivityData_DocumentRelation{
						DocumentRelation: &users.UserDocumentRelation{
							DocumentId: uint64(docId),
							Added:      activity.NewValue != "",
							Relation:   int32(relation),
						},
					},
				}

			case "Plugin.Billing.Fines":
				activity.Type = users.UserActivityType_USER_ACTIVITY_TYPE_FINE

				if activity.OldValue == "" {
					activity.OldValue = "0"
				}
				oldVal, err := strconv.Atoi(activity.OldValue)
				if err != nil {
					return err
				}

				if activity.NewValue == "" {
					activity.NewValue = "0"
				}
				newVal, err := strconv.Atoi(activity.NewValue)
				if err != nil {
					return err
				}

				amount := 0
				removed := false
				if newVal == 0 {
					amount = -oldVal
				} else if newVal == oldVal {
					removed = true
					amount = -newVal
				} else {
					amount = newVal
				}

				activity.Data = &users.UserActivityData{
					Data: &users.UserActivityData_FineChange{
						FineChange: &users.UserFineChange{
							Removed: removed,
							Amount:  int64(amount),
						},
					},
				}

			case "Plugin.Jail":
				activity.Type = users.UserActivityType_USER_ACTIVITY_TYPE_JAIL

				seconds := 0
				if activity.NewValue != "" {
					newVal, err := strconv.Atoi(activity.NewValue)
					if err != nil {
						return err
					}

					if newVal != 0 {
						seconds = newVal
					} else if newVal == 0 {
						seconds = 0
					} else {
						seconds = -1
					}
				} else {
					seconds = -1
				}

				activity.Data = &users.UserActivityData{
					Data: &users.UserActivityData_JailChange{
						JailChange: &users.UserJailChange{
							Seconds: int32(seconds),
							Admin:   false,
						},
					},
				}

			case "Plugin.Licenses":
				activity.Type = users.UserActivityType_USER_ACTIVITY_TYPE_LICENSES

				added := true
				ltype := activity.NewValue
				if activity.OldValue != "" {
					added = false
					ltype = activity.OldValue
				}

				activity.Data = &users.UserActivityData{
					Data: &users.UserActivityData_LicensesChange{
						LicensesChange: &users.UserLicenseChange{
							Added: added,
							Licenses: []*users.License{
								{
									Type:  ltype,
									Label: activity.Reason,
								},
							},
						},
					},
				}

			case "UserProps.Job":
				activity.Type = users.UserActivityType_USER_ACTIVITY_TYPE_JOB

				jsplit := strings.Split(activity.NewValue, "|")
				job := strings.ToLower(jsplit[0])
				jobLabel := jsplit[0]

				activity.Data = &users.UserActivityData{
					Data: &users.UserActivityData_JobChange{
						JobChange: &users.UserJobChange{
							Job:      &job,
							JobLabel: &jobLabel,
						},
					},
				}

			case "UserProps.Labels":
				// Labels entries are not transfered but deleted
				delStmt := tUserActivity.DELETE().WHERE(tUserActivity.ID.EQ(jet.Uint64(activity.Id))).LIMIT(1)
				if _, err := delStmt.ExecContext(ctx, db); err != nil {
					return err
				}
				continue

			case "UserProps.MugShot":
				activity.Type = users.UserActivityType_USER_ACTIVITY_TYPE_MUGSHOT

				var mugShotUrl *string
				if activity.NewValue != "" {
					mugShotUrl = &activity.NewValue
				}

				activity.Data = &users.UserActivityData{
					Data: &users.UserActivityData_MugshotChange{
						MugshotChange: &users.UserMugshotChange{
							New: mugShotUrl,
						},
					},
				}

			case "UserProps.TrafficInfractionPoints":
				activity.Type = users.UserActivityType_USER_ACTIVITY_TYPE_TRAFFIC_INFRACTION_POINTS

				oldPoints, err := strconv.Atoi(activity.OldValue)
				if err != nil {
					return err
				}
				newPoints, err := strconv.Atoi(activity.NewValue)
				if err != nil {
					return err
				}

				activity.Data = &users.UserActivityData{
					Data: &users.UserActivityData_TrafficInfractionPointsChange{
						TrafficInfractionPointsChange: &users.UserTrafficInfractionPointsChange{
							Old: uint32(oldPoints),
							New: uint32(newPoints),
						},
					},
				}

			case "UserProps.Wanted":
				activity.Type = users.UserActivityType_USER_ACTIVITY_TYPE_WANTED

				activity.Data = &users.UserActivityData{
					Data: &users.UserActivityData_WantedChange{
						WantedChange: &users.UserWantedChange{
							Wanted: activity.NewValue == "true",
						},
					},
				}

			default:
				return fmt.Errorf("unknown activity key %q", activity.Key)
			}

			if err := c.updateActivity(ctx, db, activity.Id, activity.Type, activity.Reason, activity.Data); err != nil {
				return err
			}
		}

		if len(dest) < int(limit) {
			logger.Info("no more user activity expected to be migrated")
			return nil
		}

		lastActivity := dest[len(dest)-1]
		logger.Info("migration at activity id", zap.Uint64("activity_id", lastActivity.Id))
		time.Sleep(100 * time.Millisecond)

	}
}

func (c *UserActivityMigrateCmd) updateActivity(ctx context.Context, tx *sql.DB, id uint64, aType users.UserActivityType, reason string, data *users.UserActivityData) error {
	tUserActivity := table.FivenetUserActivity

	stmt := tUserActivity.
		UPDATE(
			jet.StringColumn("fivenet_user_activity.key"),
			jet.StringColumn("fivenet_user_activity.old_value"),
			jet.StringColumn("fivenet_user_activity.new_value"),
			tUserActivity.Type,
			tUserActivity.Reason,
			tUserActivity.Data,
		).
		SET(
			jet.NULL,
			jet.NULL,
			jet.NULL,
			aType,
			reason,
			data,
		).
		WHERE(
			tUserActivity.ID.EQ(jet.Uint64(id)),
		)

	if _, err := stmt.ExecContext(ctx, tx); err != nil {
		return err
	}

	return nil
}
