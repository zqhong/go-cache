package cache

import (
	"runtime"
	"testing"
	"time"
)

func TestCacheBasic(t *testing.T) {
	c := New()

	none := c.Get("not-found")
	if none != nil {
		t.Error("should receive nil")
	}

	c.Put("str", "go-cache", 100*time.Millisecond)
	s := c.Get("str")
	if s.(string) != "go-cache" {
		t.Error("should receive 'go-cache'")
	}

	c.Put("int", 1, 100*time.Millisecond)
	i := c.Get("int")
	if i.(int) != 1 {
		t.Error("should receive 1")
	}
}

func TestCacheAutoGC(t *testing.T) {
	sign := make(chan struct{})
	go func() {
		interval := 10 * time.Millisecond
		ttl := 15 * time.Millisecond
		c := New()
		c.SetCleanupInterval(interval)
		c.Put("int", 1, ttl)
		sign <- struct{}{}
	}()

	<-sign
	runtime.GC()
}
