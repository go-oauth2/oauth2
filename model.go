package oauth2

import "time"

// 相关模型接口的定义
type (
	// ClientInfo 客户端信息模型接口
	ClientInfo interface {
		// 客户端唯一标识
		GetID() string
		// 客户端秘钥
		GetSecret() string
		// 客户端域名URL
		GetDomain() string
		// 自定义数据
		GetRetainData() interface{}
	}

	// Authorize 授权信息模型接口
	Authorize interface {
		// 客户端标识
		GetClientID() string
		// 设置客户端标识
		SetClientID(string)
		// 用户标识
		GetUserID() string
		// 设置用户标识
		SetUserID(string)
		// 重定向URI
		GetRedirectURI() string
		// 设置重定向URI
		SetRedirectURI(string)
		// 权限范围
		GetScope() string
		// 设置权限范围
		SetScope(string)
		// 创建时间
		GetCreateAt() time.Time
		// 设置创建时间
		SetCreateAt(time.Time)
		// 有效期
		GetExpiresIn() time.Duration
		// 设置有效期
		SetExpiresIn(time.Duration)
		// 授权令牌
		GetToken() string
		// 设置授权令牌
		SetToken(string)
		// 用于标识授权令牌的唯一标识码
		GetIdentifier() string
		// 设置用于标识授权令牌的唯一标识码
		SetIdentifier(string)
	}

	// TokenInfo 令牌信息模型接口
	TokenInfo interface {
		// 客户端标识
		GetClientID() string
		// 设置客户端标识
		SetClientID(string)
		// 用户标识
		GetUserID() string
		// 设置用户标识
		SetUserID(string)
		// 访问令牌
		GetAccess() TokenBasic
		// 设置访问令牌
		SetAccess(TokenBasic)
		// 更新令牌
		GetRefresh() TokenBasic
		// 设置更新令牌
		SetRefresh(TokenBasic)
		// 权限范围
		GetScope() string
		// 设置权限范围
		SetScope(string)
	}

	// TokenBasic 令牌基础模型接口
	TokenBasic interface {
		// 创建时间
		GetCreateAt() time.Time
		// 设置创建时间
		SetCreateAt(time.Time)
		// 有效期
		GetExpiresIn() time.Duration
		// 设置有效期
		SetExpiresIn(time.Duration)
		// 令牌
		GetToken() string
		// 设置令牌
		SetToken(string)
		// 用于标识令牌的唯一标识码
		GetIdentifier() string
		// 设置用于标识令牌的唯一标识码
		SetIdentifier(string)
	}
)
