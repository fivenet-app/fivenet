package model

import (
	"gorm.io/gorm"
)

type ContentType string

const (
	PlaintextContentType ContentType = "plaintext"
	MakrdownContentType  ContentType = "markdown"
)

type Document struct {
	gorm.Model

	Title       string
	Content     string
	ContentType ContentType
	Creator     uint
	Path        string
}

type AccessTypes struct {
	CanView bool
	CanEdit bool
}

type JobAccess struct {
	gorm.Model

	DocumentID uint

	AccessTypes

	JobID        uint
	Job          Job
	MinimumGrade int `gorm:"default:0"`
}

func (a JobAccess) CanEdit(citizen Citizen) bool {
	if a.Job.ID != citizen.Job.ID {
		return false
	}

	if a.MinimumGrade > citizen.JobGrade.Grade {
		return false
	}

	return a.AccessTypes.CanEdit
}

func (a JobAccess) CanView(citizen Citizen) bool {
	if a.Job.ID != citizen.Job.ID {
		return false
	}

	if a.MinimumGrade > citizen.JobGrade.Grade {
		return false
	}

	return a.AccessTypes.CanView
}

type CitizenAccess struct {
	gorm.Model

	DocumentID uint

	AccessTypes

	CitizenID uint
	Citizen   Citizen
}

func (a CitizenAccess) CanEdit(citizen Citizen) bool {
	if a.CitizenID != citizen.ID {
		return false
	}

	return a.AccessTypes.CanEdit
}

func (a CitizenAccess) CanView(citizen Citizen) bool {
	if a.CitizenID != citizen.ID {
		return false
	}

	return a.AccessTypes.CanView
}
