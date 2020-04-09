package god

import (
	"fmt"
	"sync/atomic"
)

// Counter is a simple operations counter.
type Counter struct {
	total int64
	fmt   func(int) string
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
	Register(name, c)

	return c
}

// Add will increment the counter.
func (c *Counter) Add(n int) {
	atomic.AddInt64(&c.total, int64(n))
}

// Collect implements the Metric interface.
func (c *Counter) Collect() string {
	str := c.fmt(int(atomic.LoadInt64(&c.total)))
	atomic.StoreInt64(&c.total, 0)
	return str
}
