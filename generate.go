package oauth2

import "time"

type (
	// TokenData 提供生成令牌的基础数据
	TokenData struct {
		Client     ClientInfo    // 客户端信息
		UserID     string        // 用户标识
		Scope      string        // 权限范围
		CreateAt   time.Time     // 创建时间
		ExpiresIn  time.Duration // 有效期
		Identifier string        // 唯一标识码
	}

	// AuthorizeGenerate 授权令牌生成接口
	AuthorizeGenerate interface {
		// 生成授权令牌
		Token(data *TokenData) (string, error)

		// 验证令牌的有效性
		Verify(token string, data *TokenData) (bool, error)
	}

	// TokenGenerate 访问令牌生成接口
	TokenGenerate interface {
		// 生成访问令牌
		AccessToken(data *TokenData) (string, error)

		// 生成刷新令牌
		RefreshToken(data *TokenData) (string, error)
	}
)
