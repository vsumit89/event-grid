package services

import (
	"server/internal/handlers/dtos"
	"server/internal/models"
	"server/internal/repository"
)

type IUserSvc interface {
	RegisterUser(user *dtos.CreateUserReqeust) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
}

type UserSvcOptions struct {
	Repository repository.IUserRepository
}

func NewUserSvc(options *UserSvcOptions) IUserSvc {
	return &userSvcImpl{
		userRepo: options.Repository,
	}
}

type userSvcImpl struct {
	userRepo repository.IUserRepository
}

func (u *userSvcImpl) RegisterUser(user *dtos.CreateUserReqeust) (*models.User, error) {
	return nil, nil
}

func (u *userSvcImpl) GetUserByID(id uint) (*models.User, error) {
	return nil, nil
}
