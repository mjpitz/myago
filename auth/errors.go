package auth

import (
	"errors"
)

var (
	// ErrUnauthorized is returned when no user information is found on a context.
	ErrUnauthorized = errors.New("unauthorized")
)
