package oauth2

// GrantType 定义授权模式
type GrantType byte

const (
	// AuthorizationCode 授权码模式
	AuthorizationCode GrantType = 1 << (iota + 1)
	// Implicit 简化模式
	Implicit
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
		return "implicit"
	case 1 << 3:
		return "password"
	case 1 << 4:
		return "clientcredentials"
	}
	return "unknown"
}
