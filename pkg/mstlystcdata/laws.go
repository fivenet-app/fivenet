package mstlystcdata

import (
	"context"
	"database/sql"
	"errors"
	"slices"
	"sort"
	"strings"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/laws"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"github.com/puzpuzpuz/xsync/v3"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Laws struct {
	logger *zap.Logger
	db     *sql.DB

	tracer trace.Tracer

	lawBooks *xsync.MapOf[uint64, *laws.LawBook]
}

func NewLaws(p Params) *Laws {
	ctxCancel, cancel := context.WithCancel(context.Background())

	c := &Laws{
		logger: p.Logger,
		db:     p.DB,

		tracer: p.TP.Tracer("mstlystcdata-laws"),

		lawBooks: xsync.NewMapOf[uint64, *laws.LawBook](),
	}

	p.LC.Append(fx.StartHook(func(ctxStartup context.Context) error {
		p.CronHandlers.Add("mstlystcdata.laws", func(ctx context.Context, data *cron.CronjobData) error {
			ctx, span := c.tracer.Start(ctx, "mstlystcdata-laws")
			defer span.End()

			if err := c.loadLaws(ctxCancel, 0); err != nil {
				c.logger.Error("failed to refresh laws cache", zap.Error(err))
				return err
			}

			return nil
		})

		if err := p.Cron.RegisterCronjob(ctxStartup, &cron.Cronjob{
			Name:     "mstlystcdata.laws",
			Schedule: "@always", // Every minute
		}); err != nil {
			return err
		}

		return nil
	}))

	p.LC.Append(fx.StopHook(func(_ context.Context) error {
		cancel()

		return nil
	}))

	return c
}

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
			tLawBooks.Name.ASC(),
			tLaws.Name.ASC(),
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

func (c *Laws) GetLawBooks() []*laws.LawBook {
	lawBooks := []*laws.LawBook{}
	c.lawBooks.Range(func(key uint64, value *laws.LawBook) bool {
		lawBooks = append(lawBooks, value)
		return true
	})

	sort.Slice(lawBooks, func(i, j int) bool {
		return strings.EqualFold(lawBooks[i].Name, lawBooks[j].Name)
	})

	return lawBooks
}

func (c *Laws) Refresh(ctx context.Context, lawBookId uint64) error {
	if err := c.loadLaws(ctx, lawBookId); err != nil {
		return err
	}

	return nil
}
