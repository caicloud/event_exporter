package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Exporter collects events and convert to metrics. It implements prometheus.Collector.
type Exporter struct {
	store *EventStore
}

// NewExporter return a new event exporter
func NewExporter(store *EventStore) *Exporter {
	return &Exporter{
		store: store,
	}
}

// Collect implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.scrape(ch)
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
}

func (e *Exporter) scrape(ch chan<- prometheus.Metric) {
	e.store.Scrap(ch)
}
