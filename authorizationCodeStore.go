package oauth2

import (
	"time"
)

// ACInfo 授权码信息(Authorization Code Info)
type ACInfo struct {
	ID          int64         // 唯一标识
	ClientID    string        // 客户端标识
	UserID      string        // 用户标识
	RedirectURI string        // 重定向URI
	Scope       string        // 申请的权限范围
	Code        string        // 随机码
	CreateAt    int64         // 创建时间（时间戳）
	ExpiresIn   time.Duration // 有效期(单位秒)
}

// ACStore 授权码存储接口(临时存储，提供自动GC过期的元素)(Authorization Code Store)
type ACStore interface {
	// Put 将元素放入存储，返回存储的ID
	// 如果存储发生异常，则返回错误
	Put(item ACInfo) (int64, error)

	// TakeByID 根据ID取出元素
	// 如果元素找不到或发生异常，则返回错误
	TakeByID(id int64) (*ACInfo, error)
}
