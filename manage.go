package oauth2

import (
	"context"
	"net/http"
	"time"
)

// TokenGenerateRequest provide to generate the token request parameters
// 提供以生成令牌请求参数
type TokenGenerateRequest struct {
	ClientID       string
	ClientSecret   string
	UserID         string
	RedirectURI    string
	Scope          string
	Code           string
	Refresh        string
	AccessTokenExp time.Duration
	Request        *http.Request
}

// Manager authorization management interface
// 授权管理接口
type Manager interface {
	// get the client information
	// 获取客户端信息
	GetClient(ctx context.Context, clientID string) (cli ClientInfo, err error)

	// generate the authorization token(code)
	// 生成授权令牌
	GenerateAuthToken(ctx context.Context, rt ResponseType, tgr *TokenGenerateRequest) (authToken TokenInfo, err error)

	// generate the access token
	// 生成访问令牌
	GenerateAccessToken(ctx context.Context, rt GrantType, tgr *TokenGenerateRequest) (accessToken TokenInfo, err error)

	// refreshing an access token
	// 刷新访问令牌
	RefreshAccessToken(ctx context.Context, tgr *TokenGenerateRequest) (accessToken TokenInfo, err error)

	// use the access token to delete the token information
	// 使用访问令牌删除令牌信息
	RemoveAccessToken(ctx context.Context, access string) (err error)

	// use the refresh token to delete the token information
	// 使用刷新令牌删除令牌信息
	RemoveRefreshToken(ctx context.Context, refresh string) (err error)

	// according to the access token for corresponding token information
	// 根据访问令牌获取相应的令牌信息
	LoadAccessToken(ctx context.Context, access string) (ti TokenInfo, err error)

	// according to the refresh token for corresponding token information
	// 根据刷新令牌获取相应的令牌信息
	LoadRefreshToken(ctx context.Context, refresh string) (ti TokenInfo, err error)
}
