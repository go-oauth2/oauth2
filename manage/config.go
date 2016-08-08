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

// default configs
var (
	DefaultCodeExp               = time.Minute * 10
	DefaultAuthorizeCodeTokenCfg = &Config{AccessTokenExp: time.Hour * 2, RefreshTokenExp: time.Hour * 24 * 3, IsGenerateRefresh: true}
	DefaultImplicitTokenCfg      = &Config{AccessTokenExp: time.Hour * 1}
	DefaultPasswordTokenCfg      = &Config{AccessTokenExp: time.Hour * 2, RefreshTokenExp: time.Hour * 24 * 7, IsGenerateRefresh: true}
	DefaultClientTokenCfg        = &Config{AccessTokenExp: time.Hour * 2}
	DefaultRefreshTokenCfg       = &Config{}
)
