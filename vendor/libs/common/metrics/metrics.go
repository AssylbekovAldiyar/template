package metrics

import (
	"github.com/go-kit/kit/metrics"
	kitprom "github.com/go-kit/kit/metrics/prometheus"
	stdprom "github.com/prometheus/client_golang/prometheus"
)

func NewRequestCounter(namespace, name string) metrics.Counter {
	return kitprom.NewCounterFrom(stdprom.CounterOpts{
		Namespace: namespace,
		Subsystem: name,
		Name:      "request_count",
		Help:      "Number of requests received.",
	}, []string{"method"})
}

func NewErrorCounter(namespace, name string) metrics.Counter {
	return kitprom.NewCounterFrom(stdprom.CounterOpts{
		Namespace: namespace,
		Subsystem: name,
		Name:      "error_count",
		Help:      "Number of error requests received.",
	}, []string{"method"})
}

func NewRequestLatency(namespace, name string) metrics.Histogram {
	return kitprom.NewSummaryFrom(stdprom.SummaryOpts{
		Namespace: namespace,
		Subsystem: name,
		Name:      "request_latency_microseconds",
		Help:      "Total duration of requests in microseconds.",
	}, []string{"method"})
}
