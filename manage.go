package oauth2

// TokenGenerateRequest 提供生成令牌的请求参数
type TokenGenerateRequest struct {
	ClientID          string // 客户端标识
	ClientSecret      string // 客户端密钥
	UserID            string // 用户标识
	RedirectURI       string // 重定向URI
	Scope             string // 授权范围
	Code              string // 授权码(授权码模式使用)
	IsGenerateRefresh bool   // 是否生成刷新令牌
}

// Manager OAuth2授权管理接口
type Manager interface {
	// GenerateAuthToken 生成授权令牌
	// rt 授权类型
	// tgr 生成令牌的请求参数
	GenerateAuthToken(rt ResponseType, tgr *TokenGenerateRequest) (token string, err error)

	// GenerateToken 生成访问令牌、刷新令牌
	// rt 授权模式
	// tgr 生成令牌的请求参数
	GenerateToken(rt GrantType, tgr *TokenGenerateRequest) (token, refresh string, err error)

	// RefreshToken 使用刷新令牌更新访问令牌
	// refresh 刷新令牌
	// scope 作用域
	RefreshToken(refresh, scope string) (token string, err error)

	// RevokeToken 使用访问令牌废除令牌信息
	// token 访问令牌
	RevokeToken(token string) (err error)

	// CheckToken 令牌检查，如果存在则返回令牌信息
	CheckToken(token string) (ti TokenInfo, err error)

	// CheckRefreshToken 访问令牌检查，如果存在则返回令牌信息
	CheckRefreshToken(refresh string) (ti TokenInfo, err error)
}
