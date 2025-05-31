package errors

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrRoleNotFound       = errors.New("role not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrInvalidToken       = errors.New("invalid token")
	ErrForbidden          = errors.New("forbidden")
	ErrNotFound           = errors.New("not found")
)
