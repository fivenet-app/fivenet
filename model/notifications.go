package model

import "time"

type Notification struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	ReadAt    time.Time `json:"read_at"`

	Title   string
	Content string
	Link    *string

	ExpiresAt time.Time `json:"expires_at"`
}
