package citizensstore

import (
	"context"
	"errors"
	"slices"
	"strings"

	citizenslabels "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/citizens/labels"
	"github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/timestamp"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	"github.com/fivenet-app/fivenet/v2026/pkg/dbutils"
	"github.com/fivenet-app/fivenet/v2026/pkg/utils"
	"github.com/fivenet-app/fivenet/v2026/query/fivenet/table"
	"github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/proto"
)

func (s *Store) GetUserProps(
	ctx context.Context,
	tx qrm.DB,
	userId int32,
) (*usersprops.UserProps, error) {
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
			tUserProps.AvatarFileID,
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
			tUserProps.UserID.EQ(mysql.Int32(userId)),
		).
		LIMIT(1)

	dest := &usersprops.UserProps{UserId: userId}
	if err := stmt.QueryContext(ctx, tx, dest); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	dest.UserId = userId

	labels, err := s.GetUserLabels(
		ctx,
		tx,
		mysql.AND(tCitizenLabels.UserID.EQ(mysql.Int32(userId))),
		nil,
	)
	if err != nil {
		return nil, err
	}
	dest.Labels = labels

	return dest, nil
}

func (s *Store) HandleUserPropsChanges(
	ctx context.Context,
	tx qrm.DB,
	x *usersprops.UserProps,
	in *usersprops.UserProps,
	sourceUserId *int32,
	reason string,
) ([]*usersactivity.UserActivity, error) {
	x.Default()

	tUserProps := table.FivenetUserProps

	updateSets := []mysql.ColumnAssigment{}

	if in.Wanted != nil {
		updateSets = append(updateSets, tUserProps.Wanted.SET(mysql.Bool(in.GetWanted())))
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
			tUserProps.Job.SET(mysql.String(in.GetJobName())),
			tUserProps.JobGrade.SET(mysql.Int32(in.GetJobGradeNumber())),
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
			tUserProps.TrafficInfractionPoints.SET(mysql.Uint32(in.GetTrafficInfractionPoints())),
		)

		if in.GetTrafficInfractionPoints() > 0 {
			in.TrafficInfractionPointsUpdatedAt = timestamp.Now()
		} else {
			in.TrafficInfractionPointsUpdatedAt = nil
		}
	} else {
		in.TrafficInfractionPoints = x.TrafficInfractionPoints
		in.TrafficInfractionPointsUpdatedAt = x.GetTrafficInfractionPointsUpdatedAt()
	}

	if in.OpenFines != nil {
		updateSets = append(updateSets, tUserProps.OpenFines.SET(
			mysql.IntExp(mysql.Raw(
				"CASE WHEN COALESCE(`open_fines`, 0) + $amount < 0 THEN 0 ELSE COALESCE(`open_fines`, 0) + $amount END",
				mysql.RawArgs{
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
			tUserProps.MugShot.SET(mysql.StringExp(mysql.Raw("VALUES(`mug_shot`)"))),
		)
	} else {
		in.MugshotFileId = x.MugshotFileId
	}

	if in.GetLabels() != nil {
		if in.Labels.List == nil {
			in.Labels.List = []*citizenslabels.Label{}
		}

		slices.SortFunc(in.GetLabels().GetList(), func(a, b *citizenslabels.Label) int {
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
				in.GetTrafficInfractionPointsUpdatedAt(),
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

	activities := []*usersactivity.UserActivity{}

	if x.GetWanted() != in.GetWanted() &&
		(x.Wanted == nil || in.Wanted == nil || x.GetWanted() != in.GetWanted()) {
		var wanted bool
		if in.Wanted != nil {
			wanted = in.GetWanted()
		}

		activities = append(activities, &usersactivity.UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.GetUserId(),
			Type:         usersactivity.UserActivityType_USER_ACTIVITY_TYPE_WANTED,
			Reason:       reason,
			Data: &usersactivity.UserActivityData{
				Data: &usersactivity.UserActivityData_WantedChange{
					WantedChange: &usersactivity.WantedChange{Wanted: wanted},
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

		activities = append(activities, &usersactivity.UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.GetUserId(),
			Type:         usersactivity.UserActivityType_USER_ACTIVITY_TYPE_JOB,
			Reason:       reason,
			Data: &usersactivity.UserActivityData{
				Data: &usersactivity.UserActivityData_JobChange{
					JobChange: &usersactivity.JobChange{
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

		activities = append(activities, &usersactivity.UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.GetUserId(),
			Type:         usersactivity.UserActivityType_USER_ACTIVITY_TYPE_TRAFFIC_INFRACTION_POINTS,
			Reason:       reason,
			Data: &usersactivity.UserActivityData{
				Data: &usersactivity.UserActivityData_TrafficInfractionPointsChange{
					TrafficInfractionPointsChange: &usersactivity.TrafficInfractionPointsChange{
						Old: old,
						New: new,
					},
				},
			},
		})
	}
	if x.GetMugshotFileId() != in.GetMugshotFileId() &&
		(x.MugshotFileId == nil || in.MugshotFileId == nil || x.GetMugshotFileId() != in.GetMugshotFileId()) {
		activities = append(activities, &usersactivity.UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.GetUserId(),
			Type:         usersactivity.UserActivityType_USER_ACTIVITY_TYPE_MUGSHOT,
			Reason:       reason,
			Data: &usersactivity.UserActivityData{
				Data: &usersactivity.UserActivityData_MugshotChange{
					MugshotChange: &usersactivity.MugshotChange{},
				},
			},
		})
	}
	if x.GetLabels() != in.GetLabels() && !proto.Equal(in.GetLabels(), x.GetLabels()) {
		if in.GetLabels() == nil {
			in.Labels = &citizenslabels.Labels{}
		}

		added, removed := utils.SlicesDifferenceFunc(
			x.GetLabels().GetList(),
			in.GetLabels().GetList(),
			func(in *citizenslabels.Label) int64 { return in.GetId() },
		)

		if err := s.updateUserLabels(ctx, tx, in.GetUserId(), added, removed); err != nil {
			return nil, err
		}

		activities = append(activities, &usersactivity.UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.GetUserId(),
			Type:         usersactivity.UserActivityType_USER_ACTIVITY_TYPE_LABELS,
			Reason:       reason,
			Data: &usersactivity.UserActivityData{
				Data: &usersactivity.UserActivityData_LabelsChange{
					LabelsChange: &usersactivity.LabelsChange{
						Added:   added,
						Removed: removed,
					},
				},
			},
		})
	}

	return activities, nil
}

func (s *Store) updateUserLabels(
	ctx context.Context,
	tx qrm.DB,
	userId int32,
	added []*citizenslabels.Label,
	removed []*citizenslabels.Label,
) error {
	tUserLabels := table.FivenetUserLabels
	if len(added) > 0 {
		stmt := tUserLabels.
			INSERT(
				tUserLabels.UserID,
				tUserLabels.LabelID,
				tUserLabels.ExpiresAt,
			)

		for _, label := range added {
			stmt = stmt.VALUES(userId, label.GetId(), label.GetExpiresAt())
		}

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			if !dbutils.IsDuplicateError(err) {
				return err
			}
		}
	}

	if len(removed) > 0 {
		ids := make([]mysql.Expression, len(removed))
		for i := range removed {
			ids[i] = mysql.Int64(removed[i].GetId())
		}

		stmt := tUserLabels.
			DELETE().
			WHERE(mysql.AND(
				tUserLabels.UserID.EQ(mysql.Int32(userId)),
				tUserLabels.LabelID.IN(ids...),
			)).
			LIMIT(int64(len(removed)))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}
