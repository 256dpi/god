package god

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

var traceReceivers = map[chan Event]struct{}{}
var traceMutex sync.Mutex

// Event is a trace event.
type Event struct {
	name  string
	task  string
	start time.Time
	stop  time.Time
}

// End will complete a trace event.
func (e Event) End() {
	// set stop
	e.stop = time.Now()

	// distribute event
	traceMutex.Lock()
	for rec := range traceReceivers {
		select {
		case rec <- e:
		default:
		}
	}
	traceMutex.Unlock()
}

// Trace will start a new trace event.
func Trace(name, task string) Event {
	return Event{
		name:  name,
		task:  task,
		start: time.Now(),
	}
}

func traceHandler(w http.ResponseWriter, r *http.Request) {
	// create channel
	ch := make(chan Event, 32)

	// add receiver
	traceMutex.Lock()
	traceReceivers[ch] = struct{}{}
	traceMutex.Unlock()

	// remove receiver
	defer func() {
		traceMutex.Lock()
		delete(traceReceivers, ch)
		traceMutex.Unlock()
	}()

	// write header
	w.Header().Set("Content-Type", "text/event-stream")
	w.WriteHeader(200)
	w.(http.Flusher).Flush()

	// write events
	for event := range ch {
		_, err := fmt.Fprintf(w, "%s;%s;%s;%s\n", event.name, event.task, event.start.Format(time.RFC3339Nano), event.stop.Format(time.RFC3339Nano))
		if err != nil {
			return
		}
		w.(http.Flusher).Flush()
	}
}
