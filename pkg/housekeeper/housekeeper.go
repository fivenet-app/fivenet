package housekeeper

import (
	"context"
	"database/sql"
	"fmt"
	"slices"

	"github.com/fivenet-app/fivenet/gen/go/proto/resources/common/cron"
	"github.com/fivenet-app/fivenet/pkg/croner"
	jet "github.com/go-jet/jet/v2/mysql"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var Module = fx.Module("db_housekeeper",
	fx.Provide(
		New,
	),
)

const DefaultRetentionDays = 30

type Housekeeper struct {
	logger *zap.Logger
	tracer trace.Tracer

	db *sql.DB
}

type Params struct {
	fx.In

	LC fx.Lifecycle

	Logger *zap.Logger
	DB     *sql.DB
	TP     *tracesdk.TracerProvider

	Cron         croner.ICron
	CronHandlers *croner.Handlers
}

const lastTableMapIndex = "last_key"

func New(p Params) *Housekeeper {
	h := &Housekeeper{
		logger: p.Logger.Named("housekeeper"),
		tracer: p.TP.Tracer("housekeeper"),
		db:     p.DB,
	}

	p.LC.Append(fx.StartHook(func(ctx context.Context) error {
		if err := p.Cron.RegisterCronjob(ctx, &cron.Cronjob{
			Name:     "housekeeper.run",
			Schedule: "@5minutes", // Every 5 minutes
		}); err != nil {
			return err
		}

		return nil
	}))

	p.CronHandlers.Add("housekeeper.run", func(ctx context.Context, data *cron.CronjobData) error {
		ctx, span := h.tracer.Start(ctx, "housekeeper.run")
		defer span.End()

		dest := &cron.GenericCronData{
			Attributes: map[string]string{},
		}
		if data.Data == nil {
			data.Data = &anypb.Any{}
		}

		if err := anypb.UnmarshalTo(data.Data, dest, proto.UnmarshalOptions{}); err != nil {
			h.logger.Error("failed to unmarshal housekeeper cron data", zap.Error(err))
		}

		if err := h.runHousekeeper(ctx, dest); err != nil {
			return fmt.Errorf("error during docstore workflow handling. %w", err)
		}

		if err := data.Data.MarshalFrom(dest); err != nil {
			return fmt.Errorf("failed to marshal updated housekeeper cron data. %w", err)
		}

		return nil
	})

	return h
}

func (h *Housekeeper) runHousekeeper(ctx context.Context, data *cron.GenericCronData) error {
	keys := []string{}
	for key := range tablesList {
		keys = append(keys, key)
	}
	slices.Sort(keys)
	if len(keys) == 0 {
		return nil
	}

	lastTblKey, ok := data.Attributes[lastTableMapIndex]
	if !ok {
		// Take first table
		lastTblKey = keys[0]
	} else {
		idx := slices.Index(keys, lastTblKey)
		if idx == -1 || len(keys) <= idx+1 {
			h.logger.Warn("last table key not found in keys, falling back to first table")
			lastTblKey = keys[0]
		} else {
			lastTblKey = keys[idx+1]
		}
	}

	tbl, ok := tablesList[lastTblKey]
	if !ok {
		return nil
	}

	var condition jet.BoolExpression
	if tbl.TimestampColumn != nil {
		condition = jet.AND(
			tbl.TimestampColumn.IS_NOT_NULL(),
			tbl.TimestampColumn.LT_EQ(
				jet.CURRENT_DATE().SUB(jet.INTERVAL(tbl.MinDays, jet.DAY)),
			),
		)
	} else {
		condition = jet.AND(
			tbl.DateColumn.IS_NOT_NULL(),
			tbl.DateColumn.LT_EQ(
				jet.CAST(
					jet.CURRENT_DATE().SUB(jet.INTERVAL(tbl.MinDays, jet.DAY)),
				).AS_DATE(),
			),
		)
	}

	if tbl.Condition != nil {
		condition = condition.AND(tbl.Condition)
	}

	stmt := tbl.Table.
		DELETE().
		WHERE(condition).
		LIMIT(2000)

	res, err := stmt.ExecContext(ctx, h.db)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected > 0 {
		h.logger.Info("housekeeper run deleted rows", zap.String("table", lastTblKey), zap.Int64("rows", rowsAffected))
	}

	data.Attributes[lastTableMapIndex] = lastTblKey

	return nil
}
