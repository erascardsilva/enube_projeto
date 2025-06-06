package parser

import "errors"

var (
	// ErrInvalidRecord is returned when a record has invalid format
	ErrInvalidRecord = errors.New("invalid record format")

	// ErrInvalidDate is returned when a date field is invalid
	ErrInvalidDate = errors.New("invalid date format")

	// ErrInvalidAmount is returned when an amount field is invalid
	ErrInvalidAmount = errors.New("invalid amount format")
)
