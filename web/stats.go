package web

import (
	"encoding/json"
	"fmt"
	"mtui/bridge"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"sync/atomic"
	"time"
)

var current_stats atomic.Pointer[command.StatsCommand]
var current_stats_updated atomic.Int64

func (a *Api) StatsEventListener(c chan *bridge.CommandResponse) {
	for {
		cmd := <-c
		sc := &command.StatsCommand{}
		err := json.Unmarshal(cmd.Data, sc)
		if err != nil {
			fmt.Printf("Stats-listener-error: %s\n", err.Error())
			continue
		}
		current_stats.Store(sc)
		current_stats_updated.Store(time.Now().Unix())
	}
}

type StatResponse struct {
	*command.StatsCommand
	Maintenance        bool `json:"maintenance"`
	FilebrowserEnabled bool `json:"filebrowser_enabled"`
}

func (a *Api) GetStats(w http.ResponseWriter, r *http.Request, claims *types.Claims) {

	sc := &StatResponse{
		StatsCommand:       &command.StatsCommand{},
		FilebrowserEnabled: a.app.Config.FilebrowserURL != "",
	}
	sc.Maintenance = a.app.MaintenanceMode.Load()

	last_updated := current_stats_updated.Load()
	seconds_ago := time.Now().Unix() - last_updated

	cs := current_stats.Load()
	if cs != nil && seconds_ago < 10 {
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
