package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/apis/extensions"
	"k8s.io/kubernetes/pkg/client/cache"
	k8s_clientset "k8s.io/kubernetes/pkg/client/clientset_generated/internalclientset"
	k8s_client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/kubelet/container"
	"k8s.io/kubernetes/pkg/runtime"
	deploymentutil "k8s.io/kubernetes/pkg/util/deployment"
	"k8s.io/kubernetes/pkg/watch"
)

const (
	// Resync period for the kube controller loop.
	resyncPeriod = 5 * time.Minute
)

var (
	keyFunc = cache.MetaNamespaceKeyFunc
)

// Controller listwatch events, pod, deployments and  handle them
type Controller struct {
	client          *k8s_client.Client
	clientset       *k8s_clientset.Clientset
	stopCh          chan struct{}
	stopLock        sync.Mutex
	shutdown        bool
	eventController *framework.Controller
	eventStore      cache.Store
	podController   *framework.Controller
	podStore        cache.Store
	dpController    *framework.Controller
	dpStore         cache.Store
	backoff         *Backoff
	eventsDesc      *prometheus.Desc
	rcMapperDesc    *prometheus.Desc
	dpMapperDesc    *prometheus.Desc
	mapper          map[string]*api.PodList
}

// NewController returns Controller instance or an error
func NewController(client *k8s_client.Client, clientset *k8s_clientset.Clientset, init, max time.Duration) (*Controller, error) {
	es := &Controller{
		client:    client,
		clientset: clientset,
		stopCh:    make(chan struct{}),
		backoff:   NewBackoff(init, max),
		mapper:    map[string]*api.PodList{},
		eventsDesc: prometheus.NewDesc(
			"kubernetes_events",
			"State of kubernetes events",
			[]string{"event_namespace", "event_name", "event_kind", "event_reason", "event_type", "event_subobject", "event_source"},
			nil,
		),
		rcMapperDesc: prometheus.NewDesc(
			"kubernetes_resource_mapper",
			"Resource mapper of kubernetes",
			[]string{"io_kubernetes_pod_uid", "kubernetes_pod_name", "kubernetes_namespace", "kubernetes_rc_name"},
			nil,
		),
		dpMapperDesc: prometheus.NewDesc(
			"kubernetes_resource_mapper",
			"Resource mapper of kubernetes",
			[]string{"io_kubernetes_pod_uid", "kubernetes_pod_name", "kubernetes_namespace", "kubernetes_rs_name", "kubernetes_dp_name"},
			nil,
		),
	}

	eventHandler := framework.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			es.syncEvent(obj)
		},
		UpdateFunc: func(old, cur interface{}) {
			es.syncEvent(cur)
		},
	}

	es.eventStore, es.eventController = framework.NewInformer(
		&cache.ListWatch{
			ListFunc:  eventListFunc(es.client, api.NamespaceAll),
			WatchFunc: eventWatchFunc(es.client, api.NamespaceAll),
		},
		&api.Event{}, resyncPeriod, eventHandler)

	es.podStore, es.podController = framework.NewInformer(
		&cache.ListWatch{
			ListFunc:  podListFunc(es.client, api.NamespaceAll),
			WatchFunc: podWatchFunc(es.client, api.NamespaceAll),
		},
		&api.Pod{}, resyncPeriod, framework.ResourceEventHandlerFuncs{})

	es.dpStore, es.dpController = framework.NewInformer(
		&cache.ListWatch{
			ListFunc:  dpListFunc(es.client, api.NamespaceAll),
			WatchFunc: dpWatchFunc(es.client, api.NamespaceAll),
		},
		&extensions.Deployment{}, resyncPeriod, framework.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				dp := obj.(*extensions.Deployment)
				key := fmt.Sprintf("%s/%s", dp.Namespace, dp.Name)
				pods, err := deploymentutil.ListPods(dp, func(ns string, options api.ListOptions) (*api.PodList, error) {
					return es.client.Pods(ns).List(options)
				})
				if err != nil {
					log.Infof("%v", err)
				}
				es.mapper[key] = pods
			},
			UpdateFunc: func(old, cur interface{}) {
				dp := cur.(*extensions.Deployment)
				key := fmt.Sprintf("%s/%s", dp.Namespace, dp.Name)
				pods, err := deploymentutil.ListPods(dp, func(ns string, options api.ListOptions) (*api.PodList, error) {
					return es.client.Pods(ns).List(options)
				})
				if err != nil {
					log.Infof("%v", err)
				}
				es.mapper[key] = pods
			},
			DeleteFunc: func(obj interface{}) {
				dp := obj.(*extensions.Deployment)
				key := fmt.Sprintf("%s/%s", dp.Namespace, dp.Name)
				delete(es.mapper, key)
			},
		})

	return es, nil
}

