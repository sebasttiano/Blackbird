package storage

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewMemStorage(t *testing.T) {
	tests := []struct {
		name    string
		storage MemStorage
		want    MemStorage
	}{
		{
			name:    "Create New MemStorage",
			storage: NewMemStorage(),
			want: MemStorage{
				Gauge:   make(map[string]float64),
				Counter: map[string]int64{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, tt.storage)
		})
	}
	testsValues := []struct {
		name        string
		metricType  string
		metricValue string
		want        string
	}{
		{name: "Check gauge value #1", metricType: "gauge", metricValue: "138.34", want: "138.34"},
		{name: "Check gauge value #2", metricType: "gauge", metricValue: "-138.34", want: "-138.34"},
		{name: "Check counter value #1", metricType: "counter", metricValue: "10", want: "10"},
		{name: "Check counter value #2", metricType: "counter", metricValue: "15", want: "25"},
	}

	var localStorage HandleMemStorage = &MemStorage{
		Gauge:   make(map[string]float64),
		Counter: make(map[string]int64),
	}

	for _, tt := range testsValues {
		t.Run(tt.name, func(t *testing.T) {
			err := localStorage.SetValue("TestMetric", tt.metricType, tt.metricValue)
			require.NoErrorf(t, err, "error returned by SetValue method. params: type %s, value %s", tt.metricType, tt.metricValue)
			valueBack, err := localStorage.GetValue("TestMetric", tt.metricType)
			require.NoErrorf(t, err, "error returned by GetValue method. params: type %s", tt.metricValue)
			assert.Equal(t, fmt.Sprintf("%v", tt.want), fmt.Sprintf("%v", valueBack), "returned value not equal to expected")
		})
	}
}
