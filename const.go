package oauth2

// ResponseType the type of authorization request
// 响应类型
type ResponseType string

// define the type of authorization request
// 定义授权请求的类型
const (
	Code  ResponseType = "code"
	Token ResponseType = "token"
)

func (rt ResponseType) String() string {
	return string(rt)
}

// GrantType authorization model
// 授权模型
type GrantType string

// define authorization model
// 定义授权模型
const (
	AuthorizationCode   GrantType = "authorization_code"
	PasswordCredentials GrantType = "password"
	ClientCredentials   GrantType = "client_credentials"
	Refreshing          GrantType = "refresh_token"
	Implicit            GrantType = "__implicit"
)

func (gt GrantType) String() string {
	if gt == AuthorizationCode ||
		gt == PasswordCredentials ||
		gt == ClientCredentials ||
		gt == Refreshing {
		return string(gt)
	}
	return ""
}
