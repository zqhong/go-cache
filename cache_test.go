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
	if c.Exists("not-found") {
		t.Error("The key named 'not-found' should not exist.")
	}

	c.Put("str", "go-cache", 100*time.Millisecond)
	s := c.Get("str")
	if s.(string) != "go-cache" {
		t.Error("it should receive 'go-cache'")
	}
	if !c.Exists("str") {
		t.Error("The key named 'str' should exist.")
	}

	c.Put("int", 1, 100*time.Millisecond)
	i := c.Get("int")
	if i.(int) != 1 {
		t.Error("it should receive 1")
	}

	c.Put("Test-Del", 1, 100*time.Millisecond)
	if !c.Exists("Test-Del") {
		t.Error("The key named 'Test-Del' should exist.")
	}
	c.Del("Test-Del")
	if c.Exists("Test-Del") {
		t.Error("The key named 'Test-Del' should not exist.")
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
