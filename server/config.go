package server

import "gopkg.in/oauth2.v2"

// Config 配置参数
type Config struct {
	// TokenType 令牌类型（默认为Bearer）
	TokenType string
	// AllowedResponseType 允许的授权类型（默认code）
	AllowedResponseType []oauth2.ResponseType
	// AllowedGrantType 允许的授权模式（默认authorization_code）
	AllowedGrantType []oauth2.GrantType
}

// NewConfig 创建默认的配置参数
func NewConfig() *Config {
	return &Config{
		TokenType:           "Bearer",
		AllowedResponseType: []oauth2.ResponseType{oauth2.Code},
		AllowedGrantType:    []oauth2.GrantType{oauth2.AuthorizationCodeCredentials},
	}
}
