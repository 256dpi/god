package god

import (
	"net/http"
	"net/http/pprof"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Options define the available debugging options.
type Options struct {
	// The port to run the debug interface on.
	//
	// Default: 6060.
	Port int

	// The mutex profile fraction.
	//
	// Default: 1.
	MutexProfileFraction int

	// The block profile rate.
	//
	// Default: 1.
	BlockProfileRate int

	// A custom handler for the status endpoint.
	//
	// Default: "OK" writer.
	StatusHandler http.HandlerFunc
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

	// get address
	addr := "0.0.0.0:6060"
	if opts.Port > 0 {
		addr = "0.0.0.0:" + strconv.Itoa(opts.Port)
	}

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

		// add prometheus endpoint
		mux.Handle("/metrics", promhttp.Handler())

		// add status endpoint
		mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
			if opts.StatusHandler != nil {
				opts.StatusHandler(w, r)
			} else {
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("OK"))
			}
		})

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
