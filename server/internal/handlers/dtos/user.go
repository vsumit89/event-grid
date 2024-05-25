package dtos

import (
	"fmt"
	"server/internal/models"
	"server/pkg/utils"
)

type CreateUserReqeust struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *CreateUserReqeust) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}

	if c.Password == "" {
		return fmt.Errorf("password is required")
	}

	if !utils.ValidateEmail(c.Email) {
		return fmt.Errorf("invalid email")
	}

	return nil
}

func (c CreateUserReqeust) MapToUserModel() *models.User {
	return &models.User{
		Name:     c.Name,
		Email:    c.Email,
		Password: c.Password,
	}
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresAt    string `json:"expires_at"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (l *LoginRequest) Validate() error {
	if l.Email == "" {
		return fmt.Errorf("email is required")
	}

	if l.Password == "" {
		return fmt.Errorf("password is required")
	}

	return nil
}
