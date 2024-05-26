package repository

import (
	"errors"
	"server/internal/commons"
	db "server/internal/infrastructure/database"
	"server/internal/models"
	"time"

	"gorm.io/gorm"
)

type IEventsRepository interface {
	CreateEvent(event *models.Event) (*models.Event, error)
	GetEventByID(userID, eventID uint) (*models.Event, error)
	GetEventByTitle(title string) (*models.Event, error)
	GetEventsInRange(userID uint, start, end time.Time) ([]models.Event, error)
	UpdateEvent(id uint, event *models.Event) (*models.Event, error)
	DeleteEvent(userID, eventID uint) error
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

func (e *eventsPgRepoImpl) UpdateEvent(id uint, event *models.Event) (*models.Event, error) {
	err := e.db.Model(&models.Event{}).Where("id = ?", id).Updates(event).Error
	if err != nil {
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

	err := e.db.Model(&models.Event{}).Where("start BETWEEN ? AND ? AND ? IN (SELECT user_id FROM event_attendees WHERE event_attendees.event_id = events.id)", start, end, userID).Find(&events).Error
	if err != nil {
		return nil, err
	}

	return events, nil
}
