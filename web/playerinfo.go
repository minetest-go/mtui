package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

type PlayerInfo struct {
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

	auth, err := a.app.DBContext.Auth.GetByUsername(playername)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if auth == nil {
		SendError(w, 404, "player not found")
		return
	}

	player, err := a.app.DBContext.Player.GetPlayer(playername)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	privs, err := a.app.DBContext.Privs.GetByID(*auth.ID)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	info := &PlayerInfo{
		Name:       playername,
		LastLogin:  player.ModificationDate,
		FirstLogin: player.CreationDate,
		Privs:      make([]string, 0),
		Breath:     player.Breath,
		HP:         player.HP,
	}

	for _, priv := range privs {
		info.Privs = append(info.Privs, priv.Privilege)
	}

	SendJson(w, info)
}
