package services

import (
	"server/internal/handlers/dtos"
	"server/internal/models"
	"server/internal/repository"
)

type IEventSvc interface {
	CreateEvent(eventDetails *dtos.CreateEvent) (*models.Event, error)
}

type EventSvcOptions struct {
	EventRepository repository.IEventsRepository
	UserRepository  repository.IUserRepository
}

type eventSvcImpl struct {
	EventRepository repository.IEventsRepository
	UserRepository  repository.IUserRepository
}

func NewEventSvc(opts *EventSvcOptions) IEventSvc {
	return &eventSvcImpl{
		EventRepository: opts.EventRepository,
		UserRepository:  opts.UserRepository,
	}
}

func (e *eventSvcImpl) CreateEvent(eventDetails *dtos.CreateEvent) (*models.Event, error) {
	

	return nil, nil
}
