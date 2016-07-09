package oauth2

// 提供存储接口
type (
	// ClientStore 客户端信息存储接口
	ClientStore interface {
		// GetByID 根据ID获取客户端信息
		GetByID(id string) (ClientInfo, error)
	}

	// TokenStore 令牌信息存储接口
	TokenStore interface {
		// Create 创建并存储新的令牌信息
		Create(info TokenInfo) error

		// RemoveByAccess 使用访问令牌删除令牌信息
		RemoveByAccess(access string) error

		// RemoveByRefresh 使用更新令牌删除令牌信息
		RemoveByRefresh(refresh string) error

		// 使用访问令牌获取令牌信息数据
		GetByAccess(access string) (TokenInfo, error)

		// 根据更新令牌获取令牌信息数据
		GetByRefresh(refresh string) (TokenInfo, error)
	}
)
