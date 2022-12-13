package command

import "mtui/bridge"

const (
	COMMAND_METRICS bridge.CommandType = "metrics"
)

type GameMetric struct {
	Name  string  `json:"name"`
	Type  string  `json:"type"`
	Help  string  `json:"help"`
	Value float64 `json:"value"`
}
