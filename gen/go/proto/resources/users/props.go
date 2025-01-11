package users

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/fivenet-app/fivenet/pkg/utils"
	"github.com/fivenet-app/fivenet/pkg/utils/dbutils"
	"github.com/fivenet-app/fivenet/query/fivenet/model"
	"github.com/fivenet-app/fivenet/query/fivenet/table"
	jet "github.com/go-jet/jet/v2/mysql"
	"github.com/go-jet/jet/v2/qrm"
	"google.golang.org/protobuf/encoding/protojson"
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
		))

	if err := stmt.QueryContext(ctx, tx, &list.List); err != nil {
		if !errors.Is(err, qrm.ErrNoRows) {
			return nil, err
		}
	}

	return list, nil
}

func (x *UserProps) HandleChanges(ctx context.Context, tx qrm.DB, in *UserProps, sourceUserId *int32, reason string) ([]*UserActivity, error) {
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
			grade := int32(1)
			in.JobGradeNumber = &grade
		}

		if in.Job == nil || in.JobGrade == nil {
			return nil, errors.New("invalid job")
		}

		updateSets = append(updateSets, tUserProps.Job.SET(jet.String(*in.JobName)))
		updateSets = append(updateSets, tUserProps.JobGrade.SET(jet.Int32(*in.JobGradeNumber)))
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
				tUserProps.MugShot,
			).
			VALUES(
				in.UserId,
				in.Wanted,
				in.JobName,
				in.JobGradeNumber,
				in.TrafficInfractionPoints,
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
	if *in.Wanted != *x.Wanted {
		activities = append(activities, &UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			Type:         UserActivityType_USER_ACTIVITY_TYPE_CHANGED,
			Key:          "UserProps.Wanted",
			OldValue:     strconv.FormatBool(*x.Wanted),
			NewValue:     strconv.FormatBool(*in.Wanted),
			Reason:       reason,
		})
	}
	if *in.JobName != *x.JobName || *in.JobGradeNumber != *x.JobGradeNumber {
		activities = append(activities, &UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			Type:         UserActivityType_USER_ACTIVITY_TYPE_CHANGED,
			Key:          "UserProps.Job",
			OldValue:     fmt.Sprintf("%s|%s", x.Job.Label, x.JobGrade.Label),
			NewValue:     fmt.Sprintf("%s|%s", in.Job.Label, in.JobGrade.Label),
			Reason:       reason,
		})
	}
	if *in.TrafficInfractionPoints != *x.TrafficInfractionPoints {
		activities = append(activities, &UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			Type:         UserActivityType_USER_ACTIVITY_TYPE_CHANGED,
			Key:          "UserProps.TrafficInfractionPoints",
			OldValue:     strconv.Itoa(int(*x.TrafficInfractionPoints)),
			NewValue:     strconv.Itoa(int(*in.TrafficInfractionPoints)),
			Reason:       reason,
		})
	}
	if in.MugShot != nil && (x.MugShot == nil || in.MugShot.Url != x.MugShot.Url) {
		previousUrl := ""
		if x.MugShot != nil && x.MugShot.Url != nil {
			previousUrl = *x.MugShot.Url
		}
		currentUrl := ""
		if in.MugShot != nil && in.MugShot.Url != nil {
			currentUrl = *in.MugShot.Url
		}

		activities = append(activities, &UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			Type:         UserActivityType_USER_ACTIVITY_TYPE_CHANGED,
			Key:          "UserProps.MugShot",
			OldValue:     previousUrl,
			NewValue:     currentUrl,
			Reason:       reason,
		})
	}
	if !proto.Equal(in.Labels, x.Labels) {
		added, removed := utils.SlicesDifferenceFunc(x.Labels.List, in.Labels.List,
			func(in *CitizenLabel) uint64 {
				return in.Id
			})

		if err := x.updateCitizenLabels(ctx, tx, in.UserId, added, removed); err != nil {
			return nil, err
		}

		addedOut, err := protojson.Marshal(&CitizenLabels{
			List: added,
		})
		if err != nil {
			return nil, err
		}
		removedOut, err := protojson.Marshal(&CitizenLabels{
			List: removed,
		})
		if err != nil {
			return nil, err
		}

		activities = append(activities, &UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			Type:         UserActivityType_USER_ACTIVITY_TYPE_CHANGED,
			Key:          "UserProps.Labels",
			OldValue:     string(removedOut),
			NewValue:     string(addedOut),
			Reason:       reason,
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
