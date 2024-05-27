package models

import (
	"time"
)

type Event struct {
	Base
	CreatedBy   uint      `json:"created_by,omitempty"`
	Title       string    `json:"title,omitempty"`
	Description string    `json:"description,omitempty"`
	Start       time.Time `json:"start,omitempty"`
	End         time.Time `json:"end,omitempty"`
	Attendees   []User    `gorm:"many2many:event_attendees" json:"attendees,omitempty"`
	MeetingURL  string    `json:"meeting_url,omitempty"`
}
