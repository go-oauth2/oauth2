package oauth2

import "time"

type (
	// ClientInfo The client information model interface
	ClientInfo interface {
		// The client id
		GetID() string
		// The client secret
		GetSecret() string
		// The client domain
		GetDomain() string
		// The extension data related to the client
		GetExtraData() interface{}
	}

	// TokenInfo The token information model interface
	TokenInfo interface {
		// Get client id
		GetClientID() string
		// Set client id
		SetClientID(string)
		// Get user id
		GetUserID() string
		// Set user id
		SetUserID(string)
		// Get Redirect URI
		GetRedirectURI() string
		// Set Redirect URI
		SetRedirectURI(string)
		// Get Scope of authorization
		GetScope() string
		// Set Scope of authorization
		SetScope(string)

		// Get Access Token
		GetAccess() string
		// Set Access Token
		SetAccess(string)
		// Get Create Time
		GetAccessCreateAt() time.Time
		// Set Create Time
		SetAccessCreateAt(time.Time)
		// Get The lifetime in seconds of the access token
		GetAccessExpiresIn() time.Duration
		// Set The lifetime in seconds of the access token
		SetAccessExpiresIn(time.Duration)

		// Get Refresh Token
		GetRefresh() string
		// Set Refresh Token
		SetRefresh(string)
		// Get Create Time
		GetRefreshCreateAt() time.Time
		// Set Create Time
		SetRefreshCreateAt(time.Time)
		// Get The lifetime in seconds of the access token
		GetRefreshExpiresIn() time.Duration
		// Set The lifetime in seconds of the access token
		SetRefreshExpiresIn(time.Duration)
	}
)
