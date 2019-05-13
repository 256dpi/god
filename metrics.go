package god

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type collector interface {
	string() string
	reset()
}

type metric struct {
	name      string
	collector collector
}

var metrics []metric

// Metrics will enable the collection and printing of debug metrics.
func Metrics() {
	// create ticker
	ticker := time.Tick(time.Second)

	// run goroutine that calls print every second
	go func() {
		for {
			select {
			case <-ticker:
				printMetrics()
			}
		}
	}()
}

func add(name string, collector collector) {
	// add metric
	metrics = append(metrics, metric{
		name:      name,
		collector: collector,
	})

	// sort by name
	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].name < metrics[j].name
	})
}

func printMetrics() {
	// collect strings
	s := make([]string, 0, len(metrics))
	for _, m := range metrics {
		s = append(s, fmt.Sprintf("%s: %s", m.name, m.collector.string()))
		m.collector.reset()
	}

	// print
	fmt.Println(strings.Join(s, " ｜ "))
}
