package livemapper

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/lru"
	jobs "github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/gen/go/proto/resources/livemap"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DispatchMarkerLimit = 60
)

var (
	tPlayerLocs = table.FivenetUserLocations
	tUsers      = table.Users.AS("user")
	tJobProps   = table.FivenetJobProps
)

var (
	FailedErr = status.Error(codes.Internal, "Failed to stream livemap data!")
)

type Server struct {
	LivemapperServiceServer

	tracer trace.Tracer
	ctx    context.Context
	logger *zap.Logger
	db     *sql.DB
	p      perms.Permissions
	c      *mstlystcdata.Enricher

	dispatchesCache *cache.Cache[string, []*livemap.DispatchMarker]
	usersCache      *cache.Cache[string, []*livemap.UserMarker]

	broker *utils.Broker[interface{}]

	visibleJobs []string
}

func NewServer(ctx context.Context, logger *zap.Logger, tp *tracesdk.TracerProvider, db *sql.DB, p perms.Permissions, c *mstlystcdata.Enricher, usersCacheSize int, visibleJobs []string) *Server {
	dispatchesCache := cache.NewContext(
		ctx,
		cache.AsLRU[string, []*livemap.DispatchMarker](lru.WithCapacity(len(visibleJobs))),
		cache.WithJanitorInterval[string, []*livemap.DispatchMarker](90*time.Second),
	)
	usersCache := cache.NewContext(
		ctx,
		cache.AsLRU[string, []*livemap.UserMarker](lru.WithCapacity(usersCacheSize)),
		cache.WithJanitorInterval[string, []*livemap.UserMarker](90*time.Second),
	)

	broker := utils.NewBroker[interface{}](ctx)
	go broker.Start()

	return &Server{
		ctx:    ctx,
		logger: logger,

		tracer: tp.Tracer("livemapper-cache"),
		db:     db,
		p:      p,
		c:      c,

		dispatchesCache: dispatchesCache,
		usersCache:      usersCache,

		broker: broker,

		visibleJobs: visibleJobs,
	}
}

func (s *Server) Start() {
	for {
		s.refreshCache()

		select {
		case <-s.ctx.Done():
			return
		// 3.85 seconds
		case <-time.After(3850 * time.Millisecond):
		}
	}
}

