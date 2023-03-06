package livemap

import (
	"math/rand"
	"time"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/permissions"
	"github.com/galexrt/arpanet/query"
	"github.com/galexrt/arpanet/query/arpanet/model"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
)

func init() {
	permissions.RegisterPerms([]*permissions.Perm{
		{Key: "livemap", Name: "Stream", PerJob: true},
	})
}

type Server struct {
	LivemapServiceServer
}

func NewServer() *Server {
	s := &Server{}
	go s.generateRandomMarker()
	return s
}

func (s *Server) Stream(req *StreamRequest, srv LivemapService_StreamServer) error {
	user, err := auth.GetUserFromContext(srv.Context())
	if err != nil {
		return err
	}

	// TODO Establish which jobs the user can select

	l := table.ArpanetUserLocations
	stmt := l.SELECT(l.AllColumns).
		FROM(l).
		WHERE(
			l.Job.EQ(jet.String(user.Job)).
				AND(l.Hidden.IS_FALSE()),
		)

	for {
		resp := &ServerStreamResponse{
			Dispatches: []*Marker{},
		}

		if err := stmt.QueryContext(srv.Context(), query.DB, &resp.Users); err != nil {
			return err
		}

		if err := srv.Send(resp); err != nil {
			return err
		}
		time.Sleep(3 * time.Second)
	}
}

func (s *Server) generateRandomMarker() {
	userIDs := []int32{
		// ambulance
		26061,
		4650,
		29225,
		29205,
		931,
		16235,
		6173,
		20232,
		1634,
		17800,
		27434,
		15706,
		3046,
	}

	l := table.ArpanetUserLocations
	for {
		markers := make([]*model.ArpanetUserLocations, len(userIDs))
		for i := 0; i < len(userIDs); i++ {
			xMin := -3500
			xMax := 4300
			x := float64(rand.Intn(xMax-xMin+1) + xMin)

			yMin := -3600
			yMax := 7800
			y := float64(rand.Intn(yMax-yMin+1) + yMin)

			job := "ambulance"
			hidden := false
			markers[i] = &model.ArpanetUserLocations{
				UserID: userIDs[i],
				Job:    &job,
				Hidden: &hidden,

				X: &x,
				Y: &y,
			}
		}

		stmt := l.INSERT().
			MODELS(markers).
			ON_DUPLICATE_KEY_UPDATE(
				l.X.SET(l.X.ADD(l.X)),
				l.Y.SET(l.Y.ADD(l.Y)),
			)
		_, err := stmt.Exec(query.DB)
		_ = err

		time.Sleep(2 * time.Second)
	}
}
