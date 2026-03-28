// Package demo provides functionality for generating random dispatches, user locations,
// and optional fake user data in demo mode.
//
//nolint:gosec // G404: randomness is non-cryptographic and intended for demo data only.
package demo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/fivenet-app/fivenet/v2026/pkg/config"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/fivenet-app/fivenet/v2026/services/centrum/dispatches"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	tAccounts             = table.FivenetAccounts
	tJobs                 = table.FivenetJobs
	tJobsGrades           = table.FivenetJobsGrades
	tJobProps             = table.FivenetJobProps
	tLicenses             = table.FivenetLicenses
	tLawbooks             = table.FivenetLawbooks
	tLawbooksLaws         = table.FivenetLawbooksLaws
	tRbacPermissions      = table.FivenetRbacPermissions
	tRbacRoles            = table.FivenetRbacRoles
	tRbacRolePerms        = table.FivenetRbacRolesPermissions
	tRbacJobPerms         = table.FivenetRbacJobPermissions
	tRbacAttrs            = table.FivenetRbacAttrs
	tRbacRoleAttrs        = table.FivenetRbacRolesAttrs
	tRbacJobAttrs         = table.FivenetRbacJobAttrs
	tJobColleagueActivity = table.FivenetJobColleagueActivity
	tLocs                 = table.FivenetCentrumUserLocations
	tTimeClock            = table.FivenetJobTimeclock
	tUsers                = table.FivenetUser
	tUserJobs             = table.FivenetUserJobs
	tUserLicenses         = table.FivenetUserLicenses
	tUserPhoneNumbers     = table.FivenetUserPhoneNumbers
	tUserProps            = table.FivenetUserProps
	tOwnedVehicles        = table.FivenetOwnedVehicles
	tVehicleProps         = table.FivenetVehiclesProps
	tCentrumSettings      = table.FivenetCentrumSettings
	tCentrumUnits         = table.FivenetCentrumUnits
)

var Module = fx.Module(
	"demo",
	fx.Provide(
		New,
	),
)

const (
	defaultUserLookupLimit   = 500
	minRuntimeTargetJobUsers = 10
)

var (
	xBounds = []float64{-2750, 2500}
	yBounds = []float64{-3000, 6000}
)

type user struct {
	UserID     int32  `alias:"user.userid"     sql:"primary_key"`
	Identifier string `alias:"user.identifier"`
}

type startupGenerator interface {
	Name() string
	Enabled(*Demo) bool
	Run(context.Context, *Demo) error
}

type demoCatalogGenerator struct{}

func (g demoCatalogGenerator) Name() string {
	return "catalog_seed"
}

func (g demoCatalogGenerator) Enabled(_ *Demo) bool {
	return true
}

func (g demoCatalogGenerator) Run(ctx context.Context, d *Demo) error {
	return d.seedDemoCatalog(ctx)
}

// Demo provides demo mode functionality for generating random dispatches, user locations,
// timeclock rows and optional fake users.
type Demo struct {
	logger     *zap.Logger
	db         *sql.DB
	dispatches *dispatches.DispatchDB
	cfg        *config.Config

	randMu sync.Mutex
	rng    *rand.Rand
	fake   *gofakeit.Faker

	demoJobNames  []string
	demoJobGrades map[string][]int32

	users []*user
}

// Params contains dependencies for constructing a Demo instance.
type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger     *zap.Logger
	Cfg        *config.Config
	DB         *sql.DB
	Dispatches *dispatches.DispatchDB
}

// New creates a new Demo instance if demo mode is enabled in the config.
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
	d.initRandomizers()

	d.logger.Warn("Demo mode is enabled. This will generate demo data and user locations.")

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

func (d *Demo) initRandomizers() {
	seed := d.cfg.Demo.Seed
	d.rng = rand.New(rand.NewPCG(seed, seed))
	d.fake = gofakeit.New(seed)
	d.initDemoJobCatalog()
}

func (d *Demo) startupGenerators() []startupGenerator {
	return []startupGenerator{
		demoCatalogGenerator{},
		fakeUsersGenerator{},
	}
}

func (d *Demo) runStartupGenerators(ctx context.Context) error {
	for _, generator := range d.startupGenerators() {
		if !generator.Enabled(d) {
			continue
		}

		d.logger.Info("running demo startup generator", zap.String("generator", generator.Name()))
		if err := generator.Run(ctx, d); err != nil {
			return fmt.Errorf("failed to run startup generator %s. %w", generator.Name(), err)
		}
	}

	return nil
}

func (d *Demo) lookupUsers(
	ctx context.Context,
	identifiers []string,
	limit int64,
) ([]*user, error) {
	condition := tUsers.Job.EQ(mysql.String(d.cfg.Demo.TargetJob))
	if len(identifiers) > 0 {
		idds := make([]mysql.Expression, len(identifiers))
		for i, id := range identifiers {
			idds[i] = mysql.String(id)
		}
		condition = condition.AND(tUsers.Identifier.IN(idds...))
	}

	if limit <= 0 {
		limit = defaultUserLookupLimit
	}

	stmt := tUsers.
		SELECT(
			tUsers.ID.AS("user.userid"),
			tUsers.Identifier.AS("user.identifier"),
		).
		FROM(tUsers).
		WHERE(condition).
		ORDER_BY(tUsers.ID.ASC()).
		LIMIT(limit)

	users := []*user{}
	if err := stmt.QueryContext(ctx, d.db, &users); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, fmt.Errorf("failed to query users. %w", err)
		}
	}

	return users, nil
}

