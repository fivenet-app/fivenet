package livemapper

import (
	"math/rand"
	"time"

	"github.com/galexrt/arpanet/pkg/auth"
	"github.com/galexrt/arpanet/pkg/perms"
	"github.com/galexrt/arpanet/proto/resources/livemap"
	"github.com/galexrt/arpanet/query"
	"github.com/galexrt/arpanet/query/arpanet/model"
	"github.com/galexrt/arpanet/query/arpanet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

var (
	l = table.ArpanetUserLocations
)

type Server struct {
	LivemapperServiceServer

	logger *zap.Logger
	p      perms.Permissions
}

func NewServer(logger *zap.Logger, p perms.Permissions) *Server {
	return &Server{
		logger: logger,
		p:      p,
	}
}

func (s *Server) SetLogger(logger *zap.Logger) {
	s.logger = logger
}

func (s *Server) Stream(req *StreamRequest, srv LivemapperService_StreamServer) error {
	userID := auth.GetUserIDFromContext(srv.Context())

	jobs, err := s.p.GetSuffixOfPermissionsByPrefixOfUser(userID, "LivemapperService.Stream")
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

	l := l.AS("marker")
	stmt := l.SELECT(
		l.UserID,
		l.Job,
		l.X,
		l.Y,
		l.UpdatedAt,
	).
		FROM(l).
		WHERE(
			l.Job.IN(sqlJobs...).
				AND(
					l.Hidden.IS_FALSE(),
					// TODO Ignore all columns where `updated_at` older than 5 minutes
				),
		)

	for {
		resp := &ServerStreamResponse{
			Dispatches: []*livemap.UserMarker{},
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

func (s *Server) GenerateRandomUserMarker() {
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

		stmt := l.INSERT(
			l.UserID,
			l.Job,
			l.X,
			l.Y,
			l.Hidden,
		).
			MODELS(markers).
			ON_DUPLICATE_KEY_UPDATE(
				l.X.SET(jet.RawFloat("VALUES(x)")),
				l.Y.SET(jet.RawFloat("VALUES(y)")),
			)

		_, err := stmt.Exec(query.DB)
		if err != nil {
			s.logger.Error("failed to insert/ update random location to locations table", zap.Error(err))
		}

		time.Sleep(2 * time.Second)
	}
}
