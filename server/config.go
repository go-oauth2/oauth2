package server

import "gopkg.in/oauth2.v3"

// Config configuration parameters
type Config struct {
	TokenType            string                // TokenType token type（Default is Bearer）
	AllowedResponseTypes []oauth2.ResponseType // Allow the authorization type(Default is all)
	AllowedGrantTypes    []oauth2.GrantType    // Allow the grant type(Default is all)
}

// NewConfig create to configuration instance
func NewConfig() *Config {
	return &Config{
		TokenType:            "Bearer",
		AllowedResponseTypes: []oauth2.ResponseType{oauth2.Code, oauth2.Token},
		AllowedGrantTypes: []oauth2.GrantType{
			oauth2.AuthorizationCode,
			oauth2.PasswordCredentials,
			oauth2.ClientCredentials,
			oauth2.Refreshing,
		},
	}
}
