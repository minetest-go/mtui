package web

import (
	"encoding/json"
	"fmt"
	"mtui/minetestconfig/depanalyzer"
	"mtui/modmanager"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"path"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func (a *Api) GetMods(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	mods, err := a.app.Repos.ModRepo.GetAll()
	Send(w, mods, err)
}

func (a *Api) GetMod(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	mod, err := a.app.Repos.ModRepo.GetByID(vars["id"])
	Send(w, mod, err)
}

func (a *Api) UpdateModVersion(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]
	version := vars["version"]

	m, err := a.app.Repos.ModRepo.GetByID(id)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if m == nil {
		SendError(w, 404, fmt.Errorf("not found"))
		return
	}

	err = a.app.ModManager.Update(m, version)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	// update locally
	m.Version = version

	SendJson(w, m)

	// create log entry
	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "mods",
		Message:  fmt.Sprintf("User '%s' updated the %s '%s' (%s) to version '%s'", claims.Username, m.ModType, m.Name, m.SourceType, m.Version),
	}, r)

	// send notification to engine
	err = a.app.Bridge.ExecuteCommand(command.COMMAND_NOTIFY_MODS_CHANGED, nil, nil, time.Second*5)
	if err != nil {
		// ignore error, just log
		logrus.WithError(err).Warn("mods updated notification failed")
	}
}

func (a *Api) CreateMod(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	m := &types.Mod{}
	err := json.NewDecoder(r.Body).Decode(m)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	Send(w, m, a.app.ModManager.Create(m))

	// create log entry
	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "mods",
		Message:  fmt.Sprintf("User '%s' creates the %s '%s' (%s) in version '%s'", claims.Username, m.ModType, m.Name, m.SourceType, m.Version),
	}, r)

	// send notification to engine
	err = a.app.Bridge.ExecuteCommand(command.COMMAND_NOTIFY_MODS_CHANGED, nil, nil, time.Second*5)
	if err != nil {
		// ignore error, just log
		logrus.WithError(err).Warn("mods updated notification failed")
	}
}

func (a *Api) CreateMTUIMod(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	m, err := a.app.CreateMTUIMod()
	if err == nil && m != nil {
		// create log entry
		a.app.CreateUILogEntry(&types.Log{
			Username: claims.Username,
			Event:    "mods",
			Message:  fmt.Sprintf("User '%s' creates the %s '%s' (%s) in version '%s'", claims.Username, m.ModType, m.Name, m.SourceType, m.Version),
		}, r)
	}
	Send(w, m, nil)
}

func (a *Api) CreateBeerchatMod(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	m, err := a.app.CreateBeerchatMod()
	if err == nil && m != nil {
		// create log entry
		a.app.CreateUILogEntry(&types.Log{
			Username: claims.Username,
			Event:    "mods",
			Message:  fmt.Sprintf("User '%s' creates the %s '%s' (%s) in version '%s'", claims.Username, m.ModType, m.Name, m.SourceType, m.Version),
		}, r)
	}
	Send(w, m, nil)
}

func (a *Api) CreateMapserverMod(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	m, err := a.app.CreateMapserverMod()
	if err == nil && m != nil {
		// create log entry
		a.app.CreateUILogEntry(&types.Log{
			Username: claims.Username,
			Event:    "mods",
			Message:  fmt.Sprintf("User '%s' creates the %s '%s' (%s) in version '%s'", claims.Username, m.ModType, m.Name, m.SourceType, m.Version),
		}, r)
	}
	Send(w, m, nil)
}

func (a *Api) UpdateMod(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]

	m, err := a.app.Repos.ModRepo.GetByID(id)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if m == nil {
		SendError(w, 404, fmt.Errorf("not found"))
		return
	}

	err = json.NewDecoder(r.Body).Decode(m)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	Send(w, m, a.app.Repos.ModRepo.Update(m))

	// create log entry
	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "mods",
		Message:  fmt.Sprintf("User '%s' updates the metadata of  %s '%s' (%s)", claims.Username, m.ModType, m.Name, m.SourceType),
	}, r)
}

func (a *Api) DeleteMod(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	id := vars["id"]

	m, err := a.app.Repos.ModRepo.GetByID(id)
	if err != nil {
		SendError(w, 500, err)
		return
	}
	if m == nil {
		SendError(w, 404, fmt.Errorf("not found"))
		return
	}

	err = a.app.ModManager.Remove(m)
	if err != nil {
		SendError(w, 500, err)
	}

	// create log entry
	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "mods",
		Message:  fmt.Sprintf("User '%s' deletes the %s '%s' (%s)", claims.Username, m.ModType, m.Name, m.SourceType),
	}, r)

	// send notification to engine
	err = a.app.Bridge.ExecuteCommand(command.COMMAND_NOTIFY_MODS_CHANGED, nil, nil, time.Second*5)
	if err != nil {
		// ignore error, just log
		logrus.WithError(err).Warn("mods updated notification failed")
	}
}

func (a *Api) ModsCheckUpdates(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	mods, err := a.app.Repos.ModRepo.GetAll()
	if err != nil {
		SendError(w, 500, err)
		return
	}

	updated_mods, err := modmanager.CheckUpdates(a.app.WorldDir, mods)
	Send(w, updated_mods, err)
}

func (a *Api) ModsValidate(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	ad, err := depanalyzer.AnalyzeDeps(
		path.Join(a.app.WorldDir, "worldmods"),
		path.Join(a.app.WorldDir, "game/mods"),
	)
	Send(w, ad, err)
}
