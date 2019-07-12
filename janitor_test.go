package gocache

import (
	"testing"
	"time"
)

func TestJanitorTTL(t *testing.T) {
	c := New()
	ttl := 10 * time.Millisecond
	interval := 20 * time.Millisecond
	c.SetCleanupInterval(interval)
	c.Put("int", 1, ttl)

	i := c.Get("int")
	if i.(int) != 1 {
		t.Error("it should receive 1")
	}

	time.Sleep(interval + 20*time.Microsecond)
	i = c.Get("int")
	if i != nil {
		t.Error("it should receive nil")
	}
}
