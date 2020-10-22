package filters

import (
	"testing"

	v1 "k8s.io/api/core/v1"
)

func TestEventTypeFilter_Filter(t *testing.T) {
	warningType := "Waring"
	normalType := "Normal"
	type fields struct {
		AllowedTypes []string
	}
	type args struct {
		event *v1.Event
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "warning event,warning rule",
			fields: fields{
				AllowedTypes: []string{warningType},
			},
			args: args{&v1.Event{
				Type: warningType,
			}},
			want: true,
		},
		{
			name:   "normal event,normal rule",
			fields: fields{AllowedTypes: []string{normalType}},
			args: args{&v1.Event{
				Type: normalType,
			}},
			want: true,
		},
		{
			name:   "normal event,warning rule",
			fields: fields{AllowedTypes: []string{warningType}},
			args: args{&v1.Event{
				Type: normalType,
			}},
			want: false,
		},
		{
			name:   "warning event,normal rule",
			fields: fields{AllowedTypes: []string{normalType}},
			args: args{&v1.Event{
				Type: warningType,
			}},
			want: false,
		},
		{
			name:   "warning event,all types rule",
			fields: fields{AllowedTypes: []string{normalType, warningType}},
			args: args{&v1.Event{
				Type: warningType,
			}},
			want: true,
		},
		{
			name:   "normal event,all types rule",
			fields: fields{AllowedTypes: []string{normalType, warningType}},
			args: args{&v1.Event{
				Type: normalType,
			}},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EventTypeFilter{
				AllowedTypes: tt.fields.AllowedTypes,
			}
			if got := e.Filter(tt.args.event); got != tt.want {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}
