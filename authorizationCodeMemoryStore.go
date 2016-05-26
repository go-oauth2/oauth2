package oauth2

import (
	"container/list"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

// NewACMemoryStore 创建授权码的内存存储
// gcInterval GC周期（单位秒，默认60秒执行一次）
func NewACMemoryStore(gcInterval int64) ACStore {
	if gcInterval == 0 {
		gcInterval = 60
	}
	memStore := &ACMemoryStore{
		gcInterval: time.Second * time.Duration(gcInterval),
		data:       list.New(),
	}
	go memStore.gc()
	return memStore
}

// ACMemoryStore 提供授权码的内存存储
type ACMemoryStore struct {
	sync.RWMutex
	globalID   int64
	gcInterval time.Duration
	data       *list.List
}

func (am *ACMemoryStore) gc() {
	time.AfterFunc(am.gcInterval, func() {
		defer am.gc()
		for {
			am.RLock()
			ele := am.data.Front()
			if ele == nil {
				am.RUnlock()
				break
			}
			item := ele.Value.(ACInfo)
			am.RUnlock()
			if (item.CreateAt + int64(item.ExpiresIn/time.Second)) < time.Now().Unix() {
				am.Lock()
				am.data.Remove(ele)
				am.Unlock()
				continue
			}
			break
		}
	})
}

// Put Put item
func (am *ACMemoryStore) Put(item ACInfo) (int64, error) {
	am.Lock()
	defer am.Unlock()
	atomic.AddInt64(&am.globalID, 1)
	item.ID = am.globalID
	am.data.PushBack(item)
	return item.ID, nil
}

// TakeByID Take item by ID
func (am *ACMemoryStore) TakeByID(id int64) (*ACInfo, error) {
	am.RLock()
	var takeEle *list.Element
	for ele := am.data.Back(); ele != nil; ele = ele.Prev() {
		item := ele.Value.(ACInfo)
		if item.ID == id {
			takeEle = ele
			break
		}
	}
	if takeEle == nil {
		am.RUnlock()
		return nil, errors.New("Item not found")
	}
	item := takeEle.Value.(ACInfo)
	am.RUnlock()
	am.Lock()
	am.data.Remove(takeEle)
	am.Unlock()
	return &item, nil
}
