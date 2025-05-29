package users

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/timestamp"
	"github.com/fivenet-app/fivenet/v2025/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2025/pkg/utils"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/model"
	"github.com/fivenet-app/fivenet/v2025/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/proto"
)

func GetUserProps(ctx context.Context, tx qrm.DB, userId int32, attrJobs []string) (*UserProps, error) {
	tUserProps := table.FivenetUserProps.AS("user_props")
	stmt := tUserProps.
		SELECT(
			tUserProps.UserID,
			tUserProps.UpdatedAt,
			tUserProps.Wanted,
			tUserProps.Job,
			tUserProps.JobGrade,
			tUserProps.TrafficInfractionPoints,
			tUserProps.TrafficInfractionPointsUpdatedAt,
			tUserProps.OpenFines,
			tUserProps.Mugshot,
			tUserProps.MugshotFileID,
		).
		FROM(tUserProps).
		WHERE(
			tUserProps.UserID.EQ(jet.Int32(userId)),
		).
		LIMIT(1)

	dest := &UserProps{
		UserId: userId,
	}
	if err := stmt.QueryContext(ctx, tx, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	dest.UserId = userId

	attributes, err := GetUserLabels(ctx, tx, userId, attrJobs)
	if err != nil {
		return nil, err
	}
	dest.Labels = attributes

	return dest, nil
}

func GetUserLabels(ctx context.Context, tx qrm.DB, userId int32, jobs []string) (*Labels, error) {
	list := &Labels{
		List: []*Label{},
	}

	if len(jobs) == 0 {
		return list, nil
	}

	jobsExp := make([]jet.Expression, len(jobs))
	for i := range jobs {
		jobsExp[i] = jet.String(jobs[i])
	}

	tCitizensLabelsJob := table.FivenetUserLabelsJob.AS("citizen_label")
	tUserLabels := table.FivenetUserLabels

	stmt := tUserLabels.
		SELECT(
			tCitizensLabelsJob.ID,
			tCitizensLabelsJob.Job,
			tCitizensLabelsJob.Name,
			tCitizensLabelsJob.Color,
		).
		FROM(
			tUserLabels.
				INNER_JOIN(tCitizensLabelsJob,
					tCitizensLabelsJob.ID.EQ(tUserLabels.LabelID),
				),
		).
		WHERE(jet.AND(
			tUserLabels.UserID.EQ(jet.Int32(userId)),
			tCitizensLabelsJob.Job.IN(jobsExp...),
		)).
		ORDER_BY(
			tCitizensLabelsJob.SortKey.ASC(),
		)

	if err := stmt.QueryContext(ctx, tx, &list.List); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return list, nil
}

func (x *UserProps) Default() {
	if x.Wanted == nil {
		v := false
		x.Wanted = &v
	}

	if x.TrafficInfractionPoints == nil {
		v := uint32(0)
		x.TrafficInfractionPoints = &v
	}

	if x.OpenFines == nil {
		v := int64(0)
		x.OpenFines = &v
	}

	if x.Labels == nil {
		x.Labels = &Labels{}
	}
}

func (x *UserProps) HandleChanges(ctx context.Context, tx qrm.DB, in *UserProps, sourceUserId *int32, reason string) ([]*UserActivity, error) {
	x.Default()

	tUserProps := table.FivenetUserProps

	updateSets := []jet.ColumnAssigment{}

	// Generate the update sets
	if in.Wanted != nil {
		updateSets = append(updateSets, tUserProps.Wanted.SET(jet.Bool(*in.Wanted)))
	} else {
		in.Wanted = x.Wanted
	}

	if in.JobName != nil {
		if in.JobGradeNumber == nil {
			grade := int32(0)
			in.JobGradeNumber = &grade
		}

		if in.Job == nil || in.JobGrade == nil {
			return nil, errors.New("invalid job")
		}

		updateSets = append(updateSets,
			tUserProps.Job.SET(jet.String(*in.JobName)),
			tUserProps.JobGrade.SET(jet.Int32(*in.JobGradeNumber)),
		)
	} else {
		in.JobName = x.JobName
		in.Job = x.Job
		in.JobGradeNumber = x.JobGradeNumber
		in.JobGrade = x.JobGrade
	}

	if in.TrafficInfractionPoints != nil {
		updateSets = append(updateSets, tUserProps.TrafficInfractionPoints.SET(jet.Uint32(*in.TrafficInfractionPoints)))

		// Update the timestamp if points are added
		if *in.TrafficInfractionPoints > 0 {
			in.TrafficInfractionPointsUpdatedAt = timestamp.Now()
		} else {
			// Reset the timestamp if points are "reset" (0)
			in.TrafficInfractionPointsUpdatedAt = nil
		}
	} else {
		in.TrafficInfractionPoints = x.TrafficInfractionPoints
		in.TrafficInfractionPointsUpdatedAt = x.TrafficInfractionPointsUpdatedAt
	}

	if in.OpenFines != nil {
		updateSets = append(updateSets, tUserProps.OpenFines.SET(
			jet.IntExp(jet.Raw(
				"CASE WHEN COALESCE(`open_fines`, 0) + $amount < 0 THEN 0 ELSE COALESCE(`open_fines`, 0) + $amount END",
				jet.RawArgs{
					"$amount": *in.OpenFines,
				},
			)),
		))
	} else {
		in.OpenFines = x.OpenFines
	}

	if in.Mugshot != nil {
		updateSets = append(updateSets, tUserProps.Mugshot.SET(jet.StringExp(jet.Raw("VALUES(`mug_shot`)"))))
	} else {
		in.Mugshot = x.Mugshot
	}

	if in.Labels != nil {
		if in.Labels.List == nil {
			in.Labels.List = []*Label{}
		}

		slices.SortFunc(in.Labels.List, func(a, b *Label) int {
			return strings.Compare(a.Name, b.Name)
		})
	} else {
		in.Labels = x.Labels
	}

	if len(updateSets) > 0 {
		stmt := tUserProps.
			INSERT(
				tUserProps.UserID,
				tUserProps.Wanted,
				tUserProps.Job,
				tUserProps.JobGrade,
				tUserProps.TrafficInfractionPoints,
				tUserProps.TrafficInfractionPointsUpdatedAt,
				tUserProps.OpenFines,
				tUserProps.Mugshot,
			).
			VALUES(
				in.UserId,
				in.Wanted,
				in.JobName,
				in.JobGradeNumber,
				in.TrafficInfractionPoints,
				in.TrafficInfractionPointsUpdatedAt,
				in.OpenFines,
				in.Mugshot,
			).
			ON_DUPLICATE_KEY_UPDATE(
				updateSets...,
			)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return nil, err
		}
	}

	activities := []*UserActivity{}

	// Create user activity entries

	if x.Wanted != in.Wanted && (x.Wanted == nil || in.Wanted == nil || *x.Wanted != *in.Wanted) {
		var wanted bool
		if in.Wanted != nil {
			wanted = *in.Wanted
		}

		activities = append(activities, &UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			Type:         UserActivityType_USER_ACTIVITY_TYPE_WANTED,
			Reason:       reason,
			Data: &UserActivityData{
				Data: &UserActivityData_WantedChange{
					WantedChange: &WantedChange{
						Wanted: wanted,
					},
				},
			},
		})
	}
	if (x.JobName != in.JobName && (x.JobName == nil || in.JobName == nil || *x.JobName != *in.JobName)) || (x.JobGradeNumber != in.JobGradeNumber &&
		(x.JobGradeNumber == nil || in.JobGradeNumber == nil || *x.JobGradeNumber != *in.JobGradeNumber)) {
		var jobLabel *string
		if in.Job != nil {
			jobLabel = &in.Job.Label
		}

		var gradeLabel *string
		if in.Job != nil {
			gradeLabel = &in.JobGrade.Label
		}

		activities = append(activities, &UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			Type:         UserActivityType_USER_ACTIVITY_TYPE_JOB,
			Reason:       reason,
			Data: &UserActivityData{
				Data: &UserActivityData_JobChange{
					JobChange: &JobChange{
						Job:        jobLabel,
						JobLabel:   in.JobName,
						Grade:      in.JobGradeNumber,
						GradeLabel: gradeLabel,
					},
				},
			},
		})
	}
	if x.TrafficInfractionPoints != in.TrafficInfractionPoints && (x.TrafficInfractionPoints == nil || in.TrafficInfractionPoints == nil ||
		*x.TrafficInfractionPoints != *in.TrafficInfractionPoints) {
		old := uint32(0)
		if x.TrafficInfractionPoints != nil {
			old = *x.TrafficInfractionPoints
		}
		new := uint32(0)
		if in.TrafficInfractionPoints != nil {
			new = *in.TrafficInfractionPoints
		}

		activities = append(activities, &UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			Type:         UserActivityType_USER_ACTIVITY_TYPE_TRAFFIC_INFRACTION_POINTS,
			Reason:       reason,
			Data: &UserActivityData{
				Data: &UserActivityData_TrafficInfractionPointsChange{
					TrafficInfractionPointsChange: &TrafficInfractionPointsChange{
						Old: old,
						New: new,
					},
				},
			},
		})
	}
	if x.MugshotFileId != in.MugshotFileId && (x.MugshotFileId == nil || in.MugshotFileId == nil || x.MugshotFileId != in.MugshotFileId) {
		activities = append(activities, &UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			Type:         UserActivityType_USER_ACTIVITY_TYPE_MUGSHOT,
			Reason:       reason,
			Data: &UserActivityData{
				Data: &UserActivityData_MugshotChange{
					MugshotChange: &MugshotChange{},
				},
			},
		})
	}
	if x.Labels != in.Labels && !proto.Equal(in.Labels, x.Labels) {
		if in.Labels == nil {
			in.Labels = &Labels{}
		}

		added, removed := utils.SlicesDifferenceFunc(x.Labels.List, in.Labels.List,
			func(in *Label) uint64 {
				return in.Id
			})

		if err := x.updateLabels(ctx, tx, in.UserId, added, removed); err != nil {
			return nil, err
		}

		activities = append(activities, &UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			Type:         UserActivityType_USER_ACTIVITY_TYPE_LABELS,
			Reason:       reason,
			Data: &UserActivityData{
				Data: &UserActivityData_LabelsChange{
					LabelsChange: &LabelsChange{
						Added:   added,
						Removed: removed,
					},
				},
			},
		})
	}

	return activities, nil
}

func (s *UserProps) updateLabels(ctx context.Context, tx qrm.DB, userId int32, added []*Label, removed []*Label) error {
	tUserLabels := table.FivenetUserLabels

	if len(added) > 0 {
		addedLabels := make([]*model.FivenetUserLabels, len(added))
		for i, label := range added {
			addedLabels[i] = &model.FivenetUserLabels{
				UserID:  userId,
				LabelID: label.Id,
			}
		}

		stmt := tUserLabels.
			INSERT(
				tUserLabels.UserID,
				tUserLabels.LabelID,
			).
			MODELS(addedLabels)

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return err
			}
		}
	}

	if len(removed) > 0 {
		ids := make([]jet.Expression, len(removed))

		for i := range removed {
			ids[i] = jet.Uint64(removed[i].Id)
		}

		stmt := tUserLabels.
			DELETE().
			WHERE(jet.AND(
				tUserLabels.UserID.EQ(jet.Int32(userId)),
				tUserLabels.LabelID.IN(ids...),
			)).
			LIMIT(int64(len(removed)))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}
