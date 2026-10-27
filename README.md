# god

**A small tool for simplifying debugging of go applications.**

Import the library in your app:

```go
import "github.com/256dpi/god"
```

Then add the following line as early as possible:

```go
god.Init(god.Options{})
```

This will open a debug endpoint on `0.0.0.0:6060` that serves:

- `/debug/pprof/` the pprof profiles
- `/debug/fgprof` a full (on and off CPU) profile, using fgprof
- `/debug/runtime` the runtime metrics as text, histograms as quantiles
- `/debug/flight` a snapshot of the flight recorder, if enabled
- `/metrics` a prometheus endpoint, if a `MetricsHandler` is given
- `/status` an `OK` or what a `StatusHandler` writes
- `/trace` a stream of the events recorded with `god.Trace`

The options set where it listens and what the runtime records:

```go
god.Init(god.Options{
	Host: "127.0.0.1", // listen on loopback only, e.g. inside a shared pod
	Port: 6064,
	// sample one allocation per 1 KiB instead of Go's 512 KiB
	MemoryProfileRate: 1024,
	// keep at least the last 10 seconds of execution trace
	FlightRecorder: &trace.FlightRecorderConfig{
		MinAge: 10 * time.Second,
	},
})
```

`Host` defaults to `0.0.0.0` and `Port` to `6060`. The memory profile rate
stays Go's own unless set. The flight recorder traces continuously, which costs
some CPU and the window's memory, so it is off unless configured. Fetch a
snapshot right after something went wrong and open it with `go tool trace`:

```bash
curl -o flight.trace http://localhost:6060/debug/flight
go tool trace flight.trace
```

Use the `god` utility to interact with the endpoint:

```bash
god -duration 10
```

The tool will fetch all available profiles, run the pprof servers and open a
custom frontend from which all profiles can be accessed.

It can also collect a CPU profile for profile-guided optimization:

```bash
god -pgo ./cmd/app/default.pgo -duration 30
```

The profile is merged with the one already in the file, so running it a few
times under different loads builds a representative profile. Go uses a
`default.pgo` in the main package's directory when building.
