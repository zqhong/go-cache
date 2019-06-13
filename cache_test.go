package cache

import (
	"runtime"
	"testing"
	"time"
)

func TestCache_Put(t *testing.T) {
	c := New()

	c.Put("str", "go-cache", 100*time.Millisecond)
	s := c.Get("str")
	if s.(string) != "go-cache" {
		t.Error("it should receive 'go-cache'")
	}

	c.Put("int", 1, 100*time.Millisecond)
	i := c.Get("int")
	if i.(int) != 1 {
		t.Error("it should receive 1")
	}
}

func TestCache_Get(t *testing.T) {
	c := New()

	none := c.Get("not-found")
	if none != nil {
		t.Error("should receive nil")
	}

	c.Put("str", "go-cache", 100*time.Millisecond)
	s := c.Get("str")
	if s.(string) != "go-cache" {
		t.Error("it should receive 'go-cache'")
	}
}

func TestCache_Exists(t *testing.T) {
	c := New()

	if c.Exists("not-found") {
		t.Error("The key named 'not-found' should not exist.")
	}

	c.Put("str", "go-cache", 100*time.Millisecond)
	if !c.Exists("str") {
		t.Error("The key named 'str' should exist.")
	}
}

func TestCache_Del(t *testing.T) {
	c := New()

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
