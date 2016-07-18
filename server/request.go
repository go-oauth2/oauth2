package server

import "gopkg.in/oauth2.v3"

// AuthorizeRequest The authorization request
type AuthorizeRequest struct {
	ResponseType oauth2.ResponseType
	ClientID     string
	Scope        string
	RedirectURI  string
	State        string
	UserID       string
}
