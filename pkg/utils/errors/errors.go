package errors

import "errors"

var (
	ErrEntityNotFound = errors.New("entity not found")
	ErrInternalServer = errors.New("internal server error")
	ErrBadRequest     = errors.New("bad request")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
	ErrConflict       = errors.New("conflict")
	ErrInvalidToken   = errors.New("invalid token")
	ErrCreateToken    = errors.New("error in creating token")
	ErrValidateOTP    = errors.New("failed to validate otp")
	ErrAlreadyExists  = errors.New("entity already exists")
)
