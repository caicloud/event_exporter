package main

import (
	"reflect"
	"testing"
	"time"
)

func TestBackoff_AllKeysStateSinceUpdate(t *testing.T) {
	type fields struct {
		baseDuration time.Duration
		maxDuration  time.Duration
		perItemEntry map[string]*backoffEntry
	}
	type args struct {
		eventTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]bool
	}{
		{
			name: "empty",
			fields: fields{
				baseDuration: 10 * time.Second,
				maxDuration:  300 * time.Second,
				perItemEntry: map[string]*backoffEntry{},
			},
			args: args{
				eventTime: time.Now(),
			},
			want: map[string]bool{},
		},
		{
			name: "contain",
			fields: fields{
				baseDuration: 10 * time.Second,
				maxDuration:  300 * time.Second,
				perItemEntry: map[string]*backoffEntry{"entry": {
					backoff:    10 * time.Second,
					lastUpdate: time.Now(),
				}},
			},
			args: args{
				eventTime: time.Now(),
			},
			want: map[string]bool{"entry": true},
		},
		{
			name: "expire",
			fields: fields{
				baseDuration: 10 * time.Second,
				maxDuration:  300 * time.Second,
				perItemEntry: map[string]*backoffEntry{"entry": {
					backoff:    10 * time.Second,
					lastUpdate: time.Now().Add(-20 * time.Second),
				}},
			},
			args: args{
				eventTime: time.Now(),
			},
			want: map[string]bool{"entry": false},
		},
		{
			name: "multiple",
			fields: fields{
				baseDuration: 10 * time.Second,
				maxDuration:  300 * time.Second,
				perItemEntry: map[string]*backoffEntry{
					"entry1": {
						backoff:    10 * time.Second,
						lastUpdate: time.Now().Add(-1 * time.Second),
					},
					"entry2": {
						backoff:    10 * time.Second,
						lastUpdate: time.Now().Add(-5 * time.Second),
					},
					"entry3": {
						backoff:    10 * time.Second,
						lastUpdate: time.Now().Add(-20 * time.Second),
					},
				},
			},
			args: args{
				eventTime: time.Now(),
			},
			want: map[string]bool{"entry1": true, "entry2": true, "entry3": false},
		},
	}
	for _, tt := range tests {
		p := &Backoff{
			baseDuration: tt.fields.baseDuration,
			maxDuration:  tt.fields.maxDuration,
			perItemEntry: tt.fields.perItemEntry,
		}
		if got := p.AllKeysStateSinceUpdate(tt.args.eventTime); !reflect.DeepEqual(got, tt.want) {
			t.Errorf("%q. Backoff.AllKeysStateSinceUpdate() = %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestBackoff_Next(t *testing.T) {
	type fields struct {
		baseDuration time.Duration
		maxDuration  time.Duration
		perItemEntry map[string]*backoffEntry
	}
	type args struct {
		id        string
		count     int
		eventTime time.Time
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   time.Duration
		want1  time.Time
	}{
		{
			name: "init",
			fields: fields{
				baseDuration: 10 * time.Second,
				maxDuration:  300 * time.Second,
				perItemEntry: map[string]*backoffEntry{},
			},
			args: args{
				id:        "entry",
				count:     1,
				eventTime: time.Date(2016, time.November, 11, 0, 0, 0, 0, time.Local),
			},
			want:  10 * time.Second,
			want1: time.Date(2016, time.November, 11, 0, 0, 0, 0, time.Local),
		},
		{
			name: "update",
			fields: fields{
				baseDuration: 10 * time.Second,
				maxDuration:  300 * time.Second,
				perItemEntry: map[string]*backoffEntry{"entry": {
					backoff:    10 * time.Second,
					lastUpdate: time.Now(),
				}},
			},
			args: args{
				id:        "entry",
				count:     1,
				eventTime: time.Date(2016, time.November, 11, 0, 0, 0, 0, time.Local),
			},
			want:  20 * time.Second,
			want1: time.Date(2016, time.November, 11, 0, 0, 0, 0, time.Local),
		},
		{
			name: "bound",
			fields: fields{
				baseDuration: 10 * time.Second,
				maxDuration:  300 * time.Second,
				perItemEntry: map[string]*backoffEntry{"entry": {
					backoff:    200 * time.Second,
					lastUpdate: time.Now(),
				}},
			},
			args: args{
				id:        "entry",
				count:     1,
				eventTime: time.Date(2016, time.November, 11, 0, 0, 0, 0, time.Local),
			},
			want:  300 * time.Second,
			want1: time.Date(2016, time.November, 11, 0, 0, 0, 0, time.Local),
		},
	}
	for _, tt := range tests {
		p := &Backoff{
			baseDuration: tt.fields.baseDuration,
			maxDuration:  tt.fields.maxDuration,
			perItemEntry: tt.fields.perItemEntry,
		}
		p.Next(tt.args.id, tt.args.eventTime)
		got, got1 := p.Get(tt.args.id)
		if got != tt.want {
			t.Errorf("%q. Backoff.Get() got = %v, want %v", tt.name, got, tt.want)
		}
		if !reflect.DeepEqual(got1, tt.want1) {
			t.Errorf("%q. Backoff.Get() got1 = %v, want %v", tt.name, got1, tt.want1)
		}
	}
}
