package util

import "errors"

// ErrDurationNotSet is an error for when a duration is not set
func ErrDurationNotSet() error {
	return errors.New("duration not set")
}

// ErrDurationNegative is an error for when a duration is negative
func ErrDurationNegative() error {
	return errors.New("duration cannot be negative")
}
