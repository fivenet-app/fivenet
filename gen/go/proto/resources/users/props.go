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

var tUserProps = table.FivenetUserProps

func (x *UserProps) HandleChanges(ctx context.Context, tx qrm.DB, in *UserProps, sourceUserId *int32, reason string) ([]*UserActivity, error) {
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

	if in.Attributes != nil {
		if in.Attributes.List == nil {
			in.Attributes.List = []*CitizenAttribute{}
		}

		slices.SortFunc(in.Attributes.List, func(a, b *CitizenAttribute) int {
			return strings.Compare(a.Name, b.Name)
		})
	} else {
		in.Attributes = x.Attributes
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
	if !proto.Equal(in.Attributes, x.Attributes) {
		added, removed := utils.SlicesDifferenceFunc(x.Attributes.List, in.Attributes.List,
			func(in *CitizenAttribute) uint64 {
				return in.Id
			})

		if err := x.updateCitizenAttributes(ctx, tx, in.UserId, added, removed); err != nil {
			return nil, err
		}

		addedOut, err := protojson.Marshal(&CitizenAttributes{
			List: added,
		})
		if err != nil {
			return nil, err
		}
		removedOut, err := protojson.Marshal(&CitizenAttributes{
			List: removed,
		})
		if err != nil {
			return nil, err
		}

		activities = append(activities, &UserActivity{
			SourceUserId: sourceUserId,
			TargetUserId: x.UserId,
			Type:         UserActivityType_USER_ACTIVITY_TYPE_CHANGED,
			Key:          "UserProps.Attributes",
			OldValue:     string(removedOut),
			NewValue:     string(addedOut),
			Reason:       reason,
		})
	}

	return activities, nil
}

func (s *UserProps) updateCitizenAttributes(ctx context.Context, tx qrm.DB, userId int32, added []*CitizenAttribute, removed []*CitizenAttribute) error {
	tUserCitizenAttributes := table.FivenetUserCitizenAttributes

	if len(added) > 0 {
		addedAttributes := make([]*model.FivenetUserCitizenAttributes, len(added))
		for i, attribute := range added {
			addedAttributes[i] = &model.FivenetUserCitizenAttributes{
				UserID:      userId,
				AttributeID: attribute.Id,
			}
		}

		stmt := tUserCitizenAttributes.
			INSERT(
				tUserCitizenAttributes.UserID,
				tUserCitizenAttributes.AttributeID,
			).
			MODELS(addedAttributes)

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

		stmt := tUserCitizenAttributes.
			DELETE().
			WHERE(jet.AND(
				tUserCitizenAttributes.UserID.EQ(jet.Int32(userId)),
				tUserCitizenAttributes.AttributeID.IN(ids...),
			)).
			LIMIT(int64(len(removed)))

		if _, err := stmt.ExecContext(ctx, tx); err != nil {
			return err
		}
	}

	return nil
}
