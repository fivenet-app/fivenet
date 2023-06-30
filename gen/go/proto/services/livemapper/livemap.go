package livemapper

import (
	"context"
	"database/sql"
	"strconv"
	"strings"
	"time"

	jobs "github.com/galexrt/fivenet/gen/go/proto/resources/jobs"
	"github.com/galexrt/fivenet/gen/go/proto/resources/livemap"
	"github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	"github.com/galexrt/fivenet/pkg/grpc/auth"
	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/utils"
	"github.com/galexrt/fivenet/pkg/utils/syncx"
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
	tLocs     = table.FivenetUserLocations
	tUsers    = table.Users.AS("user")
	tJobProps = table.FivenetJobProps
)

var (
	ErrStreamFailed = status.Error(codes.Internal, "errors.LivemapperService.ErrStreamFailed")
)

type Server struct {
	LivemapperServiceServer

	tracer trace.Tracer
	ctx    context.Context
	logger *zap.Logger
	db     *sql.DB
	p      perms.Permissions
	c      *mstlystcdata.Enricher

	dispatchesCache syncx.Map[string, []*livemap.DispatchMarker]
	usersCache      syncx.Map[string, []*livemap.UserMarker]

	broker *utils.Broker[interface{}]

	refreshTime time.Duration
	visibleJobs []string
}

func NewServer(ctx context.Context, logger *zap.Logger, tp *tracesdk.TracerProvider, db *sql.DB, p perms.Permissions, c *mstlystcdata.Enricher, refreshTime time.Duration, visibleJobs []string) *Server {
	broker := utils.NewBroker[interface{}](ctx)
	go broker.Start()

	return &Server{
		ctx:    ctx,
		logger: logger,

		tracer: tp.Tracer("livemapper-cache"),
		db:     db,
		p:      p,
		c:      c,

		dispatchesCache: syncx.Map[string, []*livemap.DispatchMarker]{},
		usersCache:      syncx.Map[string, []*livemap.UserMarker]{},

		broker: broker,

		refreshTime: refreshTime,
		visibleJobs: visibleJobs,
	}
}

func (s *Server) Start() {
	for {
		s.refreshCache()

		select {
		case <-s.ctx.Done():
			return
		case <-time.After(s.refreshTime):
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
		return ErrStreamFailed
	}
	playersAttr, err := s.p.Attr(userInfo, LivemapperServicePerm, LivemapperServiceStreamPerm, LivemapperServiceStreamPlayersPermField)
	if err != nil {
		return ErrStreamFailed
	}

	var dispatchesJobs []string
	if dispatchesAttr != nil {
		dispatchesJobs = dispatchesAttr.([]string)
	}
	if userInfo.SuperUser {
		dispatchesJobs = s.visibleJobs
	}

	var playersJobs map[string]int32
	if playersAttr != nil {
		playersJobs, _ = playersAttr.(map[string]int32)
	}

	if userInfo.SuperUser {
		playersJobs = map[string]int32{}
		for _, j := range s.visibleJobs {
			playersJobs[j] = -1
		}
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
	resp.JobsUsers = []*jobs.Job{}
	for job := range playersJobs {
		j := &jobs.Job{
			Name: job,
		}
		s.c.EnrichJobName(j)
		resp.JobsUsers = append(resp.JobsUsers, j)
	}

	signalCh := s.broker.Subscribe()
	defer s.broker.Unsubscribe(signalCh)

	for {
		userMarkers, _, err := s.getUserLocations(playersJobs, userInfo)
		if err != nil {
			return ErrStreamFailed
		}
		resp.Users = userMarkers

		dispatchMarkers, err := s.getUserDispatches(dispatchesJobs)
		if err != nil {
			return ErrStreamFailed
		}
		resp.Dispatches = dispatchMarkers

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

func (s *Server) getUserLocations(jobs map[string]int32, userInfo *userinfo.UserInfo) ([]*livemap.UserMarker, bool, error) {
	ds := []*livemap.UserMarker{}

	found := false
	if userInfo.SuperUser {
		found = true
	}
	for job, grade := range jobs {
		markers, ok := s.usersCache.Load(job)
		if !ok {
			continue
		}

		for i := 0; i < len(markers); i++ {
			// SuperUser returns grade as `-1`, job has access to that grade or it is the user itself
			if grade == -1 || (markers[i].User.JobGrade <= grade || markers[i].User.UserId == userInfo.UserId) {
				ds = append(ds, markers[i])
			}

			// If the user is found in the list of user markers, set found state
			if !found && (userInfo.Job == job && markers[i].Id == userInfo.UserId) {
				found = true
			}
		}
	}

	if found {
		return ds, true, nil
	}

	return nil, false, nil
}

func (s *Server) getUserDispatches(jobs []string) ([]*livemap.DispatchMarker, error) {
	ds := []*livemap.DispatchMarker{}

	for _, job := range jobs {
		markers, ok := s.dispatchesCache.Load(job)
		if !ok {
			continue
		}

		ds = append(ds, markers...)
	}

	return ds, nil
}

func (s *Server) refreshUserLocations(ctx context.Context) error {
	markers := map[string][]*livemap.UserMarker{}

	tLocs := tLocs.AS("usermarker")
	stmt := tLocs.
		SELECT(
			tLocs.Identifier,
			tLocs.Job,
			tLocs.X,
			tLocs.Y,
			tLocs.UpdatedAt,
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
		if dest[i].IconColor == "" {
			dest[i].IconColor = jobs.DefaultLivemapMarkerColor
		}

		markers[job] = append(markers[job], dest[i])
	}
	for job, v := range markers {
		s.usersCache.Store(job, v)
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
		s.dispatchesCache.Store(job, v)
	}

	return nil
}