func (es *Controller) syncEvent(obj interface{}) {
	event := obj.(*api.Event)
	eventCleaning(event)
	key, err := keyFunc(event)
	if err != nil {
		log.Errorf("cannot generate key of event: %v", event)
		return
	}
	es.backoff.Next(key, int(event.Count), event.LastTimestamp.Time)
}

func (es *Controller) syncPod(obj interface{}) {
	log.Infof("sync pod")
}

func (es *Controller) syncDp(obj interface{}) {
	log.Infof("sync dp")
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

func podListFunc(c *k8s_client.Client, ns string) func(api.ListOptions) (runtime.Object, error) {
	return func(options api.ListOptions) (runtime.Object, error) {
		return c.Pods(ns).List(options)
	}
}

func podWatchFunc(c *k8s_client.Client, ns string) func(api.ListOptions) (watch.Interface, error) {
	return func(options api.ListOptions) (watch.Interface, error) {
		return c.Pods(ns).Watch(options)
	}
}

func dpListFunc(c *k8s_client.Client, ns string) func(api.ListOptions) (runtime.Object, error) {
	return func(options api.ListOptions) (runtime.Object, error) {
		return c.Deployments(ns).List(options)
	}
}

func dpWatchFunc(c *k8s_client.Client, ns string) func(api.ListOptions) (watch.Interface, error) {
	return func(options api.ListOptions) (watch.Interface, error) {
		return c.Deployments(ns).Watch(options)
	}
}

// Run event store
func (es *Controller) Run() {
	log.Infoln("start event store...")
	go es.eventController.Run(es.stopCh)
	go es.podController.Run(es.stopCh)
	go es.dpController.Run(es.stopCh)
	<-es.stopCh
}

// Stop event store
func (es *Controller) Stop() error {
	es.stopLock.Lock()
	defer es.stopLock.Unlock()

	if !es.shutdown {
		es.shutdown = true
		close(es.stopCh)

		return nil
	}

	return fmt.Errorf("shutdown already in progress")
}

func (es *Controller) controllersInSync() bool {
	return es.eventController.HasSynced()
}

func eval(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

// GetCreatedBy return the SerializedReference in created-by annotation
func GetCreatedBy(pod *api.Pod) (*api.SerializedReference, error) {
	raw, ok := pod.Annotations["kubernetes.io/created-by"]
	if !ok {
		return nil, fmt.Errorf("no created-by annotation")
	}
	obj, err := runtime.Decode(api.Codecs.UniversalDecoder(), []byte(raw))
	if err != nil {
		return nil, err
	}
	return obj.(*api.SerializedReference), nil
}

// ScrapMapper scrap pod to rc/rs/dp map
func (es *Controller) ScrapMapper(ch chan<- prometheus.Metric) error {
	var err error
	err = nil
	for _, obj := range es.podStore.List() {
		pod := obj.(*api.Pod)
		createdBy, _ := GetCreatedBy(pod)
		if createdBy != nil && createdBy.Reference.Kind == "ReplicationController" {
			rc := createdBy.Reference.Name
			ch <- prometheus.MustNewConstMetric(
				es.rcMapperDesc, prometheus.GaugeValue, 1,
				string(pod.GetUID()), pod.Name, pod.Namespace, rc,
			)
		}
	}
	for key, podList := range es.mapper {
		name := strings.Split(key, "/")[1]
		for _, pod := range podList.Items {
			createdBy, _ := GetCreatedBy(&pod)
			if createdBy != nil && createdBy.Reference.Kind == "ReplicaSet" {
				ch <- prometheus.MustNewConstMetric(
					es.dpMapperDesc, prometheus.GaugeValue, 1,
					string(pod.GetUID()), pod.Name, pod.Namespace, createdBy.Reference.Name, name,
				)
			}
		}
	}
	return err
}

// ScrapEvents scrap events and generate metrics
func (es *Controller) ScrapEvents(ch chan<- prometheus.Metric) error {
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
			es.eventsDesc, prometheus.GaugeValue,
			eval(isHappening), event.InvolvedObject.Namespace,
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

// Scrap all
func (es *Controller) Scrap(ch chan<- prometheus.Metric) error {
	var err error
	err = es.ScrapEvents(ch)
	if err != nil {
		return err
	}
	err = es.ScrapMapper(ch)
	return err
}
