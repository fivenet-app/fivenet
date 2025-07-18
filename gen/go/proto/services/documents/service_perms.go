// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/documents/collab.proto
// source: services/documents/documents.proto

package documents

import (
	"github.com/fivenet-app/fivenet/v2025/gen/go/proto/resources/permissions"
	permkeys "github.com/fivenet-app/fivenet/v2025/gen/go/proto/services/documents/perms"
	"github.com/fivenet-app/fivenet/v2025/pkg/perms"
)

var PermsRemap = map[string]string{
	// Service: documents.CollabService
	"documents.CollabService/JoinRoom": "documents.DocumentsService/UpdateDocument",

	// Service: documents.DocumentsService
	"documents.DocumentsService/CreateDocument":          "documents.DocumentsService/UpdateDocument",
	"documents.DocumentsService/EditComment":             "documents.DocumentsService/ListDocuments",
	"documents.DocumentsService/GetComments":             "documents.DocumentsService/ListDocuments",
	"documents.DocumentsService/GetDocument":             "documents.DocumentsService/ListDocuments",
	"documents.DocumentsService/GetDocumentAccess":       "documents.DocumentsService/ListDocuments",
	"documents.DocumentsService/GetDocumentReferences":   "documents.DocumentsService/ListDocuments",
	"documents.DocumentsService/GetDocumentRelations":    "documents.DocumentsService/ListDocuments",
	"documents.DocumentsService/GetTemplate":             "documents.DocumentsService/ListTemplates",
	"documents.DocumentsService/ListDocumentPins":        "documents.DocumentsService/ListDocuments",
	"documents.DocumentsService/PostComment":             "documents.DocumentsService/ListDocuments",
	"documents.DocumentsService/RemoveDocumentReference": "documents.DocumentsService/AddDocumentReference",
	"documents.DocumentsService/RemoveDocumentRelation":  "documents.DocumentsService/AddDocumentRelation",
	"documents.DocumentsService/SetDocumentAccess":       "documents.DocumentsService/UpdateDocument",
	"documents.DocumentsService/UpdateDocumentReq":       "documents.DocumentsService/CreateDocumentReq",
	"documents.DocumentsService/UpdateTemplate":          "documents.DocumentsService/CreateTemplate",
	"documents.DocumentsService/UploadFile":              "documents.DocumentsService/UpdateDocument",
}

func init() {
	perms.AddPermsToList([]*perms.Perm{

		// Service: documents.DocumentsService
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceAddDocumentReferencePerm,
			Attrs:    []perms.Attr{},
			Order:    0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceAddDocumentRelationPerm,
			Attrs:    []perms.Attr{},
			Order:    0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceChangeDocumentOwnerPerm,
			Attrs: []perms.Attr{
				{
					Key:         permkeys.DocumentsServiceChangeDocumentOwnerAccessPermField,
					Type:        permissions.StringListAttributeType,
					ValidValues: []string{"Own", "Lower_Rank", "Same_Rank", "Any"},
				},
			},
			Order: 0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceCreateDocumentReqPerm,
			Attrs: []perms.Attr{
				{
					Key:         permkeys.DocumentsServiceCreateDocumentReqTypesPermField,
					Type:        permissions.StringListAttributeType,
					ValidValues: []string{"Access", "Closure", "Update", "Deletion", "OwnerChange"},
				},
			},
			Order: 0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceCreateOrUpdateCategoryPerm,
			Attrs:    []perms.Attr{},
			Order:    0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceCreateTemplatePerm,
			Attrs:    []perms.Attr{},
			Order:    0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceDeleteCategoryPerm,
			Attrs:    []perms.Attr{},
			Order:    0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceDeleteCommentPerm,
			Attrs: []perms.Attr{
				{
					Key:         permkeys.DocumentsServiceDeleteCommentAccessPermField,
					Type:        permissions.StringListAttributeType,
					ValidValues: []string{"Own", "Lower_Rank", "Same_Rank", "Any"},
				},
			},
			Order: 0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceDeleteDocumentPerm,
			Attrs: []perms.Attr{
				{
					Key:         permkeys.DocumentsServiceDeleteDocumentAccessPermField,
					Type:        permissions.StringListAttributeType,
					ValidValues: []string{"Own", "Lower_Rank", "Same_Rank", "Any"},
				},
			},
			Order: 0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceDeleteDocumentReqPerm,
			Attrs:    []perms.Attr{},
			Order:    0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceDeleteTemplatePerm,
			Attrs:    []perms.Attr{},
			Order:    0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceListCategoriesPerm,
			Attrs:    []perms.Attr{},
			Order:    0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceListDocumentActivityPerm,
			Attrs:    []perms.Attr{},
			Order:    0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceListDocumentReqsPerm,
			Attrs:    []perms.Attr{},
			Order:    0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceListDocumentsPerm,
			Attrs:    []perms.Attr{},
			Order:    0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceListTemplatesPerm,
			Attrs:    []perms.Attr{},
			Order:    0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceListUserDocumentsPerm,
			Attrs:    []perms.Attr{},
			Order:    0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceSetDocumentReminderPerm,
			Attrs:    []perms.Attr{},
			Order:    0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceToggleDocumentPerm,
			Attrs: []perms.Attr{
				{
					Key:         permkeys.DocumentsServiceToggleDocumentAccessPermField,
					Type:        permissions.StringListAttributeType,
					ValidValues: []string{"Own", "Lower_Rank", "Same_Rank", "Any"},
				},
			},
			Order: 0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceToggleDocumentPinPerm,
			Attrs: []perms.Attr{
				{
					Key:         permkeys.DocumentsServiceToggleDocumentPinTypesPermField,
					Type:        permissions.StringListAttributeType,
					ValidValues: []string{"JobWide"},
				},
			},
			Order: 0,
		},
		{
			Category: permkeys.DocumentsServicePerm,
			Name:     permkeys.DocumentsServiceUpdateDocumentPerm,
			Attrs: []perms.Attr{
				{
					Key:         permkeys.DocumentsServiceUpdateDocumentAccessPermField,
					Type:        permissions.StringListAttributeType,
					ValidValues: []string{"Own", "Lower_Rank", "Same_Rank", "Any"},
				},
			},
			Order: 0,
		},
	})
}
