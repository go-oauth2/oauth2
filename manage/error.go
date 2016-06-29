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

	// ErrAccessInvalid Access token expired
	ErrAccessInvalid = errors.New("access token invalid")

	// ErrAccessExpired Access token expired
	ErrAccessExpired = errors.New("access token expired")

	// ErrRefreshInvalid Refresh token invalid
	ErrRefreshInvalid = errors.New("refresh token invalid")

	// ErrRefreshExpired Refresh token expired
	ErrRefreshExpired = errors.New("refresh token expired")
)
