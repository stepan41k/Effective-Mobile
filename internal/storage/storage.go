package storage

import "errors"

var (
	ErrProfilesNotFound = errors.New("profiles not found")
	ErrNoChanges = errors.New("no changes or profile not found")
	ErrProfileNotFound = errors.New("profile not found")
)