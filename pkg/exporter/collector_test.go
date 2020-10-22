package exporter

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/caicloud/event_exporter/pkg/utils"
)

func TestVersion(t *testing.T) {
	var testcase utils.MetricsTestCases = map[string]utils.MetricsTestCase{
		"Version": {
			Target: exporterVersion,
			Want: fmt.Sprintf(`
        # HELP kube_event_exporter_version Version of the exporter
        # TYPE kube_event_exporter_version gauge
        kube_event_exporter_version{branch="UNKNOWN",build_date="UNKNOWN",build_user="Caicloud Authors",go_version="%s",version="1.0.0"} 1
`, runtime.Version()),
		},
	}
	testcase.Test(t)
}
