package repository

import (
	"errors"
	"server/internal/commons"
	db "server/internal/infrastructure/database"
	"server/internal/models"

	"gorm.io/gorm"
)

type IEventsRepository interface {
	CreateEvent(event *models.Event) (*models.Event, error)
	GetEventByID(id uint) (*models.Event, error)
	GetEventByTitle(title string) (*models.Event, error)
	UpdateEvent(id uint, event *models.Event) (*models.Event, error)
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

func (e *eventsPgRepoImpl) GetEventByID(id uint) (*models.Event, error) {
	var event models.Event
	err := e.db.Where("id = ?", id).First(&event).Error
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

func (e *eventsPgRepoImpl) DeleteEvent(id uint) error {
	return e.db.Where("id = ?", id).Delete(&models.Event{}).Error
}
