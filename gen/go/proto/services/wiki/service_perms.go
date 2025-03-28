// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/wiki/wiki.proto

package wiki

import (
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/permissions"
	permkeys "github.com/fivenet-app/fivenet/gen/go/proto/services/wiki/perms"
	"github.com/fivenet-app/fivenet/pkg/perms"
)

var PermsRemap = map[string]string{
	// Service: WikiService
	"WikiService/GetPage": "WikiService/ListPages",
}

func init() {
	perms.AddPermsToList([]*perms.Perm{

		// Service: WikiService
		{
			Category: permkeys.WikiServicePerm,
			Name:     permkeys.WikiServiceCreatePagePerm,
			Attrs: []perms.Attr{
				{
					Key:         permkeys.WikiServiceCreatePageFieldsPermField,
					Type:        permissions.StringListAttributeType,
					ValidValues: []string{"Public"},
				},
			},
		},
		{
			Category: permkeys.WikiServicePerm,
			Name:     permkeys.WikiServiceDeletePagePerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: permkeys.WikiServicePerm,
			Name:     permkeys.WikiServiceListPageActivityPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: permkeys.WikiServicePerm,
			Name:     permkeys.WikiServiceListPagesPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: permkeys.WikiServicePerm,
			Name:     permkeys.WikiServiceUpdatePagePerm,
			Attrs:    []perms.Attr{},
		},
	})
}
