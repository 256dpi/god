package god

import (
	"net/http"
	"net/http/pprof"
	"runtime"
	"strconv"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Debug will run a god compatible debug endpoint.
func Debug(port ...int) {
	// get address
	addr := "0.0.0.0:6060"
	if len(port) > 0 {
		addr = "0.0.0.0:" + strconv.Itoa(port[0])
	}

	// enable mutex profiling
	runtime.SetMutexProfileFraction(100)

	// enable block profiling
	runtime.SetBlockProfileRate(100)

	// enable debugging
	go func() {
		// create mux
		mux := http.NewServeMux()

		// add pprof endpoints
		mux.HandleFunc("/debug/pprof/", pprof.Index)
		mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
		mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
		mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
		mux.HandleFunc("/debug/pprof/trace", pprof.Trace)

		// add prometheus endpoint
		mux.Handle("/metrics", promhttp.Handler())

		// add status endpoint
		mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("OK"))
		})

		// launch debug server
		err := http.ListenAndServe(addr, mux)
		if err != nil {
			println(err.Error())
		}
	}()
}
