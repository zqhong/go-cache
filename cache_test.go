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
		t.Error("Should receive 'go-cache' here")
	}

	c.Put("int", 1, 100*time.Millisecond)
	i := c.Get("int")
	if i.(int) != 1 {
		t.Error("Should receive 1 here")
	}
}

func TestCache_Get(t *testing.T) {
	c := New()

	none := c.Get("not-found")
	if none != nil {
		t.Error("Should receive nil here")
	}

	c.Put("str", "go-cache", 100*time.Millisecond)
	s := c.Get("str")
	if s.(string) != "go-cache" {
		t.Error("Should receive 'go-cache' here")
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

func TestCache_IncrBy(t *testing.T) {
	c := New()

	err := c.IncrBy("not-found", 1)
	if err == nil {
		t.Error("Here should receive an error")
	}

	c.Put("Test-IncrBy1", 1, 1*time.Nanosecond)
	time.Sleep(2 * time.Nanosecond)
	err = c.IncrBy("Test-IncrBy1", 1)
	if err == nil {
		t.Error("Here should receive an error")
	}

	c.Put("Test-IncrBy2", "str", 100*time.Millisecond)
	err = c.IncrBy("Test-IncrBy2", 1)
	if err == nil {
		t.Error("Here should receive an error")
	}

	c.Put("Test-IncrBy3", int(1), 100*time.Millisecond)
	err = c.IncrBy("Test-IncrBy3", 100)
	if err != nil {
		t.Error("Here should receive an error")
	}
	if c.Get("Test-IncrBy3") != 101 {
		t.Error("Should receive 101 here")
	}
}

func TestCache_Incr(t *testing.T) {
	c := New()

	c.Put("Test-IncrBy", int(1), 100*time.Millisecond)
	err := c.IncrBy("Test-IncrBy", 1)
	if err != nil {
		t.Error("it should not receive an error")
	}
	if c.Get("Test-IncrBy") != 2 {
		t.Error("Should receive 2 here")
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
