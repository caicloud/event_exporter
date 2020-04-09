package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/caicloud/event_exporter/pkg/version"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/pflag"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

var (
	showVersion = pflag.Bool(
		"version", false,
		"Print version information",
	)
	listenAddress = pflag.String(
		"web.listen-address", ":9102",
		"Address to listen on for web interface and telemetry.",
	)
	metricPath = pflag.String(
		"web.telemetry-path", "/metrics",
		"Path under which to expose metrics.",
	)
	maxPreserve = pflag.Int(
		"event.max-length", 300,
		"Upper bound duration(sec) of an event to preserve",
	)
	initPreserve = pflag.Int(
		"event.init-length", 20,
		"Lower bound duration(sec) of an event to preserve",
	)
	kubeNamespace = pflag.String(
		"namespace", core_v1.NamespaceAll,
		"Optional, the namespace to watch (default all)",
	)
	kubeConfig = pflag.String(
		"kubeconfig", "",
		"Absolute path to the kubeconfig file",
	)
	apiserver = pflag.String(
		"apiserver", "",
		"The URL of the apiserver to use as a master",
	)
)

var landingPage = []byte(`<html>
<head><title>Event exporter</title></head>
<body>
<h1>Event exporter</h1>
<p><a href='` + *metricPath + `'>Metrics</a></p>
</body>
</html>
`)

func init() {
	pflag.Parse()
	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Message())
		os.Exit(0)
	}
}

func main() {
	config, err := clientcmd.BuildConfigFromFlags(*apiserver, *kubeConfig)
	if err != nil {
		klog.Fatalf("build kubeconfig: %v", err)
	}
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalln("create client:", err)
	}
	store := NewEventStore(client,
		time.Duration(*initPreserve)*time.Second,
		time.Duration(*maxPreserve)*time.Second,
		*kubeNamespace)
	go store.Run()
	exporter := NewExporter(store)
	prometheus.MustRegister(exporter)

	http.Handle(*metricPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(landingPage)
	})

	klog.Infoln("Listening on", *listenAddress)
	klog.Fatal(http.ListenAndServe(*listenAddress, nil))
}
