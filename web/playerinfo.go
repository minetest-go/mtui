package web

import (
	"encoding/json"
	"mtui/types"
	"net/http"
	"strconv"

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

	// privileged info
	Pitch float64 `json:"pitch"`
	Yaw   float64 `json:"yaw"`
	PosX  float64 `json:"posx"`
	PosY  float64 `json:"posy"`
	PosZ  float64 `json:"posz"`

	Stats map[string]float64 `json:"stats"`
}

var player_stat_fields = map[string]bool{
	"played_time":  true,
	"digged_nodes": true,
	"placed_nodes": true,
	"died":         true,
	"crafted":      true,
}

func mapPlayerInfo(auth *auth.AuthEntry, privs []*auth.PrivilegeEntry, player *player.Player, md []*player.PlayerMetadata, claims *types.Claims) *PlayerInfo {
	info := &PlayerInfo{
		Privs: make([]string, 0),
		Stats: make(map[string]float64),
	}

	if auth != nil {
		info.Name = auth.Name
		info.AuthEntry = true
		info.AuthID = *auth.ID
	}

	for _, priv := range privs {
		info.Privs = append(info.Privs, priv.Privilege)
	}

	for _, meta := range md {
		if player_stat_fields[meta.Metadata] {
			value, err := strconv.ParseFloat(meta.Value, 64)
			if err == nil {
				info.Stats[meta.Metadata] = value
			}
		}
	}

	if player != nil {
		info.PlayerEntry = true
		info.LastLogin = player.ModificationDate
		info.FirstLogin = player.CreationDate
		info.Breath = player.Breath
		info.HP = player.HP

		if claims.HasPriv("ban") {
			// fill in privileged infos
			info.Pitch = player.Pitch
			info.Yaw = player.Yaw
			info.PosX = player.PosX
			info.PosY = player.PosY
			info.PosZ = player.PosZ
		}
	}

	return info
}

func (a *Api) GetPlayerInfo(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
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

	md, err := a.app.DBContext.PlayerMetadata.GetPlayerMetadata(playername)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	info := mapPlayerInfo(auth, privs, player, md, claims)
	SendJson(w, info)
}

func (a *Api) SearchPlayer(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	s := &auth.AuthSearch{}
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	list, err := a.app.DBContext.Auth.Search(s)
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

		md, err := a.app.DBContext.PlayerMetadata.GetPlayerMetadata(auth.Name)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}

		result[i] = mapPlayerInfo(auth, privs, player, md, claims)
	}

	Send(w, result, err)
}

func (a *Api) CountPlayer(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	s := &auth.AuthSearch{}
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	count, err := a.app.DBContext.Auth.Count(s)
	Send(w, count, err)
}
