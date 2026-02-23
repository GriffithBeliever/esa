package domain

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrForbidden        = errors.New("forbidden")
	ErrConflict         = errors.New("conflict")
	ErrInvalidInput     = errors.New("invalid input")
	ErrExpired          = errors.New("expired")
	ErrAlreadyExists    = errors.New("already exists")
	ErrInvitationClosed = errors.New("invitation already responded to")
)
