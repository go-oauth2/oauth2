package manage

import "time"

// Config authorization configuration parameters
// 配置授权配置参数
type Config struct {
	// access token expiration time, 0 means it doesn't expire
	// 访问令牌的过期时间，0表示它没有过期
	AccessTokenExp time.Duration
	// refresh token expiration time, 0 means it doesn't expire
	// 刷新令牌过期时间，0表示它没有过期
	RefreshTokenExp time.Duration
	// whether to generate the refreshing token
	// 刷新令牌过期时间，0表示它没有过期
	IsGenerateRefresh bool
}

// RefreshingConfig refreshing token config
// 刷新令牌配置
type RefreshingConfig struct {
	// access token expiration time, 0 means it doesn't expire
	// 访问令牌的过期时间，0表示它没有过期
	AccessTokenExp time.Duration
	// refresh token expiration time, 0 means it doesn't expire
	// 刷新令牌过期时间，0表示它没有过期
	RefreshTokenExp time.Duration
	// whether to generate the refreshing token
	// 是否生成刷新令牌
	IsGenerateRefresh bool
	// whether to reset the refreshing create time
	// 是否重置刷新时间
	IsResetRefreshTime bool
	// whether to remove access token
	// 是否删除访问令牌
	IsRemoveAccess bool
	// whether to remove refreshing token
	// 是否删除刷新令牌
	IsRemoveRefreshing bool
}

// default configs
// 默认配置
var (
	DefaultCodeExp               = time.Minute * 10
	DefaultAuthorizeCodeTokenCfg = &Config{AccessTokenExp: time.Hour * 2, RefreshTokenExp: time.Hour * 24 * 3, IsGenerateRefresh: true}
	DefaultImplicitTokenCfg      = &Config{AccessTokenExp: time.Hour * 1}
	DefaultPasswordTokenCfg      = &Config{AccessTokenExp: time.Hour * 2, RefreshTokenExp: time.Hour * 24 * 7, IsGenerateRefresh: true}
	DefaultClientTokenCfg        = &Config{AccessTokenExp: time.Hour * 2}
	DefaultRefreshTokenCfg       = &RefreshingConfig{IsGenerateRefresh: true, IsRemoveAccess: true, IsRemoveRefreshing: true}
)
