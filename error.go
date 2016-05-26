package oauth2

import (
	"errors"
)

var (
	// ErrClientNotFound Client not found
	ErrClientNotFound = errors.New("The client is not found.")

	// ErrACNotFound Authorization code not found
	ErrACNotFound = errors.New("The authorization code is not found.")

	// ErrACInvalid Authorization code invalid
	ErrACInvalid = errors.New("The authorization code is invalid.")

	// ErrCSInvalid Client secret invalid
	ErrCSInvalid = errors.New("The client secret is invalid.")

	// ErrATNotFound Refresh token not found
	ErrATNotFound = errors.New("The access token is not found.")

	// ErrATInvalid Access token invalid
	ErrATInvalid = errors.New("The access token is invalid.")

	// ErrATExpire Access token expire
	ErrATExpire = errors.New("The access token is expire.")

	// ErrRTNotFound Refresh token not found
	ErrRTNotFound = errors.New("The refresh token is not found.")

	// ErrRTInvalid Refresh token invalid
	ErrRTInvalid = errors.New("The refresh token is invalid.")

	// ErrRTExpire Refresh token expire
	ErrRTExpire = errors.New("The refresh token is expire.")
)
