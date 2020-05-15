package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	core_v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog"
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
	client          kubernetes.Interface
	stopCh          chan struct{}
	stopLock        sync.Mutex
	shutdown        bool
	eventController cache.Controller
	eventStore      cache.Store
	backoff         *Backoff
	events          *prometheus.Desc
}

// NewEventStore returns EventStore or error
func NewEventStore(client kubernetes.Interface, init, max time.Duration, namespace string) *EventStore {
	es := &EventStore{
		client:  client,
		stopCh:  make(chan struct{}),
		backoff: NewBackoff(init, max),
		events: prometheus.NewDesc(
			"kubernetes_events",
			"State of kubernetes events",
			[]string{"event_metaname", "event_namespace", "event_name", "event_kind", "event_reason", "event_type", "event_subobject", "event_message", "event_source"},
			nil,
		),
	}

	eventHandler := cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			es.sync(obj)
		},
		UpdateFunc: func(old, cur interface{}) {
			es.sync(cur)
		},
	}

	es.eventStore, es.eventController = cache.NewInformer(
		&cache.ListWatch{
			ListFunc:  eventListFunc(es.client, namespace),
			WatchFunc: eventWatchFunc(es.client, namespace),
		},
		&core_v1.Event{}, resyncPeriod, eventHandler)

	return es
}

func (es *EventStore) sync(obj interface{}) {
	event := obj.(*core_v1.Event)
	eventCleaning(event)
	key, err := keyFunc(event)
	if err != nil {
		klog.Errorf("cannot generate key of event: %v", event)
		return
	}
	es.backoff.Next(key, event.LastTimestamp.Time)
}

// This function used to handle the not well formated k8s events type
func eventCleaning(e *core_v1.Event) {
	if e.Reason == "FailedSync" {
		if strings.Contains(e.Message, "ErrImagePull") {
			e.Reason = fmt.Sprintf("%s_%s", e.Reason, "ErrImagePull")
		} else if strings.Contains(e.Message, "ImagePullBackOff") {
			e.Reason = fmt.Sprintf("%s_%s", e.Reason, "ImagePullBackOff")
		}
	}
}

func eventListFunc(c kubernetes.Interface, ns string) func(meta_v1.ListOptions) (runtime.Object, error) {
	return func(options meta_v1.ListOptions) (runtime.Object, error) {
		return c.CoreV1().Events(ns).List(options)
	}
}

func eventWatchFunc(c kubernetes.Interface, ns string) func(meta_v1.ListOptions) (watch.Interface, error) {
	return func(options meta_v1.ListOptions) (watch.Interface, error) {
		return c.CoreV1().Events(ns).Watch(options)
	}
}

// Run event store
func (es *EventStore) Run() {
	klog.Infoln("start event store...")
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

func eval(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

// Scrap events and generate metrics
func (es *EventStore) Scrap(ch chan<- prometheus.Metric) {
	for key, isHappening := range es.backoff.AllKeysStateSinceUpdate(time.Now()) {
		obj, exists, err := es.eventStore.GetByKey(key)
		if err != nil {
			continue
		} else if !exists {
			klog.Errorf("event not found: %s", key)
			continue
		}
		event := obj.(*core_v1.Event)
		ch <- prometheus.MustNewConstMetric(
			es.events, prometheus.GaugeValue,
			eval(isHappening),
			event.ObjectMeta.Name,
			event.InvolvedObject.Namespace,
			event.InvolvedObject.Name,
			event.InvolvedObject.Kind,
			event.Reason,
			event.Type,
			event.InvolvedObject.FieldPath,
			event.Message,
			fmt.Sprintf("%s/%s", event.Source.Host, event.Source.Component),
		)
	}
	es.backoff.GC()
}
