package dtos

import (
	"errors"
	"server/internal/models"
	"time"
)

type CreateEvent struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Attendees   []string  `json:"attendees"`
	MeetingURL  string    `json:"meeting_url"`
}

func (e *CreateEvent) Validate() error {
	if e.Title == "" {
		return errors.New("title cannot be empty")
	}

	if len(e.Description) < 10 {
		return errors.New("description must be at least 10 characters")
	}

	if e.StartTime.After(e.EndTime) {
		return errors.New("start time cannot be after end time")
	}

	if len(e.Attendees) == 0 {
		return errors.New("attendees cannot be empty")
	}

	return nil
}

func (e *CreateEvent) MapToModel(userID uint) *models.Event {
	users := make([]models.User, 0)

	return &models.Event{
		Title:       e.Title,
		Description: e.Description,
		Start:       e.StartTime,
		End:         e.EndTime,
		Attendees:   users,
		MeetingURL:  e.MeetingURL,
		CreatedBy:   userID,
	}
}
