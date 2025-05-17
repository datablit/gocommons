package syncache

import (
	"time"

	"github.com/patrickmn/go-cache"
	"golang.org/x/sync/singleflight"
)

// Cache is a generic thread-safe in-memory cache with deduplicated loading.
type Cache[T any] struct {
	mem *cache.Cache
	sg  *singleflight.Group
	ttl time.Duration
}

// New creates a new Cache with the given TTL.
func New[T any](ttl time.Duration) *Cache[T] {
	return &Cache[T]{
		mem: cache.New(ttl, ttl*2),
		sg:  &singleflight.Group{},
		ttl: ttl,
	}
}

// GetOrLoad returns cached value or runs the loader (only once per key concurrently).
func (c *Cache[T]) GetOrLoad(key string, load func() (T, error)) (T, error) {
	// Try from cache
	if val, ok := c.mem.Get(key); ok {
		return val.(T), nil
	}

	// Deduplicate concurrent loads
	val, err, _ := c.sg.Do(key, func() (interface{}, error) {
		v, err := load()
		if err != nil {
			var zero T
			return zero, err
		}
		c.mem.Set(key, v, c.ttl)
		return v, nil
	})
	return val.(T), err
}

// Set manually sets a value in the cache for the given key.
func (c *Cache[T]) Set(key string, val T) {
	c.mem.Set(key, val, c.ttl)
}

// Delete removes a key from the cache.
func (c *Cache[T]) Delete(key string) {
	c.mem.Delete(key)
}

// ClearAll removes all items from the cache.
func (c *Cache[T]) ClearAll() {
	c.mem.Flush()
}
