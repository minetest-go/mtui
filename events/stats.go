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
	// just the playername and hp/breath info
	PlayerStatsEvent eventbus.EventType = "player_stats"
	// all player infos, including ip,rtt,etc
	PlayerStatsEventExtra eventbus.EventType = "player_stats_extra"
)

func statsLoop(e *eventbus.EventBus, ch chan *bridge.CommandResponse) {
	for cmd := range ch {
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
			Cache: true,
		})

		lightPlayerData := make([]*command.PlayerStats, len(stats.Players))
		for i, p := range stats.Players {
			lightPlayerData[i] = &command.PlayerStats{
				Name:   p.Name,
				HP:     p.HP,
				Breath: p.Breath,
			}
		}

		// "light" data
		e.Emit(&eventbus.Event{
			Type:  PlayerStatsEvent,
			Data:  lightPlayerData,
			Cache: true,
		})

		// full infos about players
		e.Emit(&eventbus.Event{
			Type:         PlayerStatsEventExtra,
			Data:         stats.Players,
			RequiredPriv: "ban",
			Cache:        true,
		})
	}
}
