// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/citizenstore/citizenstore.proto

package citizenstore

import (
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/permissions"
	permkeys "github.com/fivenet-app/fivenet/gen/go/proto/services/citizenstore/perms"
	"github.com/fivenet-app/fivenet/pkg/perms"
)

var PermsRemap = map[string]string{
	// Service: CitizenStoreService
	"CitizenStoreService/SetProfilePicture": "Any",
}

func init() {
	perms.AddPermsToList([]*perms.Perm{

		// Service: CitizenStoreService
		{
			Category: permkeys.CitizenStoreServicePerm,
			Name:     permkeys.CitizenStoreServiceGetUserPerm,
			Attrs: []perms.Attr{
				{
					Key:  permkeys.CitizenStoreServiceGetUserJobsPermField,
					Type: permissions.JobGradeListAttributeType,
				},
			},
		},
		{
			Category: permkeys.CitizenStoreServicePerm,
			Name:     permkeys.CitizenStoreServiceListCitizensPerm,
			Attrs: []perms.Attr{
				{
					Key:         permkeys.CitizenStoreServiceListCitizensFieldsPermField,
					Type:        permissions.StringListAttributeType,
					ValidValues: []string{"PhoneNumber", "Licenses", "UserProps.Wanted", "UserProps.Job", "UserProps.TrafficInfractionPoints", "UserProps.OpenFines", "UserProps.BloodType", "UserProps.MugShot", "UserProps.Labels", "UserProps.Email"},
				},
			},
		},
		{
			Category: permkeys.CitizenStoreServicePerm,
			Name:     permkeys.CitizenStoreServiceListUserActivityPerm,
			Attrs: []perms.Attr{
				{
					Key:         permkeys.CitizenStoreServiceListUserActivityFieldsPermField,
					Type:        permissions.StringListAttributeType,
					ValidValues: []string{"SourceUser", "Own"},
				},
			},
		},
		{
			Category: permkeys.CitizenStoreServicePerm,
			Name:     permkeys.CitizenStoreServiceManageCitizenLabelsPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: permkeys.CitizenStoreServicePerm,
			Name:     permkeys.CitizenStoreServiceSetUserPropsPerm,
			Attrs: []perms.Attr{
				{
					Key:         permkeys.CitizenStoreServiceSetUserPropsFieldsPermField,
					Type:        permissions.StringListAttributeType,
					ValidValues: []string{"Wanted", "Job", "TrafficInfractionPoints", "MugShot", "Labels"},
				},
			},
		},
	})
}
