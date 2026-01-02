package errors

import (
	"errors"
)

// Auth errors
var (
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidStore       = errors.New("store is not a PGXStore")
)
