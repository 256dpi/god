package god

import (
	"time"
)

func ExampleMetrics() {
	// enable
	Metrics()

	// counter
	counter := NewCounter("counter", nil)
	counter.Add(1)
	counter.Add(2)

	// timer
	timer := NewTimer("timer")
	timer.Add(time.Millisecond)
	timer.Add(time.Second)

	// track
	TrackInt("int", func() int64 { return 2 })
	TrackFloat("float", func() float64 { return 4.2 })

	time.Sleep(1500 * time.Millisecond)

	// Output:
	// counter: 3 c/s - float: 4.20 - int: 2 - timer: 1ms/500.5ms/1s
}
