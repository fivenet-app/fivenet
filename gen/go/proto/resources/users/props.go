package users

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/fivenet-app/fivenet/pkg/dbutils"
	"github.com/fivenet-app/fivenet/pkg/utils"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/proto"
)

func GetUserProps(ctx context.Context, tx qrm.DB, userId int32, attrJobs []string) (*UserProps, error) {
	tUserProps := table.FivenetUserProps.AS("userprops")
	stmt := tUserProps.
		SELECT(
			tUserProps.UserID,
			tUserProps.UpdatedAt,
			tUserProps.Wanted,
			tUserProps.Job,
			tUserProps.JobGrade,
			tUserProps.TrafficInfractionPoints,
			tUserProps.OpenFines,
			tUserProps.MugShot,
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

func GetUserLabels(ctx context.Context, tx qrm.DB, userId int32, jobs []string) (*CitizenLabels, error) {
	list := &CitizenLabels{
		List: []*CitizenLabel{},
	}

	if len(jobs) == 0 {
		return list, nil
	}

	jobsExp := make([]jet.Expression, len(jobs))
	for i := 0; i < len(jobs); i++ {
		jobsExp[i] = jet.String(jobs[i])
	}

	tJobCitizenLabels := table.FivenetJobCitizenLabels.AS("citizen_label")
	tUserCitizenLabels := table.FivenetUserCitizenLabels

	stmt := tUserCitizenLabels.
		SELECT(
			tJobCitizenLabels.ID,
			tJobCitizenLabels.Job,
			tJobCitizenLabels.Name,
			tJobCitizenLabels.Color,
		).
		FROM(
			tUserCitizenLabels.
				INNER_JOIN(tJobCitizenLabels,
					tJobCitizenLabels.ID.EQ(tUserCitizenLabels.AttributeID),
				),
		).
		WHERE(jet.AND(
			tUserCitizenLabels.UserID.EQ(jet.Int32(userId)),
			tJobCitizenLabels.Job.IN(jobsExp...),
		)).
		ORDER_BY(
			tJobCitizenLabels.SortKey.ASC(),
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
		x.Labels = &CitizenLabels{}
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
	} else {
		in.TrafficInfractionPoints = x.TrafficInfractionPoints
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

	if in.MugShot != nil {
		updateSets = append(updateSets, tUserProps.MugShot.SET(jet.StringExp(jet.Raw("VALUES(`mug_shot`)"))))
	} else {
		in.MugShot = x.MugShot
	}

	if in.Labels != nil {
		if in.Labels.List == nil {
			in.Labels.List = []*CitizenLabel{}
		}

		slices.SortFunc(in.Labels.List, func(a, b *CitizenLabel) int {
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
				tUserProps.OpenFines,
				tUserProps.MugShot,
			).
			VALUES(
				in.UserId,
				in.Wanted,
				in.JobName,
				in.JobGradeNumber,
				in.TrafficInfractionPoints,
				in.OpenFines,
				in.MugShot,
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
					WantedChange: &UserWantedChange{
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
					JobChange: &UserJobChange{
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
					TrafficInfractionPointsChange: &UserTrafficInfractionPointsChange{
						Old: old,
						New: new,
					},
				},
			},
		})
	}
	if x.MugShot != in.MugShot && (x.MugShot == nil || in.MugShot == nil || x.MugShot.GetUrl() != in.MugShot.GetUrl()) {
		var url *string
		if in.MugShot != nil && in.MugShot.Url != nil {
			url = in.MugShot.Url
		}

		activities = append(activities, &UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			Type:         UserActivityType_USER_ACTIVITY_TYPE_MUGSHOT,
			Reason:       reason,
			Data: &UserActivityData{
				Data: &UserActivityData_MugshotChange{
					MugshotChange: &UserMugshotChange{
						New: url,
					},
				},
			},
		})
	}
	if x.Labels != in.Labels && !proto.Equal(in.Labels, x.Labels) {
		if in.Labels == nil {
			in.Labels = &CitizenLabels{}
		}

		added, removed := utils.SlicesDifferenceFunc(x.Labels.List, in.Labels.List,
			func(in *CitizenLabel) uint64 {
				return in.Id
			})

		if err := x.updateCitizenLabels(ctx, tx, in.UserId, added, removed); err != nil {
			return nil, err
		}

		activities = append(activities, &UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			Type:         UserActivityType_USER_ACTIVITY_TYPE_LABELS,
			Reason:       reason,
			Data: &UserActivityData{
				Data: &UserActivityData_LabelsChange{
					LabelsChange: &UserLabelsChange{
						Added:   added,
						Removed: removed,
					},
				},
			},
		})
	}

	return activities, nil
}

func (s *UserProps) updateCitizenLabels(ctx context.Context, tx qrm.DB, userId int32, added []*CitizenLabel, removed []*CitizenLabel) error {
	tUserCitizenLabels := table.FivenetUserCitizenLabels

	if len(added) > 0 {
		addedLabels := make([]*model.FivenetUserCitizenLabels, len(added))
		for i, attribute := range added {
			addedLabels[i] = &model.FivenetUserCitizenLabels{
				UserID:      userId,
				AttributeID: attribute.Id,
			}
		}

		stmt := tUserCitizenLabels.
			INSERT(
				tUserCitizenLabels.UserID,
				tUserCitizenLabels.AttributeID,
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

		stmt := tUserCitizenLabels.
			DELETE().
			WHERE(jet.AND(
				tUserCitizenLabels.UserID.EQ(jet.Int32(userId)),
				tUserCitizenLabels.AttributeID.IN(ids...),
			)).
			LIMIT(int64(len(removed)))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}
