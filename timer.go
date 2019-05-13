package god

import (
	"fmt"
	"time"
)

// Timer is a simple operations Timer.
type Timer struct {
	list []time.Duration
}

// NewTimer will create and return a Timer.
func NewTimer(name string) *Timer {
	t := &Timer{}
	add(name, t)
	return t
}

// Add will add the duration to the Timer.
func (t *Timer) Add(d time.Duration) {
	t.list = append(t.list, d)
}

func (t *Timer) string() string {
	// prepare mean
	var min time.Duration
	var mean time.Duration
	var max time.Duration

	// get mean
	if len(t.list) > 0 {
		// set min and max
		min = t.list[0]
		max = t.list[0]

		// prepare sum
		var sum time.Duration

		// iterate through all values
		for _, i := range t.list {
			// add up
			sum += i

			// check min
			if i < min {
				min = i
			}

			// check max
			if i > max {
				max = i
			}
		}

		// calculate mean
		mean = sum / time.Duration(len(t.list))
	}

	return fmt.Sprintf("%s - %s - %s", min.String(), mean.String(), max.String())
}

func (t *Timer) reset() {
	t.list = nil
}
