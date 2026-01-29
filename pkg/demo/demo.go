// Package demo provides functionality for generating random dispatches and user locations
// in demo mode, allowing for testing and demonstration of the centrum service.
//
//nolint:gosec // G404: rand.IntN is not cryptographically secure, but we don't need it to be here.
package demo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand/v2"
	"time"

	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum"
	centrumdispatches "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/centrum/dispatches"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/dispatches"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	tLocs      = table.FivenetCentrumUserLocations
	tTimeClock = table.FivenetJobTimeclock
)

var Module = fx.Module(
	"demo",
	fx.Provide(
		New,
	),
)

var (
	// Define clamp bounds outside the loop for configurability.
	xBounds = []float64{-2750, 2500}
	yBounds = []float64{-3000, 6000}

	dispatchDescriptions = []string{
		"A person was seen acting suspiciously in the area.",
		"A vehicle was reported speeding.",
		"Loud noises were heard coming from a building.",
		"Possible altercation in progress.",
		"Unattended package found.",
	}
	dispatchMessages = []string{
		"Suspicious activity reported",
		"Speeding vehicle reported",
		"Noise complaint received",
		"Possible fight reported",
		"Suspicious package reported",
	}

	dispatchStatusProgression = []centrumdispatches.StatusDispatch{
		centrumdispatches.StatusDispatch_STATUS_DISPATCH_NEW,
		centrumdispatches.StatusDispatch_STATUS_DISPATCH_EN_ROUTE,
		centrumdispatches.StatusDispatch_STATUS_DISPATCH_ON_SCENE,
		centrumdispatches.StatusDispatch_STATUS_DISPATCH_COMPLETED,
	}
)

type user struct {
	UserID     int32  `alias:"user.userid"     sql:"primary_key"`
	Identifier string `alias:"user.identifier"`
}

// Demo provides demo mode functionality for generating random dispatches and user locations.
type Demo struct {
	// logger for demo operations
	logger *zap.Logger
	// database connection
	db *sql.DB
	// centrum manager for dispatch operations
	dispatches *dispatches.DispatchDB
	// application configuration
	cfg *config.Config

	// users is a list of user identifiers for demo operations
	users []*user
}

// Params contains dependencies for constructing a Demo instance.
type Params struct {
	fx.In

	// LC is the Fx lifecycle for registering hooks
	LC fx.Lifecycle

	Logger     *zap.Logger
	Cfg        *config.Config
	DB         *sql.DB
	Dispatches *dispatches.DispatchDB
}

