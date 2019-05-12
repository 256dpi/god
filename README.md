# god

**A small tool for simplifying debugging of go applications.**

Import the library in your app:

```go
import "github.com/256dpi/god"
```

Then add the following line as early as possible:

```go
god.Debug()
```

This will open a pprof and prometheus endpoint on port `6060`.

Use the `god` utility to interact with the endpoint:

```bash
god -duration 10
```

The tool will fetch the CPU profile by default but can also fetch other profiles:

```
Usage of god:
  -block
    	block profile
  -duration int
    	trace duration (default 5)
  -mem
    	memory profile
  -mutex
    	mutex profile
  -trace
    	trace profile
```
