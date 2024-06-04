package repository

// GaugeMetric модель для маппинга метрики Gauge на БД
type GaugeMetric struct {
	ID    int64   `db:"id"`
	Name  string  `db:"name"`
	Value float64 `db:"gauge"`
}

// CounterMetric модель для маппинга метрики Counter на БД
type CounterMetric struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Value int64  `db:"counter"`
}

// StoreMetrics хранит массивы с GaugeMetric и CounterMetric
type StoreMetrics struct {
	Gauge   []GaugeMetric   `json:"gauges,omitempty"`
	Counter []CounterMetric `json:"counters,omitempty"`
}
