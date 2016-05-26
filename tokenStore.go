package oauth2

import (
	"time"
)

// Token 令牌信息
type Token struct {
	ID           int64         `bson:"ID"`           // 唯一标识(自增ID)
	ClientID     string        `bson:"ClientID"`     // 客户端标识
	UserID       string        `bson:"UserID"`       // 用户标识
	AccessToken  string        `bson:"AccessToken"`  // 访问令牌
	ATCreateAt   int64         `bson:"ATCreateAt"`   // 访问令牌创建时间（时间戳）
	ATExpiresIn  time.Duration `bson:"ATExpiresIn"`  // 访问令牌有效期(单位秒)
	RefreshToken string        `bson:"RefreshToken"` // 更新令牌
	RTCreateAt   int64         `bson:"RTCreateAt"`   // 更新令牌创建时间（时间戳）
	RTExpiresIn  time.Duration `bson:"RTExpiresIn"`  // 更新令牌有效期(单位秒)
	Scope        string        `bson:"Scope"`        // 申请的权限范围
	CreateAt     int64         `bson:"CreateAt"`     // 创建时间(时间戳)
	Status       STATUS        `bson:"Status"`       // 令牌状态
}

// TokenStore 令牌存储接口(持久化存储)
type TokenStore interface {
	// Create 创建新的令牌，返回令牌ID
	// 如果创建发生异常，则返回错误
	Create(item Token) (int64, error)

	// Update 根据ID更新令牌信息
	// info 需要更新的字段信息(字段名称与结构体的字段名保持一致)
	// 如果更新发生异常，则返回错误
	Update(id int64, info map[string]interface{}) error

	// GetByAccessToken 根据访问令牌，获取令牌信息
	// 如果不存则返回nil
	GetByAccessToken(accessToken string) (*Token, error)

	// GetByRefreshToken 根据更新令牌，获取令牌信息
	// 如果不存则返回nil
	GetByRefreshToken(refreshToken string) (*Token, error)
}
