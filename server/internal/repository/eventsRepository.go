package repository

import (
	"errors"
	"fmt"
	"server/internal/commons"
	"server/internal/handlers/dtos"
	db "server/internal/infrastructure/database"
	"server/internal/models"
	"server/pkg/logger"
	"time"

	"gorm.io/gorm"
)

type IEventsRepository interface {
	CreateEvent(event *models.Event) (*models.Event, error)
	GetEventByID(userID, eventID uint) (*models.Event, error)
	GetEventByTitle(title string) (*models.Event, error)
	GetEventsInRange(userID uint, start, end time.Time) ([]models.Event, error)
	UpdateEvent(userID, eventID uint, event *dtos.EventDTO) (*models.Event, error)
	DeleteEvent(userID, eventID uint) error

	IsUserAvailableForEvent(userID uint, start, end time.Time) (bool, error)
}

type eventsPgRepoImpl struct {
	db *gorm.DB
}

func NewEventsRepository(dbSvc db.IDatabase) IEventsRepository {
	return &eventsPgRepoImpl{
		db: dbSvc.GetClient().(*gorm.DB),
	}
}

func (e *eventsPgRepoImpl) CreateEvent(event *models.Event) (*models.Event, error) {
	err := e.db.Create(event).Error
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (e *eventsPgRepoImpl) GetEventByID(userID, eventID uint) (*models.Event, error) {
	var event models.Event
	err := e.db.Where("id = ? and created_by = ?", eventID, userID).Preload("Attendees").First(&event).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, commons.ErrEventNotFound
		}
		return nil, err
	}

	return &event, nil
}

func (e *eventsPgRepoImpl) GetEventByTitle(title string) (*models.Event, error) {
	var event models.Event
	err := e.db.Where("title = ?", title).First(&event).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, commons.ErrEventNotFound
		}
		return nil, err
	}

	return &event, nil
}

func (e *eventsPgRepoImpl) UpdateEvent(userID, eventID uint, eventDetails *dtos.EventDTO) (*models.Event, error) {
	logger.Info("updating event attendees", "attendees", eventDetails.Attendees, "len", len(eventDetails.Attendees))

	event, err := e.GetEventByID(userID, eventID)
	if err != nil {
		return nil, err
	}

	if time.Now().After(event.Start) {
		return nil, fmt.Errorf("user can't update an event that has already started")
	}

	event = &models.Event{
		Base:        models.Base{ID: eventID},
		CreatedBy:   userID,
		Title:       eventDetails.Title,
		Description: eventDetails.Description,
		Start:       eventDetails.StartTime,
		End:         eventDetails.EndTime,
		MeetingURL:  eventDetails.MeetingURL,
	}

	var existingAttendees []models.User
	err = e.db.Model(event).Association("Attendees").Find(&existingAttendees)
	if err != nil {
		return nil, err
	}

	existingAttendeeIDs := make(map[uint]bool)
	for _, attendee := range existingAttendees {
		existingAttendeeIDs[attendee.ID] = true
	}

	for _, email := range eventDetails.Attendees {
		// Check if the attendee already exists
		var attendee models.User
		if err := e.db.Where("email = ?", email).First(&attendee).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				attendee = models.User{Email: email}
				err = e.db.Create(&attendee).Error
				if err != nil {
					return nil, err
				}
			}
		}

		logger.Info("updating event attendee", "attendee", attendee)

		// Add attendee to the event if not already added
		if _, ok := existingAttendeeIDs[attendee.ID]; !ok {
			logger.Info("adding event attendee", "attendee", attendee)

			err := e.db.Model(event).Association("Attendees").Append(&attendee)
			if err != nil {
				return nil, err
			}
		}

		// Remove attendee ID from existing list (for attendees already associated with event)
		delete(existingAttendeeIDs, attendee.ID)
	}

	// this will make sure that we are not removing the user who created the event
	delete(existingAttendeeIDs, userID)

	for attendeeID := range existingAttendeeIDs {
		var attendee models.User
		err = e.db.First(&attendee, attendeeID).Error
		if err != nil {
			return nil, err
		}

		err = e.db.Model(event).Association("Attendees").Delete(attendee)
		if err != nil {
			return nil, err
		}
	}

	event.Attendees = nil
	event.CreatedAt = time.Now()

	if err := e.db.Save(event).Error; err != nil {
		return nil, err
	}

	return event, nil
}

func (e *eventsPgRepoImpl) DeleteEvent(userID, eventID uint) error {
	return e.db.Where("id = ? and created_by = ?", eventID, userID).Delete(&models.Event{}).Error
}

func (e *eventsPgRepoImpl) GetEventsByUserID(userID uint) ([]models.Event, error) {
	var events []models.Event
	err := e.db.Where("created_by = ?", userID).Preload("Attendees").Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}

func (e *eventsPgRepoImpl) GetEventsInRange(userID uint, start, end time.Time) ([]models.Event, error) {
	var events []models.Event

	startTime := start.Format("2006-01-02 15:04:05") + "+05:30"
	endTime := end.Format("2006-01-02 15:04:05") + "+05:30"

	result := e.db.Model(&models.Event{}).Where("start >= ? AND start <= ? AND ? IN (SELECT user_id FROM event_attendees WHERE event_attendees.event_id = events.id)", startTime, endTime, userID).Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}

	return events, nil
}

func (e *eventsPgRepoImpl) IsUserAvailableForEvent(userID uint, start, end time.Time) (bool, error) {
	events, err := e.GetEventsInRange(userID, start, end)
	if err != nil {
		return false, err
	}

	if len(events) > 0 {
		return false, nil
	}

	return true, nil
}
