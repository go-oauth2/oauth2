package oauth2

import "time"

type (
	// TokenGenerateData 提供生成令牌的基础数据
	TokenGenerateData struct {
		Client   ClientInfo // 客户端信息
		UserID   string     // 用户标识
		Scope    string     // 权限范围
		CreateAt time.Time  // 创建时间
	}

	// AuthorizeGenerate 授权令牌生成接口
	AuthorizeGenerate interface {
		// 授权令牌
		Token(data *TokenGenerateData) (string, error)
	}

	// TokenGenerate 令牌生成接口
	TokenGenerate interface {
		// 生成令牌
		Token(data *TokenGenerateData, isGenRefresh bool) (string, string, error)
	}
)
