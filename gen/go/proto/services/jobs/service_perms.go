// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/jobs/conduct.proto
// source: services/jobs/jobs.proto
// source: services/jobs/timeclock.proto

package jobs

import (
	"github.com/fivenet-app/fivenet/gen/go/proto/resources/permissions"
	permkeys "github.com/fivenet-app/fivenet/gen/go/proto/services/jobs/perms"
	"github.com/fivenet-app/fivenet/pkg/perms"
)

var PermsRemap = map[string]string{
	// Service: JobsService
	"JobsService/GetColleagueLabels":      "JobsService/GetColleague",
	"JobsService/GetColleagueLabelsStats": "JobsService/GetColleague",
	"JobsService/GetMOTD":                 "Any",
	"JobsService/GetSelf":                 "JobsService/ListColleagues",

	// Service: JobsTimeclockService
	"JobsTimeclockService/GetTimeclockStats": "JobsTimeclockService/ListTimeclock",
}

func init() {
	perms.AddPermsToList([]*perms.Perm{

		// Service: JobsConductService
		{
			Category: permkeys.JobsConductServicePerm,
			Name:     permkeys.JobsConductServiceCreateConductEntryPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: permkeys.JobsConductServicePerm,
			Name:     permkeys.JobsConductServiceDeleteConductEntryPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: permkeys.JobsConductServicePerm,
			Name:     permkeys.JobsConductServiceListConductEntriesPerm,
			Attrs: []perms.Attr{
				{
					Key:         permkeys.JobsConductServiceListConductEntriesAccessPermField,
					Type:        permissions.StringListAttributeType,
					ValidValues: []string{"Own", "All"},
				},
			},
		},
		{
			Category: permkeys.JobsConductServicePerm,
			Name:     permkeys.JobsConductServiceUpdateConductEntryPerm,
			Attrs:    []perms.Attr{},
		},

		// Service: JobsService
		{
			Category: permkeys.JobsServicePerm,
			Name:     permkeys.JobsServiceGetColleaguePerm,
			Attrs: []perms.Attr{
				{
					Key:         permkeys.JobsServiceGetColleagueAccessPermField,
					Type:        permissions.StringListAttributeType,
					ValidValues: []string{"Own", "Lower_Rank", "Same_Rank", "Any"},
				},
				{
					Key:         permkeys.JobsServiceGetColleagueTypesPermField,
					Type:        permissions.StringListAttributeType,
					ValidValues: []string{"Note", "Labels"},
				},
			},
		},
		{
			Category: permkeys.JobsServicePerm,
			Name:     permkeys.JobsServiceListColleagueActivityPerm,
			Attrs: []perms.Attr{
				{
					Key:         permkeys.JobsServiceListColleagueActivityTypesPermField,
					Type:        permissions.StringListAttributeType,
					ValidValues: []string{"HIRED", "FIRED", "PROMOTED", "DEMOTED", "ABSENCE_DATE", "NOTE", "LABELS", "NAME"},
				},
			},
		},
		{
			Category: permkeys.JobsServicePerm,
			Name:     permkeys.JobsServiceListColleaguesPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: permkeys.JobsServicePerm,
			Name:     permkeys.JobsServiceManageColleagueLabelsPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: permkeys.JobsServicePerm,
			Name:     permkeys.JobsServiceSetJobsUserPropsPerm,
			Attrs: []perms.Attr{
				{
					Key:         permkeys.JobsServiceSetJobsUserPropsAccessPermField,
					Type:        permissions.StringListAttributeType,
					ValidValues: []string{"Own", "Lower_Rank", "Same_Rank", "Any"},
				},
				{
					Key:         permkeys.JobsServiceSetJobsUserPropsTypesPermField,
					Type:        permissions.StringListAttributeType,
					ValidValues: []string{"AbsenceDate", "Note", "Labels", "Name"},
				},
			},
		},
		{
			Category: permkeys.JobsServicePerm,
			Name:     permkeys.JobsServiceSetMOTDPerm,
			Attrs:    []perms.Attr{},
		},

		// Service: JobsTimeclockService
		{
			Category: permkeys.JobsTimeclockServicePerm,
			Name:     permkeys.JobsTimeclockServiceListInactiveEmployeesPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: permkeys.JobsTimeclockServicePerm,
			Name:     permkeys.JobsTimeclockServiceListTimeclockPerm,
			Attrs: []perms.Attr{
				{
					Key:         permkeys.JobsTimeclockServiceListTimeclockAccessPermField,
					Type:        permissions.StringListAttributeType,
					ValidValues: []string{"All"},
				},
			},
		},
	})
}
