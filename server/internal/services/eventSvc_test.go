package services_test

import (
	"errors"
	"server/internal/commons"
	"server/internal/handlers/dtos"
	"server/internal/mocks"
	"server/internal/models"
	"server/internal/services"
	"testing"
	"time"

	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
)

func TestCreateEvent(t *testing.T) {
	mockEventRepo := &mocks.MockEventRepository{}
	mockUserRepo := &mocks.MockUserRepository{}
	mockMQ := &mocks.MockMQ{}

	eventSvc := services.NewEventSvc(&services.EventSvcOptions{
		EventRepository: mockEventRepo,
		UserRepository:  mockUserRepo,
		MQ:              mockMQ,
	})

	userID := uint(1)
	eventDetails := &dtos.EventDTO{
		Title:     "Test Event",
		StartTime: time.Now().Add(1 * time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
		Attendees: []string{"test@example.com"},
	}

	mockUser := &models.User{
		Base:  models.Base{ID: userID},
		Email: "user@example.com",
	}

	mockEvent := &models.Event{
		Base:      models.Base{ID: 1},
		Title:     eventDetails.Title,
		Start:     eventDetails.StartTime,
		End:       eventDetails.EndTime,
		CreatedBy: userID,
	}

	mockUserRepo.MockGetUserByID = func(id uint) (*models.User, error) {
		if id == userID {
			return mockUser, nil
		}
		return nil, commons.ErrUserNotFound
	}

	mockEventRepo.MockIsUserAvailableForEvent = func(userID uint, start, end time.Time) (bool, error) {
		return true, nil
	}

	mockUserRepo.MockGetUsersByEmailList = func(emails []string) ([]models.User, error) {
		return []models.User{}, nil
	}

	mockUserRepo.MockCreateUsersByEmailList = func(emails []string) ([]models.User, error) {
		return []models.User{}, nil
	}

	mockEventRepo.MockCreateEvent = func(event *models.Event) (*models.Event, error) {
		return mockEvent, nil
	}

	mockMQ.MockDeclareQueue = func(queueName string) (*amqp091.Channel, error) {
		return &amqp091.Channel{}, nil
	}

	mockMQ.MockPublishWithExchange = func(ch *amqp091.Channel, body []byte, exchangeName string) error {
		return nil
	}

	event, err := eventSvc.CreateEvent(userID, eventDetails)
	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, mockEvent.Title, event.Title)
}

func TestGetEventByID(t *testing.T) {
	mockEventRepo := &mocks.MockEventRepository{}
	mockUserRepo := &mocks.MockUserRepository{}
	mockMQ := &mocks.MockMQ{}

	eventSvc := services.NewEventSvc(&services.EventSvcOptions{
		EventRepository: mockEventRepo,
		UserRepository:  mockUserRepo,
		MQ:              mockMQ,
	})

	userID := uint(1)
	eventID := uint(1)

	mockEvent := &models.Event{
		Base:      models.Base{ID: 1},
		Title:     "Test Event",
		Start:     time.Now().Add(1 * time.Hour),
		End:       time.Now().Add(2 * time.Hour),
		CreatedBy: userID,
	}

	mockEventRepo.MockGetEventByID = func(userID, eventID uint) (*models.Event, error) {
		if eventID == mockEvent.ID {
			return mockEvent, nil
		}
		return nil, commons.ErrEventNotFound
	}

	event, err := eventSvc.GetEventByID(userID, eventID)
	assert.NoError(t, err)
	assert.NotNil(t, event)
	assert.Equal(t, mockEvent.Title, event.Title)
}

func TestDeleteEvent(t *testing.T) {
	mockEventRepo := &mocks.MockEventRepository{}
	mockUserRepo := &mocks.MockUserRepository{}
	mockMQ := &mocks.MockMQ{}

	eventSvc := services.NewEventSvc(&services.EventSvcOptions{
		EventRepository: mockEventRepo,
		UserRepository:  mockUserRepo,
		MQ:              mockMQ,
	})

	userID := uint(1)
	eventID := uint(1)

	mockEventRepo.MockDeleteEvent = func(userID, eventID uint) error {
		return nil
	}

	err := eventSvc.DeleteEvent(userID, eventID)
	assert.NoError(t, err)
}

