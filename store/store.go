package store

import (
	"sync"

	"github.com/cagnosolutions/safemap"
)

type SafeMapStore struct {
	SafeMaps map[string]*safemap.SafeMap
	sync.RWMutex
}

func NewSafeMapStore(shardCount int) *SafeMapStore {
	if shardCount == 0 || shardCount%2 != 0 {
		shardCount = 16
	}
	safemap.SHARD_COUNT = shardCount
	return &SafeMapStore{
		SafeMaps: make(map[string]*safemap.SafeMap),
	}
}

func (sms *SafeMapStore) Set(key, fld string, val interface{}) bool {
	sm, ok := sms.GetSafeMap(key)
	if !ok {
		sms.Lock()
		sms.SafeMaps[key] = safemap.NewSafeMap(safemap.SHARD_COUNT)
		sm = sms.SafeMaps[key]
		sms.Unlock()
	}
	return sm.Set(fld, val)
}

func (sms *SafeMapStore) Get(key, fld string) (interface{}, bool) {
	if sm, ok := sms.GetSafeMap(key); ok {
		return sm.Get(fld)
	}
	return nil, false
}

func (sms *SafeMapStore) Del(key, fld string) bool {
	if sm, ok := sms.GetSafeMap(key); ok {
		return sm.Del(fld)
	}
	return true
}

func (sms *SafeMapStore) AddStore(key string) bool {
	if _, ok := sms.GetSafeMap(key); !ok {
		sms.Lock()
		sms.SafeMaps[key] = safemap.NewSafeMap(safemap.SHARD_COUNT)
		sms.Unlock()
	}
	sms.RLock()
	_, ok := sms.SafeMaps[key]
	sms.RUnlock()
	return ok
}

func (sms *SafeMapStore) GetSafeMap(key string) (*safemap.SafeMap, bool) {
	sms.RLock()
	sm, ok := sms.SafeMaps[key]
	sms.RUnlock()
	return sm, ok
}

func (sms *SafeMapStore) DelStore(key string) bool {
	var ok bool
	if _, ok := sms.GetSafeMap(key); ok {
		sms.Lock()
		delete(sms.SafeMaps, key)
		_, ok = sms.SafeMaps[key]
		sms.Unlock()
	}
	return !ok
}
