package tracker

import (
	"math/rand"
	"time"

	"github.com/galexrt/fivenet/query/fivenet/model"
	jet "github.com/go-jet/jet/v2/mysql"
	"go.uber.org/zap"
)

func (s *Manager) randomizeUserMarkers() {
	userIdentifiers := []string{
		"char1:1419c5a584086a76dbddae6aed1051157596a02a",
		"char1:2db159960d3cbf1a9598670d9f75bbe11ecf28e8",
		"char1:60368551546ee3b8c6d6c3d6830865ce664fac79",
		"char1:8369b091d49b5aa2fe061c654d4624f5ce5b89f1",
		"char1:c4c6e6f876a2417b3060025773fb193b24666e8b",
		"char1:c834a0f52b6d5465643ad9c8c689d2b5abe9801f",
		"char1:d4b145bb77a128e66c0cd96618a9950236cf0a70",
		"char1:eef505578bf3ee13dff3446a36e88439db2e5f5d",
		"char2:1419c5a584086a76dbddae6aed1051157596a02a",
		"char2:da4f9eacf69feb5eb5bdfad0de1e3c2dc9dee335",
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
					tLocs.Hidden.SET(jet.RawBool("VALUES(`hidden`)")),
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