func TestGetEvents(t *testing.T) {
	mockEventRepo := &mocks.MockEventRepository{}
	mockUserRepo := &mocks.MockUserRepository{}
	mockMQ := &mocks.MockMQ{}

	eventSvc := services.NewEventSvc(&services.EventSvcOptions{
		EventRepository: mockEventRepo,
		UserRepository:  mockUserRepo,
		MQ:              mockMQ,
	})

	userID := uint(1)
	filters := commons.EventFilters{
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
	}

	mockEvents := []models.Event{
		{
			Base:      models.Base{ID: 1},
			Title:     "Event 1",
			Start:     time.Now().Add(1 * time.Hour),
			End:       time.Now().Add(2 * time.Hour),
			CreatedBy: userID,
		},
		{
			Base:      models.Base{ID: 2},
			Title:     "Event 2",
			Start:     time.Now().Add(3 * time.Hour),
			End:       time.Now().Add(4 * time.Hour),
			CreatedBy: userID,
		},
	}

	mockEventRepo.MockGetEventsInRange = func(userID uint, start, end time.Time) ([]models.Event, error) {
		return mockEvents, nil
	}

	events, err := eventSvc.GetEvents(userID, filters)
	assert.NoError(t, err)
	assert.NotNil(t, events)
	assert.Len(t, events, 2)
}