// Start runs the demo runtime loops based on enabled feature toggles.
func (d *Demo) Start(ctx context.Context) error {
	if err := d.runStartupGenerators(ctx); err != nil {
		d.logger.Error("failed to run startup demo generators", zap.Error(err))
		return err
	}

	if d.cfg.Demo.Features.Locations || d.cfg.Demo.Features.Timeclock {
		users, err := d.buildRuntimeUsers(ctx)
		if err != nil {
			d.logger.Error("failed to build runtime user list", zap.Error(err))
			return err
		}
		d.users = users
	}

	if d.cfg.Demo.Features.Timeclock {
		d.createTimeclockEntries(ctx)
	}

	if d.cfg.Demo.Features.Locations {
		go d.moveUserMarkers(ctx)
	}

	if !d.cfg.Demo.Features.Dispatches {
		<-ctx.Done()
		return nil
	}

	for {
		d.logger.Info("running demo dispatch generation cycle")
		if err := d.updateDispatches(ctx); err != nil {
			d.logger.Error("failed to update dispatches", zap.Error(err))
		}

		d.generateDispatches(ctx)

		randWait := d.randIntN(300) + 30 // 30-329 seconds

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(time.Duration(randWait) * time.Second):
		}
	}
}

func (d *Demo) buildRuntimeUsers(ctx context.Context) ([]*user, error) {
	limit := int64(max(defaultUserLookupLimit, d.cfg.Demo.FakeUsers.Count*2))
	if err := d.ensureRuntimeTargetJobUsers(ctx, minRuntimeTargetJobUsers); err != nil {
		return nil, err
	}

	users, err := d.lookupUsers(ctx, nil, int64(max(int(limit), minRuntimeTargetJobUsers)))
	if err != nil {
		return nil, err
	}

	mainIdentifier := d.getMainCharacterIdentifier()
	foundMain := false
	for _, u := range users {
		if u.Identifier == mainIdentifier {
			foundMain = true
			break
		}
	}
	if !foundMain {
		mainUser, err := d.lookupUsers(ctx, []string{mainIdentifier}, 1)
		if err != nil {
			return nil, err
		}
		users = append(mainUser, users...)
	}

	seen := map[string]struct{}{}
	out := make([]*user, 0, len(users))
	for _, u := range users {
		if _, ok := seen[u.Identifier]; ok {
			continue
		}
		seen[u.Identifier] = struct{}{}
		out = append(out, u)
	}

	return out, nil
}

func (d *Demo) randIntN(n int) int {
	d.randMu.Lock()
	defer d.randMu.Unlock()
	return d.rng.IntN(n)
}

func (d *Demo) randInt32N(n int32) int32 {
	d.randMu.Lock()
	defer d.randMu.Unlock()
	return d.rng.Int32N(n)
}

func (d *Demo) randFloat64() float64 {
	d.randMu.Lock()
	defer d.randMu.Unlock()
	return d.rng.Float64()
}

func (d *Demo) randPerm(n int) []int {
	d.randMu.Lock()
	defer d.randMu.Unlock()
	return d.rng.Perm(n)
}

// moveUserMarkers periodically updates user locations.
func (d *Demo) moveUserMarkers(ctx context.Context) {
	if len(d.users) == 0 {
		d.logger.Info("demo location generator has no users to move")
		return
	}

	firstRun := true

	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(2 * time.Second):
		}

		for _, user := range d.users {
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
				WHERE(tLocs.UserID.EQ(mysql.Int32(user.UserID))).
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
			if firstRun || errors.Is(err, qrm.ErrNoRows) {
				newX = d.randFloat64()*(xBounds[1]-xBounds[0]) + xBounds[0]
				newY = d.randFloat64()*(yBounds[1]-yBounds[0]) + yBounds[0]
			} else {
				maxStep := 35.0
				deltaX := (d.randFloat64()*2 - 1) * maxStep
				deltaY := (d.randFloat64()*2 - 1) * maxStep
				newX = curr.X + deltaX
				newY = curr.Y + deltaY

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
					tLocs.UserID,
					tLocs.X,
					tLocs.Y,
					tLocs.Hidden,
				).
				VALUES(user.UserID, newX, newY, false).
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

// createTimeclockEntries creates or updates timeclock entries for demo users.
func (d *Demo) createTimeclockEntries(ctx context.Context) {
	if len(d.users) == 0 {
		d.logger.Info("demo timeclock generator has no users to seed")
		return
	}

	for _, user := range d.users {
		numPastEntries := d.randIntN(10) + 1
		for range numPastEntries {
			daysAgo := d.randIntN(14) + 1
			entryDate := time.Now().AddDate(0, 0, -daysAgo)
			startTime := entryDate.Add(time.Duration(d.randIntN(8)) * time.Hour)
			endTime := startTime.Add(time.Duration(d.randIntN(8)+1) * time.Hour)
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

		if d.randIntN(2) == 0 {
			activeStartTime := time.Now().Add(-time.Duration(d.randIntN(8)+1) * time.Hour)

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
					nil,
					nil,
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
