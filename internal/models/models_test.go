package models

import (
	"reflect"
	"testing"
)

func TestMetricSet_CastToMetrics(t *testing.T) {
	type fields struct {
		Set []*MetricsProtobuf
	}
	tests := []struct {
		name   string
		fields fields
		want   []*Metrics
	}{
		{name: "success", fields: fields{Set: []*MetricsProtobuf{
			{ID: "GetSetZip144", MType: "counter"},
			{ID: "OtherSys", MType: "gauge"}}},
			want: []*Metrics{
				{ID: "GetSetZip144", MType: "counter"},
				{ID: "OtherSys", MType: "gauge"},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MetricSet{
				Set: tt.fields.Set,
			}
			if got := m.CastToMetrics(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CastToMetrics() = %v, want %v", got, tt.want)
			}
		})
	}
}
