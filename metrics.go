package god

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"
)

type metric struct {
	name   string
	metric Metric
}

var metrics []metric
var metricsEnabled bool
var metricsStarted bool
var metricsMutex sync.Mutex

// Metric represents a periodically collected metric.
type Metric interface {
	// Collect is called every second to retrieve the current value.
	Collect() string
}

// Register will register the provided metric.
func Register(name string, m Metric) {
	// acquire mutex
	metricsMutex.Lock()
	defer metricsMutex.Unlock()

	// add metric to a copy, as the printer may be reading the current list
	list := append(metrics[:len(metrics):len(metrics)], metric{
		name:   name,
		metric: m,
	})

	// sort by name
	sort.Slice(list, func(i, j int) bool {
		return list[i].name < list[j].name
	})

	// set list
	metrics = list

	// start printing
	startMetrics()
}

func enableMetrics() {
	// acquire mutex
	metricsMutex.Lock()
	defer metricsMutex.Unlock()

	// set flag
	metricsEnabled = true

	// start printing
	startMetrics()
}

func startMetrics() {
	// print once enabled and a metric is registered
	if metricsEnabled && !metricsStarted && len(metrics) > 0 {
		metricsStarted = true
		go printMetrics()
	}
}

func printMetrics() {
	// create ticker
	ticker := time.Tick(time.Second)

	// print metrics
	for {
		// await tick
		<-ticker

		// get metrics
		metricsMutex.Lock()
		list := metrics
		metricsMutex.Unlock()

		// collect strings
		s := make([]string, 0, len(list))
		for _, m := range list {
			s = append(s, fmt.Sprintf("%s: %s", m.name, m.metric.Collect()))
		}

		// print
		fmt.Println(strings.Join(s, " - "))
	}
}
