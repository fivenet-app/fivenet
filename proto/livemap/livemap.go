package livemap

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/permissions"
	"github.com/galexrt/arpanet/query"
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
	return &Server{}
}

func (s *Server) Stream(req *StreamRequest, srv LivemapService_StreamServer) error {
	user, err := auth.GetUserFromContext(srv.Context())
	if err != nil {
		return err
	}

	// TODO use our own location table in the future
	v := query.VpcL
	net := ""
	switch user.Job {
	case "ambulance":
		net = "medic"
	case "police":
		net = "cop"
	}

	q := v.Where(v.Net.Eq(net))
	_ = q
	for {
		resp := &ServerStreamResponse{
			Dispatches: []*Marker{},
		}

		// Start
		//locations, err := q.Find()
		//if err != nil {
		//	continue
		//}
		//resp.Users = make([]*Marker, len(locations))
		//for key, loc := range locations {
		//	x, _ := strconv.ParseFloat(loc.Coordsx, 32)
		//	y, _ := strconv.ParseFloat(loc.Coordsy, 32)
		//	resp.Users[key] = &Marker{
		//		Id:    loc.PlayerID,
		//		X:     float32(x),
		//		Y:     float32(y),
		//		Name:  loc.PlayerID,
		//		Popup: loc.PlayerID,
		//	}
		//}
		// or
		resp.Users = s.generateRandomMarker()

		if err := srv.Send(resp); err != nil {
			return err
		}
		time.Sleep(3 * time.Second)
	}
}

//lint:ignore U1000 used for testing
func (s *Server) generateRandomMarker() []*Marker {
	randomMarkerCount := rand.Intn(25) + 1
	markers := make([]*Marker, randomMarkerCount)

	for i := 0; i < randomMarkerCount; i++ {
		xMin := -3500
		xMax := 4300
		x := float32(rand.Intn(xMax-xMin+1) + xMin)

		yMin := -3600
		yMax := 7800
		y := float32(rand.Intn(yMax-yMin+1) + yMin)

		markers[i] = &Marker{
			Id:   fmt.Sprintf("%d", i),
			Name: fmt.Sprintf("Test Marker %d", i),
			X:    x,
			Y:    y,
		}
	}

	return markers
}
