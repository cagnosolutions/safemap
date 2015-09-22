package safemap

import (
	"sync"

	"github.com/cagnosolutions/safemap/util"
)

var SHARD_COUNT int

type SafeMap []*Shard

type Shard struct {
	items map[string]interface{}
	sync.RWMutex
}

func NewSafeMap(shardCount int) *SafeMap {
	if shardCount == 0 || shardCount%2 != 0 {
		shardCount = 16
	}
	SHARD_COUNT = shardCount
	sm := make(SafeMap, SHARD_COUNT)
	for i := 0; i < SHARD_COUNT; i++ {
		sm[i] = &Shard{
			items: make(map[string]interface{}),
		}
	}
	return &sm
}

func (sm *SafeMap) GetShard(key string) *Shard {
	bucket := util.Sum32([]byte(key)) % uint32(SHARD_COUNT)
	return (*sm)[bucket]
}

func (sm *SafeMap) Set(key string, val interface{}) bool {
	shard := sm.GetShard(key)
	shard.Lock()
	shard.items[key] = val
	_, ok := shard.items[key]
	shard.Unlock()
	return ok
}

func (sm *SafeMap) Get(key string) (interface{}, bool) {
	shard := sm.GetShard(key)
	shard.RLock()
	val, ok := shard.items[key]
	shard.RUnlock()
	return val, ok
}

func (sm *SafeMap) Del(key string) bool {
	var ok bool
	if shard := sm.GetShard(key); shard != nil {
		shard.Lock()
		delete(shard.items, key)
		_, ok = shard.items[key]
		shard.Unlock()
	}
	return !ok
}

type EntrySet struct {
	Key string
	Val interface{}
}

func (sm *SafeMap) Iter() <-chan EntrySet {
	ch := make(chan EntrySet)
	go func() {
		for _, shard := range *sm {
			shard.RLock()
			for key, val := range shard.items {
				ch <- EntrySet{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}
