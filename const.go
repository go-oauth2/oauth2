package oauth2

// ResponseType 定义授权类型
type ResponseType string

const (
	// Code 授权码类型
	Code ResponseType = "code"
	// Token 令牌类型
	Token ResponseType = "token"
)

func (rt ResponseType) String() string {
	return string(rt)
}

// GrantType 定义授权模式
type GrantType string

const (
	// AuthorizationCodeCredentials 授权码模式
	AuthorizationCodeCredentials GrantType = "authorization_code"
	// PasswordCredentials 密码模式
	PasswordCredentials GrantType = "password"
	// ClientCredentials 客户端模式
	ClientCredentials GrantType = "clientcredentials"
	// RefreshCredentials 更新令牌模式
	RefreshCredentials GrantType = "refreshtoken"
)

func (gt GrantType) String() string {
	return string(gt)
}
