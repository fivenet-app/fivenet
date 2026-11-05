package citizensstore

import (
	"testing"

	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	usersprops "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/props"
	"github.com/stretchr/testify/assert"
)

func TestMergePropertyMutationActivitiesEmitsOneScalarActivity(t *testing.T) {
	t.Parallel()
	value := &properties.PropertyValue{}
	value.SetStringValue("new")
	scalar := UserPropertyActivities(42, nil, "manual", nil, []PropertyChange{{
		Property: &properties.Property{Key: "custom.note", Value: value},
	}})
	reason := "reviewed"
	lifecycle := UserPropertyActivities(
		42,
		nil,
		"manual",
		[]*properties.Property{{Key: "custom.note", Value: value}},
		[]PropertyChange{{
			Property:  &properties.Property{Key: "custom.note", Value: value},
			Lifecycle: &usersactivity.PropertyLifecycleChange{StateReason: &reason},
		}},
	)
	merged := mergePropertyMutationActivities(scalar, lifecycle)
	require := assert.New(t)
	if !require.Len(merged, 1) {
		return
	}
	require.Equal(
		"reviewed",
		merged[0].GetData().GetPropertyChange().GetLifecycle().GetStateReason(),
	)
}

func TestMergePropertyMutationActivitiesSuppressesNoOpAfterLifecycleDiff(t *testing.T) {
	t.Parallel()
	value := &properties.PropertyValue{}
	value.SetStringValue("unchanged")
	reason := "unchanged"
	scalar := UserPropertyActivities(
		42,
		nil,
		"manual",
		[]*properties.Property{{Key: "custom.note", Value: value}},
		[]PropertyChange{{
			Property:  &properties.Property{Key: "custom.note", Value: value},
			Lifecycle: &usersactivity.PropertyLifecycleChange{StateReason: &reason},
			Cause:     usersactivity.PropertyChangeCause_PROPERTY_CHANGE_CAUSE_MANUAL,
		}},
	)
	assert.Empty(t, mergePropertyMutationActivities(scalar, nil))
}

func TestUserWantedMutationActivitiesRecordsStateOnlyPatch(t *testing.T) {
	t.Parallel()
	wanted := true
	reason := "reviewed"
	activities := UserWantedMutationActivities(
		&usersprops.UserProps{UserId: 42, Wanted: &wanted},
		&usersprops.UserProps{UserId: 42, Wanted: &wanted},
		[]*properties.PropertyMutation{
			{Key: "user.wanted", Reason: &reason},
		},
		nil,
		"manual update",
	)
	require := assert.New(t)
	if !require.Len(activities, 1) {
		return
	}
	assert.Equal(t, int32(42), activities[0].GetTargetUserId())
	assert.True(t, activities[0].GetData().GetWantedChange().GetWanted())
	assert.Equal(t, "manual update", activities[0].GetReason())
}
