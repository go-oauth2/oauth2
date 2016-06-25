package oauth2

// 提供存储接口
type (
	// ClientStorage 客户端信息存储接口
	ClientStorage interface {
		// GetByID 根据ID获取客户端信息
		GetByID(id string) (ClientInfo, error)
	}

	// AuthorizeStorage 授权码信息存储接口
	AuthorizeStorage interface {
		// 将授权信息放入存储
		Put(info Authorize) error

		// 根据授权令牌取出授权信息
		TakeByToken(token string) (Authorize, error)
	}

	// TokenStorage 令牌信息存储接口
	TokenStorage interface {
		// Create 创建并存储新的令牌信息
		Create(info TokenInfo) error

		// UpdateByRefresh 根据刷新令牌更新令牌信息
		UpdateByRefresh(refresh string, info TokenInfo) error

		// 根据访问令牌获取令牌信息数据
		GetByAccess(access string) (TokenInfo, error)

		// 根据刷新令牌获取令牌信息数据
		GetByRefresh(refresh string) (TokenInfo, error)

		// 根据访问令牌废除令牌信息
		RevokeByAccess(access string) error

		// 将该访问令牌对应的令牌信息作过期处理
		ExpiredByAccess(access string) error

		// 将该刷新令牌对应的令牌信息作过期处理
		ExpiredByRefresh(refresh string) error
	}
)
