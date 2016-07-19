package oauth2

// TokenGenerateRequest Provide to generate the token request parameters
type TokenGenerateRequest struct {
	ClientID     string // The client information
	ClientSecret string // The client secret
	UserID       string // The user id
	RedirectURI  string // Redirect URI
	Scope        string // Scope of authorization
	Code         string // Authorization code
	Refresh      string // Refreshing token
}

// Manager Authorization management interface
type Manager interface {
	// GetClient Get the client information
	GetClient(clientID string) (cli ClientInfo, err error)

	// GenerateAuthToken Generate the authorization token(code)
	GenerateAuthToken(rt ResponseType, tgr *TokenGenerateRequest) (authToken TokenInfo, err error)

	// GenerateAccessToken Generate the access token
	GenerateAccessToken(rt GrantType, tgr *TokenGenerateRequest) (accessToken TokenInfo, err error)

	// RefreshAccessToken Refreshing an access token
	RefreshAccessToken(tgr *TokenGenerateRequest) (accessToken TokenInfo, err error)

	// RemoveAccessToken Use the access token to delete the token information
	RemoveAccessToken(access string) (err error)

	// RemoveRefreshToken Use the refresh token to delete the token information
	RemoveRefreshToken(refresh string) (err error)

	// LoadAccessToken According to the access token for corresponding token information
	LoadAccessToken(access string) (ti TokenInfo, err error)

	// LoadRefreshToken According to the refresh token for corresponding token information
	LoadRefreshToken(refresh string) (ti TokenInfo, err error)
}
