package god

import "fmt"

type tracker struct {
	fn func() string
}

// Track will track a string over time.
func Track(name string, fn func() string) {
	add(name, &tracker{
		fn: fn,
	})
}

// TrackFloat will track a float over time.
func TrackFloat(name string, fn func() float64) {
	Track(name, func() string {
		return fmt.Sprintf("%.2f", fn())
	})
}

func (c *tracker) string() string {
	return c.fn()
}

func (c *tracker) reset() {}
