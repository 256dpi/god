package god

import (
	"fmt"
	"sync"
	"time"
)

// Timer is a simple operations timer.
type Timer struct {
	list  []time.Duration
	mutex sync.Mutex
}

// NewTimer will create and return a timer.
func NewTimer(name string) *Timer {
	t := &Timer{}
	add(name, t)
	return t
}

// Add will add the duration to the timer.
func (t *Timer) Add(d time.Duration) {
	t.mutex.Lock()
	t.list = append(t.list, d)
	t.mutex.Unlock()
}

// Measure can be used to track a function using defer.
func (t *Timer) Measure() func() {
	now := time.Now()
	return func() {
		t.Add(time.Since(now))
	}
}

func (t *Timer) string() string {
	// lock mutex
	t.mutex.Lock()

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

	// unlock mutex
	t.mutex.Unlock()

	return fmt.Sprintf("%s - %s - %s", min.String(), mean.String(), max.String())
}

func (t *Timer) reset() {
	t.mutex.Lock()
	t.list = nil
	t.mutex.Unlock()
}
