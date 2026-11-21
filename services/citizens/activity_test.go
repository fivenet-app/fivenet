package citizens

import (
	"testing"

	usersactivity "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/users/activity"
	"github.com/stretchr/testify/assert"
)

func TestUserActivityTypesForPropertyKeys(t *testing.T) {
	t.Parallel()

	activityTypes := []usersactivity.UserActivityType{
		usersactivity.UserActivityType_USER_ACTIVITY_TYPE_WANTED,
		usersactivity.UserActivityType_USER_ACTIVITY_TYPE_JOB,
		usersactivity.UserActivityType_USER_ACTIVITY_TYPE_PROPERTY,
	}

	assert.Equal(t,
		[]usersactivity.UserActivityType{
			usersactivity.UserActivityType_USER_ACTIVITY_TYPE_WANTED,
			usersactivity.UserActivityType_USER_ACTIVITY_TYPE_JOB,
		},
		userActivityTypesForPropertyKeys(
			[]string{"user.wanted", "user.job", "user.risk_level"},
			activityTypes,
			[]string{"user.wanted", "user.job", "user.risk_level"},
		),
	)
	assert.Empty(t, userActivityTypesForPropertyKeys(
		[]string{"user.wanted"},
		activityTypes,
		[]string{"user.job"},
	))
}
