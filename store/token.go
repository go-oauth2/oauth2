package store

import (
	"encoding/json"
	"time"

	"github.com/satori/go.uuid"
	"github.com/tidwall/buntdb"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
)

// NewMemoryTokenStore create a token store instance based on memory
func NewMemoryTokenStore() (store oauth2.TokenStore, err error) {
	store, err = NewFileTokenStore(":memory:")
	return
}

// NewFileTokenStore create a token store instance based on file
func NewFileTokenStore(filename string) (store oauth2.TokenStore, err error) {
	db, err := buntdb.Open(filename)
	if err != nil {
		return
	}
	store = &TokenStore{db: db}
	return
}

// TokenStore token storage based on buntdb(https://github.com/tidwall/buntdb)
type TokenStore struct {
	db *buntdb.DB
}

// Create create and store the new token information
func (ts *TokenStore) Create(info oauth2.TokenInfo) (err error) {
	ct := time.Now()
	jv, err := json.Marshal(info)
	if err != nil {
		return
	}
	err = ts.db.Update(func(tx *buntdb.Tx) (err error) {
		if code := info.GetCode(); code != "" {
			_, _, err = tx.Set(code, string(jv), &buntdb.SetOptions{Expires: true, TTL: info.GetCodeExpiresIn()})
			return
		}
		basicID := uuid.NewV4().String()
		aexp := info.GetAccessExpiresIn()
		rexp := aexp
		if refresh := info.GetRefresh(); refresh != "" {
			rexp = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()).Sub(ct)
			if aexp.Seconds() > rexp.Seconds() {
				aexp = rexp
			}
			_, _, err = tx.Set(refresh, basicID, &buntdb.SetOptions{Expires: true, TTL: rexp})
			if err != nil {
				return
			}
		}
		_, _, err = tx.Set(basicID, string(jv), &buntdb.SetOptions{Expires: true, TTL: rexp})
		if err != nil {
			return
		}
		_, _, err = tx.Set(info.GetAccess(), basicID, &buntdb.SetOptions{Expires: true, TTL: aexp})
		return
	})
	return
}

// remove key
func (ts *TokenStore) remove(key string) (err error) {
	verr := ts.db.Update(func(tx *buntdb.Tx) (err error) {
		_, err = tx.Delete(key)
		return
	})
	if verr == buntdb.ErrNotFound {
		return
	}
	err = verr
	return
}

// RemoveByCode use the authorization code to delete the token information
func (ts *TokenStore) RemoveByCode(code string) (err error) {
	err = ts.remove(code)
	return
}

// RemoveByAccess use the access token to delete the token information
func (ts *TokenStore) RemoveByAccess(access string) (err error) {
	err = ts.remove(access)
	return
}

// RemoveByRefresh use the refresh token to delete the token information
func (ts *TokenStore) RemoveByRefresh(refresh string) (err error) {
	err = ts.remove(refresh)
	return
}

func (ts *TokenStore) getData(key string) (ti oauth2.TokenInfo, err error) {
	verr := ts.db.View(func(tx *buntdb.Tx) (err error) {
		jv, err := tx.Get(key)
		if err != nil {
			return
		}
		var tm models.Token
		err = json.Unmarshal([]byte(jv), &tm)
		if err != nil {
			return
		}
		ti = &tm
		return
	})
	if verr != nil {
		if verr == buntdb.ErrNotFound {
			return
		}
		err = verr
	}
	return
}

func (ts *TokenStore) getBasicID(key string) (basicID string, err error) {
	verr := ts.db.View(func(tx *buntdb.Tx) (err error) {
		v, err := tx.Get(key)
		if err != nil {
			return
		}
		basicID = v
		return
	})
	if verr != nil {
		if verr == buntdb.ErrNotFound {
			return
		}
		err = verr
	}
	return
}

// GetByCode use the authorization code for token information data
func (ts *TokenStore) GetByCode(code string) (ti oauth2.TokenInfo, err error) {
	ti, err = ts.getData(code)
	return
}

// GetByAccess use the access token for token information data
func (ts *TokenStore) GetByAccess(access string) (ti oauth2.TokenInfo, err error) {
	basicID, err := ts.getBasicID(access)
	if err != nil {
		return
	}
	ti, err = ts.getData(basicID)
	return
}

// GetByRefresh use the refresh token for token information data
func (ts *TokenStore) GetByRefresh(refresh string) (ti oauth2.TokenInfo, err error) {
	basicID, err := ts.getBasicID(refresh)
	if err != nil {
		return
	}
	ti, err = ts.getData(basicID)
	return
}
