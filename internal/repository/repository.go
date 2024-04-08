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

//type Store interface {
//	GetValue(ctx context.Context, string, metricType string) (interface{}, error)
//	GetModelValue(ctx context.Context, metric *models.Metrics) error
//	SetValue(ctx context.Context, metricName string, metricType string, metricValue string) error
//	SetModelValue(ctx context.Context, metrics []*models.Metrics) error
//	GetAllValues(ctx context.Context) *StoreMetrics
//	Save() error
//	Restore() error
//}
