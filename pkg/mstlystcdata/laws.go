package mstlystcdata

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"sort"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/laws"
	"github.com/fivenet-app/fivenet/v2025/pkg/croner"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/maruel/natural"
	"github.com/puzpuzpuz/xsync/v4"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Laws provides methods for loading, caching, and refreshing law books from the database.
type Laws struct {
	// logger for logging
	logger *zap.Logger
	// db is the database connection
	db *sql.DB

	// tracer is the OpenTelemetry tracer for this component
	tracer trace.Tracer

	// lawBooks is a concurrent map of law book IDs to LawBook structs
	lawBooks *xsync.Map[uint64, *laws.LawBook]
}

// LawsResult is the output struct for NewLaws, providing Laws and a cronjob register.
type LawsResult struct {
	fx.Out

	// Laws is the main Laws instance
	Laws *Laws
	// CronRegister is used to register cronjobs for law updates
	CronRegister croner.CronRegister `group:"cronjobregister"`
}

// NewLaws creates a new Laws instance, sets up lifecycle hooks, and returns a LawsResult.
func NewLaws(p Params) LawsResult {
	ctxCancel, cancel := context.WithCancel(context.Background())

	c := &Laws{
		logger:   p.Logger,
		db:       p.DB,
		tracer:   p.TP.Tracer("mstlystcdata.laws"),
		lawBooks: xsync.NewMap[uint64, *laws.LawBook](),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		if err := c.loadLaws(ctxCancel, 0); err != nil {
			c.logger.Error("failed to loads laws into cache", zap.Error(err))
			return err
		}
		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()
		return nil
	}))

	return LawsResult{
		Laws:         c,
		CronRegister: c,
	}
}

// RegisterCronjobs registers the law refresh cronjob with the given registry.
func (c *Laws) RegisterCronjobs(ctx context.Context, registry croner.IRegistry) error {
	if err := registry.RegisterCronjob(ctx, &cron.Cronjob{
		Name:     "mstlystcdata.laws",
		Schedule: "* * * * *", // Every minute
	}); err != nil {
		return err
	}
	return nil
}

// RegisterCronjobHandlers adds the handler for the law refresh cronjob.
func (c *Laws) RegisterCronjobHandlers(h *croner.Handlers) error {
	h.Add("mstlystcdata.laws", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := c.tracer.Start(ctx, "mstlystcdata-laws")
		defer span.End()

		if err := c.loadLaws(ctx, 0); err != nil {
			c.logger.Error("failed to refresh laws in cache", zap.Error(err))
			return err
		}
		return nil
	})
	return nil
}

// loadLaws loads law books and their laws from the database into the cache.
// If lawBookId is 0, loads all law books; otherwise, loads only the specified law book.
func (c *Laws) loadLaws(ctx context.Context, lawBookId uint64) error {
	tLawBooks := table.FivenetLawbooks.AS("lawbook")
	tLaws := table.FivenetLawbooksLaws.AS("law")

	stmt := tLawBooks.
		SELECT(
			tLawBooks.ID,
			tLawBooks.CreatedAt,
			tLawBooks.UpdatedAt,
			tLawBooks.Name,
			tLawBooks.Description,
			tLaws.ID,
			tLaws.LawbookID,
			tLaws.CreatedAt,
			tLaws.UpdatedAt,
			tLaws.Name,
			tLaws.Description,
			tLaws.Hint,
			tLaws.Fine,
			tLaws.DetentionTime,
			tLaws.StvoPoints,
		).
		FROM(tLawBooks.
			LEFT_JOIN(tLaws,
				tLaws.LawbookID.EQ(tLawBooks.ID),
			),
		).
		ORDER_BY(
			tLawBooks.SortKey.ASC(),
			tLaws.SortKey.ASC(),
		)

	if lawBookId > 0 {
		stmt = stmt.WHERE(
			tLawBooks.ID.EQ(jet.Uint64(lawBookId)),
		)
	}

	var dest []*laws.LawBook
	if err := stmt.QueryContext(ctx, c.db, &dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return err
		}
	}

	// Single lawbook update or lawbook deleted => not found, need to remove it
	if len(dest) == 0 {
		c.lawBooks.Delete(lawBookId)
	} else if lawBookId > 0 {
		c.lawBooks.Store(lawBookId, dest[0])
	} else {
		// Update cache
		found := []uint64{}
		for _, lawbook := range dest {
			c.lawBooks.Store(lawbook.Id, lawbook)
			found = append(found, lawbook.Id)
		}

		// Delete non-existing law books, based on which are in the database
		c.lawBooks.Range(func(key uint64, value *laws.LawBook) bool {
			if !slices.ContainsFunc(found, func(in uint64) bool {
				return in == key
			}) {
				c.lawBooks.Delete(key)
			}
			return true
		})
	}

	return nil
}

// GetLawBooks returns all cached law books, sorted by name using natural order.
func (c *Laws) GetLawBooks() []*laws.LawBook {
	lawBooks := []*laws.LawBook{}
	c.lawBooks.Range(func(key uint64, value *laws.LawBook) bool {
		lawBooks = append(lawBooks, value)
		return true
	})

	sort.Slice(lawBooks, func(i, j int) bool {
		return natural.Less(lawBooks[i].Name, lawBooks[j].Name)
	})

	return lawBooks
}

// Refresh reloads the specified law book (or all if lawBookId is 0) from the database.
func (c *Laws) Refresh(ctx context.Context, lawBookId uint64) error {
	if err := c.loadLaws(ctx, lawBookId); err != nil {
		return err
	}
	return nil
}
