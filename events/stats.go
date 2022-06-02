package events

import (
	"fmt"
	"mtui/bridge"
	"mtui/eventbus"
	"mtui/types"
)

const StatsEvent eventbus.EventType = "stats"

func statsLoop(e *eventbus.EventBus, ch chan *bridge.Command) {
	for {
		cmd := <-ch
		payload, err := types.ParseCommand(cmd)
		if err != nil {
			fmt.Printf("Payload error: %s\n", err.Error())
			return
		}
		switch data := payload.(type) {
		case *types.StatsCommand:
			e.Emit(&eventbus.Event{
				Type: StatsEvent,
				Data: data,
			})
		}
	}
}
