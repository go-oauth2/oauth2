package oauth2

// ResponseType Response Type
type ResponseType string

const (
	// Code Authorization code type
	Code ResponseType = "code"
	// Token Token type
	Token ResponseType = "token"
)

func (rt ResponseType) String() string {
	return string(rt)
}

// GrantType Authorization Grant
type GrantType string

const (
	// AuthorizationCode Authorization Code
	AuthorizationCode GrantType = "authorization_code"
	// PasswordCredentials Resource Owner Password Credentials
	PasswordCredentials GrantType = "password"
	// ClientCredentials Client Credentials
	ClientCredentials GrantType = "clientcredentials"
	// Refreshing Refresh Token
	Refreshing GrantType = "refreshtoken"
	// Implicit Implicit Grant
	Implicit GrantType = "__implicit"
)

func (gt GrantType) String() string {
	return string(gt)
}
