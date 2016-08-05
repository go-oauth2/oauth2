package models

// Client client model
type Client struct {
	ID     string
	Secret string
	Domain string
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
