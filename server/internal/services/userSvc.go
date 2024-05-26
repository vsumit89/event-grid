package services

import (
	"errors"
	"server/internal/commons"
	"server/internal/models"
	"server/internal/repository"
	"server/pkg/utils"
)

type IUserSvc interface {
	RegisterUser(user *models.User) (*models.User, error)
	Login(email, password string) (*models.User, error)
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

func (u *userSvcImpl) RegisterUser(user *models.User) (*models.User, error) {
	var err error

	userFromDB, err := u.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		if !errors.Is(err, commons.ErrUserNotFound) {
			return nil, err
		}
	}

	if userFromDB != nil {
		if userFromDB.Password != "" {
			return nil, commons.ErrUserAlreadyExists
		}

		user.Password, err = utils.HashPassword(user.Password)
		if err != nil {
			return nil, err
		}

		user, err = u.userRepo.UpdateUser(userFromDB.ID, user)
		if err != nil {
			return nil, err
		}

	} else {
		user.Password, err = utils.HashPassword(user.Password)
		if err != nil {
			return nil, err
		}

		user, err = u.userRepo.CreateUser(user)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (u *userSvcImpl) Login(email, password string) (*models.User, error) {
	user, err := u.userRepo.GetUserByEmail(email)
	if err != nil {
		if err == commons.ErrUserNotFound {
			return nil, commons.ErrInvalidEmailPassword
		}
		return nil, commons.ErrInternalServer
	}

	if err := utils.VerifyPassword(password, user.Password); err != nil {
		return nil, commons.ErrInvalidEmailPassword
	}

	return user, nil
}

func (u *userSvcImpl) GetUserByID(id uint) (*models.User, error) {
	user, err := u.userRepo.GetUserByID(id)
	if err != nil {
		if err == commons.ErrUserNotFound {
			return nil, commons.ErrUserNotFound
		}
		return nil, commons.ErrInternalServer
	}
	return user, nil
}
