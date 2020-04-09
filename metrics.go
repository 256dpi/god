package god

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// Metric represents a periodically collected metric.
type Metric interface {
	// Collect is called every second to retrieve the current value.
	Collect() string
}

type metric struct {
	name   string
	metric Metric
}

var metrics []metric

// Metrics will enable the collection and printing of debug metrics.
func Metrics() {
	// run printer
	go func() {
		// create ticker
		ticker := time.Tick(time.Second)

		// print metrics
		for {
			select {
			case <-ticker:
				// collect strings
				s := make([]string, 0, len(metrics))
				for _, m := range metrics {
					s = append(s, fmt.Sprintf("%s: %s", m.name, m.metric.Collect()))
				}

				// print
				fmt.Println(strings.Join(s, " ｜ "))
			}
		}
	}()
}

// Register will register the provided metric.
func Register(name string, m Metric) {
	// add metric
	metrics = append(metrics, metric{
		name:   name,
		metric: m,
	})

	// sort by name
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].name < metrics[j].name
	})
}
