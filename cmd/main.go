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

package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"

	"github.com/caicloud/event_exporter/pkg/collector"
	_ "github.com/caicloud/event_exporter/pkg/exporter"
	"github.com/caicloud/event_exporter/pkg/options"
	"github.com/caicloud/event_exporter/pkg/signal"
	"github.com/caicloud/event_exporter/pkg/version"
)

const (
	resync = time.Minute * 5
)

func main() {
	opts := options.NewOptions()
	opts.AddFlags()
	if err := opts.Parse(); err != nil {
		klog.Fatalf("failed to parse commandline args,err:%s", err.Error())
	}
	if opts.Version {
		fmt.Fprintln(os.Stdout, version.Message())
		os.Exit(0)
	}

	group, stopChan := signal.SetupStopSignalContext()

	kubeConfig, err := clientcmd.BuildConfigFromFlags(opts.KubeMasterURL, opts.KubeConfigPath)
	if err != nil {
		klog.Fatalf("failed to build kubernetes cluster configuration,err:%s", err.Error())
	}
	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if err != nil {
		klog.Fatalf("failed to build kubernetes client,err:%s", err.Error())
	}

	factory := informers.NewSharedInformerFactory(kubeClient, resync)
	eventCollector := collector.NewEventCollector(kubeClient, factory, opts)
	factory.Start(stopChan)

	group.Go(func() error {
		if err := eventCollector.Run(stopChan); err != nil {
			return fmt.Errorf("eventCollector run err:%s", err.Error())
		}
		return nil
	})

	klog.Infof("starting prometheus metrics server on http://localhost:%d", opts.Port)
	http.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.Handle("/metrics", promhttp.Handler())
	if err := http.ListenAndServe(fmt.Sprintf(":%d", opts.Port), nil); err != nil {
		klog.Fatal(err)
	}

	if err := group.Wait(); err != nil {
		klog.Fatal(err)
	}
}
