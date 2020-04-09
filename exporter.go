package main

import (
	"github.com/caicloud/nirvana/log"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "event"
	exporter  = "exporter"
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
	if err := e.store.Scrap(ch); err != nil {
		log.Errorln("Error scraping for events:", err)
	}
}
