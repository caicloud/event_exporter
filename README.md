# Kubernetes Event Exporter

Kuberentes events to Prometheus bridge.

A Collector that can list and watch Kubernetes events, and according to events' occurrence, determine how long the event lasts. The information is then translated into metrics.

# Build and Run

## Build

```shell
$ make
```

## Run

running outside Kuberentes (Exporter will search for kubeconfig in ~/.kube)

```shell
$ ./event_exporter --running-in-cluster=false --kubeconfig=$HOME/.kube/config
```

running inside Kubernetes (Exporter will use Kubernetes serviceaccount)

```shell
$ ./event_exporter
```

## Collector Flags

Name | Type | Description
---| --- | ---
event.init-length | int | Lower bound duration(sec) of an event to preserve (default 20) |
event.max-length | int | Upper bound duration(sec) of an event to preserve (default 300)

## General Flags

Name | Description
--- | ---
running-in-cluster | Optional. If this controller is running in a kubernetes cluster, use the pod secrets for creating a Kubernetes client. (default true)
log.level | Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]. (default info)
version | Print version information

## Use Docker

You can deploy this exporter using the Docker image `caicloud/event-exporter:${VERSION}`,
the available versions can be got from the [releases](https://github.com/caicloud/event_exporter/releases).

For example:

```shell
$ docker run -d -p 9102:9102 -v ~/.kube/config:/root/.kube/config caicloud/event-exporter:${VERSION} --running-in-cluster=false
```

then make requests:

```shell
$ curl localhost:9102/metrics
```

example response:

```
# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0.000140144
go_gc_duration_seconds{quantile="0.25"} 0.000140144
go_gc_duration_seconds{quantile="0.5"} 0.00040924
go_gc_duration_seconds{quantile="0.75"} 0.00040924
go_gc_duration_seconds{quantile="1"} 0.00040924
go_gc_duration_seconds_sum 0.000549384
go_gc_duration_seconds_count 2
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 18
# HELP http_request_duration_microseconds The HTTP request latencies in microseconds.
# TYPE http_request_duration_microseconds summary
http_request_duration_microseconds{handler="prometheus",quantile="0.5"} NaN
http_request_duration_microseconds{handler="prometheus",quantile="0.9"} NaN
http_request_duration_microseconds{handler="prometheus",quantile="0.99"} NaN
http_request_duration_microseconds_sum{handler="prometheus"} 0
http_request_duration_microseconds_count{handler="prometheus"} 0
# HELP http_request_size_bytes The HTTP request sizes in bytes.
# TYPE http_request_size_bytes summary
http_request_size_bytes{handler="prometheus",quantile="0.5"} NaN
http_request_size_bytes{handler="prometheus",quantile="0.9"} NaN
http_request_size_bytes{handler="prometheus",quantile="0.99"} NaN
http_request_size_bytes_sum{handler="prometheus"} 0
http_request_size_bytes_count{handler="prometheus"} 0
# HELP http_response_size_bytes The HTTP response sizes in bytes.
# TYPE http_response_size_bytes summary
http_response_size_bytes{handler="prometheus",quantile="0.5"} NaN
http_response_size_bytes{handler="prometheus",quantile="0.9"} NaN
http_response_size_bytes{handler="prometheus",quantile="0.99"} NaN
http_response_size_bytes_sum{handler="prometheus"} 0
http_response_size_bytes_count{handler="prometheus"} 0
# HELP kubernetes_build_info A metric with a constant '1' value labeled by major, minor, git version, git commit, git tree state, build date, Go version, and compiler from which Kubernetes was built, and platform on which it is running.
# TYPE kubernetes_build_info gauge
kubernetes_build_info{buildDate="1970-01-01T00:00:00Z",compiler="gc",gitCommit="$Format:%H$",gitTreeState="not a git tree",gitVersion="v1.3.2+$Format:%h$",goVersion="go1.6.2",major="1",minor="3"<Plug>PeepOpenlatform="darwin/amd64"} 1
# HELP kubernetes_events State of kubernetes events
# TYPE kubernetes_events gauge
kubernetes_events{event_kind="Pod",event_name="nginx-pc-534913751-2yzev",event_namespace="allen",event_reason="BackOff",event_source="kube-node-3/kubelet",event_subobject="spec.containers{nginx}",event_type="Normal"} 1
kubernetes_events{event_kind="Pod",event_name="nginx-pc-534913751-2yzev",event_namespace="allen",event_reason="Failed",event_source="kube-node-3/kubelet",event_subobject="spec.containers{nginx}",event_type="Warning"} 0
kubernetes_events{event_kind="Pod",event_name="nginx-pc-534913751-2yzev",event_namespace="allen",event_reason="FailedSync_ErrImagePull",event_source="kube-node-3/kubelet",event_subobject="",event_type="Warning"} 0
kubernetes_events{event_kind="Pod",event_name="nginx-pc-534913751-2yzev",event_namespace="allen",event_reason="FailedSync_ImagePullBackOff",event_source="kube-node-3/kubelet",event_subobject="",event_type="Warning"} 1
kubernetes_events{event_kind="Pod",event_name="nginx-pc-534913751-2yzev",event_namespace="allen",event_reason="Pulling",event_source="kube-node-3/kubelet",event_subobject="spec.containers{nginx}",event_type="Normal"} 0
# HELP last_scrape_duration_seconds Duration of the last scrape of metrics from event store.
# TYPE last_scrape_duration_seconds gauge
last_scrape_duration_seconds 5.1443e-05
# HELP last_scrape_error Whether the last scrape of metrics from event store resulted in an error (1 for error, 0 for success).
# TYPE last_scrape_error gauge
last_scrape_error 0
# HELP rest_client_request_latency_microseconds Request latency in microseconds. Broken down by verb and URL
# TYPE rest_client_request_latency_microseconds summary
rest_client_request_latency_microseconds{url="https://sysinfra.caicloudprivatetest.com/api/v1/events?resourceVersion=%7Bvalue%7D",verb="GET",quantile="0.5"} 155984
rest_client_request_latency_microseconds{url="https://sysinfra.caicloudprivatetest.com/api/v1/events?resourceVersion=%7Bvalue%7D",verb="GET",quantile="0.9"} 155984
rest_client_request_latency_microseconds{url="https://sysinfra.caicloudprivatetest.com/api/v1/events?resourceVersion=%7Bvalue%7D",verb="GET",quantile="0.99"} 155984
rest_client_request_latency_microseconds_sum{url="https://sysinfra.caicloudprivatetest.com/api/v1/events?resourceVersion=%7Bvalue%7D",verb="GET"} 155984
rest_client_request_latency_microseconds_count{url="https://sysinfra.caicloudprivatetest.com/api/v1/events?resourceVersion=%7Bvalue%7D",verb="GET"} 1
# HELP rest_client_request_status_codes Number of http requests, partitioned by metadata
# TYPE rest_client_request_status_codes counter
rest_client_request_status_codes{code="200",host="sysinfra.caicloudprivatetest.com",method="GET"} 2
# HELP scrapes_total Total number of times event store was scraped for metrics.
# TYPE scrapes_total counter
scrapes_total 2
```
