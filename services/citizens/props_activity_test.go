package citizens

import (
	"testing"

	properties "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/properties"
	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	citizensstore "github.com/fivenet-app/fivenet/v2026/stores/citizens"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func stringProperty(key, value string) *properties.Property {
	v := &properties.PropertyValue{}
	v.SetStringValue(value)
	return &properties.Property{Key: key, Value: v}
}

func TestUserPropertyActivities(t *testing.T) {
	t.Parallel()

	sourceID := int32(7)
	previous := []*properties.Property{stringProperty("user.rank", "cadet")}

	tests := []struct {
		name     string
		change   citizensstore.PropertyChange
		activity bool
		old, new string
		cleared  bool
	}{
		{
			"create",
			citizensstore.PropertyChange{Property: stringProperty("user.note", "watchlist")},
			true,
			"",
			"watchlist",
			false,
		},
		{
			"update",
			citizensstore.PropertyChange{Property: stringProperty("user.rank", "officer")},
			true,
			"cadet",
			"officer",
			false,
		},
		{
			"no-op",
			citizensstore.PropertyChange{Property: stringProperty("user.rank", "cadet")},
			false,
			"",
			"",
			false,
		},
		{
			"clear",
			citizensstore.PropertyChange{
				Property: &properties.Property{Key: "user.rank"},
				Delete:   true,
			},
			true,
			"cadet",
			"",
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			activities := citizensstore.UserPropertyActivities(
				42,
				&sourceID,
				"reviewed",
				previous,
				[]citizensstore.PropertyChange{test.change},
			)
			if !test.activity {
				assert.Empty(t, activities)
				return
			}
			require.Len(t, activities, 1)
			activity := activities[0]
			assert.Equal(
				t,
				usersactivity.UserActivityType_USER_ACTIVITY_TYPE_PROPERTY,
				activity.GetType(),
			)
			assert.Equal(t, test.change.Property.GetKey(), activity.GetKey())
			assert.Equal(t, int32(42), activity.GetTargetUserId())
			change := activity.GetData().GetPropertyChange()
			if test.old == "" {
				assert.Nil(t, change.GetOldValue())
			} else {
				assert.Equal(t, test.old, change.GetOldValue().GetStringValue())
			}
			if test.cleared {
				assert.Nil(t, change.GetNewValue())
			} else {
				assert.Equal(t, test.new, change.GetNewValue().GetStringValue())
			}
		})
	}
}
