package oauth2

type (
	// ClientStore The client information storage interface
	ClientStore interface {
		// GetByID According to the ID for the client information
		GetByID(id string) (ClientInfo, error)
	}

	// TokenStore The token information storage interface
	TokenStore interface {
		// Create Create and store the new token information
		Create(info TokenInfo) error

		// RemoveByAccess Use the access token to delete the token information(Along with the refresh token)
		RemoveByAccess(access string) error

		// RemoveByRefresh Use the refresh token to delete the token information
		RemoveByRefresh(refresh string) error

		// Use the access token for token information data
		GetByAccess(access string) (TokenInfo, error)

		// Use the refresh token for token information data
		GetByRefresh(refresh string) (TokenInfo, error)
	}
)
