package auth_errors

import "errors"

var (
	ErrUserDoesNotExist  = errors.New("user does not exist")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidPassword   = errors.New("Invalid password provided")
)