func TestUpdateEvent(t *testing.T) {
	mockEventRepo := &mocks.MockEventRepository{}
	mockUserRepo := &mocks.MockUserRepository{}
	mockMQ := &mocks.MockMQ{}

	eventSvc := services.NewEventSvc(&services.EventSvcOptions{
		EventRepository: mockEventRepo,
		UserRepository:  mockUserRepo,
		MQ:              mockMQ,
	})

	userID := uint(1)
	eventID := uint(1)

	eventDetails := &dtos.EventDTO{
		Title:     "Updated Event",
		StartTime: time.Now().Add(1 * time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
		Attendees: []string{"test@example.com"},
	}

	mockEvent := &models.Event{
		Base: models.Base{
			ID: eventID,
		},
		Title:     eventDetails.Title,
		Start:     eventDetails.StartTime,
		End:       eventDetails.EndTime,
		CreatedBy: userID,
	}

	mockEventRepo.MockUpdateEvent = func(userID, eventID uint, event *dtos.EventDTO) (*models.Event, error) {
		return mockEvent, nil
	}

	mockMQ.MockDeclareQueue = func(queueName string) (*amqp091.Channel, error) {
		return &amqp091.Channel{}, nil
	}

	mockMQ.MockPublishWithExchange = func(ch *amqp091.Channel, body []byte, exchangeName string) error {
		return nil
	}

	updatedEvent, err := eventSvc.UpdateEvent(userID, eventID, eventDetails)
	assert.NoError(t, err)
	assert.NotNil(t, updatedEvent)
	assert.Equal(t, mockEvent.Title, updatedEvent.Title)
}

func TestCreateEvent_FailureCases(t *testing.T) {
	mockEventRepo := &mocks.MockEventRepository{}
	mockUserRepo := &mocks.MockUserRepository{}
	mockMQ := &mocks.MockMQ{}

	eventSvc := services.NewEventSvc(&services.EventSvcOptions{
		EventRepository: mockEventRepo,
		UserRepository:  mockUserRepo,
		MQ:              mockMQ,
	})

	userID := uint(1)
	eventDetails := &dtos.EventDTO{
		Title:     "Test Event",
		StartTime: time.Now().Add(1 * time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
		Attendees: []string{"test@example.com"},
	}

	t.Run("User not found", func(t *testing.T) {
		mockUserRepo.MockGetUserByID = func(id uint) (*models.User, error) {
			return nil, commons.ErrUserNotFound
		}

		event, err := eventSvc.CreateEvent(userID, eventDetails)
		assert.Error(t, err)
		assert.Nil(t, event)
		assert.Equal(t, commons.ErrUserNotFound, err)
	})

	t.Run("User not available for event", func(t *testing.T) {
		mockUser := &models.User{
			Base:  models.Base{ID: userID},
			Email: "user@example.com",
		}

		mockUserRepo.MockGetUserByID = func(id uint) (*models.User, error) {
			return mockUser, nil
		}

		mockEventRepo.MockIsUserAvailableForEvent = func(userID uint, start, end time.Time) (bool, error) {
			return false, nil
		}

		event, err := eventSvc.CreateEvent(userID, eventDetails)
		assert.Error(t, err)
		assert.Nil(t, event)
		assert.Equal(t, commons.ErrInvalidEventTime, err)
	})

	t.Run("Error creating event", func(t *testing.T) {
		mockUser := &models.User{
			Base:  models.Base{ID: userID},
			Email: "user@example.com",
		}

		mockUserRepo.MockGetUserByID = func(id uint) (*models.User, error) {
			return mockUser, nil
		}

		mockEventRepo.MockIsUserAvailableForEvent = func(userID uint, start, end time.Time) (bool, error) {
			return true, nil
		}

		mockUserRepo.MockGetUsersByEmailList = func(emails []string) ([]models.User, error) {
			return []models.User{}, nil
		}

		mockUserRepo.MockCreateUsersByEmailList = func(emails []string) ([]models.User, error) {
			return []models.User{}, nil
		}

		mockEventRepo.MockCreateEvent = func(event *models.Event) (*models.Event, error) {
			return nil, errors.New("error creating event")
		}

		event, err := eventSvc.CreateEvent(userID, eventDetails)
		assert.Error(t, err)
		assert.Nil(t, event)
		assert.Equal(t, "error creating event", err.Error())
	})
}

func TestGetEventByID_FailureCases(t *testing.T) {
	mockEventRepo := &mocks.MockEventRepository{}
	mockUserRepo := &mocks.MockUserRepository{}
	mockMQ := &mocks.MockMQ{}

	eventSvc := services.NewEventSvc(&services.EventSvcOptions{
		EventRepository: mockEventRepo,
		UserRepository:  mockUserRepo,
		MQ:              mockMQ,
	})

	userID := uint(1)
	eventID := uint(1)

	t.Run("Event not found", func(t *testing.T) {
		mockEventRepo.MockGetEventByID = func(userID, eventID uint) (*models.Event, error) {
			return nil, commons.ErrEventNotFound
		}

		event, err := eventSvc.GetEventByID(userID, eventID)
		assert.Error(t, err)
		assert.Nil(t, event)
		assert.Equal(t, commons.ErrEventNotFound, err)
	})
}

func TestDeleteEvent_FailureCases(t *testing.T) {
	mockEventRepo := &mocks.MockEventRepository{}
	mockUserRepo := &mocks.MockUserRepository{}
	mockMQ := &mocks.MockMQ{}

	eventSvc := services.NewEventSvc(&services.EventSvcOptions{
		EventRepository: mockEventRepo,
		UserRepository:  mockUserRepo,
		MQ:              mockMQ,
	})

	userID := uint(1)
	eventID := uint(1)

	t.Run("Error deleting event", func(t *testing.T) {
		mockEventRepo.MockDeleteEvent = func(userID, eventID uint) error {
			return errors.New("error deleting event")
		}

		err := eventSvc.DeleteEvent(userID, eventID)
		assert.Error(t, err)
		assert.Equal(t, "error deleting event", err.Error())
	})
}

func TestGetEvents_FailureCases(t *testing.T) {
	mockEventRepo := &mocks.MockEventRepository{}
	mockUserRepo := &mocks.MockUserRepository{}
	mockMQ := &mocks.MockMQ{}

	eventSvc := services.NewEventSvc(&services.EventSvcOptions{
		EventRepository: mockEventRepo,
		UserRepository:  mockUserRepo,
		MQ:              mockMQ,
	})

	userID := uint(1)
	filters := commons.EventFilters{
		StartTime: time.Now(),
		EndTime:   time.Now().Add(24 * time.Hour),
	}

	t.Run("Error fetching events", func(t *testing.T) {
		mockEventRepo.MockGetEventsInRange = func(userID uint, start, end time.Time) ([]models.Event, error) {
			return nil, errors.New("error fetching events")
		}

		events, err := eventSvc.GetEvents(userID, filters)
		assert.Error(t, err)
		assert.Nil(t, events)
		assert.Equal(t, "error fetching events", err.Error())
	})
}

func TestUpdateEvent_FailureCases(t *testing.T) {
	mockEventRepo := &mocks.MockEventRepository{}
	mockUserRepo := &mocks.MockUserRepository{}
	mockMQ := &mocks.MockMQ{}

	eventSvc := services.NewEventSvc(&services.EventSvcOptions{
		EventRepository: mockEventRepo,
		UserRepository:  mockUserRepo,
		MQ:              mockMQ,
	})

	userID := uint(1)
	eventID := uint(1)

	eventDetails := &dtos.EventDTO{
		Title:     "Updated Event",
		StartTime: time.Now().Add(1 * time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
		Attendees: []string{"test@example.com"},
	}

	t.Run("Error updating event", func(t *testing.T) {
		mockEventRepo.MockUpdateEvent = func(userID, eventID uint, event *dtos.EventDTO) (*models.Event, error) {
			return nil, errors.New("error updating event")
		}

		updatedEvent, err := eventSvc.UpdateEvent(userID, eventID, eventDetails)
		assert.Error(t, err)
		assert.Nil(t, updatedEvent)
		assert.Equal(t, "error updating event", err.Error())
	})
}
