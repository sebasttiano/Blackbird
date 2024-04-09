package repository

type GaugeMetric struct {
	ID    int64   `db:"id"`
	Name  string  `db:"name"`
	Value float64 `db:"gauge"`
}

type CounterMetric struct {
	ID    int64  `db:"id"`
	Name  string `db:"name"`
	Value int64  `db:"counter"`
}

type StoreMetrics struct {
	Gauge   []GaugeMetric
	Counter []CounterMetric
}
