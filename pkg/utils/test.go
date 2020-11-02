package utils

import (
	"bytes"
	"log"
	"testing"

	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/prometheus/client_golang/prometheus/testutil/promlint"
)

// MetricsTestCase can be used to unit test a Prometheus Collector implementation. It takes
// a initialized Prometheus Collector and check if the metric it defines is standard and
// produces the expected output.
type MetricsTestCase struct {
	// Target is the Prometheus Collector implementation to test
	Target prometheus.Collector
	// Want is the expected output in the metric API response as the result of Target Collector
	Want string
	// Metrics is the names of the metrics to test; leave empty to test all metrics.
	Metrics []string
}

// Test runs the tests. It returns an error if the Collector does not work properly. It returns
// a list of link errors and no errors if the Collector works but has non-standard definition.
func (mtc MetricsTestCase) Test() ([]promlint.Problem, error) {
	want := bytes.NewBufferString(mtc.Want)
	if err := testutil.CollectAndCompare(mtc.Target, want, mtc.Metrics...); err != nil {
		return nil, errors.Wrap(err, "output verification failed")
	}
	return nil, nil
}

// MetricsTestCases is a alias for a set of MetricsTestCase. It provide a convenient and standard
// way to unit test multiple MetricsTestCase. The keys of the map are the names of the test cases.
type MetricsTestCases map[string]MetricsTestCase

// Test runs all given test cases under the given testing.T object.
func (mtc MetricsTestCases) Test(t *testing.T) {
	for name, tc := range mtc {
		t.Run(name, func(tt *testing.T) {
			problems, err := tc.Test()
			if err != nil {
				tt.Error(err)
			}
			for _, problem := range problems {
				log.Printf("non-standard metric '%s': %s", problem.Metric, problem.Text)
			}
		})
	}
}
