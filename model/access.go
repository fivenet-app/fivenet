package model

import (
	"time"

	"gorm.io/gorm"
)

type AccessTypes struct {
	CanView bool `json:"canView"`
	CanEdit bool `json:"canEdit"`
}

const TableNameAccessJob = "arpanet_access_jobs"

type AccessJob struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time

	DocumentID uint

	AccessTypes

	JobID        uint
	Job          Job
	MinimumGrade int `gorm:"default:0"`
}

func (*AccessJob) TableName() string {
	return TableNameAccessJob
}

func (a AccessJob) CanEdit(citizen Citizen) bool {
	//if a.Job.ID != citizen.Job.ID {
	//	return false
	//}

	if a.MinimumGrade > citizen.JobGrade {
		return false
	}

	return a.AccessTypes.CanEdit
}

func (a AccessJob) CanView(citizen Citizen) bool {
	//if a.Job.ID != citizen.Job.ID {
	//	return false
	//}

	if a.MinimumGrade > citizen.JobGrade {
		return false
	}

	return a.AccessTypes.CanView
}

const TableNameAccessCitizen = "arpanet_access_citizens"

type AccessCitizen struct {
	gorm.Model

	DocumentID uint

	AccessTypes

	CitizenID uint
	Citizen   Citizen
}

// TableName Document's table name
func (*AccessCitizen) TableName() string {
	return TableNameAccessCitizen
}

func (a AccessCitizen) CanEdit(citizen Citizen) bool {
	//if a.CitizenID != citizen.ID {
	//	return false
	//}

	return a.AccessTypes.CanEdit
}

func (a AccessCitizen) CanView(citizen Citizen) bool {
	//if a.CitizenID != citizen.ID {
	//	return false
	//}

	return a.AccessTypes.CanView
}
