package cache

import (
	"time"
)

type janitor struct {
	stop     chan struct{}
}

var janitorInterval = make(chan time.Duration)

func (j *janitor) stopJanitor() {
	j.stop <- struct{}{}
}

func (j *janitor) process(c *cache) {
	interval := <- janitorInterval
	ticker := time.NewTicker(interval)

Loop:
	for {
		select {
		case <-ticker.C:
			c.cleanup()
		case <-j.stop:
			ticker.Stop()
			break Loop
		}
	}
}
