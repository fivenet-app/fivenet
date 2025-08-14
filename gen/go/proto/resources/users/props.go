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

func GetUserProps(
	ctx context.Context,
	tx qrm.DB,
	userId int32,
	attrJobs []string,
) (*UserProps, error) {
	tUserProps := table.FivenetUserProps.AS("user_props")
	tFiles := table.FivenetFiles.AS("mugshot")

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
			tUserProps.MugshotFileID,
			tFiles.ID.AS("mugshot.mugshot_file_id"),
			tFiles.FilePath,
		).
		FROM(
			tUserProps.
				LEFT_JOIN(tFiles,
					tFiles.ID.EQ(tUserProps.MugshotFileID),
				),
		).
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

	if x.GetLabels() == nil {
		x.Labels = &Labels{}
	}
}

func (x *UserProps) HandleChanges(
	ctx context.Context,
	tx qrm.DB,
	in *UserProps,
	sourceUserId *int32,
	reason string,
) ([]*UserActivity, error) {
	x.Default()

	tUserProps := table.FivenetUserProps

	updateSets := []jet.ColumnAssigment{}

	// Generate the update sets
	if in.Wanted != nil {
		updateSets = append(updateSets, tUserProps.Wanted.SET(jet.Bool(in.GetWanted())))
	} else {
		in.Wanted = x.Wanted
	}

	if in.JobName != nil {
		if in.JobGradeNumber == nil {
			grade := int32(0)
			in.JobGradeNumber = &grade
		}

		if in.GetJob() == nil || in.GetJobGrade() == nil {
			return nil, errors.New("invalid job")
		}

		updateSets = append(updateSets,
			tUserProps.Job.SET(jet.String(in.GetJobName())),
			tUserProps.JobGrade.SET(jet.Int32(in.GetJobGradeNumber())),
		)
	} else {
		in.JobName = x.JobName
		in.Job = x.GetJob()
		in.JobGradeNumber = x.JobGradeNumber
		in.JobGrade = x.GetJobGrade()
	}

	if in.TrafficInfractionPoints != nil {
		updateSets = append(
			updateSets,
			tUserProps.TrafficInfractionPoints.SET(jet.Uint32(in.GetTrafficInfractionPoints())),
		)

		// Update the timestamp if points are added
		if in.GetTrafficInfractionPoints() > 0 {
			in.TrafficInfractionPointsUpdatedAt = timestamp.Now()
		} else {
			// Reset the timestamp if points are "reset" (0)
			in.TrafficInfractionPointsUpdatedAt = nil
		}
	} else {
		in.TrafficInfractionPoints = x.TrafficInfractionPoints
		in.TrafficInfractionPointsUpdatedAt = x.GetTrafficInfractionPointsUpdatedAt()
	}

	if in.OpenFines != nil {
		updateSets = append(updateSets, tUserProps.OpenFines.SET(
			jet.IntExp(jet.Raw(
				"CASE WHEN COALESCE(`open_fines`, 0) + $amount < 0 THEN 0 ELSE COALESCE(`open_fines`, 0) + $amount END",
				jet.RawArgs{
					"$amount": in.GetOpenFines(),
				},
			)),
		))
	} else {
		in.OpenFines = x.OpenFines
	}

	if in.MugshotFileId != nil {
		updateSets = append(
			updateSets,
			tUserProps.MugShot.SET(jet.StringExp(jet.Raw("VALUES(`mug_shot`)"))),
		)
	} else {
		in.MugshotFileId = x.MugshotFileId
	}

	if in.GetLabels() != nil {
		if in.Labels.List == nil {
			in.Labels.List = []*Label{}
		}

		slices.SortFunc(in.GetLabels().GetList(), func(a, b *Label) int {
			return strings.Compare(a.GetName(), b.GetName())
		})
	} else {
		in.Labels = x.GetLabels()
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
				tUserProps.MugshotFileID,
			).
			VALUES(
				in.GetUserId(),
				in.Wanted,
				in.JobName,
				in.JobGradeNumber,
				in.TrafficInfractionPoints,
				in.TrafficInfractionPointsUpdatedAt,
				in.OpenFines,
				in.MugshotFileId,
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

	if x.GetWanted() != in.GetWanted() &&
		(x.Wanted == nil || in.Wanted == nil || x.GetWanted() != in.GetWanted()) {
		var wanted bool
		if in.Wanted != nil {
			wanted = in.GetWanted()
		}

		activities = append(activities, &UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.GetUserId(),
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
	if (x.GetJobName() != in.GetJobName() && (x.JobName == nil || in.JobName == nil || x.GetJobName() != in.GetJobName())) ||
		(x.GetJobGradeNumber() != in.GetJobGradeNumber() &&
			(x.JobGradeNumber == nil || in.JobGradeNumber == nil || x.GetJobGradeNumber() != in.GetJobGradeNumber())) {
		var jobLabel *string
		if in.GetJob() != nil {
			jobLabel = &in.Job.Label
		}

		var gradeLabel *string
		if in.GetJob() != nil {
			gradeLabel = &in.JobGrade.Label
		}

		activities = append(activities, &UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.GetUserId(),
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
	if x.GetTrafficInfractionPoints() != in.GetTrafficInfractionPoints() &&
		(x.TrafficInfractionPoints == nil || in.TrafficInfractionPoints == nil ||
			x.GetTrafficInfractionPoints() != in.GetTrafficInfractionPoints()) {
		old := uint32(0)
		if x.TrafficInfractionPoints != nil {
			old = x.GetTrafficInfractionPoints()
		}
		new := uint32(0)
		if in.TrafficInfractionPoints != nil {
			new = in.GetTrafficInfractionPoints()
		}

		activities = append(activities, &UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.GetUserId(),
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
	if x.GetMugshotFileId() != in.GetMugshotFileId() &&
		(x.MugshotFileId == nil || in.MugshotFileId == nil || x.GetMugshotFileId() != in.GetMugshotFileId()) {
		activities = append(activities, &UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.GetUserId(),
			Type:         UserActivityType_USER_ACTIVITY_TYPE_MUGSHOT,
			Reason:       reason,
			Data: &UserActivityData{
				Data: &UserActivityData_MugshotChange{
					MugshotChange: &MugshotChange{},
				},
			},
		})
	}
	if x.GetLabels() != in.GetLabels() && !proto.Equal(in.GetLabels(), x.GetLabels()) {
		if in.GetLabels() == nil {
			in.Labels = &Labels{}
		}

		added, removed := utils.SlicesDifferenceFunc(
			x.GetLabels().GetList(),
			in.GetLabels().GetList(),
			func(in *Label) uint64 {
				return in.GetId()
			},
		)

		if err := x.updateLabels(ctx, tx, in.GetUserId(), added, removed); err != nil {
			return nil, err
		}

		activities = append(activities, &UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.GetUserId(),
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

func (s *UserProps) updateLabels(
	ctx context.Context,
	tx qrm.DB,
	userId int32,
	added []*Label,
	removed []*Label,
) error {
	tUserLabels := table.FivenetUserLabels

	if len(added) > 0 {
		addedLabels := make([]*model.FivenetUserLabels, len(added))
		for i, label := range added {
			addedLabels[i] = &model.FivenetUserLabels{
				UserID:  userId,
				LabelID: label.GetId(),
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
			ids[i] = jet.Uint64(removed[i].GetId())
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
