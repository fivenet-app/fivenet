// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/docstore/docstore.proto

package docstore

import "github.com/galexrt/arpanet/pkg/perms"

var PermsRemap = map[string]string{
	// Service: DocStoreService
	"DocStoreService/EditDocumentComment":     "DocStoreService/PostDocumentComment",
	"DocStoreService/GetDocumentReferences":   "DocStoreService/GetDocument",
	"DocStoreService/GetDocumentRelations":    "DocStoreService/GetDocument",
	"DocStoreService/GetTemplate":             "DocStoreService/ListTemplates",
	"DocStoreService/RemoveDocumentReference": "DocStoreService/AddDocumentReference",
	"DocStoreService/RemoveDocumentRelation":  "DocStoreService/AddDocumentRelation",
}

func (s *Server) GetPermsRemap() map[string]string {
	return PermsRemap
}

const (
	DocStoreServicePermKey = "DocStoreService"
)

func init() {
	perms.AddPermsToList([]*perms.Perm{
		// Service: DocStoreService
		{
			Key:  DocStoreServicePermKey,
			Name: "AddDocumentReference",
		},
		{
			Key:  DocStoreServicePermKey,
			Name: "AddDocumentRelation",
		},
		{
			Key:  DocStoreServicePermKey,
			Name: "CreateDocument",
		},
		{
			Key:  DocStoreServicePermKey,
			Name: "PostDocumentComment",
		},
		{
			Key:  DocStoreServicePermKey,
			Name: "FindDocuments",
		},
		{
			Key:  DocStoreServicePermKey,
			Name: "GetDocument",
		},
		{
			Key:  DocStoreServicePermKey,
			Name: "GetDocumentAccess",
		},
		{
			Key:  DocStoreServicePermKey,
			Name: "GetDocumentComments",
		},
		{
			Key:  DocStoreServicePermKey,
			Name: "GetDocument",
		},
		{
			Key:  DocStoreServicePermKey,
			Name: "GetDocument",
		},
		{
			Key:  DocStoreServicePermKey,
			Name: "ListTemplates",
		},
		{
			Key:  DocStoreServicePermKey,
			Name: "GetUserDocuments",
		},
		{
			Key:  DocStoreServicePermKey,
			Name: "ListTemplates",
		},
		{
			Key:  DocStoreServicePermKey,
			Name: "PostDocumentComment",
		},
		{
			Key:  DocStoreServicePermKey,
			Name: "AddDocumentReference",
		},
		{
			Key:  DocStoreServicePermKey,
			Name: "AddDocumentRelation",
		},
		{
			Key:  DocStoreServicePermKey,
			Name: "SetDocumentAccess",
		},
		{
			Key:  DocStoreServicePermKey,
			Name: "UpdateDocument",
		},
	})
}
