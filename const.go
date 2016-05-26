package oauth2

const (
	// DefaultRandomCodeLen 默认随机码的长度
	DefaultRandomCodeLen = 6
	// DefaultACExpiresIn 默认授权码模式的授权码有效期(10分钟)
	DefaultACExpiresIn = 60 * 10
	// DefaultATExpiresIn 默认授权码模式的访问令牌有效期(7天)
	DefaultATExpiresIn = 60 * 60 * 24 * 7
	// DefaultRTExpiresIn 默认授权码模式的更新令牌有效期(30天)
	DefaultRTExpiresIn = 60 * 60 * 24 * 30
	// DefaultIATExpiresIn 默认简化模式的访问令牌有效期(1小时)
	DefaultIATExpiresIn = 60 * 60
	// DefaultCCATExpiresIn 默认客户端模式的访问令牌有效期(1天)
	DefaultCCATExpiresIn = 60 * 60 * 24
)

// STATUS 提供一些状态标识
type STATUS byte

const (
	// Deleted 删除状态
	Deleted STATUS = iota
	// Actived 激活状态
	Actived
	// Blocked 冻结状态
	Blocked
	// Expired 过期状态
	Expired
)
