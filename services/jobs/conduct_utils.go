package jobs

import (
	"context"

	jobsconduct "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/jobs/conduct"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
)

func (s *Server) getConductEntry(
	ctx context.Context,
	id int64,
	withFiles bool,
) (*jobsconduct.ConductEntry, error) {
	tColleague := table.FivenetUser.AS("target_user")
	tCreator := tColleague.AS("creator")

	stmt := tConduct.
		SELECT(
			tConduct.ID,
			tConduct.CreatedAt,
			tConduct.UpdatedAt,
			tConduct.DeletedAt,
			tConduct.Job,
			tConduct.Type,
			tConduct.Draft,
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
				LEFT_JOIN(tColleague,
					tColleague.ID.EQ(tConduct.TargetUserID),
				).
				LEFT_JOIN(tCreator,
					tCreator.ID.EQ(tConduct.CreatorID),
				),
		).
		WHERE(tConduct.ID.EQ(mysql.Int64(id))).
		LIMIT(1)

	dest := &jobsconduct.ConductEntry{}
	if err := stmt.QueryContext(ctx, s.db, dest); err != nil {
		return nil, err
	}

	if withFiles {
		files, err := s.fHandler.ListFilesForParentID(ctx, id)
		if err != nil {
			return nil, err
		}
		dest.Files = files
	}

	return dest, nil
}
