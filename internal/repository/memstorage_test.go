package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMemStorage(t *testing.T) {
	tests := []struct {
		name    string
		storage MemStorage
		want    MemStorage
	}{
		{
			name:    "Create New MemStorage",
			storage: *NewMemStorage(),
			want: MemStorage{
				Gauge:   make(map[string]float64),
				Counter: make(map[string]int64),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.storage)
		})
	}
	testsValuesGauge := []struct {
		name        string
		metricName  string
		metricValue float64
		want        float64
	}{
		{name: "Check gauge value #1", metricName: "gauge1", metricValue: 138.34, want: 138.34},
		{name: "Check gauge value #2", metricName: "gauge2", metricValue: -138.34, want: -138.34},
	}

	testsValuesCounter := []struct {
		name        string
		metricName  string
		metricValue int64
		want        int64
	}{
		{name: "Check counter value #1", metricName: "counter1", metricValue: 10, want: 10},
		{name: "Check counter value #2", metricName: "counter1", metricValue: 15, want: 25},
	}

	var localStorage = &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}

	for _, tt := range testsValuesGauge {
		t.Run(tt.name, func(t *testing.T) {
			err := localStorage.SetGauge(context.TODO(), &GaugeMetric{Name: tt.metricName, Value: tt.metricValue})
			require.NoErrorf(t, err, "error returned by SetGauge method. params: value %f", tt.metricValue)
			m := GaugeMetric{Name: tt.metricName}
			err = localStorage.GetGauge(context.TODO(), &m)
			require.NoErrorf(t, err, "error returned by GetGauge method. params: type %f", tt.metricValue)
			assert.Equal(t, fmt.Sprintf("%v", tt.want), fmt.Sprintf("%v", m.Value), "returned value not equal to expected")
		})
	}
	for _, tt := range testsValuesCounter {
		t.Run(tt.name, func(t *testing.T) {
			err := localStorage.SetCounter(context.TODO(), &CounterMetric{Name: tt.metricName, Value: tt.metricValue})
			require.NoErrorf(t, err, "error returned by SetCounter method. params: value %d", tt.metricValue)
			m := CounterMetric{Name: tt.metricName}
			err = localStorage.GetCounter(context.TODO(), &m)
			require.NoErrorf(t, err, "error returned by GetCounter method. params: type %d", tt.metricValue)
			assert.Equal(t, fmt.Sprintf("%v", tt.want), fmt.Sprintf("%v", m.Value), "returned value not equal to expected")
		})

	}
}
