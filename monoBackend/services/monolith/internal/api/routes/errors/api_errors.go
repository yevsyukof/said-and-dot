package api_errors

import "errors"

var (
	ErrUserDoesNotExist  = errors.New("user does not exist")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidPassword   = errors.New("invalid password provided")
	ErrDatabaseError     = errors.New("internal server error")
)
