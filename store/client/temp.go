package client

import (
	"errors"

	"gopkg.in/oauth2.v2"
	"gopkg.in/oauth2.v2/models"
)

// NewTempStore 创建客户端临时存储实例
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

// TempStore 客户端信息的临时存储
type TempStore struct {
	data map[string]*models.Client
}

// GetByID 获取客户端信息
func (ts *TempStore) GetByID(id string) (cli oauth2.ClientInfo, err error) {
	if c, ok := ts.data[id]; ok {
		cli = c
		return
	}
	err = errors.New("not found")
	return
}
