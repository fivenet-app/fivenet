//
// Code generated by go-jet DO NOT EDIT.
//
// WARNING: Changes to this file may cause incorrect behavior
// and will be lost if the code is regenerated
//

package model

type FivenetDocumentsAccess struct {
	ID           uint64  `sql:"primary_key" json:"id"`
	TargetID     uint64  `json:"target_id"`
	UserID       *int32  `json:"user_id"`
	Job          *string `json:"job"`
	MinimumGrade *int32  `json:"minimum_grade"`
	Access       int16   `json:"access"`
}
