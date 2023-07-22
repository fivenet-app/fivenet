package tracker

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/gen/go/proto/resources/livemap"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/pkg/utils/syncx"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const (
	DispatchMarkerLimit = 60
)

var (
	tLocs     = table.FivenetUserLocations
	tUsers    = table.Users.AS("user")
	tJobProps = table.FivenetJobProps
)

// TODO keep log of on duty (non hidden) player locations

// TODO keep track of all markers in database (in the future also "Sperrzonen" and Panicbuttons)

type Tracker struct {
	ctx    context.Context
	logger *zap.Logger
	tracer trace.Tracer
	db     *sql.DB
	c      *mstlystcdata.Enricher

	dispatchesCache syncx.Map[string, []*livemap.DispatchMarker]
	usersCache      syncx.Map[string, []*livemap.UserMarker]

	broker *utils.Broker[interface{}]

	refreshTime time.Duration
	visibleJobs []string
}

func New(ctx context.Context, logger *zap.Logger, tp *tracesdk.TracerProvider, db *sql.DB, c *mstlystcdata.Enricher, refreshTime time.Duration, visibleJobs []string) *Tracker {
	broker := utils.NewBroker[interface{}](ctx)
	go broker.Start()

	return &Tracker{
		ctx:    ctx,
		logger: logger,
		tracer: tp.Tracer("tracker-cache"),
		db:     db,
		c:      c,

		dispatchesCache: syncx.Map[string, []*livemap.DispatchMarker]{},
		usersCache:      syncx.Map[string, []*livemap.UserMarker]{},

		broker: broker,

		refreshTime: refreshTime,
		visibleJobs: visibleJobs,
	}
}

func (s *Tracker) Start() {
	for {
		s.refreshCache()

		select {
		case <-s.ctx.Done():
			s.broker.Stop()
			return
		case <-time.After(s.refreshTime):
		}
	}
}

func (s *Tracker) refreshCache() {
	ctx, span := s.tracer.Start(s.ctx, "livemap-refresh-cache")
	defer span.End()

	if err := s.refreshUserLocations(ctx); err != nil {
		s.logger.Error("failed to refresh livemap users cache", zap.Error(err))
	}
	if err := s.refreshDispatches(ctx); err != nil {
		s.logger.Error("failed to refresh livemap dispatches cache", zap.Error(err))
	}

	s.broker.Publish(nil)
}

func (s *Tracker) refreshUserLocations(ctx context.Context) error {
	markers := map[string][]*livemap.UserMarker{}

	tLocs := tLocs.AS("genericmarker")
	stmt := tLocs.
		SELECT(
			tLocs.Identifier,
			tLocs.Job,
			tLocs.X,
			tLocs.Y,
			tLocs.UpdatedAt,
			tUsers.ID.AS("user.id"),
			tUsers.ID.AS("genericmarker.id"),
			tUsers.Identifier,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
			tJobProps.LivemapMarkerColor.AS("genericmarker.icon_color"),
		).
		FROM(
			tLocs.
				INNER_JOIN(tUsers,
					tLocs.Identifier.EQ(tUsers.Identifier),
				).
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tUsers.Job),
				),
		).
		WHERE(jet.AND(
			tLocs.Hidden.IS_FALSE(),
			tLocs.UpdatedAt.GT_EQ(jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(60, jet.MINUTE))),
		))

	var dest []*livemap.UserMarker
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	for i := 0; i < len(dest); i++ {
		s.c.EnrichJobInfo(dest[i].User)

		job := dest[i].User.Job
		if _, ok := markers[job]; !ok {
			markers[job] = []*livemap.UserMarker{}
		}
		if dest[i].Marker.IconColor == "" {
			dest[i].Marker.IconColor = jobs.DefaultLivemapMarkerColor
		}

		markers[job] = append(markers[job], dest[i])
	}
	for job, v := range markers {
		s.usersCache.Store(job, v)
	}

	return nil
}

func (s *Tracker) refreshDispatches(ctx context.Context) error {
	if len(s.visibleJobs) == 0 {
		s.logger.Warn("empty livemap jobs in config, no dispatches can be found because of that")
		return nil
	}

	gksphoneJobM := table.GksphoneJobMessage
	stmt := gksphoneJobM.
		SELECT(
			gksphoneJobM.ID,
			gksphoneJobM.Name,
			gksphoneJobM.Number,
			gksphoneJobM.Message,
			gksphoneJobM.Gps,
			gksphoneJobM.Owner,
			gksphoneJobM.Jobm,
			gksphoneJobM.Anon,
			gksphoneJobM.Time,
		).
		FROM(
			gksphoneJobM,
		).
		WHERE(
			jet.AND(
				gksphoneJobM.Jobm.REGEXP_LIKE(jet.String("\\[\"("+strings.Join(s.visibleJobs, "|")+")\"\\]")),
				gksphoneJobM.Time.GT_EQ(jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(20, jet.MINUTE))),
			),
		).
		ORDER_BY(
			gksphoneJobM.Owner.ASC(),
			gksphoneJobM.Time.DESC(),
		).
		LIMIT(DispatchMarkerLimit)

	var dest []*model.GksphoneJobMessage
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return err
	}

	markers := map[string][]*livemap.DispatchMarker{}
	for _, v := range dest {
		gps, _ := strings.CutPrefix(*v.Gps, "GPS: ")
		gpsSplit := strings.Split(gps, ", ")
		x, _ := strconv.ParseFloat(gpsSplit[0], 32)
		y, _ := strconv.ParseFloat(gpsSplit[1], 32)

		var icon string
		var iconColor string
		if v.Owner == 0 {
			icon = "dispatch-open.svg"
			iconColor = "96E6B3"
		} else {
			icon = "dispatch-closed.svg"
			iconColor = "DA3E52"
		}

		var name string
		if v.Anon != nil && *v.Anon == "1" {
			name = "Anonym"
		} else {
			name = *v.Name
		}

		var message string
		if v.Message != nil && *v.Message != "" {
			message = *v.Message
		} else {
			message = "N/A"
		}

		// Remove the "json" leftovers (the data looks like this, e.g., `["ambulance"]`)
		job := strings.TrimSuffix(strings.TrimPrefix(*v.Jobm, "[\""), "\"]")
		if _, ok := markers[job]; !ok {
			markers[job] = []*livemap.DispatchMarker{}
		}
		marker := &livemap.DispatchMarker{
			Marker: &livemap.GenericMarker{
				Id:        v.ID,
				X:         float32(x),
				Y:         float32(y),
				Icon:      icon,
				IconColor: iconColor,
				Name:      name,
				Popup:     message,
				UpdatedAt: timestamp.New(v.Time),
			},
			Job: job,
		}
		if v.Owner == 0 {
			marker.Active = true
		}

		s.c.EnrichJobName(marker)
		markers[job] = append(markers[job], marker)
	}

	for job, v := range markers {
		s.dispatchesCache.Store(job, v)
	}

	return nil
}
