package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/caicloud/nirvana/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	"github.com/spf13/pflag"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	flags       = pflag.NewFlagSet("", pflag.ExitOnError)
	showVersion = flags.Bool(
		"version", false,
		"Print version information",
	)
	listenAddress = flags.String(
		"web.listen-address", ":9102",
		"Address to listen on for web interface and telemetry.",
	)
	metricPath = flags.String(
		"web.telemetry-path", "/metrics",
		"Path under which to expose metrics.",
	)
	maxPreserve = flags.Int(
		"event.max-length", 300,
		"Upper bound duration(sec) of an event to preserve",
	)
	initPreserve = flags.Int(
		"event.init-length", 20,
		"Lower bound duration(sec) of an event to preserve",
	)
	inCluster = flags.Bool(
		"running-in-cluster", true,
		`Optional, if this controller is running in a kubernetes cluster, use the
		pod secrets for creating a Kubernetes client.`,
	)
	kubeconfig = flags.String("kubeconfig", "", "absolute path to the kubeconfig file")
)

var landingPage = []byte(`<html>
<head><title>Event exporter</title></head>
<body>
<h1>Event exporter</h1>
<p><a href='` + *metricPath + `'>Metrics</a></p>
</body>
</html>
`)

func serveHTTP() {
}

// func init() {
// 	prometheus.MustRegister(version.NewCollector("event_exporter"))
// }

func main() {
	flags.AddGoFlagSet(flag.CommandLine)
	flags.Parse(os.Args)

	// Workaround of noisy log, see https://github.com/kubernetes/kubernetes/issues/17162
	flag.CommandLine.Parse([]string{})

	if *showVersion {
		fmt.Fprintln(os.Stdout, version.Print("event_exporter"))
		os.Exit(0)
	}
	var client kubernetes.Interface
	var config *rest.Config
	var err error
	if *inCluster {
		config, err = rest.InClusterConfig()
		if err != nil {
			log.Fatalf("error get incluster config: %v", err)
		}
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
		if err != nil {
			log.Fatalf("error connecting to the client: %v", err)
		}
	}

	client, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln("failed to create client:", err)
	}
	store, err := NewEventStore(client,
		time.Duration(*initPreserve)*time.Second,
		time.Duration(*maxPreserve)*time.Second)
	if err != nil {
		log.Fatalln("error create event store:", err)
	}
	go store.Run()
	exporter := NewExporter(store)
	prometheus.MustRegister(exporter)
	log.Infoln("Starting event_exporter", version.Info())
	log.Infoln("Build context", version.BuildContext())

	http.Handle(*metricPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(landingPage)
	})

	log.Infoln("Listening on", *listenAddress)
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
