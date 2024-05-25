package services

import "server/internal/models"

type IUserSvc interface {
	CreateUser(user *models.User) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
}
