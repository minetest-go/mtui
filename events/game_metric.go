package events

import (
	"encoding/json"
	"fmt"
	"mtui/bridge"
	"mtui/types/command"
)

func metricLoop(ch chan *bridge.CommandResponse) {
	for cmd := range ch {
		metrics := make([]*command.GameMetric, 0)
		err := json.Unmarshal(cmd.Data, &metrics)
		if err != nil {
			fmt.Printf("Payload error: %s\n", err.Error())
			return
		}

		fmt.Printf("Got %d metrics\n", len(metrics))
	}
}
