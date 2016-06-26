package manage

import "errors"

var (
	// ErrNilValue Nil Value
	ErrNilValue = errors.New("nil value")

	// ErrClientNotFound Client not Found
	ErrClientNotFound = errors.New("client not found")

	// ErrClientInvalid Client invalid
	ErrClientInvalid = errors.New("client invalid")

	// ErrAuthTokenInvalid Authorize token invalid
	ErrAuthTokenInvalid = errors.New("authorize token invalid")

	// ErrExpired Expired
	ErrExpired = errors.New("expired")

	// ErrForbidden Forbidden
	ErrForbidden = errors.New("forbidden")
)
