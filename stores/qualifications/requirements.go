package qualificationsstore

import (
	"context"
	"slices"

	resqualifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/qualifications"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
)

func (s *Store) HandleQualificationRequirementsChanges(
	ctx context.Context,
	tx qrm.DB,
	qualificationId int64,
	reqs []*resqualifications.QualificationRequirement,
) error {
	current, err := s.GetQualificationRequirements(ctx, qualificationId)
	if err != nil {
		return err
	}

	toCreate, toDelete := compareQualificationRequirements(current, reqs)
	tQReqs := table.FivenetQualificationsRequirements

	for _, req := range toDelete {
		stmt := tQReqs.
			DELETE().
			WHERE(mysql.AND(tQReqs.ID.EQ(mysql.Int64(req.GetId())))).
			LIMIT(1)
		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	for _, req := range toCreate {
		stmt := tQReqs.
			INSERT(
				tQReqs.QualificationID,
				tQReqs.TargetQualificationID,
			).
			VALUES(
				qualificationId,
				req.GetTargetQualificationId(),
			)
		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}

func compareQualificationRequirements(
	current, in []*resqualifications.QualificationRequirement,
) ([]*resqualifications.QualificationRequirement, []*resqualifications.QualificationRequirement) {
	toCreate := []*resqualifications.QualificationRequirement{}
	toDelete := []*resqualifications.QualificationRequirement{}

	if current == nil {
		return in, toDelete
	}

	if len(current) == 0 {
		if len(in) == 0 {
			toDelete = current
		} else {
			toCreate = in
		}
	} else {
		foundTracker := []int{}
		for _, cq := range current {
			var found *resqualifications.QualificationRequirement
			var foundIdx int
			for i, qj := range in {
				if cq.GetTargetQualificationId() != qj.GetTargetQualificationId() {
					continue
				}
				found = qj
				foundIdx = i
				break
			}
			if found == nil {
				toDelete = append(toDelete, cq)
				continue
			}

			foundTracker = append(foundTracker, foundIdx)
		}

		for i, uj := range in {
			if idx := slices.Index(foundTracker, i); idx == -1 {
				toCreate = append(toCreate, uj)
			}
		}
	}

	return toCreate, toDelete
}
