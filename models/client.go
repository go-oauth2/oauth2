package models

// NewClient Create to client model instance
func NewClient() *Client {
	return &Client{}
}

// Client Client model
type Client struct {
	ID     string // The client id
	Secret string // The client secret
	Domain string // The client domain
}

// GetID The client id
func (c *Client) GetID() string {
	return c.ID
}

// GetSecret The client domain
func (c *Client) GetSecret() string {
	return c.Secret
}

// GetDomain The client domain
func (c *Client) GetDomain() string {
	return c.Domain
}

// GetExtraData The extension data related to the client
func (c *Client) GetExtraData() interface{} {
	return nil
}
