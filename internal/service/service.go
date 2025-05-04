package service

import "errors"

var (
	ErrNoChanges        = errors.New("no changes")
	ErrProfilesNotFound = errors.New("profiles not found")
	ErrProfileNotFound  = errors.New("profile not found")
)
