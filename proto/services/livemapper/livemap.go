package livemapper

import (
	"context"
	"database/sql"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/galexrt/arpanet/proto/resources/livemap"
	"github.com/galexrt/arpanet/query/arpanet/model"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

var (
	locs  = table.ArpanetUserLocations
	users = table.Users.AS("user")
)

type Server struct {
	LivemapperServiceServer

	logger *zap.Logger
	db     *sql.DB
	p      perms.Permissions
}

func NewServer(logger *zap.Logger, db *sql.DB, p perms.Permissions) *Server {
	return &Server{
		logger: logger,
		db:     db,
		p:      p,
	}
}

func (s *Server) SetLogger(logger *zap.Logger) {
	s.logger = logger
}

func (s *Server) Stream(req *StreamRequest, srv LivemapperService_StreamServer) error {
	userId := auth.GetUserIDFromContext(srv.Context())

	jobs, err := s.p.GetSuffixOfPermissionsByPrefixOfUser(userId, "LivemapperService.Stream")
	if err != nil {
		return err
	}

	if len(jobs) == 0 {
		return nil
	}

	sqlJobs := make([]jet.Expression, len(jobs))
	for k := 0; k < len(jobs); k++ {
		sqlJobs[k] = jet.String(jobs[k])
	}

	locs := locs.AS("usermarker")
	stmt := locs.
		SELECT(
			locs.UserID.AS("usermarker.userid"),
			locs.Job,
			locs.X,
			locs.Y,
			locs.UpdatedAt,
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
					locs.UserID.EQ(users.ID),
				),
		).
		WHERE(
			locs.Job.IN(sqlJobs...).
				AND(
					locs.Hidden.IS_FALSE(),
				).
				AND(
					locs.UpdatedAt.GT_EQ(jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(5, jet.MINUTE))),
				),
		)

	for {
		resp := &StreamResponse{}

		if err := stmt.QueryContext(srv.Context(), s.db, &resp.Users); err != nil {
			return err
		}

		ds, err := s.getDispatches(srv.Context(), jobs)
		if err != nil {
			return err
		}
		resp.Dispatches = ds

		if err := srv.Send(resp); err != nil {
			return err
		}
		time.Sleep(3 * time.Second)
	}
}

func (s *Server) getDispatches(ctx context.Context, jobs []string) ([]*livemap.DispatchMarker, error) {
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
				d.Jobm.REGEXP_LIKE(jet.String("\\[\"("+strings.Join(jobs, "|")+")\"\\]")),
				d.Time.GT_EQ(jet.CURRENT_TIMESTAMP().SUB(jet.INTERVAL(20, jet.MINUTE))),
			),
		).
		ORDER_BY(d.Time.DESC()).
		LIMIT(50)

	var dest []*model.GksphoneJobMessage
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	ds := make([]*livemap.DispatchMarker, len(dest))

	for i, v := range dest {
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

		ds[i] = &livemap.DispatchMarker{
			X:         float32(x),
			Y:         float32(y),
			Id:        v.ID,
			Icon:      icon,
			IconColor: iconColor,
			Name:      *v.Name,
			Popup:     *v.Message,
		}
	}

	return ds, nil
}

func (s *Server) GenerateRandomUserMarker() {
	userIds := []int32{
		26061,
		4650,
		29225,
		931,
		6173,
		16235,
	}

	markers := make([]*model.ArpanetUserLocations, len(userIds))

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
			markers[i] = &model.ArpanetUserLocations{
				UserID: userIds[i],
				Job:    &job,
				Hidden: &hidden,

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
				locs.UserID,
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
		time.Sleep(2 * time.Second)
	}
}
