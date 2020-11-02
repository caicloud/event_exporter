package exporter

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/caicloud/event_exporter/pkg/utils"
)

func TestVersion(t *testing.T) {
	var testcase utils.MetricsTestCases = map[string]utils.MetricsTestCase{
		"Build_info": {
			Target: exporterVersion,
			Want: fmt.Sprintf(`
        # HELP event_exporter_build_info A metric with a constant '1' value labeled by version, branch,build_user,build_date and go_version from which event_exporter was built
        # TYPE event_exporter_build_info gauge
        event_exporter_build_info{branch="UNKNOWN",build_date="UNKNOWN",build_user="Caicloud Authors",go_version="%s",version="1.0.0"} 1
`, runtime.Version()),
		},
	}
	testcase.Test(t)
}
