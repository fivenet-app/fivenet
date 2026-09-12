package notifications

import (
	"testing"

	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	resourcesnotifications "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/notifications"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyUserInfoChanged(t *testing.T) {
	t.Parallel()

	current := &pbuserinfo.UserInfo{
		Job:      "police",
		JobGrade: 1,
	}

	applyUserInfoChanged(current, &pbuserinfo.UserInfoChanged{
		NewJob:      new("ems"),
		NewJobGrade: new(int32(3)),
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
	assert.False(t, current.GetCanBeConfigAdmin())
}

func TestApplyAccountGroupsChangedClearsNil(t *testing.T) {
	t.Parallel()

	current := &pbuserinfo.UserInfo{
		Groups: &accounts.AccountGroups{Groups: []string{"old"}},
	}

	applyAccountGroupsChanged(current, &pbuserinfo.AccountGroupsChanged{})

	assert.Nil(t, current.GetGroups())
	assert.False(t, current.GetCanBeSuperuser())
	assert.False(t, current.GetCanBeConfigAdmin())
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
	assert.False(t, current.GetCanBeConfigAdmin())
	assert.False(t, current.GetSuperuser())
	assert.Equal(t, "ems", current.GetJob())
	assert.Equal(t, int32(2), current.GetJobGrade())
}

func TestBuildSubjectsAccountOnly(t *testing.T) {
	t.Parallel()

	s := &Server{}

	baseSubjects, additionalSubjects, err := s.buildSubjects(t.Context(), &pbuserinfo.UserInfo{
		AccountId: 42,
	})
	require.NoError(t, err)

	assert.Equal(t, []string{
		"notifi.account.42",
		"notifi.sys",
	}, baseSubjects)
	assert.Empty(t, additionalSubjects)
}

func TestValidPreferenceScope(t *testing.T) {
	t.Parallel()

	assert.True(t, validPreferenceScope(&resourcesnotifications.NotificationPreference{}))
	assert.True(t, validPreferenceScope(&resourcesnotifications.NotificationPreference{
		Category: resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
	}))
	assert.True(t, validPreferenceScope(&resourcesnotifications.NotificationPreference{
		Category: resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_JOBS,
		Kind:     resourcesnotifications.NotificationKind_NOTIFICATION_KIND_JOBS_GROUP_LEADERSHIP_ADDED,
	}))
	assert.False(t, validPreferenceScope(&resourcesnotifications.NotificationPreference{
		Category: resourcesnotifications.NotificationCategory_NOTIFICATION_CATEGORY_DOCUMENT,
		Kind:     resourcesnotifications.NotificationKind_NOTIFICATION_KIND_JOBS_GROUP_LEADERSHIP_ADDED,
	}))
	assert.False(t, validPreferenceScope(&resourcesnotifications.NotificationPreference{
		Category: resourcesnotifications.NotificationCategory(99),
	}))
	assert.False(t, validPreferenceScope(nil))
}

func TestHasDeliveryPreferenceOverride(t *testing.T) {
	t.Parallel()

	assert.False(t, hasDeliveryPreferenceOverride(nil))
	assert.False(t, hasDeliveryPreferenceOverride(&resourcesnotifications.NotificationPreference{}))
	assert.True(t, hasDeliveryPreferenceOverride(&resourcesnotifications.NotificationPreference{InboxEnabled: new(true)}))
}
