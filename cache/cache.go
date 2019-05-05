package cache

import (
	"sync"
	"time"
)

type CacheGetter func(key string) (*CacheItem, error)

type CacheItem struct {
	Content   interface{}
	ExpiresAt uint64
}

type Cache struct {
	pool map[string]*CacheItem
	lock sync.Mutex
}

func (c *Cache) Get(key string, onFailed CacheGetter) (*CacheItem, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	item, ok := c.pool[key]

	if ok {
		now := time.Now()
		return item, nil
	}
	if item, err := onFailed(key); err == nil {
		c.pool[key] = item
		return item, nil
	} else {
		return nil, err
	}
}
