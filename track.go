package god

import "fmt"

type tracker struct {
	fn func() float64
}

// Track will track the length of a value.
func Track(name string, fn func() float64) {
	add(name, &tracker{
		fn: fn,
	})
}

func (c *tracker) string() string {
	return fmt.Sprintf("%.2f", c.fn())
}

func (c *tracker) reset() {}
