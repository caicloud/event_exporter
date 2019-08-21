package main

import (
	"time"

	"github.com/caicloud/nirvana/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "event"
	exporter  = "exporter"
)

// Exporter collects events and convert to metrics. It implements prometheus.Collector.
type Exporter struct {
	store           *EventStore
	duration, error prometheus.Gauge
	totalScrapes    prometheus.Counter
	up              prometheus.Gauge
}

// NewExporter return a new event exporter
func NewExporter(store *EventStore) *Exporter {
	return &Exporter{
		store: store,
		duration: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "last_scrape_duration_seconds",
			Help: "Duration of the last scrape of metrics from event store.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "scrapes_total",
			Help: "Total number of times event store was scraped for metrics.",
		}),
		error: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "last_scrape_error",
			Help: "Whether the last scrape of metrics from event store resulted in an error (1 for error, 0 for success).",
		}),
	}
}

// Collect implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	ch <- e.duration
	ch <- e.totalScrapes
	ch <- e.error
	e.scrape(ch)
	return
}

// Describe implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	metricCh := make(chan prometheus.Metric)
	doneCh := make(chan struct{})
	go func() {
		for m := range metricCh {
			ch <- m.Desc()
		}
		close(doneCh)
	}()
	e.Collect(metricCh)
	close(metricCh)
	<-doneCh
	return
}

func (e *Exporter) scrape(ch chan<- prometheus.Metric) {
	e.totalScrapes.Inc()
	var err error
	defer func(begun time.Time) {
		e.duration.Set(time.Since(begun).Seconds())
		if err == nil {
			e.error.Set(0)
		} else {
			e.error.Set(1)
		}
	}(time.Now())
	if err = e.store.Scrap(ch); err != nil {
		log.Errorln("Error scraping for events:", err)
	}
}
