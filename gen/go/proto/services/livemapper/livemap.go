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
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DispatchMarkerLimit = 60
)

var (
	locs     = table.FivenetUserLocations
	users    = table.Users.AS("user")
	jobProps = table.FivenetJobProps
)

var (
	FailedErr = status.Error(codes.Internal, "Failed to stream livemap data!")
)

type Server struct {
	LivemapperServiceServer

	ctx    context.Context
	logger *zap.Logger
	db     *sql.DB
	p      perms.Permissions
	c      *mstlystcdata.Enricher

	dispatchesCache *cache.Cache[string, []*livemap.DispatchMarker]
	usersCache      *cache.Cache[string, []*livemap.UserMarker]

	broker *utils.Broker[interface{}]
}

func NewServer(ctx context.Context, logger *zap.Logger, db *sql.DB, p perms.Permissions, c *mstlystcdata.Enricher) *Server {
	dispatchesCache := cache.NewContext(
		ctx,
		cache.AsLRU[string, []*livemap.DispatchMarker](lru.WithCapacity(32)),
		cache.WithJanitorInterval[string, []*livemap.DispatchMarker](120*time.Second),
	)
	usersCache := cache.NewContext(
		ctx,
		cache.AsLRU[string, []*livemap.UserMarker](lru.WithCapacity(32)),
		cache.WithJanitorInterval[string, []*livemap.UserMarker](120*time.Second),
	)

	broker := utils.NewBroker[interface{}](ctx)
	go broker.Start()

	return &Server{
		ctx:    ctx,
		logger: logger,
		db:     db,
		p:      p,
		c:      c,

		dispatchesCache: dispatchesCache,
		usersCache:      usersCache,
		broker:          broker,
	}
}

func (s *Server) Start() {
	for {
		select {
		case <-s.ctx.Done():
			return
		// 3.85 seconds
		case <-time.After(3850 * time.Millisecond):
			if err := s.refreshUserLocations(); err != nil {
				s.logger.Error("failed to refresh livemap users cache", zap.Error(err))
			}
			if err := s.refreshDispatches(); err != nil {
				s.logger.Error("failed to refresh livemap dispatches cache", zap.Error(err))
			}
			s.broker.Publish(nil)
		}
	}
}

func (s *Server) Stream(req *StreamRequest, srv LivemapperService_StreamServer) error {
	userInfo := auth.GetUserInfoFromContext(srv.Context())

	dispatchesAttr, err := s.p.Attr(userInfo.CharID, userInfo.Job, userInfo.JobGrade, LivemapperServicePerm, LivemapperServiceStreamPerm, LivemapperServiceStreamDispatchesPermField)
	if err != nil {
		return FailedErr
	}
	var dispatchesJobs []string
	if dispatchesAttr != nil {
		dispatchesJobs = dispatchesAttr.([]string)
	}

	playersAttr, err := s.p.Attr(userInfo.CharID, userInfo.Job, userInfo.JobGrade, LivemapperServicePerm, LivemapperServiceStreamPerm, LivemapperServiceStreamDispatchesPermField)
	if err != nil {
		return FailedErr
	}
	var playersJobs []string
	if playersAttr != nil {
		playersJobs = playersAttr.([]string)
	}

	if len(dispatchesJobs) == 0 && len(playersJobs) == 0 {
		return nil
	}

	resp := &StreamResponse{}

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

		userMarkers, err := s.getUserLocations(playersJobs, userInfo.CharID, userInfo.Job)
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

func (s *Server) refreshUserLocations() error {
	markers := map[string][]*livemap.UserMarker{}

	locs := locs.AS("usermarker")
	stmt := locs.
		SELECT(
			locs.Identifier,
			locs.Job,
			locs.X,
			locs.Y,
			locs.UpdatedAt,
			users.ID.AS("user.id"),
			users.ID.AS("usermarker.id"),
			users.Identifier,
			users.Job,
			users.JobGrade,
			users.Firstname,
			users.Lastname,
			jobProps.LivemapMarkerColor.AS("usermarker.icon_color"),
		).
		FROM(
			locs.
				INNER_JOIN(users,
					locs.Identifier.EQ(users.Identifier),
				).
				LEFT_JOIN(jobProps,
					jobProps.Job.EQ(users.Job),
				),
		).
		WHERE(
			locs.Hidden.IS_FALSE().
				AND(
					locs.UpdatedAt.GT_EQ(jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(120, jet.MINUTE))),
				),
		)

	var dest []*livemap.UserMarker
	if err := stmt.QueryContext(s.ctx, s.db, &dest); err != nil {
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

func (s *Server) refreshDispatches() error {
	if len(config.C.Game.LivemapJobs) == 0 {
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
				d.Jobm.REGEXP_LIKE(jet.String("\\[\"("+strings.Join(config.C.Game.LivemapJobs, "|")+")\"\\]")),
				d.Time.GT_EQ(jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(20, jet.MINUTE))),
			),
		).
		ORDER_BY(
			d.Owner.ASC(),
			d.Time.DESC(),
		).
		LIMIT(DispatchMarkerLimit)

	var dest []*model.GksphoneJobMessage
	if err := stmt.QueryContext(s.ctx, s.db, &dest); err != nil {
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
