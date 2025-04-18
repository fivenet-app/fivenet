// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/auth/auth.proto

package auth

import (
	permkeys "github.com/fivenet-app/fivenet/gen/go/proto/services/auth/perms"
	"github.com/fivenet-app/fivenet/pkg/perms"
)

func init() {
	perms.AddPermsToList([]*perms.Perm{

		// Service: AuthService
		{
			Category: permkeys.AuthServicePerm,
			Name:     permkeys.AuthServiceChooseCharacterPerm,
			Attrs:    []perms.Attr{},
			Order:    0,
		},
	})
}
