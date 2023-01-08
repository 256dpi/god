package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/256dpi/god"
)

func main() {
	// init
	god.Init(god.Options{})

	// block events
	go func() {
		for {
			<-time.After(time.Millisecond)
		}
	}()

	// mutex events
	var m sync.Mutex
	for i := 0; i < 2; i++ {
		go func() {
			for {
				m.Lock()
				time.Sleep(time.Millisecond)
				m.Unlock()
			}
		}()
	}

	// trace
	go func() {
		for {
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			func() {
				e := god.Trace("foo", "bar")
				defer e.End()
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			}()
		}
	}()
	go func() {
		for {
			time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
			func() {
				e := god.Trace("foo", "baz")
				defer e.End()
				time.Sleep(time.Duration(rand.Intn(1500)) * time.Millisecond)
			}()
		}
	}()

	// metrics
	counter := god.NewCounter("counter", nil)
	timer := god.NewTimer("timer")
	god.TrackInt("int", func() int64 { return 2 })
	god.TrackFloat("float", func() float64 { return 4.2 })
	go func() {
		for {
			d := time.Duration(rand.Intn(500)) * time.Millisecond
			time.Sleep(d)
			counter.Add(1)
			timer.Add(d)
		}
	}()

	select {}
}
