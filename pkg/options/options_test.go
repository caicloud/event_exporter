package options

import (
	"os"
	"reflect"
	"testing"
)

func TestOptionsParse(t *testing.T) {
	defaultEventTypes := []string{"Warning"}
	defaultPort := 9102
	defaultVersion := false
	defaultKubeMasterURL := ""
	defaultKubeConfigPath := ""

	tests := []struct {
		Name     string
		Args     []string
		Expected *Options
	}{
		{
			Name: "version command line argument",
			Args: []string{"./event_exporter", "--version"},
			Expected: &Options{
				EventType:      defaultEventTypes,
				Port:           defaultPort,
				KubeConfigPath: defaultKubeConfigPath,
				KubeMasterURL:  defaultKubeMasterURL,
				Version:        true,
			},
		},
		{
			Name: "exporter kubernetes config path and port",
			Args: []string{"./event_exporter",
				"--port=8090",
				"--kubeConfigPath=/Users/admin/.kube/config",
			},
			Expected: &Options{
				KubeConfigPath: "/Users/admin/.kube/config",
				KubeMasterURL:  defaultKubeMasterURL,
				EventType:      defaultEventTypes,
				Port:           8090,
				Version:        defaultVersion,
			},
		},
		{
			Name: "exporter event types",
			Args: []string{"./event_exporter",
				"--eventType=Normal",
				"--eventType=Warning",
			},
			Expected: &Options{
				KubeMasterURL:  defaultKubeMasterURL,
				KubeConfigPath: defaultKubeConfigPath,
				EventType:      []string{"Normal", "Warning"},
				Port:           defaultPort,
				Version:        defaultVersion,
			},
		},
		{
			Name: "default config",
			Args: []string{"./event_exporter"},
			Expected: &Options{
				KubeMasterURL:  "",
				KubeConfigPath: "",
				EventType:      defaultEventTypes,
				Port:           defaultPort,
				Version:        defaultVersion,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			opts := NewOptions()
			opts.AddFlags()
			os.Args = test.Args
			opts.Parse()
			opts.flag = nil
			if !reflect.DeepEqual(opts, test.Expected) {
				t.Errorf("test error for case:%s", test.Name)
			}
		})
	}
}
