# Introduction
This will help you to deploy event-exporter to a kubernetes cluster and gather metrics on kubernetes events.

# How to deploy to a kubernetes cluster
## Assumptions
1. A kubernetes cluster is up and running
2. You have privilage to deploy to that cluster

Run following command to deploy event-exporter to a cluster

`curl https://github.com/caicloud/event_exporter/tree/master/examples/deploy.yaml | kubectl --context {add your cluster context} -n {add namespace here} apply -f -`

# How to see event metrics in Prometheus

## Assumptions
Prometheus is currently running and scraping pods in kubernetes
[Kube-state-metrics](https://github.com/kubernetes/kube-state-metrics) is deployed to your cluster

1. Open your prometheus instance
2. Search for metrics 'kubernetes_events'

# Setup alerts in alert manager

Please see references
1. https://prometheus.io/docs/alerting/overview/
2. https://itnext.io/prometheus-with-alertmanager-f2a1f7efabd6

