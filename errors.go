package bento

import "errors"

// Define package-level errors
var (
	// ErrInvalidConfig indicates missing required configuration fields
	ErrInvalidConfig = errors.New("invalid configuration: missing required fields")

	// ErrInvalidEmail indicates an invalid email address format
	ErrInvalidEmail = errors.New("invalid email address")

	// ErrInvalidIPAddress indicates an invalid IP address format
	ErrInvalidIPAddress = errors.New("invalid IP address")

	// ErrInvalidRequest indicates invalid request parameters
	ErrInvalidRequest = errors.New("invalid request parameters")

	// ErrAPIResponse indicates an unexpected API response
	ErrAPIResponse = errors.New("unexpected API response")

	// ErrInvalidName indicates an invalid name format
	ErrInvalidName = errors.New("invalid name format")

	// ErrInvalidSegmentID indicates an invalid segment ID
	ErrInvalidSegmentID = errors.New("invalid segment ID")

	// ErrInvalidContent indicates invalid content
	ErrInvalidContent = errors.New("invalid content")

	// ErrInvalidTags indicates invalid tags format
	ErrInvalidTags = errors.New("invalid tags format")

	// ErrInvalidBatchSize indicates invalid batch size
	ErrInvalidBatchSize = errors.New("invalid batch size")

	// ErrInvalidKeyLength indicates an invalid key length
	ErrInvalidKeyLength = errors.New("invalid key length")
)