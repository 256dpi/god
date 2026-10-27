package god

import (
	"bytes"
	"net/http"
	"runtime/trace"
)

func flightHandler(cfg trace.FlightRecorderConfig) http.HandlerFunc {
	// start recorder
	recorder := trace.NewFlightRecorder(cfg)
	err := recorder.Start()
	if err != nil {
		println(err.Error())
	}

	return func(w http.ResponseWriter, r *http.Request) {
		// snapshot window, fails if inactive or already snapshotting
		var buf bytes.Buffer
		_, err := recorder.WriteTo(&buf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		// write snapshot
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="flight.trace"`)
		_, _ = w.Write(buf.Bytes())
	}
}
