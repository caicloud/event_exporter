package collector

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/caicloud/event_exporter/pkg/utils"
)

func TestEventHandler(t *testing.T) {
	event := &v1.Event{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "prometheus-data-prometheus-0.163ff24070ae83e5",
			Namespace: "kube-system",
		},
		InvolvedObject: v1.ObjectReference{
			Kind:      "PersistentVolumeClaim",
			Namespace: "kube-system",
			Name:      "prometheus-data-prometheus-0",
		},
		Reason: "ProvisioningFailed",
		Source: v1.EventSource{
			Host:      "",
			Component: "persistentvolume-controller",
		},
		Count: 120,
		Type:  "Warning",
	}
	EventHandler(event)
	var testcases utils.MetricsTestCases = map[string]utils.MetricsTestCase{
		"EventCount": {
			Target: eventCount,
			Want: `
# HELP kube_event_count Number of kubernetes event happened
# TYPE kube_event_count gauge
kube_event_count{involved_object_kind="PersistentVolumeClaim",involved_object_name="prometheus-data-prometheus-0",involved_object_namespace="kube-system",name="prometheus-data-prometheus-0.163ff24070ae83e5",namespace="kube-system",reason="ProvisioningFailed",source="/persistentvolume-controller",type="Warning"} 120
`,
		},
		"EventTotal": {
			Target: eventTotal,
			Want: `
# HELP kube_event_unique_events_total Total number of kubernetes unique event happened
# TYPE kube_event_unique_events_total counter
kube_event_unique_events_total{involved_object_kind="PersistentVolumeClaim",involved_object_name="prometheus-data-prometheus-0",involved_object_namespace="kube-system",name="prometheus-data-prometheus-0.163ff24070ae83e5",namespace="kube-system",reason="ProvisioningFailed",source="/persistentvolume-controller",type="Warning"} 1
`,
		},
	}
	testcases.Test(t)
}

func TestDeleteMetric(t *testing.T) {
	event := &v1.Event{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "prometheus-data-prometheus-0.163ff24070ae83e5",
			Namespace: "kube-system",
		},
		InvolvedObject: v1.ObjectReference{
			Kind:      "PersistentVolumeClaim",
			Namespace: "kube-system",
			Name:      "prometheus-data-prometheus-0",
		},
		Reason: "ProvisioningFailed",
		Source: v1.EventSource{
			Host:      "",
			Component: "persistentvolume-controller",
		},
		Count: 120,
		Type:  "Warning",
	}
	EventHandler(event)
	DeleteMetric(event)
	var testcases utils.MetricsTestCases = map[string]utils.MetricsTestCase{
		"EventCount": {
			Target: eventCount,
			Want:   "",
		},
		"EventTotal": {
			Target: eventTotal,
			Want:   "",
		},
	}
	testcases.Test(t)
}

func TestMultiMetric(t *testing.T) {
	eventBefore := &v1.Event{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "prometheus-data-prometheus-0.163ff24070ae83e5",
			Namespace: "kube-system",
		},
		InvolvedObject: v1.ObjectReference{
			Kind:      "PersistentVolumeClaim",
			Namespace: "kube-system",
			Name:      "prometheus-data-prometheus-0",
		},
		Reason: "ProvisioningFailed",
		Source: v1.EventSource{
			Host:      "",
			Component: "persistentvolume-controller",
		},
		Count: 120,
		Type:  "Warning",
	}
	eventAfter := &v1.Event{
		ObjectMeta: meta_v1.ObjectMeta{
			Name:      "prometheus-data-prometheus-0.163ff24070ae83e5",
			Namespace: "kube-system",
		},
		InvolvedObject: v1.ObjectReference{
			Kind:      "PersistentVolumeClaim",
			Namespace: "kube-system",
			Name:      "prometheus-data-prometheus-0",
		},
		Reason: "ProvisioningFailed",
		Source: v1.EventSource{
			Host:      "",
			Component: "persistentvolume-controller",
		},
		Count: 140,
		Type:  "Warning",
	}
	EventHandler(eventBefore)
	DeleteMetric(eventBefore)
	EventHandler(eventAfter)
	var testcases utils.MetricsTestCases = map[string]utils.MetricsTestCase{
		"EventCount": {
			Target: eventCount,
			Want: `
		# HELP kube_event_count Number of kubernetes event happened
		# TYPE kube_event_count gauge
		kube_event_count{involved_object_kind="PersistentVolumeClaim",involved_object_name="prometheus-data-prometheus-0",involved_object_namespace="kube-system",name="prometheus-data-prometheus-0.163ff24070ae83e5",namespace="kube-system",reason="ProvisioningFailed",source="/persistentvolume-controller",type="Warning"} 140
`,
		},
		"EventTotal": {
			Target: eventTotal,
			Want: `
        # HELP kube_event_unique_events_total Total number of kubernetes unique event happened
        # TYPE kube_event_unique_events_total counter
        kube_event_unique_events_total{involved_object_kind="PersistentVolumeClaim",involved_object_name="prometheus-data-prometheus-0",involved_object_namespace="kube-system",name="prometheus-data-prometheus-0.163ff24070ae83e5",namespace="kube-system",reason="ProvisioningFailed",source="/persistentvolume-controller",type="Warning"} 1
`,
		},
	}
	testcases.Test(t)
}
