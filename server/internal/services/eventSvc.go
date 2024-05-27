package services

import (
	"errors"
	"server/internal/commons"
	"server/internal/handlers/dtos"
	"server/internal/infrastructure/mq"
	"server/internal/models"
	"server/internal/repository"
	"time"
)

type EventFilters struct {
	StartTime time.Time
	EndTime   time.Time
}

type IEventSvc interface {
	CreateEvent(userID uint, eventDetails *dtos.EventDTO) (*models.Event, error)
	GetEventByID(userID, eventID uint) (*models.Event, error)
	DeleteEvent(userID, eventID uint) error
	GetEvents(userID uint, filters EventFilters) ([]models.Event, error)
	UpdateEvent(userID, eventID uint, eventDetails *dtos.EventDTO) (*models.Event, error)
}

type EventSvcOptions struct {
	EventRepository repository.IEventsRepository
	UserRepository  repository.IUserRepository
	MQ              mq.IMessageQueue
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

func (e *eventSvcImpl) CreateEvent(userID uint, eventDetails *dtos.EventDTO) (*models.Event, error) {
	event := eventDetails.MapToModel(userID)

	currentUser, err := e.UserRepository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	isAvailable, err := e.EventRepository.IsUserAvailableForEvent(userID, event.Start, event.End)
	if err != nil {
		return nil, err
	}

	if !isAvailable {
		return nil, commons.ErrInvalidEventTime
	}

	// gets all the users who already exists based on their emails
	users, err := e.UserRepository.GetUsersByEmailList(eventDetails.Attendees)
	if err != nil {
		if !errors.Is(err, commons.ErrUserNotFound) {
			return nil, err
		}
	}

	users = append(users, *currentUser)
	usersToBeCreated := make([]string, 0)

	userEmails := make(map[string]bool)

	// Create a map of existing user emails
	for _, user := range users {
		userEmails[user.Email] = true
	}

	for i := 0; i < len(eventDetails.Attendees); i++ {
		if !userEmails[eventDetails.Attendees[i]] {
			usersToBeCreated = append(usersToBeCreated, eventDetails.Attendees[i])
		}
	}

	var newUsers []models.User
	if len(usersToBeCreated) > 0 {
		newUsers, err = e.UserRepository.CreateUsersByEmailList(usersToBeCreated)
		if err != nil {
			return nil, err
		}
	}

	event.Attendees = append(users, newUsers...)

	event, err = e.EventRepository.CreateEvent(event)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func (e *eventSvcImpl) GetEventByID(userID, eventID uint) (*models.Event, error) {
	return e.EventRepository.GetEventByID(userID, eventID)
}

func (e *eventSvcImpl) DeleteEvent(userID, eventID uint) error {
	return e.EventRepository.DeleteEvent(userID, eventID)
}

func (e *eventSvcImpl) GetEvents(userID uint, filters EventFilters) ([]models.Event, error) {
	events, err := e.EventRepository.GetEventsInRange(userID, filters.StartTime, filters.EndTime)
	if err != nil {
		return nil, err
	}
	return events, nil
}

func (e *eventSvcImpl) UpdateEvent(userID, eventID uint, eventDetails *dtos.EventDTO) (*models.Event, error) {
	return e.EventRepository.UpdateEvent(userID, eventID, eventDetails)
}
