package mocks

import "server/internal/models"

type MockUserRepository struct {
	MockCreateUser             func(user *models.User) (*models.User, error)
	MockGetUserByID            func(id uint) (*models.User, error)
	MockGetUserByEmail         func(email string) (*models.User, error)
	MockUpdateUser             func(id uint, user *models.User) (*models.User, error)
	MockGetUsers               func(userID uint, query string, limit int) ([]models.User, error)
	MockGetUsersByEmailList    func([]string) ([]models.User, error)
	MockCreateUsersByEmailList func([]string) ([]models.User, error)
}

func (m *MockUserRepository) CreateUser(user *models.User) (*models.User, error) {
	return m.MockCreateUser(user)
}

func (m *MockUserRepository) GetUserByID(id uint) (*models.User, error) {
	return m.MockGetUserByID(id)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*models.User, error) {
	return m.MockGetUserByEmail(email)
}

func (m *MockUserRepository) UpdateUser(id uint, user *models.User) (*models.User, error) {
	return m.MockUpdateUser(id, user)
}

func (m *MockUserRepository) GetUsers(userID uint, query string, limit int) ([]models.User, error) {
	return m.MockGetUsers(userID, query, limit)
}

func (m *MockUserRepository) GetUsersByEmailList(emails []string) ([]models.User, error) {
	return m.MockGetUsersByEmailList(emails)
}

func (m *MockUserRepository) CreateUsersByEmailList(emails []string) ([]models.User, error) {
	return m.MockCreateUsersByEmailList(emails)
}
