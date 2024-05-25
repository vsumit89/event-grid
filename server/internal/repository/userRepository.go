package repository

import (
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
