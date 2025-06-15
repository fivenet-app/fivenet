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
	"github.com/fivenet-app/fivenet/v2025/services/centrum/centrummanager"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var tLocs = table.FivenetCentrumUserLocations

var (
	// Define clamp bounds outside the loop for configurability
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

type Demo struct {
	logger *zap.Logger
	db     *sql.DB
	cs     *centrummanager.Manager
	cfg    *config.Config

	users []string
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger         *zap.Logger
	Cfg            *config.Config
	DB             *sql.DB
	CentrumManager *centrummanager.Manager
}

func New(p Params) *Demo {
	if !p.Cfg.Demo.Enabled {
		return nil
	}

	ctxCancel, cancel := context.WithCancel(context.Background())

	d := &Demo{
		logger: p.Logger.Named("demo"),
		db:     p.DB,
		cs:     p.CentrumManager,
		cfg:    p.Cfg,
	}

	d.logger.Warn("Demo mode is enabled. This will generate random dispatches and user locations!!!")

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
	d.users = append(users, d.cfg.Demo.Users...)
	utils.RemoveSliceDuplicates(d.users)

	go d.moveUserMarkers(ctx)

	for {
		d.logger.Info("Running demo dispatch generation cycle...")
		if err := d.updateDispatches(ctx); err != nil {
			d.logger.Error("failed to update dispatches", zap.Error(err))
		}

		if err := d.generateDispatches(ctx); err != nil {
			d.logger.Error("failed to generate dispatches", zap.Error(err))
		}

		randWait := rand.Intn(300) + 30 // Random wait between 30 and 270 seconds

		select {
		case <-ctx.Done():
			return nil

		case <-time.After(time.Duration(randWait) * time.Second):
		}
	}
}

func (d *Demo) generateDispatches(ctx context.Context) error {
	numDispatches := rand.Intn(2) // Up to 2 dispatches per run

	for range numDispatches {
		x := rand.Float64()*(xBounds[1]-xBounds[0]) + xBounds[0]
		y := rand.Float64()*(yBounds[1]-yBounds[0]) + yBounds[0]
		desc := dispatchDescriptions[rand.Intn(len(dispatchDescriptions))]
		msg := dispatchMessages[rand.Intn(len(dispatchMessages))]
		if _, err := d.cs.CreateDispatch(ctx, &centrum.Dispatch{
			Job:         d.cfg.Demo.TargetJob,
			Jobs:        []string{d.cfg.Demo.TargetJob},
			Message:     msg,
			Description: &desc,
			X:           x,
			Y:           y,
			Anon:        true,
		}); err != nil {
			d.logger.Error("failed to create dispatch", zap.Error(err))
		}
	}

	return nil
}

func (d *Demo) updateDispatches(ctx context.Context) error {
	dsps, _ := d.cs.ListDispatches(ctx, d.cfg.Demo.TargetJob)

	if len(dsps) == 0 {
		return nil
	}

	// Shuffle indices to pick random dispatches
	perm := rand.Perm(len(dsps))
	numToUpdate := min(len(dsps), 2)

	for i := range numToUpdate {
		dsp := dsps[perm[i]]

		// Randomize position slightly
		x := dsp.X + rand.Float64()*700 - 350
		y := dsp.Y + rand.Float64()*700 - 350

		// Pick a new status based on the status progression
		currStatus := dsp.Status.Status
		newStatusValue := currStatus
		for i, s := range dispatchStatusProgression {
			if s == currStatus && i+1 < len(dispatchStatusProgression) {
				newStatusValue = dispatchStatusProgression[i+1]
				break
			}
		}
		newStatus := &centrum.DispatchStatus{
			DispatchId: dsp.Id,
			Status:     newStatusValue,

			X: &x,
			Y: &y,
		}
		if _, err := d.cs.UpdateDispatchStatus(ctx, dsp.Job, dsp.Id, newStatus); err != nil {
			d.logger.Error("failed to update dispatch status", zap.Error(err))
		}
	}

	return nil
}

func (d *Demo) moveUserMarkers(ctx context.Context) error {
	firstRun := true

	for {
		select {
		case <-ctx.Done():
			return nil

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
				d.logger.Error("failed to select user location", zap.String("user", user), zap.Error(err))
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
				d.logger.Error("failed to update user location", zap.String("user", user), zap.Error(err))
				continue
			}
		}

		if firstRun {
			firstRun = false
		}
	}
}
