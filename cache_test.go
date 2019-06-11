package cache

import (
	"testing"
	"time"
)

func TestCache_Basic(t *testing.T) {
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
