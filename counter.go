package god

import (
	"fmt"
)

// Counter is a simple operations counter.
type Counter struct {
	total int
}

// NewCounter will create and return a Counter.
func NewCounter(name string) *Counter {
	c := &Counter{}
	add(name, c)
	return c
}

// Add will increment the counter.
func (c *Counter) Add(n int) {
	c.total += n
}

func (c *Counter) string() string {
	return fmt.Sprintf("%.2f c/s", float64(c.total))
}

func (c *Counter) reset() {
	c.total = 0
}
