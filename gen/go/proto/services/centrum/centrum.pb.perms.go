// Code generated by protoc-gen-customizer. DO NOT EDIT.
// source: services/centrum/centrum.proto

package centrum

import (
	"github.com/galexrt/fivenet/gen/go/proto/resources/permissions"
	"github.com/galexrt/fivenet/pkg/perms"
)

const (
	CentrumServicePerm perms.Category = "CentrumService"

	CentrumServiceAssignDispatchPerm        perms.Name = "AssignDispatch"
	CentrumServiceAssignUnitPerm            perms.Name = "AssignUnit"
	CentrumServiceAssignUnitAccessPermField perms.Key  = "Access"
	CentrumServiceCreateDispatchPerm        perms.Name = "CreateDispatch"
	CentrumServiceCreateOrUpdateUnitPerm    perms.Name = "CreateOrUpdateUnit"
	CentrumServiceDeleteUnitPerm            perms.Name = "DeleteUnit"
	CentrumServiceListDispatchesPerm        perms.Name = "ListDispatches"
	CentrumServiceListUnitsPerm             perms.Name = "ListUnits"
	CentrumServiceStreamPerm                perms.Name = "Stream"
	CentrumServiceTakeDispatchPerm          perms.Name = "TakeDispatch"
	CentrumServiceUpdateDispatchPerm        perms.Name = "UpdateDispatch"
	CentrumServiceUpdateUnitStatusPerm      perms.Name = "UpdateUnitStatus"
)

var PermsRemap = map[string]string{

	// Service: CentrumService
	"CentrumService/ListDispatchActivity": "CentrumService/Stream",
	"CentrumService/ListUnitActivity":     "CentrumService/ListUnits",
	"CentrumService/UpdateDispatchStatus": "CentrumService/TakeDispatch",
}

func (s *Server) GetPermsRemap() map[string]string {
	return PermsRemap
}

func init() {
	perms.AddPermsToList([]*perms.Perm{
		// Service: CentrumService
		{
			Category: CentrumServicePerm,
			Name:     CentrumServiceAssignDispatchPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: CentrumServicePerm,
			Name:     CentrumServiceAssignUnitPerm,
			Attrs: []perms.Attr{
				{
					Key:           CentrumServiceAssignUnitAccessPermField,
					Type:          permissions.StringListAttributeType,
					ValidValues:   []string{"Own", "Lower_Rank", "Same_Rank"},
					DefaultValues: []string{"Own"},
				},
			},
		},
		{
			Category: CentrumServicePerm,
			Name:     CentrumServiceCreateDispatchPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: CentrumServicePerm,
			Name:     CentrumServiceCreateOrUpdateUnitPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: CentrumServicePerm,
			Name:     CentrumServiceDeleteUnitPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: CentrumServicePerm,
			Name:     CentrumServiceListDispatchesPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: CentrumServicePerm,
			Name:     CentrumServiceListUnitsPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: CentrumServicePerm,
			Name:     CentrumServiceStreamPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: CentrumServicePerm,
			Name:     CentrumServiceTakeDispatchPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: CentrumServicePerm,
			Name:     CentrumServiceUpdateDispatchPerm,
			Attrs:    []perms.Attr{},
		},
		{
			Category: CentrumServicePerm,
			Name:     CentrumServiceUpdateUnitStatusPerm,
			Attrs:    []perms.Attr{},
		},
	})
}
