package oauth2

// TokenGenerateRequest 提供生成令牌的请求参数
type TokenGenerateRequest struct {
	ClientID          string // 客户端标识
	ClientSecret      string // 客户端密钥
	UserID            string // 用户标识
	RedirectURI       string // 重定向URI
	Scope             string // 授权范围
	Code              string // 授权码(授权码模式使用)
	Refresh           string // 刷新令牌
	IsGenerateRefresh bool   // 是否生成更新令牌
}

// Manager OAuth2授权管理接口
type Manager interface {
	// GetClient 获取客户端信息
	// clientID 客户端标识
	GetClient(clientID string) (cli ClientInfo, err error)

	// GenerateAuthToken 生成授权令牌
	// rt 授权类型
	// tgr 生成令牌的请求参数
	GenerateAuthToken(rt ResponseType, tgr *TokenGenerateRequest) (authToken TokenInfo, err error)

	// GenerateAccessToken 生成访问令牌、更新令牌
	// rt 授权模式
	// tgr 生成令牌的请求参数
	GenerateAccessToken(rt GrantType, tgr *TokenGenerateRequest) (accessToken TokenInfo, err error)

	// RefreshAccessToken 更新访问令牌
	// refresh 更新令牌
	// scope 作用域
	RefreshAccessToken(refresh, scope string) (accessToken TokenInfo, err error)

	// RemoveAccessToken 删除访问令牌
	// access 访问令牌
	RemoveAccessToken(access string) (err error)

	// RemoveRefreshToken 删除更新令牌
	// refresh 更新令牌
	RemoveRefreshToken(refresh string) (err error)

	// LoadAccessToken 加载访问令牌信息
	// access 访问令牌
	LoadAccessToken(access string) (ti TokenInfo, err error)

	// LoadRefreshToken 加载更新令牌信息
	// refresh 更新令牌
	LoadRefreshToken(refresh string) (ti TokenInfo, err error)
}
