package store

import (
	"errors"
	"sync"

	"gopkg.in/oauth2.v3"
)

// ClientStore client storage service
type ClientStore interface {
	// GetByID according to the ID for the client information
	GetByID(id string) (cli oauth2.ClientInfo, err error)
	// Set set client information
	Set(id string, cli oauth2.ClientInfo) (err error)
}

// NewClientStore create a client storage service
func NewClientStore() ClientStore {
	return &clientStore{
		data: make(map[string]oauth2.ClientInfo),
	}
}

// clientStore client information storage
type clientStore struct {
	sync.RWMutex
	data map[string]oauth2.ClientInfo
}

// GetByID according to the ID for the client information
func (cs *clientStore) GetByID(id string) (cli oauth2.ClientInfo, err error) {
	cs.RLock()
	defer cs.RUnlock()
	if c, ok := cs.data[id]; ok {
		cli = c
		return
	}
	err = errors.New("not found")
	return
}

// Set set client information
func (cs *clientStore) Set(id string, cli oauth2.ClientInfo) (err error) {
	cs.Lock()
	defer cs.Unlock()
	cs.data[id] = cli
	return
}
