// The module implements the core logic of the cache, including: PUT, GET func

package cache

import (
	"sync"
	"time"
)

// Cache exposed structure for users to use
type Cache struct {
	*cache
}

type cache struct {
	mapping sync.Map
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

// New return *Cache
func New() *Cache {
	c := &cache{}
	C := &Cache{c}
	return C
}
