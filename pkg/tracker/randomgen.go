package tracker

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/galexrt/fivenet/query/fivenet/model"
	"github.com/galexrt/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

func (s *Tracker) GenerateRandomUserMarker() {
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
				Job:        job,
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
		func() {
			ctx, span := s.tracer.Start(s.ctx, "livemap-gen-users")
			defer span.End()

			if counter >= 15 {
				resetMarkers()
				counter = 0
			} else {
				moveMarkers()
			}

			stmt := tLocs.
				INSERT(
					tLocs.Identifier,
					tLocs.Job,
					tLocs.X,
					tLocs.Y,
					tLocs.Hidden,
				).
				MODELS(markers).
				ON_DUPLICATE_KEY_UPDATE(
					tLocs.X.SET(jet.RawFloat("VALUES(`x`)")),
					tLocs.Y.SET(jet.RawFloat("VALUES(`y`)")),
				)

			_, err := stmt.ExecContext(ctx, s.db)
			if err != nil {
				s.logger.Error("failed to insert/ update random location to locations table", zap.Error(err))
			}

			counter++
			time.Sleep(3 * time.Second)
		}()
	}
}

func (s *Tracker) GenerateRandomDispatchMarker() {
	userIdentifiers := []int32{
		26061,
		43612,
		41857,
		35061,
	}

	markers := make([]*model.FivenetCentrumDispatches, len(userIdentifiers))

	jMessage := table.GksphoneJobMessage
	resetMarkers := func() {
		ctx, span := s.tracer.Start(s.ctx, "livemap-gen-dispatches")
		defer span.End()

		xMin := -3300
		xMax := 4300
		yMin := -3300
		yMax := 5000
		for i := 0; i < len(markers); i++ {
			x := float64(rand.Intn(xMax-xMin+1) + xMin)
			y := float64(rand.Intn(yMax-yMin+1) + yMin)

			message := fmt.Sprintf("TEST %d", i)
			description := fmt.Sprintf("Message of Dispatch %d", i)
			job := "ambulance"

			anon := false
			markers[i] = &model.FivenetCentrumDispatches{
				Job:         job,
				Message:     message,
				Description: &description,
				X:           &x,
				Y:           &y,
				Anon:        anon,
				UserID:      userIdentifiers[i],
			}
		}

		stmt := jMessage.
			INSERT(
				jMessage.Name,
				jMessage.Number,
				jMessage.Message,
				jMessage.Photo,
				jMessage.Gps,
				jMessage.Owner,
				jMessage.Jobm,
				jMessage.Time,
				jMessage.Anon,
			).
			MODELS(markers)

		_, err := stmt.ExecContext(ctx, s.db)
		if err != nil {
			s.logger.Error("failed to insert random dispatch", zap.Error(err))
		}
	}

	resetMarkers()

	counter := 0
	for {
		if counter >= 20 {
			resetMarkers()
			counter = 0
		}

		counter++
		time.Sleep(3 * time.Second)
	}
}
