package server

import (
	"time"

	"gopkg.in/oauth2.v3"
)

// AuthorizeRequest authorization request
type AuthorizeRequest struct {
	ResponseType   oauth2.ResponseType
	ClientID       string
	Scope          string
	RedirectURI    string
	State          string
	UserID         string
	AccessTokenExp time.Duration
}
