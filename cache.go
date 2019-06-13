// The module implements the core logic of the cache, including: PUT, GET func

package cache

import (
	"errors"
	"runtime"
	"sync"
	"time"
)

const (
	// DefaultCleanupInterval indicates the default cleaning interval, which can be modified by SetCleanupInterval
	DefaultCleanupInterval = 30 * time.Second
)

// Cache exposed structure for users to use
type Cache struct {
	*cache
}

type cache struct {
	sync.RWMutex
	mapping sync.Map
	janitor *janitor
}

type item struct {
	Payload interface{}
	Expired time.Time
}

// Put item in Cache with its ttl
func (c *cache) Put(key interface{}, payload interface{}, ttl time.Duration) {
	c.mapping.Store(key, &item{
		Payload: payload,
		Expired: time.Now().Add(ttl),
	})
}

// Get item in Cache, and drop when it expired
func (c *cache) Get(key interface{}) interface{} {
	v, ok := c.mapping.Load(key)
	if !ok {
		return nil
	}
	i := v.(*item)
	if time.Since(i.Expired) > 0 {
		c.Del(key)
		return nil
	}
	return i.Payload
}

// Exists method is used to check if a given key exists
func (c *cache) Exists(key interface{}) bool {
	_, ok := c.mapping.Load(key)
	return ok
}

// Delete the given key
func (c *cache) Del(key interface{}) {
	c.mapping.Delete(key)
}

// Increment the integer value of a key by the given amount
func (c *cache) IncrBy(key interface{}, increment int) error {
	c.RLock()
	defer c.RUnlock()

	v, ok := c.mapping.Load(key)
	if !ok {
		return errors.New("key not exists")
	}
	i := v.(*item)
	if time.Since(i.Expired) > 0 {
		c.Del(key)
		return errors.New("key is expired")
	}

	itm := i.Payload
	switch itm.(type) {
	case int:
		itm = itm.(int) + increment
	case int32:
		itm = itm.(int32) + int32(increment)
	case int64:
		itm = itm.(int64) + int64(increment)
	case uint:
		itm = itm.(uint) + uint(increment)
	case uint32:
		itm = itm.(uint32) + uint32(increment)
	case uint64:
		itm = itm.(uint64) + uint64(increment)
	default:
		return errors.New("the value of item is not of type (u)int, (u)int32, and (u)int64")
	}

	i.Payload = itm
	c.mapping.Store(key, i)

	return nil
}

// INCR key Increment the integer value of a key by one
func (c *cache) Incr(key interface{}) error {
	return c.IncrBy(key, 1)
}

func (c *cache) SetCleanupInterval(interval time.Duration) {
	c.janitor.stopJanitor()
	go c.janitor.process(c)
	c.janitor.interval <- interval
}

func (c *cache) cleanup() {
	c.mapping.Range(func(k, v interface{}) bool {
		key := k.(string)
		itm := v.(*item)
		if time.Since(itm.Expired) > 0 {
			c.mapping.Delete(key)
		}
		return true
	})
}

// New return *Cache
func New() *Cache {
	j := &janitor{
		stop:     make(chan struct{}),
		interval: make(chan time.Duration),
	}
	c := &cache{janitor: j}
	C := &Cache{c}
	go j.process(c)
	j.interval <- DefaultCleanupInterval

	runtime.SetFinalizer(C, func(c *Cache) {
		c.janitor.stopJanitor()
	})

	return C
}
