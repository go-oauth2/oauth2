package models

// Client client model
// 客户端客户端模型
type Client struct {
	ID     string
	Secret string
	Domain string
	UserID string
}

// GetID client id
// 获取客户端ID
func (c *Client) GetID() string {
	return c.ID
}

// GetSecret client secret
// 获取客户端秘钥
func (c *Client) GetSecret() string {
	return c.Secret
}

// GetDomain client domain
// 获取域名
func (c *Client) GetDomain() string {
	return c.Domain
}

// GetUserID user id
// 获取用户ID
func (c *Client) GetUserID() string {
	return c.UserID
}
