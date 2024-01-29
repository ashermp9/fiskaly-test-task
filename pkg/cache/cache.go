package cache

import (
	"sync"
)

type Cache[K comparable, V any] struct {
	items map[K]V
	mu    sync.Mutex
	locks map[K]*sync.Mutex
}

func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		items: make(map[K]V),
		locks: make(map[K]*sync.Mutex),
	}
}
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, ok := c.items[key]
	return val, ok
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = value
}

func (c *Cache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

func (c *Cache[K, V]) Lock(key K) {
	c.mu.Lock()
	if _, exists := c.locks[key]; !exists {
		c.locks[key] = &sync.Mutex{}
	}
	lock := c.locks[key]
	c.mu.Unlock()
	lock.Lock()
}

func (c *Cache[K, V]) Unlock(key K) {
	c.mu.Lock()
	if lock, exists := c.locks[key]; exists {
		c.mu.Unlock()
		lock.Unlock()
	} else {
		c.mu.Unlock()
	}
}
