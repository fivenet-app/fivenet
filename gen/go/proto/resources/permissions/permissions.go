package permissions

import (
	timestamp "github.com/galexrt/fivenet/gen/go/proto/resources/timestamp"
	"github.com/galexrt/fivenet/query/fivenet/model"
)

func (x *Role) GetJobGrade() int32 {
	return x.GetGrade()
}

func (x *Role) SetJobLabel(label string) {
	x.JobLabel = label
}

func (x *Role) SetJobGradeLabel(label string) {
	x.JobGradeLabel = label
}

func ConvertFromPerm(p *model.FivenetPermissions) *Permission {
	var createdAt *timestamp.Timestamp
	if p.CreatedAt != nil {
		createdAt = timestamp.New(*p.CreatedAt)
	}

	return &Permission{
		Id:        p.ID,
		CreatedAt: createdAt,
		Category:  p.Category,
		Name:      p.Name,
		GuardName: p.GuardName,
	}
}

func ConvertFromRole(p *model.FivenetRoles) *Role {
	var createdAt *timestamp.Timestamp
	if p.CreatedAt != nil {
		createdAt = timestamp.New(*p.CreatedAt)
	}

	return &Role{
		Id:        p.ID,
		CreatedAt: createdAt,
		Job:       p.Job,
		Grade:     p.Grade,
	}
}
