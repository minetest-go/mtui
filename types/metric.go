package types

type MetricType struct {
	Name string `json:"name"`
	Type string `json:"type"`
	Help string `json:"help"`
}

type Metric struct {
	Timestamp int64   `json:"timestamp"`
	Name      string  `json:"name"`
	Value     float64 `json:"value"`
}
