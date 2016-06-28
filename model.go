package oauth2

import "time"

// 相关模型接口的定义
type (
	// ClientInfo 客户端信息模型接口
	ClientInfo interface {
		// 客户端ID
		GetID() string
		// 客户端秘钥
		GetSecret() string
		// 客户端域名URL
		GetDomain() string
		// 自定义数据
		GetRetainData() interface{}
	}

	// TokenInfo 令牌信息模型接口
	TokenInfo interface {
		// 客户端ID
		GetClientID() string
		// 设置客户端ID
		SetClientID(string)
		// 用户ID
		GetUserID() string
		// 设置用户ID
		SetUserID(string)
		// 重定向URI
		GetRedirectURI() string
		// 设置重定向URI
		SetRedirectURI(string)
		// 权限范围
		GetScope() string
		// 设置权限范围
		SetScope(string)

		// 令牌
		GetToken() string
		// 设置令牌
		SetToken(string)
		// 令牌创建时间
		GetTokenCreateAt() time.Time
		// 设置令牌创建时间
		SetTokenCreateAt(time.Time)
		// 令牌有效期
		GetTokenExpiresIn() time.Duration
		// 设置令牌有效期
		SetTokenExpiresIn(time.Duration)

		// 更新令牌
		GetRefresh() string
		// 设置更新令牌
		SetRefresh(string)
		// 更新令牌创建时间
		GetRefreshCreateAt() time.Time
		// 设置更新令牌创建时间
		SetRefreshCreateAt(time.Time)
		// 更新令牌有效期
		GetRefreshExpiresIn() time.Duration
		// 设置更新令牌有效期
		SetRefreshExpiresIn(time.Duration)
	}
)
