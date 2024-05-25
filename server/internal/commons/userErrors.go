package commons

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrInvalidEmailPassword = errors.New("invalid email or password")

	ErrUserAlreadyExists = errors.New("user already exists")
)
