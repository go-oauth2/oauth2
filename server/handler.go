package server

import (
	"net/http"

	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/errors"
)

// ClientInfoHandler Get client info from request
type ClientInfoHandler func(r *http.Request) (clientID, clientSecret string, err error)

// ClientAuthorizedHandler Check the client allows to use this authorization grant type
type ClientAuthorizedHandler func(clientID string, grantType oauth2.GrantType) (allowed bool, err error)

// ClientScopeHandler Check the client allows to use scope
type ClientScopeHandler func(clientID, scope string) (allowed bool, err error)

// UserAuthorizationHandler Get user id from request authorization
type UserAuthorizationHandler func(w http.ResponseWriter, r *http.Request) (userID string, err error)

// PasswordAuthorizationHandler Get user id from username and password
type PasswordAuthorizationHandler func(username, password string) (userID string, err error)

// RefreshingScopeHandler Check the scope of the refreshing token
type RefreshingScopeHandler func(newScope, oldScope string) (allowed bool)

// ResponseErrorHandler Response error handing
type ResponseErrorHandler func(re *errors.Response)

// InternalErrorHandler Internal error handing
type InternalErrorHandler func(err error)

// ClientFormHandler Get client data from form
func ClientFormHandler(r *http.Request) (clientID, clientSecret string, err error) {
	clientID = r.Form.Get("client_id")
	clientSecret = r.Form.Get("client_secret")
	if clientID == "" || clientSecret == "" {
		err = errors.ErrInvalidRequest
	}
	return
}

// ClientBasicHandler Get client data from basic authorization
func ClientBasicHandler(r *http.Request) (clientID, clientSecret string, err error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		err = errors.ErrInvalidClient
		return
	}
	clientID = username
	clientSecret = password
	return
}
