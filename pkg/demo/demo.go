// Package demo provides functionality for generating random dispatches and user locations
// in demo mode, allowing for testing and demonstration of the centrum service.
//
//nolint:gosec // G404: rand.Intn is not cryptographically secure, but we don't need it to be here.
package demo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/centrum"
	"github.com/fivenet-app/fivenet/v2025/pkg/config"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	"github.com/fivenet-app/fivenet/v2025/services/centrum/dispatches"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var tLocs = table.FivenetCentrumUserLocations

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

	dispatchStatusProgression = []centrum.StatusDispatch{
		centrum.StatusDispatch_STATUS_DISPATCH_NEW,
		centrum.StatusDispatch_STATUS_DISPATCH_EN_ROUTE,
		centrum.StatusDispatch_STATUS_DISPATCH_ON_SCENE,
		centrum.StatusDispatch_STATUS_DISPATCH_COMPLETED,
	}
)

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
	users []string
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

// Start runs the main demo loop, periodically generating and updating dispatches and moving user markers.
func (d *Demo) Start(ctx context.Context) error {
	tUsers := tables.User()

	stmt := tUsers.
		SELECT(
			tUsers.Identifier,
		).
		FROM(tUsers).
		WHERE(
			tUsers.Job.EQ(jet.String(d.cfg.Demo.TargetJob)),
		).
		ORDER_BY(tUsers.ID.ASC()).
		LIMIT(20)

	var users []string
	if err := stmt.QueryContext(ctx, d.db, &users); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return fmt.Errorf("failed to query users. %w", err)
		}
	}
	users = append(users, d.cfg.Demo.Users...)
	utils.RemoveSliceDuplicates(users)
	d.users = users

	go d.moveUserMarkers(ctx)

	for {
		d.logger.Info("Running demo dispatch generation cycle...")
		if err := d.updateDispatches(ctx); err != nil {
			d.logger.Error("failed to update dispatches", zap.Error(err))
		}

		d.generateDispatches(ctx)

		randWait := rand.Intn(300) + 30 // Random wait between 30 and 270 seconds

		select {
		case <-ctx.Done():
			return nil

		case <-time.After(time.Duration(randWait) * time.Second):
		}
	}
}

// generateDispatches creates up to 2 random dispatches per run with random positions and messages.
func (d *Demo) generateDispatches(ctx context.Context) {
	numDispatches := rand.Intn(2) // Up to 2 dispatches per run

	for range numDispatches {
		x := rand.Float64()*(xBounds[1]-xBounds[0]) + xBounds[0]
		y := rand.Float64()*(yBounds[1]-yBounds[0]) + yBounds[0]
		desc := dispatchDescriptions[rand.Intn(len(dispatchDescriptions))]
		msg := dispatchMessages[rand.Intn(len(dispatchMessages))]
		if _, err := d.dispatches.Create(ctx, &centrum.Dispatch{
			Jobs: &centrum.JobList{
				Jobs: []*centrum.Job{
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
		newStatus := &centrum.DispatchStatus{
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
				WHERE(tLocs.Identifier.EQ(jet.String(user))).
				LIMIT(1)

			err := stmt.QueryContext(ctx, d.db, &curr)
			if err != nil && !errors.Is(err, qrm.ErrNoRows) {
				d.logger.Error(
					"failed to select user location",
					zap.String("user", user),
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
					tLocs.X.SET(jet.Float(newX)),
					tLocs.Y.SET(jet.Float(newY)),
					tLocs.Hidden.SET(jet.Bool(false)),
				)

			if _, err := insertStmt.ExecContext(ctx, d.db); err != nil {
				d.logger.Error(
					"failed to update user location",
					zap.String("user", user),
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
