# Kubernetes Event Exporter

Kuberentes events to Prometheus bridge.

A Collector that can list and watch Kubernetes events, according events' occurrence to determine how long the event lasts, then translate to metrics.

# Building and Running

## Build
```
make
```
## Running
running outside Kuberentes(It will search for kubeconfig in ~/.kube)

```
./event_exporter --running-in-cluster=false
```

running in Kubernetes(It will use Kubernetes serviceaccount)

```
./event_exporter
```

## Collector Flags

Name | Type | Description
---| --- | ---
event.init-length | int | Lower bound duration(sec) of an event to preserve (default 20) |
event.max-length | int | Upper bound duration(sec) of an event to preserve (default 300)

## General Flags

Name | Description
--- | ---
running-in-cluster | Optional, if this controller is running in a kubernetes cluster, use the pod secrets for creating a Kubernetes client. (default true)
log.level | Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]. (default info)
version | Print version information

#Using Docker

You can deploy this exporter using the cargo.caicloud.io/sysinfra/event-exporter Docker image.

For example:

```
docker pull prom/mysqld-exporter

docker run -d -p 9102:9102 -v ~/.kube/config:~/.kube/config cargo.caicloud.io/sysinfra/event-exporter --running-in-cluster=false
```
