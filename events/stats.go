package events

import (
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
		payload, err := command.ParseCommand(cmd)
		if err != nil {
			fmt.Printf("Payload error: %s\n", err.Error())
			return
		}
		switch data := payload.(type) {
		case *command.StatsCommand:
			e.Emit(&eventbus.Event{
				Type: StatsEvent,
				Data: map[string]float64{
					"max_lag":      data.MaxLag,
					"player_count": data.PlayerCount,
					"time_of_day":  data.TimeOfDay,
				},
			})

			e.Emit(&eventbus.Event{
				Type:         PlayerStatsEvent,
				Data:         data.Players,
				RequiredPriv: "ban",
			})
		}
	}
}
