package livemapper

import (
	"database/sql"
	"math/rand"
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
	stmt := locs.SELECT(
		locs.UserID,
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
		resp := &StreamResponse{
			Dispatches: []*livemap.DispatchMarker{},
		}

		if err := stmt.QueryContext(srv.Context(), s.db, &resp.Users); err != nil {
			return err
		}

		if err := srv.Send(resp); err != nil {
			return err
		}
		time.Sleep(3 * time.Second)
	}
}

func (s *Server) GenerateRandomUserMarker() {
	userIds := []int32{
		1,
		2,
		3,
		4,
		5,
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

		stmt := locs.INSERT(
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
