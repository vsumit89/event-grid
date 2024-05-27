package mocks

import (
	"server/internal/models"
	"time"
)

type MockEventRepository struct {
	MockCreateEvent      func(event *models.Event) (*models.Event, error)
	MockGetEventByID     func(userID, eventID uint) (*models.Event, error)
	MockGetEventByTitle  func(title string) (*models.Event, error)
	MockGetEventsInRange func(userID uint, start, end time.Time) ([]models.Event, error)
	MockUpdateEvent      func(userID, eventID uint, event *models.Event) (*models.Event, error)
	MockDeleteEvent      func(userID, eventID uint) error
}

func (m *MockEventRepository) CreateEvent(event *models.Event) (*models.Event, error) {
	return m.MockCreateEvent(event)
}

func (m *MockEventRepository) GetEventByID(userID, eventID uint) (*models.Event, error) {
	return m.MockGetEventByID(userID, eventID)
}

func (m *MockEventRepository) GetEventByTitle(title string) (*models.Event, error) {
	return m.MockGetEventByTitle(title)
}

func (m *MockEventRepository) GetEventsInRange(userID uint, start, end time.Time) ([]models.Event, error) {
	return m.MockGetEventsInRange(userID, start, end)
}

func (m *MockEventRepository) UpdateEvent(userID, eventID uint, event *models.Event) (*models.Event, error) {
	return m.MockUpdateEvent(userID, eventID, event)
}

func (m *MockEventRepository) DeleteEvent(userID, eventID uint) error {
	return m.MockDeleteEvent(userID, eventID)
}
