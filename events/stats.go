package events

import (
	"encoding/json"
	"fmt"
	"mtui/bridge"
	"mtui/eventbus"
	"mtui/types/command"
)

const (
	StatsEvent       eventbus.EventType = "stats"
	PlayerStatsEvent eventbus.EventType = "player_stats"
)

func statsLoop(e *eventbus.EventBus, ch chan *bridge.CommandResponse) {
	for {
		cmd := <-ch
		stats := &command.StatsCommand{}
		err := json.Unmarshal(cmd.Data, stats)
		if err != nil {
			fmt.Printf("Payload error: %s\n", err.Error())
			return
		}
		e.Emit(&eventbus.Event{
			Type: StatsEvent,
			Data: map[string]float64{
				"max_lag":      stats.MaxLag,
				"player_count": stats.PlayerCount,
				"time_of_day":  stats.TimeOfDay,
			},
		})

		e.Emit(&eventbus.Event{
			Type:         PlayerStatsEvent,
			Data:         stats.Players,
			RequiredPriv: "ban",
		})
	}
}
