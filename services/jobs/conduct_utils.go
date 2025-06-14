package jobs

import (
	"context"

	jobs "github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/jobs"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils/tables"
	jet "github.com/go-jet/jet/v2/mysql"
)

func (s *Server) getConductEntry(ctx context.Context, id uint64) (*jobs.ConductEntry, error) {
	tColleague := tables.User().AS("target_user")
	tCreator := tColleague.AS("creator")

	stmt := tConduct.
		SELECT(
			tConduct.ID,
			tConduct.CreatedAt,
			tConduct.UpdatedAt,
			tConduct.DeletedAt,
			tConduct.Job,
			tConduct.Type,
			tConduct.Message,
			tConduct.ExpiresAt,
			tConduct.TargetUserID,
			tColleague.ID,
			tColleague.Firstname,
			tColleague.Lastname,
			tColleague.Dateofbirth,
			tColleague.PhoneNumber,
			tConduct.CreatorID,
			tCreator.ID,
			tCreator.Firstname,
			tCreator.Lastname,
			tCreator.Dateofbirth,
			tCreator.PhoneNumber,
		).
		FROM(
			tConduct.
				INNER_JOIN(tColleague,
					tColleague.ID.EQ(tConduct.TargetUserID),
				).
				LEFT_JOIN(tCreator,
					tCreator.ID.EQ(tConduct.CreatorID),
				),
		).
		WHERE(tConduct.ID.EQ(jet.Uint64(id))).
		LIMIT(1)

	dest := &jobs.ConductEntry{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		return nil, err
	}

	return dest, nil
}
