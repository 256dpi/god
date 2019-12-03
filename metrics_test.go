package god

import (
	"time"
)

func ExampleMetrics() {
	Metrics()

	counter := NewCounter("counter", nil)
	counter.Add(1)
	counter.Add(2)

	timer := NewTimer("timer")
	timer.Add(time.Millisecond)
	timer.Add(time.Second)

	TrackFloat("track", func() float64 { return 2 })

	time.Sleep(1500 * time.Millisecond)

	// Output:
	// counter: 3 c/s ｜ timer: 1ms - 500.5ms - 1s ｜ track: 2.00
}
