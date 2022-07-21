package events

import (
	"encoding/json"
	"fmt"
	"mtui/bridge"
	"mtui/eventbus"
	"mtui/types/command"
)

const (
	// server stats
	StatsEvent eventbus.EventType = "stats"
	// all player infos, including ip,rtt,etc
	PlayerStatsEvent eventbus.EventType = "player_stats"
	// just the playername and hp/breath info
	PlayerStatsEventLight eventbus.EventType = "player_stats_light"
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

		lightPlayerData := make([]*command.PlayerStats, len(stats.Players))
		for i, p := range stats.Players {
			lightPlayerData[i] = &command.PlayerStats{
				Name:   p.Name,
				HP:     p.HP,
				Breath: p.Breath,
			}
		}

		e.Emit(&eventbus.Event{
			Type: PlayerStatsEventLight,
			Data: lightPlayerData,
		})
	}
}
