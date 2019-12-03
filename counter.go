package god

import (
	"fmt"
	"sync"
)

// Counter is a simple operations counter.
type Counter struct {
	total int
	fmt   func(int) string
	mutex sync.Mutex
}

// NewCounter will create and return a counter.
func NewCounter(name string, formatter func(total int) string) *Counter {
	// set default formatter
	if formatter == nil {
		formatter = func(total int) string {
			return fmt.Sprintf("%d c/s", total)
		}
	}

	// create counter
	c := &Counter{
		fmt: formatter,
	}

	// add counter
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

	return c.fmt(total)
}

func (c *Counter) reset() {
	c.mutex.Lock()
	c.total = 0
	c.mutex.Unlock()
}
