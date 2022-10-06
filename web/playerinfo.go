package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

type PlayerInfo struct {
	AuthEntry   bool  `json:"auth_entry"`
	AuthID      int64 `json:"auth_id"`
	PlayerEntry bool  `json:"player_entry"`

	Name       string   `json:"name"`
	Privs      []string `json:"privs"`
	LastLogin  int64    `json:"last_login"`
	FirstLogin int64    `json:"first_login"`
	Breath     int      `json:"breath"`
	HP         int      `json:"health"`
}

func (a *Api) GetPlayerInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playername := vars["playername"]

	info := &PlayerInfo{
		Name:  playername,
		Privs: make([]string, 0),
	}

	auth, err := a.app.DBContext.Auth.GetByUsername(playername)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if auth != nil {
		info.AuthEntry = true
		info.AuthID = *auth.ID
		privs, err := a.app.DBContext.Privs.GetByID(*auth.ID)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
		for _, priv := range privs {
			info.Privs = append(info.Privs, priv.Privilege)
		}
	}

	player, err := a.app.DBContext.Player.GetPlayer(playername)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if player != nil {
		info.PlayerEntry = true
		info.LastLogin = player.ModificationDate
		info.FirstLogin = player.CreationDate
		info.Breath = player.Breath
		info.HP = player.HP
	}

	SendJson(w, info)
}
