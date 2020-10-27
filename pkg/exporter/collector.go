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

package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/caicloud/event_exporter/pkg/version"
)

var (
	exporterVersion = promauto.NewGauge(prometheus.GaugeOpts{
		Namespace: "event_exporter",
		Subsystem: "",
		Name:      "build_info",
		Help:      "A metric with a constant '1' value labeled by version, branch,build_user,build_date and go_version from which event_exporter was built",
		ConstLabels: prometheus.Labels{
			"version":    version.Version,
			"branch":     version.Branch,
			"build_user": version.BuildUser,
			"build_date": version.BuildDate,
			"go_version": version.GoVersion,
		},
	})
)

func init() {
	exporterVersion.Set(1)
}
