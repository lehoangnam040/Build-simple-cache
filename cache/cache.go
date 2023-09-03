package cache

import (
	"sync"
	"time"
)

func cleaner[V any](c *Cache[V], resolution time.Duration) {
	ticker := time.NewTicker(resolution)

	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-c.done:
			ticker.Stop()
			return
		}
	}
}

type Cache[V any] struct {
	kv   map[string]item[V]
	mu   sync.RWMutex
	done chan struct{}
}

type item[V any] struct {
	deadline int64 // Unix nano
	value    V
}

func New[V any]() *Cache[V] {
	c := Cache[V]{
		kv: map[string]item[V]{},
	}
	go cleaner[V](&c, time.Millisecond)
	return &c
}

func (c *Cache[V]) cleanup() {
	now := time.Now().UnixNano()
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, item := range c.kv {
		if item.deadline < now {
			delete(c.kv, key)
		}
	}
}

func (c *Cache[V]) Set(key string, value V, ttl time.Duration) {
	item := item[V]{
		deadline: time.Now().Unix() + int64(ttl),
		value:    value,
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.kv[key] = item
}

func (c *Cache[V]) Get(key string) (V, bool) {
	c.mu.Lock()
	v, exists := c.kv[key]
	c.mu.Unlock()
	return v.value, exists
}
