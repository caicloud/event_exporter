package version

import (
	"fmt"
	"runtime"
	"testing"
)

func TestMessage(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Version",
			want: fmt.Sprintf(`event_exporter:1.0.0 (Branch: UNKNOWN, Revision: UNKNOWN)
build user: Caicloud Authors
build date: UNKNOWN
go version: %s
version   : 1.0.0
`, runtime.Version()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Message(); got != tt.want {
				t.Errorf("Message() = %v, want %v", got, tt.want)
			}
		})
	}
}
