package store

import (
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
)

// NewTestClientStore create to client information store instance
func NewTestClientStore(clients ...*models.Client) oauth2.ClientStore {
	data := map[string]*models.Client{
		"1": &models.Client{
			ID:     "1",
			Secret: "11",
			Domain: "http://localhost",
		},
	}
	for _, cli := range clients {
		data[cli.ID] = cli
	}
	return &TestClientStore{
		data: data,
	}
}

// TestClientStore client information store
type TestClientStore struct {
	data map[string]*models.Client
}

// GetByID according to the ID for the client information
func (ts *TestClientStore) GetByID(id string) (cli oauth2.ClientInfo, err error) {
	if c, ok := ts.data[id]; ok {
		cli = c
	}
	return
}
