package types

type MetricType struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Help string `json:"help"`
}

func (m *MetricType) Columns(action string) []string {
	return []string{"name", "type", "help"}
}

func (m *MetricType) Table() string {
	return "metric_type"
}

func (m *MetricType) Scan(action string, r func(dest ...any) error) error {
	return r(&m.Name, &m.Type, &m.Help)
}

func (m *MetricType) Values(action string) []any {
	return []any{m.Name, m.Type, m.Help}
}

type Metric struct {
	Timestamp int64   `json:"timestamp"`
	Name      string  `json:"name"`
	Value     float64 `json:"value"`
}
