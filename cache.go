// The module implements the core logic of the cache, including: PUT, GET func

package cache

import (
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
	mapping sync.Map
	janitor *janitor
}

type item struct {
	Payload interface{}
	Expired time.Time
}

// Put item in Cache with its ttl
func (c *cache) Put(k interface{}, x interface{}, d time.Duration) {
	c.mapping.Store(k, &item{
		Payload: x,
		Expired: time.Now().Add(d),
	})
}

// Get item in Cache, and drop when it expired
func (c *cache) Get(k interface{}) interface{} {
	v, ok := c.mapping.Load(k)
	if ok == false {
		return nil
	}
	i := v.(*item)
	if time.Since(i.Expired) > 0 {
		c.mapping.Delete(k)
		return nil
	}
	return i.Payload
}

func (c *cache) SetCleanupInterval(interval time.Duration) {
	c.janitor.stopJanitor()
	go c.janitor.process(c)
	janitorInterval <- interval
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
	}
	c := &cache{janitor: j}
	C := &Cache{c}
	go j.process(c)
	janitorInterval <- DefaultCleanupInterval

	runtime.SetFinalizer(C, func(c *Cache) {
		c.janitor.stopJanitor()
	})

	return C
}
