package livemapper

import (
	"context"
	"database/sql"
	"math/rand"
	"strconv"
	"strings"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/Code-Hex/go-generics-cache/policy/lru"
	"github.com/galexrt/fivenet/pkg/auth"
	"github.com/galexrt/fivenet/pkg/config"
	"github.com/galexrt/fivenet/pkg/mstlystcdata"
	"github.com/galexrt/fivenet/pkg/perms"
	"github.com/galexrt/fivenet/pkg/utils"
	jobs "github.com/galexrt/fivenet/proto/resources/jobs"
	"github.com/galexrt/fivenet/proto/resources/livemap"
	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

var (
	locs  = table.FivenetUserLocations
	users = table.Users.AS("user")
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
	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			case <-time.After(4 * time.Second):
				if err := s.refreshUserLocations(); err != nil {
					s.logger.Error("failed to refresh livemap users cache", zap.Error(err))
				}
				if err := s.refreshDispatches(); err != nil {
					s.logger.Error("failed to refresh livemap dispatches cache", zap.Error(err))
				}
				s.broker.Publish(nil)
			}
		}
	}()
}

func (s *Server) Stream(req *StreamRequest, srv LivemapperService_StreamServer) error {
	userId := auth.GetUserIDFromContext(srv.Context())

	js, err := s.p.GetSuffixOfPermissionsByPrefixOfUser(userId, "LivemapperService.Stream")
	if err != nil {
		return err
	}

	if len(js) == 0 {
		return nil
	}

	resp := &StreamResponse{}
	resp.Jobs = make([]*jobs.Job, len(js))
	for i := 0; i < len(resp.Jobs); i++ {
		resp.Jobs[i] = &jobs.Job{
			Name: js[i],
		}
		s.c.EnrichJobName(resp.Jobs[i])
	}

	signalCh := s.broker.Subscribe()
	defer s.broker.Unsubscribe(signalCh)

	for {
		dispatchMarkers, err := s.getUserDispatches(js)
		if err != nil {
			return err
		}
		resp.Dispatches = dispatchMarkers

		userMarkers, err := s.getUserLocations(js)
		if err != nil {
			return err
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

func (s *Server) getUserLocations(jobs []string) ([]*livemap.UserMarker, error) {
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
			users.ID.AS("usermarker.userid"),
			users.ID,
			users.Identifier,
			users.Job,
			users.JobGrade,
			users.Firstname,
			users.Lastname,
			jet.String("5C7AFF").AS("usermarker.icon_color"),
		).
		FROM(
			locs.
				LEFT_JOIN(users,
					locs.Identifier.EQ(users.Identifier),
				),
		).
		WHERE(
			locs.Hidden.IS_FALSE().
				AND(
					locs.UpdatedAt.GT_EQ(jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(10, jet.MINUTE))),
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
		markers[job] = append(markers[job], dest[i])
	}
	for job, v := range markers {
		s.usersCache.Set(job, v, cache.WithExpiration(5*time.Minute))
	}

	return nil
}

func (s *Server) refreshDispatches() error {
	d := table.GksphoneJobMessage
	stmt := d.
		SELECT(
			d.ID,
			d.Name,
			d.Number,
			d.Message,
			d.Photo,
			d.Gps,
			d.Owner,
			d.Jobm,
			d.Anon,
		).
		FROM(
			d,
		).
		WHERE(
			jet.AND(
				d.Jobm.REGEXP_LIKE(jet.String("\\[\"("+strings.Join(config.C.FiveM.LivemapJobs, "|")+")\"\\]")),
				d.Time.GT_EQ(jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(20, jet.MINUTE))),
			),
		).
		ORDER_BY(d.Time.DESC()).
		LIMIT(40)

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

		// Remove the "json" leftovers (in the gksphone table it looks like `["ambulance"]`)
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
			Name:      *v.Name,
			Popup:     *v.Message,
			Job:       job,
		}

		s.c.EnrichJobName(marker)
		markers[job] = append(markers[job], marker)
	}

	for job, v := range markers {
		s.dispatchesCache.Set(job, v, cache.WithExpiration(5*time.Minute))
	}

	return nil
}

func (s *Server) GenerateRandomUserMarker() {
	userIdentifiers := []string{
		"char1:fcee377a1fda007a8d2cc764a0a272e04d8c5d57",
		"char1:0ff2f772f2527a0626cac48670cbc20ddbdc09fb",
		"char2:d9793ddb457316fb3951d1b1092526183270a307",
		"char2:d7abbfba01625bec803788ee42da86461c96e0bd",
		"char1:ad4fb9f44bb784dd30effcc743a9c169db4d625d",
	}

	markers := make([]*model.FivenetUserLocations, len(userIdentifiers))

	resetMarkers := func() {
		xMin := -3300
		xMax := 4300
		yMin := -3300
		yMax := 5000
		for i := 0; i < len(markers); i++ {
			x := float64(rand.Intn(xMax-xMin+1) + xMin)
			y := float64(rand.Intn(yMax-yMin+1) + yMin)

			job := "ambulance"
			hidden := false
			markers[i] = &model.FivenetUserLocations{
				Identifier: userIdentifiers[i],
				Job:        &job,
				Hidden:     &hidden,

				X: &x,
				Y: &y,
			}
		}
	}

	moveMarkers := func() {
		xMin := -100
		xMax := 100
		yMin := -100
		yMax := 100

		for i := 0; i < len(markers); i++ {
			curX := *markers[i].X
			curY := *markers[i].Y

			newX := curX + float64(rand.Intn(xMax-xMin+1)+xMin)
			newY := curY + float64(rand.Intn(yMax-yMin+1)+yMin)

			markers[i].X = &newX
			markers[i].Y = &newY
		}
	}

	resetMarkers()

	counter := 0
	for {
		if counter >= 15 {
			resetMarkers()
			counter = 0
		} else {
			moveMarkers()
		}

		stmt := locs.
			INSERT(
				locs.Identifier,
				locs.Job,
				locs.X,
				locs.Y,
				locs.Hidden,
			).
			MODELS(markers).
			ON_DUPLICATE_KEY_UPDATE(
				locs.X.SET(jet.RawFloat("VALUES(x)")),
				locs.Y.SET(jet.RawFloat("VALUES(y)")),
			)

		_, err := stmt.Exec(s.db)
		if err != nil {
			s.logger.Error("failed to insert/ update random location to locations table", zap.Error(err))
		}

		counter++
		time.Sleep(3 * time.Second)
	}
}
