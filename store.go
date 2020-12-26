package oauth2

import "context"

type (
	// ClientStore the client information storage interface
	// 客户端信息存储接口
	ClientStore interface {
		// according to the ID for the client information
		// 根据客户信息ID
		GetByID(ctx context.Context, id string) (ClientInfo, error)
	}

	// TokenStore the token information storage interface
	// 令牌信息存储接口
	TokenStore interface {
		// create and store the new token information
		// 创建并存储新的令牌信息
		Create(ctx context.Context, info TokenInfo) error

		// delete the authorization code
		// 删除授权码
		RemoveByCode(ctx context.Context, code string) error

		// use the access token to delete the token information
		// 使用访问令牌删除令牌信息
		RemoveByAccess(ctx context.Context, access string) error

		// use the refresh token to delete the token information
		// 使用刷新令牌删除令牌信息
		RemoveByRefresh(ctx context.Context, refresh string) error

		// use the authorization code for token information data
		// 将授权码用于令牌信息数据
		GetByCode(ctx context.Context, code string) (TokenInfo, error)

		// use the access token for token information data
		// 将访问令牌用于令牌信息数据
		GetByAccess(ctx context.Context, access string) (TokenInfo, error)

		// use the refresh token for token information data
		// 将访问令牌用于令牌信息数据
		GetByRefresh(ctx context.Context, refresh string) (TokenInfo, error)
	}
)
