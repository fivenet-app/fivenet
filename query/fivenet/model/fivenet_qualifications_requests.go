//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

import (
	"time"
)

type FivenetQualificationsRequests struct {
	CreatedAt       *time.Time `json:"created_at"`
	DeletedAt       *time.Time `json:"deleted_at"`
	QualificationID uint64     `json:"qualification_id"`
	UserID          int32      `json:"user_id"`
	UserComment     *string    `json:"user_comment"`
	Approved        *bool      `json:"approved"`
	ApprovedAt      *time.Time `json:"approved_at"`
	ApproverComment *string    `json:"approver_comment"`
	ApproverID      *int32     `json:"approver_id"`
	ApproverJob     *string    `json:"approver_job"`
}
