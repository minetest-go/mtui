package web

import (
	"encoding/json"
	"fmt"
	"mtui/bridge"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"sync/atomic"
)

var current_stats atomic.Pointer[command.StatsCommand]

func (a *Api) StatsEventListener(c chan *bridge.CommandResponse) {
	for {
		cmd := <-c
		sc := &command.StatsCommand{}
		err := json.Unmarshal(cmd.Data, sc)
		if err != nil {
			fmt.Printf("Tan-listener-error: %s\n", err.Error())
			continue
		}
		current_stats.Store(sc)
	}
}

type StatResponse struct {
	*command.StatsCommand
	Maintenance bool `json:"maintenance"`
}

func (a *Api) GetStats(w http.ResponseWriter, r *http.Request, claims *types.Claims) {

	sc := &StatResponse{
		StatsCommand: &command.StatsCommand{},
	}
	sc.Maintenance = a.app.MaintenanceMode.Load()

	cs := current_stats.Load()
	if cs != nil {
		sc.MaxLag = cs.MaxLag
		sc.PlayerCount = cs.PlayerCount
		sc.TimeOfDay = cs.TimeOfDay
		sc.Players = make([]*command.PlayerStats, cs.PlayerCount)

		for i, ps := range cs.Players {
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
	}

	Send(w, sc, nil)
}
