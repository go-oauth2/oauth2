package oauth2

import "errors"

var (
	// ErrNotFound Not Found
	ErrNotFound = errors.New("not found")

	// ErrInvalid Invalid
	ErrInvalid = errors.New("invalid")

	// ErrExpired Expired
	ErrExpired = errors.New("expired")

	// ErrForbidden Forbidden
	ErrForbidden = errors.New("forbidden")

	// ErrNilValue Nil Value
	ErrNilValue = errors.New("nil value")
)
