package bento

import "errors"

// Define package-level errors
var ErrInvalidConfig = errors.New("invalid configuration: missing required fields")
var ErrInvalidEmail = errors.New("invalid email address")
var ErrInvalidIPAddress = errors.New("invalid IP address")
var ErrInvalidRequest = errors.New("invalid request parameters")
var ErrAPIResponse = errors.New("unexpected API response")
var ErrInvalidName = errors.New("invalid name format")
var ErrInvalidSegmentID = errors.New("invalid segment ID")
var ErrInvalidContent = errors.New("invalid content")
var ErrInvalidTags = errors.New("invalid tags format")
var ErrInvalidBatchSize = errors.New("invalid batch size")
var ErrInvalidKeyLength = errors.New("invalid key length")
