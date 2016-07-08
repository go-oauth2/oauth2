package server

import "errors"

var (
	// ErrRequestMethodInvalid Request method invalid
	ErrRequestMethodInvalid = errors.New("request method invalid")

	// ErrResponseTypeInvalid Response type invalid
	ErrResponseTypeInvalid = errors.New("response type invalid")

	// ErrGrantTypeInvalid Grant type invalid
	ErrGrantTypeInvalid = errors.New("grant type invalid")

	// ErrClientInvalid Client invalid
	ErrClientInvalid = errors.New("client invalid")

	// ErrUserInvalid User invalid
	ErrUserInvalid = errors.New("user invalid")

	// ErrAuthorizationFormInvalid Authorization form invalid
	ErrAuthorizationFormInvalid = errors.New("authorization form invalid")

	// ErrAuthorizationHeaderInvalid Authorization header invalid
	ErrAuthorizationHeaderInvalid = errors.New("authorization header invalid")

	// ErrRefreshInvalid Refresh token invalid
	ErrRefreshInvalid = errors.New("refresh token invalid")
)
