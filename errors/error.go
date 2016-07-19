package errors

import "errors"

var (
	// ErrUnauthorizedClient unauthorized client
	ErrUnauthorizedClient = errors.New("unauthorized_client")

	// ErrAccessDenied access denied
	ErrAccessDenied = errors.New("access_denied")

	// ErrUnsupportedResponseType unsupported response type
	ErrUnsupportedResponseType = errors.New("unsupported_response_type")

	// ErrInvalidScope invalid scope
	ErrInvalidScope = errors.New("invalid_scope")

	// ErrInvalidRequest invalid request
	ErrInvalidRequest = errors.New("invalid_request")

	// ErrInvalidClient invalid client
	ErrInvalidClient = errors.New("invalid_client")

	// ErrInvalidGrant invalid grant
	ErrInvalidGrant = errors.New("invalid_grant")

	// ErrUnsupportedGrantType unsupported grant type
	ErrUnsupportedGrantType = errors.New("unsupported_grant_type")

	// ErrServerError server error
	ErrServerError = errors.New("server_error")
)

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
