package domain

import "errors"

var (
	ErrNotFound      = errors.New("resource not found")
	ErrInternal      = errors.New("internal server error")
	ErrInvalidInput  = errors.New("invalid input data")
	ErrThirdParty    = errors.New("external service error")
	ErrAlreadyExists = errors.New("record already exists")
)
