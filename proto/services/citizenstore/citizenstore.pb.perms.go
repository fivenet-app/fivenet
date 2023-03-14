// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/citizenstore/citizenstore.proto

package citizenstore

import "github.com/galexrt/arpanet/pkg/perms"

var PermsRemap = map[string]string{
	// Service: CitizenStoreService
	"CitizenStoreService/GetUser": "CitizenStoreService/FindUsers",
}

func (s *Server) GetPermsRemap() map[string]string {
	return PermsRemap
}

const (
	CitizenStoreServicePermKey = "CitizenStoreService"
)

func init() {
	perms.AddPermsToList([]*perms.Perm{
		// Service: CitizenStoreService
		{
			Key:    CitizenStoreServicePermKey,
			Name:   "FindUsers",
			Fields: []string{"Licenses", "UserProps"},
		},
		{
			Key:  CitizenStoreServicePermKey,
			Name: "FindUsers",
		},
		{
			Key:    CitizenStoreServicePermKey,
			Name:   "GetUserActivity",
			Fields: []string{"CauseUser"},
		},
		{
			Key:  CitizenStoreServicePermKey,
			Name: "GetUserDocuments",
		},
		{
			Key:    CitizenStoreServicePermKey,
			Name:   "SetUserProps",
			Fields: []string{"Wanted"},
		},
	})
}
