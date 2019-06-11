package cache

import (
	"time"
)

type janitor struct {
	interval time.Duration
	stop     chan bool
}

func (j *janitor) stopJanitor() {
	j.stop <- true
}

func (j *janitor) process(c *cache) {

	ticker := time.NewTicker(j.interval)

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
