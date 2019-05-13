package god

import (
	"fmt"
	"sync"
)

// Counter is a simple operations counter.
type Counter struct {
	total int
	mutex sync.Mutex
}

// NewCounter will create and return a counter.
func NewCounter(name string) *Counter {
	c := &Counter{}
	add(name, c)
	return c
}

// Add will increment the counter.
func (c *Counter) Add(n int) {
	c.mutex.Lock()
	c.total += n
	c.mutex.Unlock()
}

func (c *Counter) string() string {
	// get value
	c.mutex.Lock()
	total := c.total
	c.mutex.Unlock()

	return fmt.Sprintf("%.2f c/s", float64(total))
}

func (c *Counter) reset() {
	c.mutex.Lock()
	c.total = 0
	c.mutex.Unlock()
}
