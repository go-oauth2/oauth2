package store

import (
	"encoding/json"
	"time"

	"github.com/tidwall/buntdb"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/models"
	"gopkg.in/oauth2.v3/utils/uuid"
)

// NewMemoryTokenStore create a token store instance based on memory
func NewMemoryTokenStore() (oauth2.TokenStore, error) {
	return NewFileTokenStore(":memory:")
}

// NewFileTokenStore create a token store instance based on file
func NewFileTokenStore(filename string) (oauth2.TokenStore, error) {
	db, err := buntdb.Open(filename)
	if err != nil {
		return nil, err
	}
	return &TokenStore{db: db}, nil
}

// TokenStore token storage based on buntdb(https://github.com/tidwall/buntdb)
type TokenStore struct {
	db *buntdb.DB
}

// Create create and store the new token information
func (ts *TokenStore) Create(info oauth2.TokenInfo) error {
	ct := time.Now()
	jv, err := json.Marshal(info)
	if err != nil {
		return err
	}
	err = ts.db.Update(func(tx *buntdb.Tx) error {
		if code := info.GetCode(); code != "" {
			_, _, err := tx.Set(code, string(jv), &buntdb.SetOptions{Expires: true, TTL: info.GetCodeExpiresIn()})
			return err
		}

		basicID := uuid.Must(uuid.NewRandom()).String()
		aexp := info.GetAccessExpiresIn()
		rexp := aexp
		if refresh := info.GetRefresh(); refresh != "" {
			rexp = info.GetRefreshCreateAt().Add(info.GetRefreshExpiresIn()).Sub(ct)
			if aexp.Seconds() > rexp.Seconds() {
				aexp = rexp
			}
			if _, _, err = tx.Set(refresh, basicID, &buntdb.SetOptions{Expires: true, TTL: rexp}); err != nil {
				return err
			}
		}
		if _, _, err := tx.Set(basicID, string(jv), &buntdb.SetOptions{Expires: true, TTL: rexp}); err != nil {
			return err
		}
		_, _, err := tx.Set(info.GetAccess(), basicID, &buntdb.SetOptions{Expires: true, TTL: aexp})
		return err
	})

	return err
}

// remove key
func (ts *TokenStore) remove(key string) error {
	err := ts.db.Update(func(tx *buntdb.Tx) error {
		_, err := tx.Delete(key)
		return err
	})
	if err == buntdb.ErrNotFound {
		return nil
	}
	return err
}

// RemoveByCode use the authorization code to delete the token information
func (ts *TokenStore) RemoveByCode(code string) error {
	return ts.remove(code)
}

// RemoveByAccess use the access token to delete the token information
func (ts *TokenStore) RemoveByAccess(access string) error {
	return ts.remove(access)
}

// RemoveByRefresh use the refresh token to delete the token information
func (ts *TokenStore) RemoveByRefresh(refresh string) error {
	return ts.remove(refresh)
}

func (ts *TokenStore) getData(key string) (oauth2.TokenInfo, error) {
	var tm models.Token
	err := ts.db.View(func(tx *buntdb.Tx) error {
		jv, err := tx.Get(key)
		if err != nil {
			return err
		}
		return json.Unmarshal([]byte(jv), &tm)
	})
	//TODO: The next 2 if's are here to make existing tests pass.
	// The rest of the function can actually be condensed down to the following:
	// if err == buntdb.ErrNotFound { err = nil }
	// return &tm, err
	if err == buntdb.ErrNotFound { //TODO: See above
		return nil, nil
	}
	if err != nil { // TODO: See above
		return nil, err
	}
	return &tm, nil
}

func (ts *TokenStore) getBasicID(key string) (string, error) {
	var v string
	err := ts.db.View(func(tx *buntdb.Tx) error {
		var err error
		v, err = tx.Get(key)
		if err != nil {
			return err
		}
		return nil
	})
	if err == buntdb.ErrNotFound {
		err = nil
	}
	return v, err
}

// GetByCode use the authorization code for token information data
func (ts *TokenStore) GetByCode(code string) (oauth2.TokenInfo, error) {
	return ts.getData(code)
}

// GetByAccess use the access token for token information data
func (ts *TokenStore) GetByAccess(access string) (oauth2.TokenInfo, error) {
	basicID, err := ts.getBasicID(access)
	if err != nil {
		return nil, err
	}
	return ts.getData(basicID)
}

// GetByRefresh use the refresh token for token information data
func (ts *TokenStore) GetByRefresh(refresh string) (oauth2.TokenInfo, error) {
	basicID, err := ts.getBasicID(refresh)
	if err != nil {
		return nil, err
	}
	return ts.getData(basicID)
}
