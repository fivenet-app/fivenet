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

type FivenetDocumentsWorkflowUsers struct {
	DocumentID            uint64     `json:"document_id"`
	UserID                int32      `json:"user_id"`
	ManualReminderTime    *time.Time `json:"manual_reminder_time"`
	ManualReminderMessage *string    `json:"manual_reminder_message"`
}
