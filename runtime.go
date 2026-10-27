package god

import (
	"fmt"
	"math"
	"net/http"
	rtmetrics "runtime/metrics"
)

func runtimeHandler(w http.ResponseWriter, r *http.Request) {
	// read all metrics
	descs := rtmetrics.All()
	samples := make([]rtmetrics.Sample, len(descs))
	for i := range samples {
		samples[i].Name = descs[i].Name
	}
	rtmetrics.Read(samples)

	// write metrics
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	for _, sample := range samples {
		switch sample.Value.Kind() {
		case rtmetrics.KindUint64:
			_, _ = fmt.Fprintf(w, "%s %d\n", sample.Name, sample.Value.Uint64())
		case rtmetrics.KindFloat64:
			_, _ = fmt.Fprintf(w, "%s %g\n", sample.Name, sample.Value.Float64())
		case rtmetrics.KindFloat64Histogram:
			hist := sample.Value.Float64Histogram()
			_, _ = fmt.Fprintf(w, "%s count=%d p50=%g p90=%g p99=%g max=%g\n", sample.Name, count(hist),
				quantile(hist, 0.5), quantile(hist, 0.9), quantile(hist, 0.99), quantile(hist, 1))
		}
	}
}

func count(hist *rtmetrics.Float64Histogram) uint64 {
	var total uint64
	for _, n := range hist.Counts {
		total += n
	}
	return total
}

// quantile returns the upper bound of the bucket holding the quantile, or its
// lower bound if the bucket is unbounded above.
func quantile(hist *rtmetrics.Float64Histogram, q float64) float64 {
	// get target
	total := count(hist)
	if total == 0 {
		return 0
	}
	target := uint64(math.Ceil(q * float64(total)))

	// find bucket
	var seen uint64
	for i, n := range hist.Counts {
		seen += n
		if n > 0 && seen >= target {
			if math.IsInf(hist.Buckets[i+1], 1) {
				return hist.Buckets[i]
			}
			return hist.Buckets[i+1]
		}
	}

	return 0
}
