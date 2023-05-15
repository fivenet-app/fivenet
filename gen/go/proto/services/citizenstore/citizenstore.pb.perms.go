// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/citizenstore/citizenstore.proto

package citizenstore

import "github.com/galexrt/fivenet/pkg/perms"

const (
	CitizenStoreServicePerm perms.Category = "CitizenStoreService"

	CitizenStoreServiceGetUserPerm                     perms.Name = "GetUser"
	CitizenStoreServiceGetUserJobsPermField            perms.Key  = "Jobs"
	CitizenStoreServiceListCitizensPerm                perms.Name = "ListCitizens"
	CitizenStoreServiceListCitizensFieldsPermField     perms.Key  = "Fields"
	CitizenStoreServiceListUserActivityPerm            perms.Name = "ListUserActivity"
	CitizenStoreServiceListUserActivityFieldsPermField perms.Key  = "Fields"
	CitizenStoreServiceSetUserPropsPerm                perms.Name = "SetUserProps"
	CitizenStoreServiceSetUserPropsFieldsPermField     perms.Key  = "Fields"
)

func init() {
	perms.AddPermsToList([]*perms.Perm{
		// Service: CitizenStoreService
		{
			Category: CitizenStoreServicePerm,
			Name:     CitizenStoreServiceGetUserPerm,
			Attrs: []perms.Attr{
				{
					Key:  CitizenStoreServiceGetUserJobsPermField,
					Type: perms.JobGradeListAttributeType,
				},
			},
		},
		{
			Category: CitizenStoreServicePerm,
			Name:     CitizenStoreServiceListCitizensPerm,
			Attrs: []perms.Attr{
				{
					Key:  CitizenStoreServiceListCitizensFieldsPermField,
					Type: perms.StringListAttributeType, ValidValues: []string{"PhoneNumber", "Licenses", "UserProps.Wanted", "UserProps.Job"},
				},
			},
		},
		{
			Category: CitizenStoreServicePerm,
			Name:     CitizenStoreServiceListUserActivityPerm,
			Attrs: []perms.Attr{
				{
					Key:  CitizenStoreServiceListUserActivityFieldsPermField,
					Type: perms.StringListAttributeType, ValidValues: []string{"SourceUser"},
				},
			},
		},
		{
			Category: CitizenStoreServicePerm,
			Name:     CitizenStoreServiceSetUserPropsPerm,
			Attrs: []perms.Attr{
				{
					Key:  CitizenStoreServiceSetUserPropsFieldsPermField,
					Type: perms.StringListAttributeType, ValidValues: []string{"Wanted", "Job"},
				},
			},
		},
	})
}
