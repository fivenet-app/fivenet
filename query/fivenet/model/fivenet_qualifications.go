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

type FivenetQualifications struct {
	ID                 uint64     `sql:"primary_key" json:"id"`
	CreatedAt          *time.Time `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at"`
	Job                string     `json:"job"`
	Weight             *uint32    `json:"weight"`
	Closed             *bool      `json:"closed"`
	Draft              *bool      `json:"draft"`
	Public             *bool      `json:"public"`
	Abbreviation       string     `json:"abbreviation"`
	Title              string     `json:"title"`
	Description        *string    `json:"description"`
	ContentType        int16      `json:"content_type"`
	Content            *string    `json:"content"`
	CreatorID          *int32     `json:"creator_id"`
	CreatorJob         string     `json:"creator_job"`
	DiscordSyncEnabled *bool      `json:"discord_sync_enabled"`
	DiscordSettings    *string    `json:"discord_settings"`
	ExamMode           *int16     `json:"exam_mode"`
	ExamSettings       *string    `json:"exam_settings"`
	LabelSyncEnabled   *bool      `json:"label_sync_enabled"`
	LabelSyncFormat    *string    `json:"label_sync_format"`
}
