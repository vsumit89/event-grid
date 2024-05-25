package dtos

import "fmt"

type CreateUserReqeust struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (c *CreateUserReqeust) Validate() error {
	if c.Name == "" {
		return fmt.Errorf("name is required")
	}
	if c.Email == "" {
		return fmt.Errorf("email is required")
	}
	if c.Password == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}
