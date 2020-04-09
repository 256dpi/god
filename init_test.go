package god

import (
	"fmt"
	"sync"
	"time"
)

func ExampleInit() {
	// init
	Init(Options{})

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

	fmt.Println("Hello")

	select {}

	// Output:
	// Hello
}
