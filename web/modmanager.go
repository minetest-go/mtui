package web

import (
	"encoding/json"
	"mtui/modmanager"
	"mtui/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Api) GetMods(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	SendJson(w, a.app.ModManager.Mods())
}

func (a *Api) ScanWorldDir(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	err := a.app.ModManager.Scan()
	Send(w, a.app.ModManager.Mods(), err)
}

func (a *Api) UpdateModVersion(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	name := vars["name"]
	version := vars["version"]

	var m *modmanager.Mod
	for _, lm := range a.app.ModManager.Mods() {
		if lm.Name == name {
			m = lm
		}
	}

	if m == nil {
		SendError(w, 404, "not found")
		return
	}

	err := a.app.ModManager.Update(m, version)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// update locally
	m.Version = version

	SendJson(w, m)
}

func (a *Api) CreateMod(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	m := &modmanager.Mod{}
	err := json.NewDecoder(r.Body).Decode(m)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	Send(w, m, a.app.ModManager.Create(m))
}

func (a *Api) DeleteMod(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	name := vars["name"]

	var m *modmanager.Mod
	for _, lm := range a.app.ModManager.Mods() {
		if lm.Name == name {
			m = lm
		}
	}

	if m == nil {
		SendError(w, 404, "not found")
		return
	}

	err := a.app.ModManager.Remove(m)
	if err != nil {
		SendError(w, 500, err.Error())
	}
}
