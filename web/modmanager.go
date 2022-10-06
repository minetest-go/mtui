package web

import (
	"encoding/json"
	"mtui/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Api) GetMods(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	list, err := a.app.ModManager.Mods()
	Send(w, list, err)
}

func (a *Api) ScanWorldDir(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	err := a.app.ModManager.Scan()
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	a.GetMods(w, r, claims)
}

func (a *Api) UpdateModVersion(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	version := vars["version"]

	m, err := a.app.ModManager.Mod(id)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if m == nil {
		SendError(w, 404, "not found")
		return
	}

	err = a.app.ModManager.Update(m, version)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	// update locally
	m.Version = version

	SendJson(w, m)
}

func (a *Api) CreateMod(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	m := &types.Mod{}
	err := json.NewDecoder(r.Body).Decode(m)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	Send(w, m, a.app.ModManager.Create(m))
}

func (a *Api) DeleteMod(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]

	m, err := a.app.ModManager.Mod(id)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	if m == nil {
		SendError(w, 404, "not found")
		return
	}

	err = a.app.ModManager.Remove(m)
	if err != nil {
		SendError(w, 500, err.Error())
	}
}

func (a *Api) ModStatus(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]

	m, err := a.app.ModManager.Mod(id)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	status, err := a.app.ModManager.Status(m)
	Send(w, status, err)
}
