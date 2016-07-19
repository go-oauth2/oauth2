package client

import (
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
)

// NewTempStore Create to client information temporary store instance
func NewTempStore(clients ...*models.Client) oauth2.ClientStore {
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
	return &TempStore{
		data: data,
	}
}

// TempStore Client information store
type TempStore struct {
	data map[string]*models.Client
}

// GetByID According to the ID for the client information
func (ts *TempStore) GetByID(id string) (cli oauth2.ClientInfo, err error) {
	if c, ok := ts.data[id]; ok {
		cli = c
	}
	return
}
