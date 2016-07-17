package manage

import (
	"github.com/LyricTian/errors"
)

var (
	// ErrNilValue Nil Value
	ErrNilValue = errors.Errorf("nil value")

	// ErrInvalidUserID invalid user id
	ErrInvalidUserID = errors.Errorf("invalid user id")

	// ErrInvalidRedirectURI invalid redirect uri
	ErrInvalidRedirectURI = errors.Errorf("invalid redirect uri")

	// ErrInvalidClient invalid client
	ErrInvalidClient = errors.Errorf("invalid client")

	// ErrInvalidAuthorizeCode invalid authorize code
	ErrInvalidAuthorizeCode = errors.Errorf("invalid authorize code")

	// ErrInvalidAccessToken invalid access token
	ErrInvalidAccessToken = errors.Errorf("invalid access token")

	// ErrInvalidRefreshToken  invalid refresh token
	ErrInvalidRefreshToken = errors.Errorf("invalid refresh token")

	// ErrExpiredAccessToken expired access token
	ErrExpiredAccessToken = errors.Errorf("expired access token")

	// ErrExpiredRefreshToken expired refresh token
	ErrExpiredRefreshToken = errors.Errorf("expired refresh token")
)
