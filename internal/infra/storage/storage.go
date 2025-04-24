package storage

import "errors"

var (
	ErrMemberNotFound    = errors.New("user not found")
	ErrEventNotFound     = errors.New("event not found")
	ErrRegNotFound       = errors.New("registration not found")
	ErrInfoNotFound      = errors.New("orchestra info not found")
	ErrLocationNotFound  = errors.New("location not found")
	ErrEventTypeNotFound = errors.New("event type not found")
	ErrMemberExists      = errors.New("member already exists")
	ErrRegAlreadyExists  = errors.New("registration already exists")
	ErrEventFull         = errors.New("event full")
	ErrEmailDuplicate    = errors.New("email already exists")
	ErrPhoneDuplicate    = errors.New("phone number already exists")
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrInvalidPhone      = errors.New("invalid phone number format")
)
