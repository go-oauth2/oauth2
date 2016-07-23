package errors

import "errors"

var (
	// ErrNilValue Nil Value
	ErrNilValue = errors.New("nil value")

	// ErrInvalidRedirectURI invalid redirect uri
	ErrInvalidRedirectURI = errors.New("invalid redirect uri")

	// ErrInvalidAuthorizeCode invalid authorize code
	ErrInvalidAuthorizeCode = errors.New("invalid authorize code")

	// ErrInvalidAccessToken invalid access token
	ErrInvalidAccessToken = errors.New("invalid access token")

	// ErrInvalidRefreshToken  invalid refresh token
	ErrInvalidRefreshToken = errors.New("invalid refresh token")

	// ErrExpiredAccessToken expired access token
	ErrExpiredAccessToken = errors.New("expired access token")

	// ErrExpiredRefreshToken expired refresh token
	ErrExpiredRefreshToken = errors.New("expired refresh token")
)
