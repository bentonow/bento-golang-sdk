package bento

import "errors"

var (
	ErrInvalidConfig    = errors.New("invalid configuration: missing required fields")
	ErrInvalidEmail     = errors.New("invalid email address")
	ErrInvalidIPAddress = errors.New("invalid IP address")
	ErrInvalidRequest   = errors.New("invalid request parameters")
	ErrAPIResponse      = errors.New("unexpected API response")
	ErrInvalidName      = errors.New("invalid name format")
	ErrInvalidSegmentID = errors.New("invalid segment ID")
	ErrInvalidContent   = errors.New("invalid content")
	ErrInvalidTags      = errors.New("invalid tags format")
	ErrInvalidBatchSize = errors.New("invalid batch size")
)
