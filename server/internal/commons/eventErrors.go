package commons

import "errors"

var (
	ErrEventNotFound = errors.New("event not found")

	ErrInvalidEventID = errors.New("invalid event id")
)
