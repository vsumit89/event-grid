package models

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Start       time.Time `json:"start,omitempty"`
	End         time.Time `json:"end,omitempty"`
	Attendees   []User    `gorm:"many2many:event_attendees" json:"attendees,omitempty"`
	CreatedBy   uint      `json:"created_by,omitempty"`
}

type EventAttendees struct {
	gorm.Model
	EventID uint
	UserID  uint
}
