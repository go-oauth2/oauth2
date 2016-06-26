package manage

import "time"

// Config 授权配置参数
type Config struct {
	TokenExpiresIn   time.Duration // 令牌有效期
	RefreshExpiresIn time.Duration // 刷新令牌有效期
}

// TokenGenerateData 提供生成令牌的相应参数
type TokenGenerateData struct {
	ClientID          string // 客户端标识
	ClientSecret      string // 客户端密钥
	UserID            string // 用户标识
	RedirectURI       string // 重定向URI
	Scope             string // 授权范围
	Code              string // 授权码(授权码模式使用)
	IsGenerateRefresh bool   // 是否生成刷新令牌
}
