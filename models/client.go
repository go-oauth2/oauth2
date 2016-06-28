package models

// Client 客户端信息
type Client struct {
	ClientID string `bson:"ClientID"` // 客户端ID
	Secret   string `bson:"Secret"`   // 密钥
	Domain   string `bson:"Domain"`   // 域名url
}

// GetID 客户端ID
func (c *Client) GetID() string {
	return c.ClientID
}

// GetSecret 客户端秘钥
func (c *Client) GetSecret() string {
	return c.Secret
}

// GetDomain 域名URL
func (c *Client) GetDomain() string {
	return c.Domain
}

// GetRetainData 自定义数据
func (c *Client) GetRetainData() interface{} {
	return nil
}
