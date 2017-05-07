package manage

import "time"

// Config authorization configuration parameters
type Config struct {
	// access token expiration time
	AccessTokenExp time.Duration
	// refresh token expiration time
	RefreshTokenExp time.Duration
	// whether to generate the refreshing token
	IsGenerateRefresh bool
}

// RefreshingConfig refreshing token config
type RefreshingConfig struct {
	// access token expiration time
	AccessTokenExp time.Duration
	// refresh token expiration time
	RefreshTokenExp time.Duration
	// whether to generate the refreshing token
	IsGenerateRefresh bool
	// whether to reset the refreshing create time
	IsResetRefreshTime bool
	// whether to remove access token
	IsRemoveAccess bool
	// whether to remove refreshing token
	IsRemoveRefreshing bool
}

// default configs
var (
	DefaultCodeExp               = time.Minute * 10
	DefaultAuthorizeCodeTokenCfg = &Config{AccessTokenExp: time.Hour * 2, RefreshTokenExp: time.Hour * 24 * 3, IsGenerateRefresh: true}
	DefaultImplicitTokenCfg      = &Config{AccessTokenExp: time.Hour * 1}
	DefaultPasswordTokenCfg      = &Config{AccessTokenExp: time.Hour * 2, RefreshTokenExp: time.Hour * 24 * 7, IsGenerateRefresh: true}
	DefaultClientTokenCfg        = &Config{AccessTokenExp: time.Hour * 2}
	DefaultRefreshTokenCfg       = &RefreshingConfig{IsRemoveAccess: true, IsRemoveRefreshing: true}
)
