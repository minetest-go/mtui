package types

type MetricType struct {
	Name string `json:"name" gorm:"primarykey;column:name"`
	Type string `json:"type" gorm:"column:type"`
	Help string `json:"help" gorm:"column:help"`
}

func (m *MetricType) TableName() string {
	return "metric_type"
}

type Metric struct {
	Timestamp int64   `json:"timestamp" gorm:"primarykey;column:timestamp"`
	Name      string  `json:"name" gorm:"primarykey;column:name"`
	Value     float64 `json:"value" gorm:"column:value"`
}

func (m *Metric) TableName() string {
	return "metric"
}

type MetricSearch struct {
	Name          *string `json:"name"`
	FromTimestamp *int64  `json:"from_timestamp"`
	ToTimestamp   *int64  `json:"to_timestamp"`
	Limit         *int    `json:"limit"`
}
