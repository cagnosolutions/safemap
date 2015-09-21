package store

import (
    "sync"
)

type SafeMapStore struct {
    Stores map[string]*SafeMap
    sync.RWMutex
}

func NewSafeMapStore(shardCount int) *SafeMapStore {
    if shardCount == 0 || shardCount%2 != 0 {
		shardCount = 16
	}
	SHARD_COUNT = shardCount
    return &SafeMapStore{
        Stores: make(map[string]*SafeMap)
    }
}

func (sms *SafeMapStore) Set(key, fld string val interface{}) {
    m, ok := sms.GetStore(key)
    if !ok {
        sms.Lock()
        sms.Stores[storeName] = NewSafeMap(SHARD_COUNT)
        sms.Unlock()
    }
    m.Set(fld, val)
}

func (sms *SafeMapStore) Get(key, fld string) (interface{}, bool) {
    if m, ok := sms.GetStore(key); ok {
        return m.Get(fld)
    }
    return nil, false
}

func (sms *SafeMapStore) Del(key, fld string) {
    if m, ok := sms.GetStore(key); ok {
        m.Del(fld)
    }
}

func (sms *SafeMapStore) AddStore(key string) {
    if _, ok := sms.GetStore(key); !ok {
        sms.Lock()
        sms.Stores[key] = NewSafeMap(SHARD_COUNT)
        sms.Unlock()
    }
}

func (sms *SafeMapStore) GetStore(key string) (*SafeMap, bool) {
    sms.RLock()
    m, ok := sms.Stores[key]
    sms.RUnlock()
    return m, ok
}

func (sms *SafeMapStore) DelStore(key string) {
    if _, ok := sms.GetStore(key); ok {
        sms.Lock()
        delete(sms.Stores, key)
        sms.Unlock()
    }
}
