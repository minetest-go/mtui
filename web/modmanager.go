package web

import (
	"encoding/json"
	"mtui/minetestconfig"
	"mtui/minetestconfig/depanalyzer"
	"mtui/types"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

func (a *Api) GetMods(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	mods, err := a.app.Repos.ModRepo.GetAll()
	Send(w, mods, err)
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

func (a *Api) ModsValidate(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	ad, err := depanalyzer.AnalyzeDeps(
		path.Join(a.app.WorldDir, "worldmods"),
		path.Join(a.app.WorldDir, "game/mods"),
	)
	Send(w, ad, err)
}

func (a *Api) GetSettingTypes(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	modst, err := minetestconfig.GetAllSettingTypes(path.Join(a.app.WorldDir, "worldmods"))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	gamest, err := minetestconfig.GetAllSettingTypes(path.Join(a.app.WorldDir, "game/mods"))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	serversettings, err := minetestconfig.GetServerSettingTypes()
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	sts := minetestconfig.SettingTypes{}
	for k, s := range modst {
		sts[k] = s
	}
	for k, s := range gamest {
		sts[k] = s
	}
	for k, s := range serversettings {
		sts[k] = s
	}

	Send(w, sts, nil)
}
