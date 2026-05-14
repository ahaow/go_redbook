package service

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrInvalidLogin = errors.New("invalid login")
	ErrForbidden    = errors.New("forbidden")
)
