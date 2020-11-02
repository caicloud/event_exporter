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

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
)

var (
	eventCount = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "kube_event",
		Subsystem: "",
		Name:      "count",
		Help:      "Number of kubernetes event happened",
	}, []string{"name", "involved_object_namespace", "namespace", "involved_object_name", "involved_object_kind", "reason", "type", "source"})
	eventTotal = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: "kube_event",
		Subsystem: "",
		Name:      "unique_events_total",
		Help:      "Total number of kubernetes unique event happened",
	}, []string{"name", "involved_object_namespace", "namespace", "involved_object_name", "involved_object_kind", "reason", "type", "source"})
)

func increaseUniqueEventTotal(event *v1.Event) {
	eventTotal.With(prometheus.Labels{
		"name":                      event.ObjectMeta.Name,
		"namespace":                 event.Namespace,
		"involved_object_namespace": event.InvolvedObject.Namespace,
		"involved_object_name":      event.InvolvedObject.Name,
		"involved_object_kind":      event.InvolvedObject.Kind,
		"reason":                    event.Reason,
		"type":                      event.Type,
		"source":                    fmt.Sprintf("%s/%s", event.Source.Host, event.Source.Component),
	}).Inc()
}

func updateEventCount(event *v1.Event) {
	eventCount.With(prometheus.Labels{
		"name":                      event.ObjectMeta.Name,
		"namespace":                 event.Namespace,
		"involved_object_namespace": event.InvolvedObject.Namespace,
		"involved_object_name":      event.InvolvedObject.Name,
		"involved_object_kind":      event.InvolvedObject.Kind,
		"reason":                    event.Reason,
		"type":                      event.Type,
		"source":                    fmt.Sprintf("%s/%s", event.Source.Host, event.Source.Component),
	}).Set(float64(event.Count))
}

func delEventCountMetric(event *v1.Event) {
	ret := eventCount.Delete(prometheus.Labels{
		"name":                      event.ObjectMeta.Name,
		"namespace":                 event.Namespace,
		"involved_object_namespace": event.InvolvedObject.Namespace,
		"involved_object_name":      event.InvolvedObject.Name,
		"involved_object_kind":      event.InvolvedObject.Kind,
		"reason":                    event.Reason,
		"type":                      event.Type,
		"source":                    fmt.Sprintf("%s/%s", event.Source.Host, event.Source.Component),
	})
	if ret {
		klog.Infof("event %s has been removed from Prometheus", event.ObjectMeta.Name)
	}
}

func delEventTotalMetric(event *v1.Event) {
	ret := eventTotal.Delete(prometheus.Labels{
		"name":                      event.ObjectMeta.Name,
		"namespace":                 event.Namespace,
		"involved_object_namespace": event.InvolvedObject.Namespace,
		"involved_object_name":      event.InvolvedObject.Name,
		"involved_object_kind":      event.InvolvedObject.Kind,
		"reason":                    event.Reason,
		"type":                      event.Type,
		"source":                    fmt.Sprintf("%s/%s", event.Source.Host, event.Source.Component),
	})
	if ret {
		klog.Infof("event %s has been removed from Prometheus", event.ObjectMeta.Name)
	}
}

func EventHandler(event *v1.Event) {
	increaseUniqueEventTotal(event)
	updateEventCount(event)
}

func DeleteMetric(event *v1.Event) {
	delEventCountMetric(event)
	delEventTotalMetric(event)
}
