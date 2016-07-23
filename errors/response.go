package errors

import "errors"

// Response error response
type Response struct {
	Error       error  `json:"error"`
	Description string `json:"error_description,omitempty"`
	URI         string `json:"error_uri,omitempty"`
	StatusCode  int    `json:"-"`
}

var (
	// ErrInvalidRequest invalid request
	ErrInvalidRequest = errors.New("invalid_request")

	// ErrUnauthorizedClient unauthorized client
	ErrUnauthorizedClient = errors.New("unauthorized_client")

	// ErrAccessDenied access denied
	ErrAccessDenied = errors.New("access_denied")

	// ErrUnsupportedResponseType unsupported response type
	ErrUnsupportedResponseType = errors.New("unsupported_response_type")

	// ErrInvalidScope invalid scope
	ErrInvalidScope = errors.New("invalid_scope")

	// ErrServerError server error
	ErrServerError = errors.New("server_error")

	// ErrTemporarilyUnavailable temporarily unavailable
	ErrTemporarilyUnavailable = errors.New("temporarily_unavailable")

	// ErrInvalidClient invalid client
	ErrInvalidClient = errors.New("invalid_client")

	// ErrInvalidGrant invalid grant
	ErrInvalidGrant = errors.New("invalid_grant")

	// ErrUnsupportedGrantType unsupported grant type
	ErrUnsupportedGrantType = errors.New("unsupported_grant_type")
)

// Descriptions error description
var Descriptions = map[error]string{
	ErrInvalidRequest:          "The request is missing a required parameter, includes an invalid parameter value, includes a parameter more than once, or is otherwise malformed",
	ErrUnauthorizedClient:      "The client is not authorized to request an authorization code using this method",
	ErrAccessDenied:            "The resource owner or authorization server denied the request",
	ErrUnsupportedResponseType: "The authorization server does not support obtaining an authorization code using this method",
	ErrInvalidScope:            "The requested scope is invalid, unknown, or malformed",
	ErrServerError:             "The authorization server encountered an unexpected condition that prevented it from fulfilling the request",
	ErrTemporarilyUnavailable:  "The authorization server is currently unable to handle the request due to a temporary overloading or maintenance of the server",
	ErrInvalidClient:           "Client authentication failed",
	ErrInvalidGrant:            "The provided authorization grant (e.g., authorization code, resource owner credentials) or refresh token is invalid, expired, revoked, does not match the redirection URI used in the authorization request, or was issued to another client",
	ErrUnsupportedGrantType:    "The authorization grant type is not supported by the authorization server",
}