// New creates a new Demo instance if demo mode is enabled in the config.
// Registers lifecycle hooks for starting and stopping demo mode.
func New(p Params) *Demo {
	if !p.Cfg.Demo.Enabled {
		return nil
	}

	ctxCancel, cancel := context.WithCancel(context.Background())

	d := &Demo{
		logger:     p.Logger.Named("demo"),
		db:         p.DB,
		dispatches: p.Dispatches,
		cfg:        p.Cfg,
	}

	d.logger.Warn(
		"Demo mode is enabled. This will generate random dispatches and user locations!!!",
	)

	p.LC.Append(fx.StartHook(func(_ context.Context) error {
		go d.Start(ctxCancel)

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return d
}

func (d *Demo) lookupUsers(ctx context.Context, identifiers []string) ([]*user, error) {
	tUsers := table.FivenetUser

	condition := tUsers.Job.EQ(mysql.String(d.cfg.Demo.TargetJob))
	if len(identifiers) > 0 {
		idds := make([]mysql.Expression, len(identifiers))
		for i, id := range identifiers {
			idds[i] = mysql.String(id)
		}
		condition = condition.AND(tUsers.Identifier.IN(idds...))
	}

	stmt := tUsers.
		SELECT(
			tUsers.ID.AS("user.userid"),
			tUsers.Identifier.AS("user.identifier"),
		).
		FROM(tUsers).
		WHERE(condition).
		ORDER_BY(tUsers.ID.ASC()).
		LIMIT(20)

	users := []*user{}
	if err := stmt.QueryContext(ctx, d.db, &users); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query users. %w", err)
		}
	}

	return users, nil
}

// Start runs the main demo loop, periodically generating and updating dispatches and moving user markers.
func (d *Demo) Start(ctx context.Context) error {
	users, err := d.lookupUsers(ctx, nil)
	if err != nil {
		d.logger.Error("failed to lookup users", zap.Error(err))
		return err
	}

	demoUsers, err := d.lookupUsers(ctx, d.cfg.Demo.Users)
	if err != nil {
		d.logger.Error("failed to lookup demo users", zap.Error(err))
		return err
	}
	users = append(users, demoUsers...)

	utils.RemoveSliceDuplicatesFn(users, func(u *user) string {
		return u.Identifier
	})
	d.users = users

	d.createTimeclockEntries(ctx)

	go d.moveUserMarkers(ctx)

	for {
		d.logger.Info("Running demo dispatch generation cycle...")
		if err := d.updateDispatches(ctx); err != nil {
			d.logger.Error("failed to update dispatches", zap.Error(err))
		}

		d.generateDispatches(ctx)

		randWait := rand.IntN(300) + 30 // Random wait between 30 and 270 seconds

		select {
		case <-ctx.Done():
			return nil

		case <-time.After(time.Duration(randWait) * time.Second):
		}
	}
}

// generateDispatches creates up to 2 random dispatches per run with random positions and messages.
func (d *Demo) generateDispatches(ctx context.Context) {
	numDispatches := rand.IntN(2) // Up to 2 dispatches per run

	for range numDispatches {
		x := rand.Float64()*(xBounds[1]-xBounds[0]) + xBounds[0]
		y := rand.Float64()*(yBounds[1]-yBounds[0]) + yBounds[0]
		desc := dispatchDescriptions[rand.IntN(len(dispatchDescriptions))]
		msg := dispatchMessages[rand.IntN(len(dispatchMessages))]
		if _, err := d.dispatches.Create(ctx, &centrumdispatches.Dispatch{
			Jobs: &centrum.JobList{
				Jobs: []*centrum.JobListEntry{
					{
						Name: d.cfg.Demo.TargetJob,
					},
				},
			},
			Message:     msg,
			Description: &desc,
			X:           x,
			Y:           y,
			Anon:        true,
		}); err != nil {
			d.logger.Error("failed to create dispatch", zap.Error(err))
		}
	}
}

// updateDispatches randomly updates the status and position of up to 2 existing dispatches.
func (d *Demo) updateDispatches(ctx context.Context) error {
	dsps := d.dispatches.List(ctx, []string{d.cfg.Demo.TargetJob})
	if len(dsps) == 0 {
		return nil
	}

	// Shuffle indices to pick random dispatches
	perm := rand.Perm(len(dsps))
	numToUpdate := min(len(dsps), 2)

	for i := range numToUpdate {
		dsp := dsps[perm[i]]

		// Randomize position slightly
		x := dsp.GetX() + rand.Float64()*700 - 350
		y := dsp.GetY() + rand.Float64()*700 - 350

		// Pick a new status based on the status progression
		currStatus := dsp.GetStatus().GetStatus()
		newStatusValue := currStatus
		for i, s := range dispatchStatusProgression {
			if s == currStatus && i+1 < len(dispatchStatusProgression) {
				newStatusValue = dispatchStatusProgression[i+1]
				break
			}
		}
		newStatus := &centrumdispatches.DispatchStatus{
			DispatchId: dsp.GetId(),
			Status:     newStatusValue,

			X: &x,
			Y: &y,

			CreatorJob: &d.cfg.Demo.TargetJob,
		}
		if _, err := d.dispatches.UpdateStatus(ctx, dsp.GetId(), newStatus); err != nil {
			d.logger.Error("failed to update dispatch status", zap.Error(err))
		}
	}

	return nil
}

// moveUserMarkers periodically updates user locations, randomizing or moving them within bounds.
func (d *Demo) moveUserMarkers(ctx context.Context) {
	firstRun := true

	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(2 * time.Second):
		}

		for _, user := range d.users {
			// Fetch the current location for the user
			var curr struct {
				X float64
				Y float64
			}
			stmt := tLocs.
				SELECT(
					tLocs.X.AS("x"),
					tLocs.Y.AS("y"),
				).
				FROM(tLocs).
				WHERE(tLocs.Identifier.EQ(mysql.String(user.Identifier))).
				LIMIT(1)

			err := stmt.QueryContext(ctx, d.db, &curr)
			if err != nil && !errors.Is(err, qrm.ErrNoRows) {
				d.logger.Error(
					"failed to select user location",
					zap.String("user", user.Identifier),
					zap.Error(err),
				)
				continue
			}

			var newX, newY float64

			if firstRun {
				// Always randomize initial position on first run
				newX = rand.Float64()*(xBounds[1]-xBounds[0]) + xBounds[0]
				newY = rand.Float64()*(yBounds[1]-yBounds[0]) + yBounds[0]
			} else if errors.Is(err, qrm.ErrNoRows) {
				// No previous location, randomize
				newX = rand.Float64()*(xBounds[1]-xBounds[0]) + xBounds[0]
				newY = rand.Float64()*(yBounds[1]-yBounds[0]) + yBounds[0]
			} else {
				// Move a small random step from the current position
				maxStep := 35.0 // max movement per tick
				deltaX := (rand.Float64()*2 - 1) * maxStep
				deltaY := (rand.Float64()*2 - 1) * maxStep
				newX = curr.X + deltaX
				newY = curr.Y + deltaY

				// Clamp to bounds using the slices
				if newX < xBounds[0] {
					newX = xBounds[0]
				} else if newX > xBounds[1] {
					newX = xBounds[1]
				}
				if newY < yBounds[0] {
					newY = yBounds[0]
				} else if newY > yBounds[1] {
					newY = yBounds[1]
				}
			}

			insertStmt := tLocs.
				INSERT(
					tLocs.Identifier,
					tLocs.X,
					tLocs.Y,
					tLocs.Hidden,
				).
				VALUES(
					user,
					newX,
					newY,
					false,
				).
				ON_DUPLICATE_KEY_UPDATE(
					tLocs.X.SET(mysql.Float(newX)),
					tLocs.Y.SET(mysql.Float(newY)),
					tLocs.Hidden.SET(mysql.Bool(false)),
				)

			if _, err := insertStmt.ExecContext(ctx, d.db); err != nil {
				d.logger.Error(
					"failed to update user location",
					zap.String("user", user.Identifier),
					zap.Error(err),
				)
				continue
			}
		}

		if firstRun {
			firstRun = false
		}
	}
}

