package web

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/minetest-go/mtdb/auth"
	"github.com/minetest-go/mtdb/player"
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

func mapPlayerInfo(auth *auth.AuthEntry, privs []*auth.PrivilegeEntry, player *player.Player) *PlayerInfo {
	info := &PlayerInfo{
		Privs: make([]string, 0),
	}

	if auth != nil {
		info.Name = auth.Name
		info.AuthEntry = true
		info.AuthID = *auth.ID
	}

	for _, priv := range privs {
		info.Privs = append(info.Privs, priv.Privilege)
	}

	if player != nil {
		info.PlayerEntry = true
		info.LastLogin = player.ModificationDate
		info.FirstLogin = player.CreationDate
		info.Breath = player.Breath
		info.HP = player.HP
	}

	return info
}

func (a *Api) GetPlayerInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playername := vars["playername"]

	auth, err := a.app.DBContext.Auth.GetByUsername(playername)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	privs, err := a.app.DBContext.Privs.GetByID(*auth.ID)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	player, err := a.app.DBContext.Player.GetPlayer(playername)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	info := mapPlayerInfo(auth, privs, player)
	SendJson(w, info)
}

func (a *Api) SearchPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	namelike := fmt.Sprintf("%%%s%%", vars["namelike"])

	list, err := a.app.DBContext.Auth.Search(&auth.AuthSearch{Usernamelike: &namelike})
	result := make([]*PlayerInfo, len(list))
	for i, auth := range list {
		privs, err := a.app.DBContext.Privs.GetByID(*auth.ID)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		player, err := a.app.DBContext.Player.GetPlayer(auth.Name)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		result[i] = mapPlayerInfo(auth, privs, player)
	}

	Send(w, result, err)
}
