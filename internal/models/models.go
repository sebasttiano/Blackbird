// Package models хранит различные модельки.
package models

// Metrics модель для парсинга запросов связанных с gauge и counter метриками
type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

type MetricsProtobuf struct {
	ID    string   `json:"id"`                     // имя метрики
	MType string   `json:"type"`                   // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,string,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"`        // значение метрики в случае передачи gauge
}

type MetricSet struct {
	Set []*MetricsProtobuf `json:"metrics,omitempty"`
}

func (m *MetricSet) CastToMetrics() []*Metrics {

	metrics := make([]*Metrics, 0, len(m.Set))
	for _, metric := range m.Set {
		nm := Metrics(*metric)
		metrics = append(metrics, &nm)
	}
	return metrics
}
