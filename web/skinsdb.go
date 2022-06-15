package web

import (
	"io"
	"io/ioutil"
	"mtui/types"
	"net/http"
	"path"
)

func getPlayerSkinFile(worlddir, playername string) string {
	return path.Join(worlddir, "worldmods", "skinsdb", "textures", "player_"+playername+".png")
}

func (a *Api) GetSkin(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	b, err := ioutil.ReadFile(getPlayerSkinFile(a.app.WorldDir, claims.Username))
	if err == nil && b != nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Add("Content-Type", "image/png")
		w.Write(b)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (a *Api) SetSkin(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	if r.ContentLength > 1024*100 {
		SendError(w, 500, "file size > 100kb")
		return
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	err = ioutil.WriteFile(getPlayerSkinFile(a.app.WorldDir, claims.Username), b, 0755)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
}
