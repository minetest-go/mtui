package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"mtui/public"
	"mtui/types"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/luanti-golang/mtdb/auth"
	"github.com/luanti-golang/mtdb/player"
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

var errPlayerNotFound = errors.New("player not found")

// getPlayerInfo assembles the player info for a single player
func (a *Api) getPlayerInfo(playername string, claims *types.Claims) (*PlayerInfo, error) {
	auth, err := a.app.DBContext.Auth.GetByUsername(playername)
	if err != nil {
		return nil, err
	}
	if auth == nil {
		return nil, errPlayerNotFound
	}
	return a.mapAuthEntry(auth, claims)
}

// searchPlayers returns the player infos for all matching auth entries
func (a *Api) searchPlayers(s *auth.AuthSearch, claims *types.Claims) ([]*PlayerInfo, error) {
	list, err := a.app.DBContext.Auth.Search(s)
	if err != nil {
		return nil, err
	}
	result := make([]*PlayerInfo, len(list))
	for i, entry := range list {
		result[i], err = a.mapAuthEntry(entry, claims)
		if err != nil {
			return nil, err
		}
	}
	return result, nil
}

func (a *Api) mapAuthEntry(entry *auth.AuthEntry, claims *types.Claims) (*PlayerInfo, error) {
	privs, err := a.app.DBContext.Privs.GetByID(*entry.ID)
	if err != nil {
		return nil, err
	}

	player, err := a.app.DBContext.Player.GetPlayer(entry.Name)
	if err != nil {
		return nil, err
	}

	md, err := a.app.DBContext.PlayerMetadata.GetPlayerMetadata(entry.Name)
	if err != nil {
		return nil, err
	}

	return mapPlayerInfo(entry, privs, player, md, claims), nil
}

func (a *Api) GetPlayerInfo(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	info, err := a.getPlayerInfo(vars["playername"], claims)
	if err == errPlayerNotFound {
		SendError(w, 404, err)
		return
	}
	Send(w, info, err)
}

func (a *Api) SearchPlayer(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	s := &auth.AuthSearch{}
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	result, err := a.searchPlayers(s, claims)
	Send(w, result, err)
}

func (a *Api) CountPlayer(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	s := &auth.AuthSearch{}
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	count, err := a.app.DBContext.Auth.Count(s)
	Send(w, count, err)
}

func (a *Api) GetPlayerSkin(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playername := vars["playername"]

	// default skin
	skin, err := public.Webapp.ReadFile(public.Prefix + "/pics/character.png")
	if err != nil {
		SendError(w, 500, err)
		return
	}

	// check if custom skin is set
	entry, err := a.app.DBContext.ModStorage.Get("skinsdb", []byte(playername))
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if entry != nil && entry.Value != nil {
		skin_path := path.Join(a.app.WorldDir, "worldmods", "skinsdb", "textures", fmt.Sprintf("%s.png", entry.Value))
		fi, _ := os.Stat(skin_path)
		if fi != nil {
			// file exists, read contents
			skin, err = os.ReadFile(skin_path)
			if err != nil {
				SendError(w, 500, err)
				return
			}
		}
	}

	w.Header().Add("Content-Type", "image/png")
	w.Write(skin)
}
