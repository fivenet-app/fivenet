// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/jobs/qualifications.proto
// source: services/jobs/conduct.proto
// source: services/jobs/jobs.proto
// source: services/jobs/timeclock.proto
// source: services/jobs/requests.proto

package jobs

import (
	"github.com/galexrt/fivenet/gen/go/proto/resources/permissions"
	permkeys "github.com/galexrt/fivenet/gen/go/proto/services/jobs/perms"
	"github.com/galexrt/fivenet/pkg/perms"
)

var PermsRemap = map[string]string{

	// Service: JobsRequestsService
	"JobsRequestsService/ListRequestComments": "JobsRequestsService/ListRequests",
	"JobsRequestsService/ListRequestTypes":    "JobsRequestsService/ListRequests",
	"JobsRequestsService/PostRequestComment":  "JobsRequestsService/CreateRequest",

	// Service: JobsTimeclockService
	"JobsTimeclockService/GetTimeclockStats": "JobsTimeclockService/TimeclockList",
}

func (s *Server) GetPermsRemap() map[string]string {
	return PermsRemap
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
					Key:           permkeys.JobsConductServiceListConductEntriesAccessPermField,
					Type:          permissions.StringListAttributeType,
					ValidValues:   []string{"Own", "All"},
					DefaultValues: []string{"Own"},
				},
			},
		},
		{
			Category: permkeys.JobsConductServicePerm,
			Name:     permkeys.JobsConductServiceUpdateConductEntryPerm,
			Attrs:    []perms.Attr{},
		},

		// Service: JobsRequestsService
		{
			Category: permkeys.JobsRequestsServicePerm,
			Name:     permkeys.JobsRequestsServiceCreateOrUpdateRequestTypePerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: permkeys.JobsRequestsServicePerm,
			Name:     permkeys.JobsRequestsServiceCreateRequestPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: permkeys.JobsRequestsServicePerm,
			Name:     permkeys.JobsRequestsServiceDeleteRequestPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: permkeys.JobsRequestsServicePerm,
			Name:     permkeys.JobsRequestsServiceDeleteRequestCommentPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: permkeys.JobsRequestsServicePerm,
			Name:     permkeys.JobsRequestsServiceDeleteRequestTypePerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: permkeys.JobsRequestsServicePerm,
			Name:     permkeys.JobsRequestsServiceListRequestsPerm,
			Attrs: []perms.Attr{
				{
					Key:           permkeys.JobsRequestsServiceListRequestsAccessPermField,
					Type:          permissions.StringListAttributeType,
					ValidValues:   []string{"Own", "All"},
					DefaultValues: []string{"Own"},
				},
			},
		},
		{
			Category: permkeys.JobsRequestsServicePerm,
			Name:     permkeys.JobsRequestsServiceUpdateRequestPerm,
			Attrs:    []perms.Attr{},
		},

		// Service: JobsService
		{
			Category: permkeys.JobsServicePerm,
			Name:     permkeys.JobsServiceListColleaguesPerm,
			Attrs:    []perms.Attr{},
		},

		// Service: JobsTimeclockService
		{
			Category: permkeys.JobsTimeclockServicePerm,
			Name:     permkeys.JobsTimeclockServiceListTimeclockPerm,
			Attrs: []perms.Attr{
				{
					Key:           permkeys.JobsTimeclockServiceListTimeclockAccessPermField,
					Type:          permissions.StringListAttributeType,
					ValidValues:   []string{"All"},
					DefaultValues: []string{},
				},
			},
		},
		{
			Category: permkeys.JobsTimeclockServicePerm,
			Name:     permkeys.JobsTimeclockServiceTimeclockListPerm,
			Attrs:    []perms.Attr{},
		},
	})
}
