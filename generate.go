package oauth2

import "time"

type (
	// TokenGenerateBasic 提供生成令牌的基础数据
	TokenGenerateBasic struct {
		Client   ClientInfo // 客户端信息
		UserID   string     // 用户标识
		CreateAt time.Time  // 创建时间
	}

	// AuthorizeTokenGenerate 授权令牌生成接口
	AuthorizeTokenGenerate interface {
		// 授权令牌
		Token(data *TokenGenerateBasic) (string, error)
	}

	// TokenGenerate 令牌生成接口
	TokenGenerate interface {
		// 生成令牌
		Token(data *TokenGenerateBasic, isGenRefresh bool) (string, string, error)
	}
)
