package services_test

import (
	"server/internal/commons"
	"server/internal/mocks"
	"server/internal/models"
	"server/internal/services"
	"server/pkg/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterUser(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{}
	userSvc := services.NewUserSvc(&services.UserSvcOptions{
		Repository: mockUserRepo,
	})

	user := &models.User{
		Email:    "test@example.com",
		Password: "password123",
	}

	t.Run("User already exists with password", func(t *testing.T) {
		mockUserRepo.MockGetUserByEmail = func(email string) (*models.User, error) {
			return &models.User{Email: email, Password: "hashedpassword"}, nil
		}

		createdUser, err := userSvc.RegisterUser(user)
		assert.Error(t, err)
		assert.Nil(t, createdUser)
		assert.Equal(t, commons.ErrUserAlreadyExists, err)
	})

	t.Run("User exists without password", func(t *testing.T) {
		mockUserRepo.MockGetUserByEmail = func(email string) (*models.User, error) {
			return &models.User{Email: email}, nil
		}

		mockUserRepo.MockUpdateUser = func(id uint, user *models.User) (*models.User, error) {
			user.ID = id
			return user, nil
		}

		createdUser, err := userSvc.RegisterUser(user)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, user.Email, createdUser.Email)
	})

	t.Run("New user registration", func(t *testing.T) {
		mockUserRepo.MockGetUserByEmail = func(email string) (*models.User, error) {
			return nil, commons.ErrUserNotFound
		}

		mockUserRepo.MockCreateUser = func(user *models.User) (*models.User, error) {
			user.ID = 1
			return user, nil
		}

		createdUser, err := userSvc.RegisterUser(user)
		assert.NoError(t, err)
		assert.NotNil(t, createdUser)
		assert.Equal(t, uint(1), createdUser.ID)
		assert.Equal(t, user.Email, createdUser.Email)
	})
}

func TestLogin(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{}
	userSvc := services.NewUserSvc(&services.UserSvcOptions{
		Repository: mockUserRepo,
	})

	email := "test@example.com"
	password := "password123"
	hashedPassword, _ := utils.HashPassword(password)

	t.Run("User not found", func(t *testing.T) {
		mockUserRepo.MockGetUserByEmail = func(email string) (*models.User, error) {
			return nil, commons.ErrUserNotFound
		}

		loggedInUser, err := userSvc.Login(email, password)
		assert.Error(t, err)
		assert.Nil(t, loggedInUser)
		assert.Equal(t, commons.ErrInvalidEmailPassword, err)
	})

	t.Run("Invalid password", func(t *testing.T) {
		mockUserRepo.MockGetUserByEmail = func(email string) (*models.User, error) {
			return &models.User{Email: email, Password: hashedPassword}, nil
		}

		loggedInUser, err := userSvc.Login(email, "wrongpassword")
		assert.Error(t, err)
		assert.Nil(t, loggedInUser)
		assert.Equal(t, commons.ErrInvalidEmailPassword, err)
	})

	t.Run("Successful login", func(t *testing.T) {
		mockUserRepo.MockGetUserByEmail = func(email string) (*models.User, error) {
			return &models.User{Email: email, Password: hashedPassword}, nil
		}

		loggedInUser, err := userSvc.Login(email, password)
		assert.NoError(t, err)
		assert.NotNil(t, loggedInUser)
		assert.Equal(t, email, loggedInUser.Email)
	})
}

func TestGetUserByID(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{}
	userSvc := services.NewUserSvc(&services.UserSvcOptions{
		Repository: mockUserRepo,
	})

	userID := uint(1)

	t.Run("User not found", func(t *testing.T) {
		mockUserRepo.MockGetUserByID = func(id uint) (*models.User, error) {
			return nil, commons.ErrUserNotFound
		}

		user, err := userSvc.GetUserByID(userID)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, commons.ErrUserNotFound, err)
	})

	t.Run("Internal server error", func(t *testing.T) {
		mockUserRepo.MockGetUserByID = func(id uint) (*models.User, error) {
			return nil, commons.ErrInternalServer
		}

		user, err := userSvc.GetUserByID(userID)
		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, commons.ErrInternalServer, err)
	})

	t.Run("Successful fetch", func(t *testing.T) {
		mockUser := &models.User{
			Base:  models.Base{ID: userID},
			Email: "test@example.com",
		}

		mockUserRepo.MockGetUserByID = func(id uint) (*models.User, error) {
			return mockUser, nil
		}

		user, err := userSvc.GetUserByID(userID)
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, mockUser, user)
	})
}

func TestGetUsers(t *testing.T) {
	mockUserRepo := &mocks.MockUserRepository{}
	userSvc := services.NewUserSvc(&services.UserSvcOptions{
		Repository: mockUserRepo,
	})

	userID := uint(1)
	query := "test"
	limit := 10

	t.Run("Successful fetch", func(t *testing.T) {
		mockUsers := []models.User{
			{Base: models.Base{ID: 1}, Email: "test1@example.com"},
			{Base: models.Base{ID: 2}, Email: "test2@example.com"},
		}

		mockUserRepo.MockGetUsers = func(userID uint, query string, limit int) ([]models.User, error) {
			return mockUsers, nil
		}

		users, err := userSvc.GetUsers(userID, query, limit)
		assert.NoError(t, err)
		assert.NotNil(t, users)
		assert.Equal(t, 2, len(users))
	})
}
