package mocks

import (
	"server/internal/commons"
	"server/internal/handlers/dtos"
	"server/internal/models"
)

type MockEventSvc struct {
	MockCreateEvent  func(userID uint, eventDetails *dtos.EventDTO) (*models.Event, error)
	MockGetEventByID func(userID, eventID uint) (*models.Event, error)
	MockDeleteEvent  func(userID, eventID uint) error
	MockGetEvents    func(userID uint, filters commons.EventFilters) ([]models.Event, error)
	MockUpdateEvent  func(userID, eventID uint, eventDetails *dtos.EventDTO) (*models.Event, error)
}

func (m *MockEventSvc) CreateEvent(userID uint, eventDetails *dtos.EventDTO) (*models.Event, error) {
	return m.MockCreateEvent(userID, eventDetails)
}

func (m *MockEventSvc) GetEventByID(userID, eventID uint) (*models.Event, error) {
	return m.MockGetEventByID(userID, eventID)
}

func (m *MockEventSvc) DeleteEvent(userID, eventID uint) error {
	return m.MockDeleteEvent(userID, eventID)
}

func (m *MockEventSvc) GetEvents(userID uint, filters commons.EventFilters) ([]models.Event, error) {
	return m.MockGetEvents(userID, filters)
}

func (m *MockEventSvc) UpdateEvent(userID, eventID uint, eventDetails *dtos.EventDTO) (*models.Event, error) {
	return m.MockUpdateEvent(userID, eventID, eventDetails)
}
