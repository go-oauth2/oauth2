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

	// ErrRefreshInvalid Refresh token invalid
	ErrRefreshInvalid = errors.New("refresh token invalid")

	// ErrRefreshExpired Refresh token expired
	ErrRefreshExpired = errors.New("refresh token expired")

	// ErrTokenInvalid Token invalid
	ErrTokenInvalid = errors.New("token invalid")

	// ErrTokenExpired Token expired
	ErrTokenExpired = errors.New("token expired")
)
