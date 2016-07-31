package oauth2

import (
	"time"
)

// TokenGenerateRequest Provide to generate the token request parameters
type TokenGenerateRequest struct {
	ClientID       string        // The client information
	ClientSecret   string        // The client secret
	UserID         string        // The user id
	RedirectURI    string        // Redirect URI
	Scope          string        // Scope of authorization
	Code           string        // Authorization code
	Refresh        string        // Refreshing token
	AccessTokenExp time.Duration // Access token expiration time (in seconds)
}

// Manager Authorization management interface
type Manager interface {
	// Check the interface implementation
	CheckInterface() (err error)

	// Get the client information
	GetClient(clientID string) (cli ClientInfo, err error)

	// Generate the authorization token(code)
	GenerateAuthToken(rt ResponseType, tgr *TokenGenerateRequest) (authToken TokenInfo, err error)

	// Generate the access token
	GenerateAccessToken(rt GrantType, tgr *TokenGenerateRequest) (accessToken TokenInfo, err error)

	// Refreshing an access token
	RefreshAccessToken(tgr *TokenGenerateRequest) (accessToken TokenInfo, err error)

	// Use the access token to delete the token information
	RemoveAccessToken(access string) (err error)

	// Use the refresh token to delete the token information
	RemoveRefreshToken(refresh string) (err error)

	// According to the access token for corresponding token information
	LoadAccessToken(access string) (ti TokenInfo, err error)

	// According to the refresh token for corresponding token information
	LoadRefreshToken(refresh string) (ti TokenInfo, err error)
}
