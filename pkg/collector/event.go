/*
Copyright 2020 CaiCloud, Inc. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package collector

import (
	"fmt"
	"sync"
	"time"

	v1api "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	coreV1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"

	"github.com/caicloud/event_exporter/pkg/filters"
	"github.com/caicloud/event_exporter/pkg/options"
)

type EventCollector struct {
	kc                kubernetes.Interface
	factory           informers.SharedInformerFactory
	eventLister       coreV1.EventLister
	eventListerSynced cache.InformerSynced
	queue             workqueue.RateLimitingInterface
	filters           []filters.EventFilter
	cache             map[string]v1api.Event
	locker            sync.Mutex
}

func NewEventCollector(kc kubernetes.Interface, factory informers.SharedInformerFactory, o *options.Options) *EventCollector {
	event := factory.Core().V1().Events()
	eventCollector := &EventCollector{
		kc:                kc,
		factory:           factory,
		eventLister:       event.Lister(),
		eventListerSynced: event.Informer().HasSynced,
		queue:             workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter()),
		cache:             make(map[string]v1api.Event),
		locker:            sync.Mutex{},
	}
	eventCollector.addFilter(o)
	event.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			eventCollector.enqueueEvent(obj)
		},
		UpdateFunc: func(old, new interface{}) {
			newObj := new.(*v1api.Event)
			oldObj := old.(*v1api.Event)
			if newObj.ResourceVersion == oldObj.ResourceVersion {
				return
			}
			eventCollector.enqueueEvent(newObj)
		},
		DeleteFunc: func(obj interface{}) {
			eventCollector.enqueueEvent(obj)
		},
	})
	return eventCollector
}

func (ec *EventCollector) Run(stopCh <-chan struct{}) error {
	defer runtime.HandleCrash()
	defer ec.queue.ShutDown()

	klog.Info("starting eventCollector")
	if ok := cache.WaitForCacheSync(stopCh, ec.eventListerSynced); !ok {
		return fmt.Errorf("failed to wait for caches to sync")
	}
	klog.Info("started")
	go wait.Until(ec.runWorker, time.Second, stopCh)

	<-stopCh
	klog.Info("shutting down")
	return nil
}

func (ec *EventCollector) runWorker() {
	for ec.processNextItem() {
	}
}

func (ec *EventCollector) processNextItem() bool {
	key, quit := ec.queue.Get()
	if quit {
		return false
	}
	defer ec.queue.Done(key)

	err := ec.syncEvent(key.(string))
	if err == nil {
		ec.queue.Forget(key)
		return true
	}

	runtime.HandleError(fmt.Errorf("syncing %v failed with : %v", key, err))
	ec.queue.AddRateLimited(key)
	return true
}

func (ec *EventCollector) enqueueEvent(obj interface{}) {
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(fmt.Errorf("couldn't get key for object %+v:%v", obj, err))
		return
	}
	ec.queue.Add(key)
}

func (ec *EventCollector) syncEvent(key string) error {
	namespace, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		runtime.HandleError(fmt.Errorf("invalid resource key: %s", key))
		return nil
	}
	event, err := ec.eventLister.Events(namespace).Get(name)
	if err != nil {
		if errors.IsNotFound(err) {
			klog.Infof("event %s has been deleted", key)
			ec.locker.Lock()
			defer ec.locker.Unlock()
			if deletedEvent, ok := ec.cache[name]; ok {
				DeleteMetric(&deletedEvent)
				delete(ec.cache, name)
			}
			return nil
		}
		return err
	}
	if ec.eventFilter(event) {
		ec.locker.Lock()
		ec.cache[event.ObjectMeta.Name] = *event
		ec.locker.Unlock()
		klog.Infof(
			"event name: %s,count: %d,involvedObject_namespace: %s,involvedObject_kind: %s,involvedObject_name: %s,reason: %s,type: %s",
			event.Name,
			event.Count,
			event.InvolvedObject.Namespace,
			event.InvolvedObject.Kind,
			event.InvolvedObject.Name,
			event.Reason,
			event.Type,
		)
		EventHandler(event)
	}
	return nil
}

func (ec *EventCollector) addFilter(o *options.Options) {
	typeFilter := filters.NewEventTypeFilter(o.EventType)
	ec.filters = append(ec.filters, typeFilter)
}

func (ec *EventCollector) eventFilter(event *v1api.Event) bool {
	for _, filter := range ec.filters {
		if !filter.Filter(event) {
			return false
		}
	}
	return true
}
