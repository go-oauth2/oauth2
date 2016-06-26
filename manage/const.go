package manage

// ResponseType 授权类型
type ResponseType byte

const (
	// Code 授权码类型
	Code ResponseType = 1 << (iota + 1)
	// Token 令牌类型
	Token
)

func (rt ResponseType) String() string {
	switch rt {
	case 1 << 1:
		return "code"
	case 1 << 2:
		return "token"
	}
	return "unknown"
}

// GrantType 定义授权模式
type GrantType byte

const (
	// AuthorizationCode 授权码模式
	AuthorizationCode GrantType = 1 << (iota + 1)
	// PasswordCredentials 密码模式
	PasswordCredentials
	// ClientCredentials 客户端模式
	ClientCredentials
)

func (gt GrantType) String() string {
	switch gt {
	case 1 << 1:
		return "authorization_code"
	case 1 << 2:
		return "password"
	case 1 << 3:
		return "clientcredentials"
	}
	return "unknown"
}
