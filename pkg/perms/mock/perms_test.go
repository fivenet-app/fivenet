package mock

import (
	"testing"

	"github.com/galexrt/fivenet/pkg/grpc/auth/userinfo"
	"github.com/stretchr/testify/assert"
)

// TestPermsMock test the most commonly used funcs to ensure we can safely
// mock permissions
func TestPermsMock(t *testing.T) {
	p := NewMock()
	ps, err := p.GetAllPermissions()
	assert.NoError(t, err)
	assert.Empty(t, ps)

	user1Perms := []string{
		"test-1",
		"test-2",
		"test-perjob-ambulance",
		"test-perjob-doj",
		"test-perjob-police",
	}
	for _, v := range user1Perms {
		p.AddUserPerm(1, v)
	}

	user1AllPerms, err := p.GetPermissionsOfUser(&userinfo.UserInfo{
		UserId:    1,
		Job:       "ambulance",
		JobGrade:  20,
		Group:     "user",
		SuperUser: false,
	})
	assert.NoError(t, err)
	assert.Len(t, user1AllPerms, 5)

	/*
		// Check if the permission check is working as expected
		can := p.Can(1, "test-2")
		assert.True(t, can)
		can = p.Can(1, "test-non-existant")
		assert.False(t, can)
		// We only match on full permission names
		can = p.Can(1, "test-perjob")
		assert.False(t, can)
		can = p.Can(1, "test-perjob")
		assert.False(t, can)
	*/

	// Make sure the perm counter has counted the 3 "can perms checks"
	pc := p.Counter.GetUser(1)
	assert.Len(t, pc, 3)
	assert.Equal(t, pc["test-2"], 1)
	assert.Equal(t, pc["test-non-existant"], 1)
	assert.Equal(t, pc["test-perjob"], 2)
}
