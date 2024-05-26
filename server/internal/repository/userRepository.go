package repository

import (
	"errors"
	"server/internal/commons"
	db "server/internal/infrastructure/database"
	"server/internal/models"

	"gorm.io/gorm"
)

// IUserRepository is an interface for the user repository
// it defines the methods that the user repository should implement
// which helps other packages to interact with the user repository
type IUserRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	GetUserByEmail(email string) (*models.User, error)
	UpdateUser(id uint, user *models.User) (*models.User, error)

	GetUsersByEmailList([]string) ([]models.User, error)

	CreateUsersByEmailList([]string) ([]models.User, error)

	GetUsers(userID uint, query string, limit int) ([]models.User, error)
}

type userPgRepoImpl struct {
	db *gorm.DB
}

func NewUserRepository(dbSvc db.IDatabase) IUserRepository {
	return &userPgRepoImpl{
		db: dbSvc.GetClient().(*gorm.DB),
	}
}

func (u *userPgRepoImpl) CreateUser(user *models.User) (*models.User, error) {
	err := u.db.Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userPgRepoImpl) GetUserByID(id uint) (*models.User, error) {
	var user models.User
	err := u.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *userPgRepoImpl) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := u.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, commons.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (u *userPgRepoImpl) UpdateUser(id uint, user *models.User) (*models.User, error) {
	err := u.db.Model(&models.User{}).Where("id = ?", id).Updates(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userPgRepoImpl) GetUsersByEmailList(emails []string) ([]models.User, error) {
	var users []models.User
	err := u.db.Where("email IN ?", emails).Find(&users).Error
	if err != nil {
		return nil, commons.ErrUserNotFound
	}

	return users, nil
}

func (u *userPgRepoImpl) CreateUsersByEmailList(emails []string) ([]models.User, error) {
	users := make([]models.User, 0)

	for _, email := range emails {
		users = append(users, models.User{Email: email})
	}

	err := u.db.Create(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *userPgRepoImpl) GetUsers(userID uint, query string, limit int) ([]models.User, error) {
	var users []models.User
	err := u.db.Where("email LIKE ? OR name LIKE ? AND id <> ?", "%"+query+"%", "%"+query+"%", userID).Order("name ASC").Order("email ASC").Limit(limit).Find(&users).Error
	if err != nil {
		return nil, commons.ErrUserNotFound
	}

	return users, nil
}
