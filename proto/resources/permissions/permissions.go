package permissions

import (
	timestamp "github.com/galexrt/fivenet/proto/resources/timestamp"
	"github.com/galexrt/fivenet/query/fivenet/model"
)

func ConvertFromPerm(p *model.FivenetPermissions) *Permission {
	var createdAt *timestamp.Timestamp
	if p.CreatedAt != nil {
		createdAt = timestamp.New(*p.CreatedAt)
	}
	var updatedAt *timestamp.Timestamp
	if p.UpdatedAt != nil {
		updatedAt = timestamp.New(*p.UpdatedAt)
	}

	return &Permission{
		Id:          p.ID,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Name:        p.Name,
		GuardName:   p.GuardName,
		Description: p.Description,
	}
}

func ConvertFromRole(p *model.FivenetRoles) *Role {
	var createdAt *timestamp.Timestamp
	if p.CreatedAt != nil {
		createdAt = timestamp.New(*p.CreatedAt)
	}
	var updatedAt *timestamp.Timestamp
	if p.UpdatedAt != nil {
		updatedAt = timestamp.New(*p.UpdatedAt)
	}

	return &Role{
		Id:          p.ID,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
		Name:        p.Name,
		GuardName:   p.GuardName,
		Description: p.Description,
	}
}
