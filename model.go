package oauth2

import (
	"time"
)

type (
	// ClientInfo the client information model interface
	// 客户端信息模型接口
	ClientInfo interface {
		GetID() string
		GetSecret() string
		GetDomain() string
		GetUserID() string
	}

	// ClientPasswordVerifier the password handler interface
	// 密码处理程序接口
	ClientPasswordVerifier interface {
		VerifyPassword(string) bool
	}

	// TokenInfo the token information model interface
	// 令牌信息模型接口
	TokenInfo interface {
		New() TokenInfo

		GetClientID() string
		SetClientID(string)
		GetUserID() string
		SetUserID(string)
		GetRedirectURI() string
		SetRedirectURI(string)
		GetScope() string
		SetScope(string)

		GetCode() string
		SetCode(string)
		GetCodeCreateAt() time.Time
		SetCodeCreateAt(time.Time)
		GetCodeExpiresIn() time.Duration
		SetCodeExpiresIn(time.Duration)

		GetAccess() string
		SetAccess(string)
		GetAccessCreateAt() time.Time
		SetAccessCreateAt(time.Time)
		GetAccessExpiresIn() time.Duration
		SetAccessExpiresIn(time.Duration)

		GetRefresh() string
		SetRefresh(string)
		GetRefreshCreateAt() time.Time
		SetRefreshCreateAt(time.Time)
		GetRefreshExpiresIn() time.Duration
		SetRefreshExpiresIn(time.Duration)
	}
)
