package notifications

import (
	"testing"

	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/assert"
)

func TestApplyUserInfoChanged(t *testing.T) {
	t.Parallel()

	current := &pbuserinfo.UserInfo{
		Job:      "police",
		JobGrade: 1,
	}

	applyUserInfoChanged(current, &pbuserinfo.UserInfoChanged{
		NewJob:      func() *string { v := "ems"; return &v }(),
		NewJobGrade: func() *int32 { v := int32(3); return &v }(),
	})

	assert.Equal(t, "ems", current.GetJob())
	assert.Equal(t, int32(3), current.GetJobGrade())
}

func TestApplyUserInfoChangedIgnoresNil(t *testing.T) {
	t.Parallel()

	current := &pbuserinfo.UserInfo{
		Job:      "police",
		JobGrade: 1,
	}

	applyUserInfoChanged(current, nil)

	assert.Equal(t, "police", current.GetJob())
	assert.Equal(t, int32(1), current.GetJobGrade())
}

func TestApplyAccountGroupsChanged(t *testing.T) {
	t.Parallel()

	current := &pbuserinfo.UserInfo{
		Groups: &accounts.AccountGroups{Groups: []string{"old"}},
	}

	applyAccountGroupsChanged(current, &pbuserinfo.AccountGroupsChanged{
		NewGroups:      &accounts.AccountGroups{Groups: []string{"supporter", "donator"}},
		CanBeSuperuser: true,
	})

	assert.Equal(t, []string{"supporter", "donator"}, current.GetGroups().GetGroups())
	assert.True(t, current.GetCanBeSuperuser())
}

func TestApplyAccountGroupsChangedClearsNil(t *testing.T) {
	t.Parallel()

	current := &pbuserinfo.UserInfo{
		Groups: &accounts.AccountGroups{Groups: []string{"old"}},
	}

	applyAccountGroupsChanged(current, &pbuserinfo.AccountGroupsChanged{})

	assert.Nil(t, current.GetGroups())
	assert.False(t, current.GetCanBeSuperuser())
	assert.False(t, current.GetSuperuser())
}

func TestApplyAccountGroupsChangedRestoresOriginalJobWhenRevokingSuperuser(t *testing.T) {
	t.Parallel()

	current := &pbuserinfo.UserInfo{
		Job:       "ems-super",
		JobGrade:  7,
		Superuser: true,
		OriginalJob: &pbuserinfo.OriginalJob{
			Job:      "ems",
			JobGrade: 2,
		},
	}

	applyAccountGroupsChanged(current, &pbuserinfo.AccountGroupsChanged{
		CanBeSuperuser: false,
	})

	assert.False(t, current.GetCanBeSuperuser())
	assert.False(t, current.GetSuperuser())
	assert.Equal(t, "ems", current.GetJob())
	assert.Equal(t, int32(2), current.GetJobGrade())
}
