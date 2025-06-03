package cache

import "sync"

type MemCache struct {
	cache sync.Map
}

func (m *MemCache) Get(key string) (any, error) {
	val, ok := m.cache.Load(key)
	if !ok {
		return nil, nil
	}
	return val, nil
}

func (m *MemCache) Set(key string, value any, expire int64) {
	m.cache.Store(key, value)
}

func NewMemCache() Cache {
	return &MemCache{
		cache: sync.Map{},
	}
}
