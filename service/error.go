package service

import "errors"

var (
	ErrUserNotFound         = errors.New("User Not Found")
	ErrEmailExists          = errors.New("Email Already Exists")
	ErrInvalidEmailPassword = errors.New("Invalid Email or Password")
	ErrDataNotFound         = errors.New("Data Not Found")
	ErrInternalServer       = errors.New("Failed Processing Request")
	ErrGenerateToken        = errors.New("Failed Generate Token")
	ErrForbidden            = errors.New("Forbidden")
	ErrClosed               = errors.New("Job Already Closed")
)
