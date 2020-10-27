# Kubernetes Event Exporter
[![Build Status](https://travis-ci.org/caicloud/event_exporter.svg?branch=master)](https://travis-ci.org/caicloud/event_exporter)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)

Kubernetes events to Prometheus bridge.

A Collector that can list and watch Kubernetes events, and according to events' occurrence, determine how long the event lasts. The information is then translated into metrics.

# Metrics Overview

1. `kube_event_count` Number of kubernetes event that happened in the past an hour.The metric value is the same as the count property of event in the cluster.
   ```
   kube_event_count{involved_object_kind="Deployment",involved_object_name="event-exporter",involved_object_namespace="default",name="event-exporter.1640452bd04fc7bf",namespace="default",reason="ScalingReplicaSet",source="/deployment-controller",type="Normal"} 1
   ```
2. `kube_event_uinque_events_total` Total number of kubernetes unique event that happened in the past an hour
   ```
   kube_event_unique_events_total{involved_object_kind="Deployment",involved_object_name="event-exporter",involved_object_namespace="default",name="event-exporter.1640452bd04fc7bf",namespace="default",reason="ScalingReplicaSet",source="/deployment-controller",type="Normal"} 1
   ```
3. `event_exporter_version`Information of the event exporter that was built
   ```
   event_exporter_build_info{branch="v1.0",build_date="2020-10-22T10:11:29Z",build_user="Caicloud Authors",go_version="go1.13.15",version="v1.0.0"} 1
   ```

# Getting Started

## Build

```shell
$ VERSION=v1.0.0 REGISTRY=docker.io make build
```
If you want to get more information about flag options,please refer to `Makefile` in our repository
## Run

running outside Kubernetes (Exporter will search for kubeconfig in ~/.kube)

```shell
$ ./event_exporter  --kubeConfigPath=$HOME/.kube/config
```

running inside Kubernetes (Exporter will use Kubernetes serviceaccount)

```shell
$ ./event_exporter
```

## General Flags

Name  | Example| Description
--- | --- | ---
kubeMasterURL|--kubeMasterURL=<APIServer-URL>|Optional. The URL of kubernetes apiserver to use as a master
kubeConfigPath| --kubeConfigPath=$HOME/.kube/config|Optional. The path of kubernetes configuration file 
eventType |--eventType=Warning --eventType=Normal |Optional.  List of allowed event types. The default value is `Warning` type
port| --port=9012|Optional. Port to expose event metrics on
version | --version| Print version information 

## Use Kubernetes

You can deploy this exporter by using the  image `caicloud/event-exporter:${VERSION}` in k8s cluster,
the available versions can be got from the [releases](https://github.com/caicloud/event_exporter/releases).

### Deploy

```shell
kubectl apply -f deploy.yml
```

Then check the pod status:

```shell
kubectl get pods | grep event
```

```
# HELP go_gc_duration_seconds A summary of the pause duration of garbage collection cycles.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 1.6811e-05
go_gc_duration_seconds{quantile="0.25"} 2.6e-05
go_gc_duration_seconds{quantile="0.5"} 3.0795e-05
go_gc_duration_seconds{quantile="0.75"} 8.0126e-05
go_gc_duration_seconds{quantile="1"} 0.000186691
go_gc_duration_seconds_sum 0.001432397
go_gc_duration_seconds_count 24
# HELP go_goroutines Number of goroutines that currently exist.
# TYPE go_goroutines gauge
go_goroutines 27
# HELP go_info Information about the Go environment.
# TYPE go_info gauge
go_info{version="go1.13.15"} 1
# HELP go_memstats_alloc_bytes Number of bytes allocated and still in use.
# TYPE go_memstats_alloc_bytes gauge
go_memstats_alloc_bytes 6.29132e+06
# HELP go_memstats_alloc_bytes_total Total number of bytes allocated, even if freed.
# TYPE go_memstats_alloc_bytes_total counter
go_memstats_alloc_bytes_total 5.6787848e+07
# HELP go_memstats_buck_hash_sys_bytes Number of bytes used by the profiling bucket hash table.
# TYPE go_memstats_buck_hash_sys_bytes gauge
go_memstats_buck_hash_sys_bytes 1.452877e+06
# HELP go_memstats_frees_total Total number of frees.
# TYPE go_memstats_frees_total counter
go_memstats_frees_total 236938
# HELP go_memstats_gc_cpu_fraction The fraction of this program's available CPU time used by the GC since the program started.
# TYPE go_memstats_gc_cpu_fraction gauge
go_memstats_gc_cpu_fraction 2.924731798616038e-06
# HELP go_memstats_gc_sys_bytes Number of bytes used for garbage collection system metadata.
# TYPE go_memstats_gc_sys_bytes gauge
go_memstats_gc_sys_bytes 2.377728e+06
# HELP go_memstats_heap_alloc_bytes Number of heap bytes allocated and still in use.
# TYPE go_memstats_heap_alloc_bytes gauge
go_memstats_heap_alloc_bytes 6.29132e+06
# HELP go_memstats_heap_idle_bytes Number of heap bytes waiting to be used.
# TYPE go_memstats_heap_idle_bytes gauge
go_memstats_heap_idle_bytes 5.8359808e+07
# HELP go_memstats_heap_inuse_bytes Number of heap bytes that are in use.
# TYPE go_memstats_heap_inuse_bytes gauge
go_memstats_heap_inuse_bytes 7.766016e+06
# HELP go_memstats_heap_objects Number of allocated objects.
# TYPE go_memstats_heap_objects gauge
go_memstats_heap_objects 21220
# HELP go_memstats_heap_released_bytes Number of heap bytes released to OS.
# TYPE go_memstats_heap_released_bytes gauge
go_memstats_heap_released_bytes 5.7688064e+07
# HELP go_memstats_heap_sys_bytes Number of heap bytes obtained from system.
# TYPE go_memstats_heap_sys_bytes gauge
go_memstats_heap_sys_bytes 6.6125824e+07
# HELP go_memstats_last_gc_time_seconds Number of seconds since 1970 of last garbage collection.
# TYPE go_memstats_last_gc_time_seconds gauge
go_memstats_last_gc_time_seconds 1.6033609023805106e+09
# HELP go_memstats_lookups_total Total number of pointer lookups.
# TYPE go_memstats_lookups_total counter
go_memstats_lookups_total 0
# HELP go_memstats_mallocs_total Total number of mallocs.
# TYPE go_memstats_mallocs_total counter
go_memstats_mallocs_total 258158
# HELP go_memstats_mcache_inuse_bytes Number of bytes in use by mcache structures.
# TYPE go_memstats_mcache_inuse_bytes gauge
go_memstats_mcache_inuse_bytes 13888
# HELP go_memstats_mcache_sys_bytes Number of bytes used for mcache structures obtained from system.
# TYPE go_memstats_mcache_sys_bytes gauge
go_memstats_mcache_sys_bytes 16384
# HELP go_memstats_mspan_inuse_bytes Number of bytes in use by mspan structures.
# TYPE go_memstats_mspan_inuse_bytes gauge
go_memstats_mspan_inuse_bytes 66096
# HELP go_memstats_mspan_sys_bytes Number of bytes used for mspan structures obtained from system.
# TYPE go_memstats_mspan_sys_bytes gauge
go_memstats_mspan_sys_bytes 81920
# HELP go_memstats_next_gc_bytes Number of heap bytes when next garbage collection will take place.
# TYPE go_memstats_next_gc_bytes gauge
go_memstats_next_gc_bytes 1.07428e+07
# HELP go_memstats_other_sys_bytes Number of bytes used for other system allocations.
# TYPE go_memstats_other_sys_bytes gauge
go_memstats_other_sys_bytes 1.772971e+06
# HELP go_memstats_stack_inuse_bytes Number of bytes in use by the stack allocator.
# TYPE go_memstats_stack_inuse_bytes gauge
go_memstats_stack_inuse_bytes 983040
# HELP go_memstats_stack_sys_bytes Number of bytes obtained from system for stack allocator.
# TYPE go_memstats_stack_sys_bytes gauge
go_memstats_stack_sys_bytes 983040
# HELP go_memstats_sys_bytes Number of bytes obtained from system.
# TYPE go_memstats_sys_bytes gauge
go_memstats_sys_bytes 7.2810744e+07
# HELP go_threads Number of OS threads created.
# TYPE go_threads gauge
go_threads 18
# HELP event_exporter_build_info A metric with a constant '1' value labeled by version, branch,build_user,build_date and go_version from which event_exporter was built
# TYPE event_exporter_build_info gauge
event_exporter_info{branch="v1.0",build_date="2020-10-22T10:11:29Z",build_user="Caicloud Authors",go_version="go1.13.15",version="v1.0.0"} 1
# HELP kube_event_count Number of kubernetes event happened
# TYPE kube_event_count gauge
kube_event_count{involved_object_kind="Deployment",involved_object_name="event-exporter",involved_object_namespace="default",name="event-exporter.1640452bd04fc7bf",namespace="default",reason="ScalingReplicaSet",source="/deployment-controller",type="Normal"} 1
kube_event_count{involved_object_kind="Deployment",involved_object_name="event-exporter",involved_object_namespace="default",name="event-exporter.164045435014f51c",namespace="default",reason="ScalingReplicaSet",source="/deployment-controller",type="Normal"} 1
kube_event_count{involved_object_kind="Deployment",involved_object_name="event-exporter",involved_object_namespace="default",name="event-exporter.164045638ee80ccb",namespace="default",reason="ScalingReplicaSet",source="/deployment-controller",type="Normal"} 1
kube_event_count{involved_object_kind="Deployment",involved_object_name="event-exporter",involved_object_namespace="default",name="event-exporter.164045efda48031f",namespace="default",reason="ScalingReplicaSet",source="/deployment-controller",type="Normal"} 1
kube_event_count{involved_object_kind="Deployment",involved_object_name="my-nginx",involved_object_namespace="default",name="my-nginx.1640456cf4c9fbad",namespace="default",reason="ScalingReplicaSet",source="/deployment-controller",type="Normal"} 1
kube_event_count{involved_object_kind="PersistentVolumeClaim",involved_object_name="prometheus-data-prometheus-0",involved_object_namespace="kube-system",name="prometheus-data-prometheus-0.163ff24070ae83e5",namespace="kube-system",reason="ProvisioningFailed",source="/persistentvolume-controller",type="Warning"} 6303
# HELP kube_event_unique_events_total Total number of kubernetes unique event happened
# TYPE kube_event_unique_events_total counter
kube_event_unique_events_total{involved_object_kind="Deployment",involved_object_name="event-exporter",involved_object_namespace="default",name="event-exporter.1640452bd04fc7bf",namespace="default",reason="ScalingReplicaSet",source="/deployment-controller",type="Normal"} 1
kube_event_unique_events_total{involved_object_kind="Deployment",involved_object_name="event-exporter",involved_object_namespace="default",name="event-exporter.164045435014f51c",namespace="default",reason="ScalingReplicaSet",source="/deployment-controller",type="Normal"} 1
kube_event_unique_events_total{involved_object_kind="Deployment",involved_object_name="event-exporter",involved_object_namespace="default",name="event-exporter.164045638ee80ccb",namespace="default",reason="ScalingReplicaSet",source="/deployment-controller",type="Normal"} 1
kube_event_unique_events_total{involved_object_kind="Deployment",involved_object_name="event-exporter",involved_object_namespace="default",name="event-exporter.164045efda48031f",namespace="default",reason="ScalingReplicaSet",source="/deployment-controller",type="Normal"} 1
kube_event_unique_events_total{involved_object_kind="Deployment",involved_object_name="my-nginx",involved_object_namespace="default",name="my-nginx.1640456cf4c9fbad",namespace="default",reason="ScalingReplicaSet",source="/deployment-controller",type="Normal"} 1
kube_event_unique_events_total{involved_object_kind="PersistentVolumeClaim",involved_object_name="prometheus-data-prometheus-0",involved_object_namespace="kube-system",name="prometheus-data-prometheus-0.163ff24070ae83e5",namespace="kube-system",reason="ProvisioningFailed",source="/persistentvolume-controller",type="Warning"} 10
# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 2.69
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1.048576e+06
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 10
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 3.4009088e+07
# HELP process_start_time_seconds Start time of the process since unix epoch in seconds.
# TYPE process_start_time_seconds gauge
process_start_time_seconds 1.60335836753e+09
# HELP process_virtual_memory_bytes Virtual memory size in bytes.
# TYPE process_virtual_memory_bytes gauge
process_virtual_memory_bytes 6.90274304e+08
# HELP process_virtual_memory_max_bytes Maximum amount of virtual memory available in bytes.
# TYPE process_virtual_memory_max_bytes gauge
process_virtual_memory_max_bytes -1
# HELP promhttp_metric_handler_requests_in_flight Current number of scrapes being served.
# TYPE promhttp_metric_handler_requests_in_flight gauge
promhttp_metric_handler_requests_in_flight 1
# HELP promhttp_metric_handler_requests_total Total number of scrapes by HTTP status code.
# TYPE promhttp_metric_handler_requests_total counter
promhttp_metric_handler_requests_total{code="200"} 174
promhttp_metric_handler_requests_total{code="500"} 0
promhttp_metric_handler_requests_total{code="503"} 0
```
# License
event_exporter is licensed under the Apache License, Version 2.0. See LICENSE for the full license text.
