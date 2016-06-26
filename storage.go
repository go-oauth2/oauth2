package oauth2

// 提供存储接口
type (
	// ClientStorage 客户端信息存储接口
	ClientStorage interface {
		// GetByID 根据ID获取客户端信息
		GetByID(id string) (ClientInfo, error)
	}

	// TokenStorage 令牌信息存储接口
	TokenStorage interface {
		// Create 创建并存储新的令牌信息
		Create(info TokenInfo) error

		// UpdateByRefresh 根据刷新令牌更新令牌信息
		UpdateByRefresh(refresh string, info TokenInfo) error

		// DeleteByToken 根据令牌删除令牌信息
		DeleteByToken(val string) error

		// 根据令牌取出令牌信息数据(获取并删除)
		TakeByToken(val string) (TokenInfo, error)

		// 根据令牌获取令牌信息数据
		GetByToken(val string) (TokenInfo, error)

		// 根据刷新令牌获取令牌信息数据
		GetByRefresh(refresh string) (TokenInfo, error)

		// 将该令牌对应的令牌信息作过期处理
		ExpiredByToken(val string) error

		// 将该刷新令牌对应的令牌信息作过期处理
		ExpiredByRefresh(refresh string) error
	}
)