func (s *Server) refreshCache() {
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

func (s *Server) Stream(req *StreamRequest, srv LivemapperService_StreamServer) error {
	userInfo := auth.MustGetUserInfoFromContext(srv.Context())

	dispatchesAttr, err := s.p.Attr(userInfo, LivemapperServicePerm, LivemapperServiceStreamPerm, LivemapperServiceStreamDispatchesPermField)
	if err != nil {
		return FailedErr
	}
	var dispatchesJobs []string
	if dispatchesAttr != nil {
		dispatchesJobs = dispatchesAttr.([]string)
	}

	playersAttr, err := s.p.Attr(userInfo, LivemapperServicePerm, LivemapperServiceStreamPerm, LivemapperServiceStreamDispatchesPermField)
	if err != nil {
		return FailedErr
	}
	var playersJobs []string
	if playersAttr != nil {
		playersJobs = playersAttr.([]string)
	}

	resp := &StreamResponse{}

	if len(dispatchesJobs) == 0 && len(playersJobs) == 0 {
		if err := srv.Send(resp); err != nil {
			return err
		}

		return nil
	}

	// Add jobs to list visible jobs list
	resp.JobsDispatches = make([]*jobs.Job, len(dispatchesJobs))
	for i := 0; i < len(dispatchesJobs); i++ {
		resp.JobsDispatches[i] = &jobs.Job{
			Name: dispatchesJobs[i],
		}
		s.c.EnrichJobName(resp.JobsDispatches[i])
	}
	resp.JobsUsers = make([]*jobs.Job, len(playersJobs))
	for i := 0; i < len(playersJobs); i++ {
		resp.JobsUsers[i] = &jobs.Job{
			Name: playersJobs[i],
		}
		s.c.EnrichJobName(resp.JobsUsers[i])
	}

	signalCh := s.broker.Subscribe()
	defer s.broker.Unsubscribe(signalCh)

	for {
		dispatchMarkers, err := s.getUserDispatches(dispatchesJobs)
		if err != nil {
			return FailedErr
		}
		resp.Dispatches = dispatchMarkers

		userMarkers, err := s.getUserLocations(playersJobs, userInfo.UserId, userInfo.Job)
		if err != nil {
			return FailedErr
		}
		resp.Users = userMarkers

		if err := srv.Send(resp); err != nil {
			return err
		}

		select {
		case <-srv.Context().Done():
			return nil
		case <-signalCh:
		}
	}
}

func (s *Server) getUserLocations(jobs []string, userId int32, userJob string) ([]*livemap.UserMarker, error) {
	ds := []*livemap.UserMarker{}

	for _, job := range jobs {
		markers, ok := s.usersCache.Get(job)
		if !ok {
			continue
		}

		ds = append(ds, markers...)
	}

	return ds, nil
}

func (s *Server) getUserDispatches(jobs []string) ([]*livemap.DispatchMarker, error) {
	ds := []*livemap.DispatchMarker{}

	for _, job := range jobs {
		markers, ok := s.dispatchesCache.Get(job)
		if !ok {
			continue
		}

		ds = append(ds, markers...)
	}

	return ds, nil
}

func (s *Server) refreshUserLocations(ctx context.Context) error {
	markers := map[string][]*livemap.UserMarker{}

	locs := tPlayerLocs.AS("usermarker")
	stmt := locs.
		SELECT(
			locs.Identifier,
			locs.Job,
			locs.X,
			locs.Y,
			locs.UpdatedAt,
			tUsers.ID.AS("user.id"),
			tUsers.ID.AS("usermarker.id"),
			tUsers.Identifier,
			tUsers.Job,
			tUsers.JobGrade,
			tUsers.Firstname,
			tUsers.Lastname,
			tJobProps.LivemapMarkerColor.AS("usermarker.icon_color"),
		).
		FROM(
			locs.
				INNER_JOIN(tUsers,
					locs.Identifier.EQ(tUsers.Identifier),
				).
				LEFT_JOIN(tJobProps,
					tJobProps.Job.EQ(tUsers.Job),
				),
		).
		WHERE(
			locs.Hidden.IS_FALSE().
				AND(
					locs.UpdatedAt.GT_EQ(jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(60, jet.MINUTE))),
				),
		)

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
		if dest[i].IconColor == "" {
			dest[i].IconColor = jobs.DefaultLivemapMarkerColor
		}

		markers[job] = append(markers[job], dest[i])
	}
	for job, v := range markers {
		s.usersCache.Set(job, v, cache.WithExpiration(10*time.Minute))
	}

	return nil
}

func (s *Server) refreshDispatches(ctx context.Context) error {
	if len(s.visibleJobs) == 0 {
		s.logger.Warn("empty livemap jobs in config, no dispatches can be found because of that")
		return nil
	}

	d := table.GksphoneJobMessage
	stmt := d.
		SELECT(
			d.ID,
			d.Name,
			d.Number,
			d.Message,
			d.Gps,
			d.Owner,
			d.Jobm,
			d.Anon,
			d.Time,
		).
		FROM(
			d,
		).
		WHERE(
			jet.AND(
				d.Jobm.REGEXP_LIKE(jet.String("\\[\"("+strings.Join(s.visibleJobs, "|")+")\"\\]")),
				d.Time.GT_EQ(jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(20, jet.MINUTE))),
			),
		).
		ORDER_BY(
			d.Owner.ASC(),
			d.Time.DESC(),
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

		// Remove the "json" leftovers (in the gksphone table it looks like, e.g., `["ambulance"]`)
		job := strings.TrimSuffix(strings.TrimPrefix(*v.Jobm, "[\""), "\"]")
		if _, ok := markers[job]; !ok {
			markers[job] = []*livemap.DispatchMarker{}
		}
		marker := &livemap.DispatchMarker{
			X:         float32(x),
			Y:         float32(y),
			Id:        v.ID,
			Icon:      icon,
			IconColor: iconColor,
			Name:      name,
			Popup:     message,
			Job:       job,
			UpdatedAt: timestamp.New(v.Time),
		}
		if v.Owner == 0 {
			marker.Active = true
		}

		s.c.EnrichJobName(marker)
		markers[job] = append(markers[job], marker)
	}

	for job, v := range markers {
		s.dispatchesCache.Set(job, v, cache.WithExpiration(5*time.Minute))
	}

	return nil
}
