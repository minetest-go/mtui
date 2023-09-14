package web

import (
	"encoding/base64"
	"fmt"
	"io"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

func getPlayerSkinFile(worlddir, playername string, skin_id int64) string {
	return path.Join(worlddir, "worldmods", "skinsdb", "textures", fmt.Sprintf("player_%s_%d.png", playername, skin_id))
}

func (a *Api) GetSkin(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	skin_id, _ := strconv.ParseInt(vars["id"], 10, 64)

	b, err := os.ReadFile(getPlayerSkinFile(a.app.WorldDir, claims.Username, skin_id))
	if err == nil && b != nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "image/png")
		w.Write(b)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (a *Api) SetSkin(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	skin_id, _ := strconv.ParseInt(vars["id"], 10, 64)

	if r.ContentLength > 1024*10 {
		SendError(w, 500, "file size > 10kb")
		return
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	err = os.WriteFile(getPlayerSkinFile(a.app.WorldDir, claims.Username, skin_id), b, 0666)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// update ingame skin with base64-encoded png
	png_str := base64.StdEncoding.EncodeToString(b)
	skin_name := fmt.Sprintf("player_%s_%d", claims.Username, skin_id)

	req := &command.LuaRequest{
		Code: `
			local skin_name = "` + skin_name + `"
			local skin = skins.get(skin_name)
			if skin then
				skin:set_texture("[png:` + png_str + `")
			end

			local player = minetest.get_player_by_name("` + claims.Username + `")
			if player then
				skins.update_player_skin(player)
			end
			return true
		`,
	}
	resp := &command.LuaResponse{}
	err = a.app.Bridge.ExecuteCommand(command.COMMAND_LUA, req, resp, time.Second*5)
	SendLuaResponse(w, err, resp)

	// create log entry
	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "skin",
		Message:  fmt.Sprintf("User '%s' uploaded a new skin in slot %d with %d bytes", claims.Username, skin_id, len(b)),
	}, r)
}

func (a *Api) RemoveSkin(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	skin_id, _ := strconv.ParseInt(vars["id"], 10, 64)

	err := os.Remove(getPlayerSkinFile(a.app.WorldDir, claims.Username, skin_id))
	Send(w, true, err)

	// create log entry
	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "skin",
		Message:  fmt.Sprintf("User '%s' removed the skin in slot %d", claims.Username, skin_id),
	}, r)
}
