# Introduction
This will help you to deploy event-exporter to a kubernetes cluster and gather metrics on kubernetes events.

# How to deploy to a kubernetes cluster
## Assumptions
1. A kubernetes cluster is up and running
2. You have privilage to deploy to that cluster

Run following commands to deploy event-exporter to a cluster

1. Download deploy.yaml file locally
2. Update namespace in ClusterRoleBinding object (line 18)
3. `kubectl --context {add your cluster context} -n {add namespace here} apply -f deploy.yaml`


# How to see event metrics in Prometheus

## Assumptions
Prometheus is currently running and scraping pods in kubernetes
[Kube-state-metrics](https://github.com/kubernetes/kube-state-metrics) is deployed to your cluster

1. Configure your scrape_config section in `prometheus.yml` file 
    ```yaml
    scrape_configs:
      # The job name is added as a label `job=<job_name>` to any timeseries scraped from this config.
      - job_name: 'prometheus'
    
        # metrics_path defaults to '/metrics'
        # scheme defaults to 'http'.
    
        static_configs:
        - targets: ['localhost:9090']
      - job_name: 'event_exporter'
    
        # metrics_path defaults to '/metrics'
        # scheme defaults to 'http'.
    
        static_configs:
        - targets: ['event-exporter:9102']
    ```
2. Starting your prometheus instance
3. Search for metrics `kube_event_count` 、 `kube_event_uinque_events_total`、`event_exporter_build_info`

# Setup alerts in alert manager

Please see references
1. https://prometheus.io/docs/alerting/overview/
2. https://itnext.io/prometheus-with-alertmanager-f2a1f7efabd6

