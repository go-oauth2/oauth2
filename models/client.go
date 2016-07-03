package models

// Client 客户端信息
type Client struct {
	ID     string // 客户端ID
	Secret string // 密钥
	Domain string // 域名url
}

// GetID 客户端ID
func (c *Client) GetID() string {
	return c.ID
}

// GetSecret 客户端秘钥
func (c *Client) GetSecret() string {
	return c.Secret
}

// GetDomain 域名URL
func (c *Client) GetDomain() string {
	return c.Domain
}

// GetUserData 用户数据
func (c *Client) GetUserData() interface{} {
	return nil
}
