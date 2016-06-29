package oauth2

import "time"

type (
	// GenerateBasic 提供生成令牌的基础数据
	GenerateBasic struct {
		Client   ClientInfo // 客户端信息
		UserID   string     // 用户标识
		CreateAt time.Time  // 创建时间
	}

	// AuthorizeGenerate 授权令牌生成接口
	AuthorizeGenerate interface {
		// 授权令牌
		Token(data *GenerateBasic) (code string, err error)
	}

	// AccessGenerate 访问令牌生成接口
	AccessGenerate interface {
		// 访问令牌、更新令牌
		Token(data *GenerateBasic, isGenRefresh bool) (access, refresh string, err error)
	}
)
