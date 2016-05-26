package oauth2

// Client 客户端的验证信息接口
type Client interface {
	// ID 客户端唯一标识
	ID() string
	// Secret 客户端秘钥
	Secret() string
	// Domain 客户端域名
	Domain() string
	// RetainData 保留数据
	RetainData() interface{}
}

// ClientStore 客户端存储接口（持久化存储）
type ClientStore interface {
	// GetByID 根据ID获取客户端信息;
	// 如果客户端不存在则返回nil
	GetByID(id string) (Client, error)
}

// DefaultClient 默认的客户端信息
type DefaultClient struct {
	ClientID     string `bson:"_id"`    // 客户端唯一标识
	ClientSecret string `bson:"Secret"` // 客户端秘钥
	ClientDomain string `bson:"Domain"` // 客户端域名
}

// ID Get ClientID
func (dc DefaultClient) ID() string {
	return dc.ClientID
}

// Secret Get ClientSecret
func (dc DefaultClient) Secret() string {
	return dc.ClientSecret
}

// Domain Get ClientDomain
func (dc DefaultClient) Domain() string {
	return dc.ClientDomain
}

// RetainData Get retain data
func (dc DefaultClient) RetainData() interface{} {
	return dc
}
