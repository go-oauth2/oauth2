package models

// Client client model
type Client struct {
	ID           string
	Secret       string
	Domain       string
	Public       bool
	UserID       string
	RedirectURIs []string
}

// GetID client id
func (c *Client) GetID() string {
	return c.ID
}

// GetSecret client secret
func (c *Client) GetSecret() string {
	return c.Secret
}

// GetDomain client domain
func (c *Client) GetDomain() string {
	return c.Domain
}

// GetRedirectURIs client redirect uris
func (c *Client) GetRedirectURIs() []string {
	return c.RedirectURIs
}

// IsPublic public
func (c *Client) IsPublic() bool {
	return c.Public
}

// GetUserID user id
func (c *Client) GetUserID() string {
	return c.UserID
}
