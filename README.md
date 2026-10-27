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

This will open a pprof endpoint (and a prometheus endpoint if a
`MetricsHandler` is given) on `0.0.0.0:6060`. The options set where it listens
and how much the runtime profiles:

```go
god.Init(god.Options{
	Host: "127.0.0.1", // listen on loopback only, e.g. inside a shared pod
	Port: 6064,
	// sample one allocation per 512 KiB (Go's default) instead of god's 1 KiB
	MemoryProfileRate: 512 * 1024,
})
```

`Host` defaults to `0.0.0.0`, `Port` to `6060` and `MemoryProfileRate` to
`1024`, which records allocations in detail at a cost to processes that
allocate a lot.

Use the `god` utility to interact with the endpoint:

```bash
god -duration 10
```

The tool will fetch all available profiles, run the pprof servers and open a
custom frontend from which all profiles can be accessed.
