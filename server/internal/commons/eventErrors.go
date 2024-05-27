package commons

import "errors"

var (
	ErrEventNotFound = errors.New("event not found")

	ErrInvalidEventID = errors.New("invalid event id")

	ErrInvalidEventTime = errors.New("user is part of the other event at that time")
)
