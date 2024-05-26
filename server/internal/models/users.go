package models

type User struct {
	Base
	Name     string  `gorm:"type:varchar(50)" json:"name"`
	Email    string  `gorm:"not null" json:"email"`
	Password string  `json:"password,omitempty"`
	Events   []Event `gorm:"many2many:event_attendees" json:"events"`
}
