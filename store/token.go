package store

import (
	"container/list"
	"strconv"
	"sync"
	"time"

	"gopkg.in/oauth2.v3"
)

// NewMemoryTokenStore Create a token store instance based on memory
// gcInterval Perform garbage collection intervals(The default is 30 seconds)
func NewMemoryTokenStore(gcInterval time.Duration) oauth2.TokenStore {
	if gcInterval == 0 {
		gcInterval = time.Second * 30
	}
	store := &MemoryTokenStore{
		gcInterval: gcInterval,
		basicList:  list.New(),
		data:       make(map[string]oauth2.TokenInfo),
		access:     make(map[string]string),
		refresh:    make(map[string]string),
	}
	go store.gc()
	return store
}

// MemoryTokenStore Memory storage for token
type MemoryTokenStore struct {
	gcInterval time.Duration
	globalID   int64
	lock       sync.RWMutex
	data       map[string]oauth2.TokenInfo
	access     map[string]string
	refresh    map[string]string
	basicList  *list.List
	listLock   sync.RWMutex
}

func (mts *MemoryTokenStore) gc() {
	time.AfterFunc(mts.gcInterval, func() {
		defer mts.gc()
		rmeles := make([]*list.Element, 0, 32)
		mts.listLock.RLock()
		ele := mts.basicList.Front()
		mts.listLock.RUnlock()
		for ele != nil {
			if rm := mts.gcElement(ele); rm {
				rmeles = append(rmeles, ele)
			}
			mts.listLock.RLock()
			ele = ele.Next()
			mts.listLock.RUnlock()
		}

		for _, e := range rmeles {
			mts.listLock.Lock()
			mts.basicList.Remove(e)
			mts.listLock.Unlock()
		}
	})
}

func (mts *MemoryTokenStore) gcElement(ele *list.Element) (rm bool) {
	basicID := ele.Value.(string)
	mts.lock.RLock()
	ti, ok := mts.data[basicID]
	mts.lock.RUnlock()
	if !ok {
		rm = true
		return
	}
	ct := time.Now()
	if refresh := ti.GetRefresh(); refresh != "" &&
		ti.GetRefreshCreateAt().Add(ti.GetRefreshExpiresIn()).Before(ct) {
		mts.lock.RLock()
		delete(mts.access, ti.GetAccess())
		delete(mts.refresh, refresh)
		delete(mts.data, basicID)
		mts.lock.RUnlock()
		rm = true
	} else if ti.GetAccessCreateAt().Add(ti.GetAccessExpiresIn()).Before(ct) {
		mts.lock.RLock()
		delete(mts.access, ti.GetAccess())
		if refresh := ti.GetRefresh(); refresh == "" {
			delete(mts.data, basicID)
			rm = true
		}
		mts.lock.RUnlock()
	}
	return
}

func (mts *MemoryTokenStore) getBasicID(id int64, info oauth2.TokenInfo) string {
	return info.GetClientID() + "_" + strconv.FormatInt(id, 10)
}

// Create Create and store the new token information
func (mts *MemoryTokenStore) Create(info oauth2.TokenInfo) (err error) {
	mts.lock.Lock()
	defer mts.lock.Unlock()
	mts.globalID++
	basicID := mts.getBasicID(mts.globalID, info)
	mts.data[basicID] = info
	mts.access[info.GetAccess()] = basicID
	if refresh := info.GetRefresh(); refresh != "" {
		mts.refresh[refresh] = basicID
	}

	mts.listLock.Lock()
	mts.basicList.PushBack(basicID)
	mts.listLock.Unlock()
	return
}

// RemoveByAccess Use the access token to delete the token information
func (mts *MemoryTokenStore) RemoveByAccess(access string) (err error) {
	mts.lock.RLock()
	v, ok := mts.access[access]
	if !ok {
		mts.lock.RUnlock()
		return
	}
	info := mts.data[v]
	mts.lock.RUnlock()

	mts.lock.Lock()
	defer mts.lock.Unlock()
	delete(mts.access, access)
	if refresh := info.GetRefresh(); refresh == "" {
		delete(mts.data, v)
	}
	return
}

// RemoveByRefresh Use the refresh token to delete the token information
func (mts *MemoryTokenStore) RemoveByRefresh(refresh string) (err error) {
	mts.lock.Lock()
	defer mts.lock.Unlock()
	delete(mts.refresh, refresh)

	return
}

// GetByAccess Use the access token for token information data
func (mts *MemoryTokenStore) GetByAccess(access string) (ti oauth2.TokenInfo, err error) {
	mts.lock.RLock()
	v, ok := mts.access[access]
	if !ok {
		mts.lock.RUnlock()
		return
	}
	ti = mts.data[v]
	mts.lock.RUnlock()
	return
}

// GetByRefresh Use the refresh token for token information data
func (mts *MemoryTokenStore) GetByRefresh(refresh string) (ti oauth2.TokenInfo, err error) {
	mts.lock.RLock()
	v, ok := mts.refresh[refresh]
	if !ok {
		mts.lock.RUnlock()
		return
	}
	ti = mts.data[v]
	mts.lock.RUnlock()
	return
}
