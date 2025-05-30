// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/completor/completor.proto

package completor

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	permkeys "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/completor/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
)

var PermsRemap = map[string]string{
	// Service: completor.CompletorService
	"completor.CompletorService/CompleteJobs": "Any",
	"completor.CompletorService/ListLawBooks": "Any",
}

func init() {
	perms.AddPermsToList([]*perms.Perm{

		// Service: completor.CompletorService
		{
			Category: permkeys.CompletorServicePerm,
			Name:     permkeys.CompletorServiceCompleteCitizenLabelsPerm,
			Attrs: []perms.Attr{
				{
					Key:  permkeys.CompletorServiceCompleteCitizenLabelsJobsPermField,
					Type: permissions.JobListAttributeType,
				},
			},
			Order: 0,
		},
		{
			Category: permkeys.CompletorServicePerm,
			Name:     permkeys.CompletorServiceCompleteCitizensPerm,
			Attrs:    []perms.Attr{},
			Order:    0,
		},
		{
			Category: permkeys.CompletorServicePerm,
			Name:     permkeys.CompletorServiceCompleteDocumentCategoriesPerm,
			Attrs: []perms.Attr{
				{
					Key:  permkeys.CompletorServiceCompleteDocumentCategoriesJobsPermField,
					Type: permissions.JobListAttributeType,
				},
			},
			Order: 0,
		},
	})
}
