package web

import (
	"encoding/json"
	"fmt"
	"mtui/bridge"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"sync"
)

var stats_mutex = &sync.RWMutex{}
var current_stats *command.StatsCommand

func (a *Api) StatsEventListener(c chan *bridge.CommandResponse) {
	for {
		cmd := <-c
		sc := &command.StatsCommand{}
		err := json.Unmarshal(cmd.Data, sc)
		if err != nil {
			fmt.Printf("Tan-listener-error: %s\n", err.Error())
			continue
		}
		stats_mutex.Lock()
		current_stats = sc
		stats_mutex.Unlock()
	}
}

func (a *Api) GetStats(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	stats_mutex.RLock()
	defer stats_mutex.RUnlock()

	sc := &command.StatsCommand{}
	sc.MaxLag = current_stats.MaxLag
	sc.PlayerCount = current_stats.PlayerCount
	sc.TimeOfDay = current_stats.TimeOfDay
	sc.Players = make([]*command.PlayerStats, current_stats.PlayerCount)

	for i, ps := range current_stats.Players {
		// infos for all users
		p := &command.PlayerStats{
			Name:   ps.Name,
			HP:     ps.HP,
			Breath: ps.Breath,
		}

		if claims != nil && claims.HasPriv("ban") {
			// infos for provileged users
			p.Pos = ps.Pos
			p.Info = ps.Info
		}

		sc.Players[i] = p
	}

	Send(w, sc, nil)
}
