package god

import (
	"net"
	"net/http"
	"net/http/pprof"
	"runtime"
	"runtime/trace"
	"strconv"
	"strings"
	"time"

	"github.com/felixge/fgprof"
)

// Options define the available debugging options.
type Options struct {
	// The host to run the debug interface on.
	//
	// Default: "0.0.0.0".
	Host string

	// The port to run the debug interface on.
	//
	// Default: 6060.
	Port int

	// The memory profile rate, see runtime.MemProfileRate.
	//
	// Default: Go's own, one sample per 512 KiB allocated.
	MemoryProfileRate int

	// The mutex profile fraction.
	//
	// Default: 1.
	MutexProfileFraction int

	// The block profile rate.
	//
	// Default: 1.
	BlockProfileRate int

	// The flight recorder configuration. If set, an execution trace of the
	// recent past is kept in memory and served as a snapshot at
	// /debug/flight. Tracing continuously costs some CPU and the window's
	// memory.
	//
	// Default: nil (disabled).
	FlightRecorder *trace.FlightRecorderConfig

	// Custom handler for the status endpoint.
	//
	// Default: "OK" writer.
	StatusHandler http.HandlerFunc

	// Custom handler for the metrics endpoint.
	MetricsHandler http.HandlerFunc
}

// Init will run a god compatible debug endpoint.
func Init(opts Options) {
	// set defaults
	if opts.MutexProfileFraction == 0 {
		opts.MutexProfileFraction = 1
	}
	if opts.BlockProfileRate == 0 {
		opts.BlockProfileRate = 1
	}

	// set memory profile rate
	if opts.MemoryProfileRate > 0 {
		runtime.MemProfileRate = opts.MemoryProfileRate
	}

	// print metrics once registered
	enableMetrics()

	// start flight recorder
	var flight http.HandlerFunc
	if opts.FlightRecorder != nil {
		flight = flightHandler(*opts.FlightRecorder)
	}

	// get address
	host := "0.0.0.0"
	if opts.Host != "" {
		host = opts.Host
	}
	port := 6060
	if opts.Port > 0 {
		port = opts.Port
	}
	addr := net.JoinHostPort(host, strconv.Itoa(port))

	// enable debugging
	go func() {
		// create mux
		mux := http.NewServeMux()

		// add pprof endpoints
		mux.HandleFunc("/debug/pprof/", profile(opts))
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

		// add fgprof endpoint
		mux.Handle("/debug/fgprof", fgprof.Handler())

		// add flight recorder endpoint
		if flight != nil {
			mux.HandleFunc("/debug/flight", flight)
		}

		// add runtime metrics endpoint
		mux.HandleFunc("/debug/runtime", runtimeHandler)

		// add prometheus endpoint
		if opts.MetricsHandler != nil {
			mux.Handle("/metrics", opts.MetricsHandler)
		}

		// add status endpoint
		mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
			if opts.StatusHandler != nil {
				opts.StatusHandler(w, r)
			} else {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("OK"))
			}
		})

		// add trace endpoint
		mux.HandleFunc("/trace", traceHandler)

		// launch debug server
		err := http.ListenAndServe(addr, mux)
		if err != nil {
			println(err.Error())
		}
	}()
}

func profile(opts Options) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get profile
		seg := strings.Split(r.URL.Path, "/")
		name := seg[len(seg)-1]

		// get seconds
		sec, err := strconv.ParseInt(r.FormValue("seconds"), 10, 64)
		if sec <= 0 || err != nil {
			sec = 30
		}

		// build temporary mutex profile
		if name == "mutex" {
			runtime.SetMutexProfileFraction(opts.MutexProfileFraction)
			defer runtime.SetMutexProfileFraction(0)
			time.Sleep(time.Duration(sec) * time.Second)
		}

		// build temporary block profile
		if name == "block" {
			runtime.SetBlockProfileRate(opts.BlockProfileRate)
			defer runtime.SetBlockProfileRate(0)
			time.Sleep(time.Duration(sec) * time.Second)
		}

		// call index
		pprof.Index(w, r)
	}
}