// createTimeclockEntries creates or updates timeclock table entries for the users in d.users.
// It creates multiple past entries of random lengths and sometimes an active entry.
func (d *Demo) createTimeclockEntries(ctx context.Context) {
	for _, user := range d.users {
		// Create multiple past entries
		// Randomly create 1 to 10 past entries
		numPastEntries := rand.IntN(10) + 1
		for range numPastEntries {
			daysAgo := rand.IntN(14) + 1 // Randomly pick a day within the last 14 days
			entryDate := time.Now().AddDate(0, 0, -daysAgo)
			startTime := entryDate.Add(
				time.Duration(rand.IntN(8)) * time.Hour,
			) // Random start time within the day
			endTime := startTime.Add(
				time.Duration(rand.IntN(8)+1) * time.Hour,
			) // Random end time after start time
			spentTime := endTime.Sub(startTime).Hours()

			insertStmt := tTimeClock.
				INSERT(
					tTimeClock.Job,
					tTimeClock.UserID,
					tTimeClock.Date,
					tTimeClock.StartTime,
					tTimeClock.EndTime,
					tTimeClock.SpentTime,
				).
				VALUES(
					d.cfg.Demo.TargetJob,
					user.UserID,
					entryDate,
					startTime,
					endTime,
					spentTime,
				)

			if _, err := insertStmt.ExecContext(ctx, d.db); err != nil {
				d.logger.Error(
					"failed to create past timeclock entry",
					zap.String("user", user.Identifier),
					zap.Error(err),
				)
				continue
			}
		}

		// Sometimes create an active entry
		if rand.IntN(2) == 0 { // 50% chance to create an active entry
			activeStartTime := time.Now().
				Add(-time.Duration(rand.IntN(8)+1) * time.Hour)
				// Random start time within the last 8 hours

			insertStmt := tTimeClock.
				INSERT(
					tTimeClock.Job,
					tTimeClock.UserID,
					tTimeClock.Date,
					tTimeClock.StartTime,
					tTimeClock.EndTime,
					tTimeClock.SpentTime,
				).
				VALUES(
					d.cfg.Demo.TargetJob,
					user.UserID,
					time.Now(),
					activeStartTime,
					nil, // Active entry has no end time
					nil, // Spent time is nil for active entries
				).
				ON_DUPLICATE_KEY_UPDATE(
					tTimeClock.StartTime.SET(mysql.TimestampT(activeStartTime)),
					tTimeClock.EndTime.SET(mysql.TimestampExp(mysql.NULL)),
					tTimeClock.SpentTime.SET(mysql.FloatExp(mysql.NULL)),
				)

			if _, err := insertStmt.ExecContext(ctx, d.db); err != nil {
				d.logger.Error(
					"failed to create active timeclock entry",
					zap.String("user", user.Identifier),
					zap.Error(err),
				)
				continue
			}
		}
	}
}
