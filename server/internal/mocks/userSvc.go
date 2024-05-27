package mocks

import "server/internal/models"

type MockUserSvc struct {
	MockRegisterUser func(user *models.User) (*models.User, error)
	MockGetUserById  func(id uint) (*models.User, error)
	MockLogin        func(email, password string) (*models.User, error)
	MockGetUsers     func(userID uint, query string, limit int) ([]models.User, error)
}

func (m *MockUserSvc) RegisterUser(user *models.User) (*models.User, error) {
	return m.MockRegisterUser(user)
}

func (m *MockUserSvc) GetUserById(id uint) (*models.User, error) {
	return m.MockGetUserById(id)
}

func (m *MockUserSvc) Login(email, password string) (*models.User, error) {
	return m.MockLogin(email, password)
}

func (m *MockUserSvc) GetUsers(userID uint, query string, limit int) ([]models.User, error) {
	return m.MockGetUsers(userID, query, limit)
}
