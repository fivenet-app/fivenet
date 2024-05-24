package jobs

import (
	"context"

	jobs "github.com/fivenet-app/fivenet/gen/go/proto/resources/jobs"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (s *Server) getConductEntry(ctx context.Context, id uint64) (*jobs.ConductEntry, error) {
	tUser := tUser.AS("target_user")
	tCreator := tUser.AS("creator")
	stmt := tConduct.
		SELECT(
			tConduct.ID,
			tConduct.CreatedAt,
			tConduct.UpdatedAt,
			tConduct.Job,
			tConduct.Type,
			tConduct.Message,
			tConduct.ExpiresAt,
			tConduct.TargetUserID,
			tUser.ID,
			tUser.Firstname,
			tUser.Lastname,
			tUser.Dateofbirth,
			tUser.PhoneNumber,
			tConduct.CreatorID,
			tCreator.ID,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
		).
		FROM(
			tConduct.
				INNER_JOIN(tUser,
					tUser.ID.EQ(tConduct.TargetUserID),
				).
				LEFT_JOIN(tCreator,
					tCreator.ID.EQ(tConduct.CreatorID),
				),
		).
		WHERE(tConduct.ID.EQ(jet.Uint64(id))).
		LIMIT(1)

	var dest jobs.ConductEntry
	if err := stmt.QueryContext(ctx, s.db, &dest); err != nil {
		return nil, err
	}

	return &dest, nil
}
