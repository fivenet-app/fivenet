package userinfo

import (
	"testing"

	accounts "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/accounts"
	pbuserinfo "github.com/fivenet-app/fivenet/v2026/gen/go/proto/resources/userinfo"
	"github.com/stretchr/testify/assert"
)

func TestCheckAndSetSuperuserAcceptsConfigAdminMembership(t *testing.T) {
	t.Parallel()

	retriever := &Retriever{
		configAdminGroups: []string{"config-admin"},
	}

	userInfo := &pbuserinfo.UserInfo{
		Groups: &accounts.AccountGroups{
			Groups: []string{"config-admin"},
		},
		License: "license-42",
	}

	retriever.checkAndSetSuperuser(userInfo)

	assert.True(t, userInfo.GetCanBeSuperuser())
	assert.False(t, userInfo.GetSuperuser())
}
