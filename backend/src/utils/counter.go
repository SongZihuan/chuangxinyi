package utils

import "sync"

type Counter struct {
	count int64
	mutex sync.Mutex
}

func NewCounter() Counter {
	return Counter{
		count: 0,
	}
}

func (c *Counter) Add(d int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.count += d
}

func (c *Counter) Sub(d int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.count -= d
}

func (c *Counter) Set(d int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.count = d
}

func (c *Counter) Get() int64 {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.count
}
