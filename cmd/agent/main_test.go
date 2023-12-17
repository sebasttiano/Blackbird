package main

import "testing"

func TestGetMetrics(t *testing.T) {
	type args struct {
		pollInterval   int
		reportInterval int64
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test pollInterval",
			args: args{pollInterval: 2, reportInterval: 10},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			GetMetrics(tt.args.pollInterval, tt.args.reportInterval)
		})
	}
}
