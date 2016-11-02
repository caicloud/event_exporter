package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	k8s_client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/runtime"
	"k8s.io/kubernetes/pkg/watch"
)

const (
	// Resync period for the kube controller loop.
	resyncPeriod = 5 * time.Minute
)

var (
	keyFunc = cache.MetaNamespaceKeyFunc
)

// EventStore stores events, handle event lasts time
type EventStore struct {
	client          *k8s_client.Client
	stopCh          chan struct{}
	stopLock        sync.Mutex
	shutdown        bool
	eventController *framework.Controller
	eventStore      cache.Store
	backoff         *Backoff
	events          *prometheus.Desc
}

// NewEventStore returns EventStore or error
func NewEventStore(client *k8s_client.Client, init, max time.Duration) (*EventStore, error) {
	es := &EventStore{
		client:  client,
		stopCh:  make(chan struct{}),
		backoff: NewBackoff(init, max),
		events: prometheus.NewDesc(
			"kubernetes_events",
			"State of kubernetes events",
			[]string{"event_namespace", "event_name", "event_kind", "event_reason", "event_type", "event_subobject", "event_source"},
			nil,
		),
	}

	eventHandler := framework.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			es.sync(obj)
		},
		UpdateFunc: func(old, cur interface{}) {
			es.sync(cur)
		},
	}

	es.eventStore, es.eventController = framework.NewInformer(
		&cache.ListWatch{
			ListFunc:  eventListFunc(es.client, api.NamespaceAll),
			WatchFunc: eventWatchFunc(es.client, api.NamespaceAll),
		},
		&api.Event{}, resyncPeriod, eventHandler)

	return es, nil
}

func (es *EventStore) sync(obj interface{}) {
	event := obj.(*api.Event)
	eventCleaning(event)
	key, err := keyFunc(event)
	if err != nil {
		log.Errorf("cannot generate key of event: %v", event)
		return
	}
	es.backoff.Next(key, int(event.Count), event.LastTimestamp.Time)
}

// This function used to handle the not well formated k8s events type
func eventCleaning(e *api.Event) {
	if e.Reason == container.FailedSync {
		if strings.Contains(e.Message, "ErrImagePull") {
			e.Reason = fmt.Sprintf("%s_%s", e.Reason, "ErrImagePull")
		} else if strings.Contains(e.Message, "ImagePullBackOff") {
			e.Reason = fmt.Sprintf("%s_%s", e.Reason, "ImagePullBackOff")
		}
	}
}

func eventListFunc(c *k8s_client.Client, ns string) func(api.ListOptions) (runtime.Object, error) {
	return func(options api.ListOptions) (runtime.Object, error) {
		return c.Events(ns).List(options)
	}
}

func eventWatchFunc(c *k8s_client.Client, ns string) func(api.ListOptions) (watch.Interface, error) {
	return func(options api.ListOptions) (watch.Interface, error) {
		return c.Events(ns).Watch(options)
	}
}

// Run event store
func (es *EventStore) Run() {
	log.Infoln("start event store...")
	go es.eventController.Run(es.stopCh)
	<-es.stopCh
}

// Stop event store
func (es *EventStore) Stop() error {
	es.stopLock.Lock()
	defer es.stopLock.Unlock()

	if !es.shutdown {
		es.shutdown = true
		close(es.stopCh)

		return nil
	}

	return fmt.Errorf("shutdown already in progress")
}

func (es *EventStore) controllersInSync() bool {
	return es.eventController.HasSynced()
}

func eval(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

// Scrap events and generate metrics
func (es *EventStore) Scrap(ch chan<- prometheus.Metric) error {
	var err error
	err = nil
	for key, isHappening := range es.backoff.AllKeysStateSinceUpdate(time.Now()) {
		obj, exists, err := es.eventStore.GetByKey(key)
		if err != nil {
			continue
		} else if !exists {
			err = fmt.Errorf("event not found: %s", key)
			continue
		}
		event := obj.(*api.Event)
		ch <- prometheus.MustNewConstMetric(
			es.events, prometheus.GaugeValue,
			eval(isHappening),
			event.InvolvedObject.Namespace,
			event.InvolvedObject.Name,
			event.InvolvedObject.Kind,
			event.Reason,
			event.Type,
			event.InvolvedObject.FieldPath,
			fmt.Sprintf("%s/%s", event.Source.Host, event.Source.Component),
		)
	}
	es.backoff.GC()
	return err
}
