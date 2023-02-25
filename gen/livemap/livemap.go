package livemap

import (
	"fmt"
	"math/rand"
	"time"

	"go.uber.org/zap"
)

type Server struct {
	LivemapServiceServer
	logger *zap.Logger
}

func NewServer(logger *zap.Logger) *Server {
	return &Server{
		logger: logger,
	}
}

func (s *Server) Stream(req *StreamRequest, srv LivemapService_StreamServer) error {
	for {
		randomMarkerCount := rand.Intn(25-1+1) + 1
		resp := StreamResponse{
			Users: make([]*Marker, randomMarkerCount),
		}

		for i := 0; i < randomMarkerCount; i++ {
			xMin := -3500
			xMax := 4300
			x := float32(rand.Intn(xMax-xMin+1) + xMin)

			yMin := -3600
			yMax := 7800
			y := float32(rand.Intn(yMax-yMin+1) + yMin)

			resp.Users = append(resp.Users, &Marker{
				Name: fmt.Sprintf("Test Marker %d", i),
				X:    x,
				Y:    y,
			})

			time.Sleep(2 * time.Second)
		}

		if err := srv.Send(&resp); err != nil {
			s.logger.Error("failed to send livemap stream", zap.Error(err))
			return err
		}
	}
}
