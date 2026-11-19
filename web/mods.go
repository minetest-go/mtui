package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"mtui/luanticonfig/depanalyzer"
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

var errModNotFound = errors.New("not found")

func (a *Api) UpdateModVersion(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	m, err := a.updateModVersion(claims, vars["id"], vars["version"], r)
	if err == errModNotFound {
		SendError(w, 404, err)
		return
	}
	Send(w, m, err)
}

// updateModVersion updates the mod to the given version and notifies the engine
func (a *Api) updateModVersion(claims *types.Claims, id, version string, r *http.Request) (*types.Mod, error) {
	m, err := a.app.Repos.ModRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if m == nil {
		return nil, errModNotFound
	}

	err = a.app.ModManager.Update(m, version)
	if err != nil {
		return nil, err
	}

	// update locally
	m.Version = version

	// create log entry
	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "mods",
		Message:  fmt.Sprintf("User '%s' updated the %s '%s' (%s) to version '%s'", claims.Username, m.ModType, m.Name, m.SourceType, m.Version),
	}, r)

	go a.notifyModsChanged()
	return m, nil
}

func (a *Api) CreateMod(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	m := &types.Mod{}
	err := json.NewDecoder(r.Body).Decode(m)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	a.logModCreate(claims, m, r)
	Send(w, m, nil)

	go func() {
		err := a.createMod(m)
		if err != nil {
			logrus.WithError(err).Error("mod create")
		}
	}()
}

func (a *Api) logModCreate(claims *types.Claims, m *types.Mod, r *http.Request) {
	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "mods",
		Message:  fmt.Sprintf("User '%s' creates the %s '%s' (%s) in version '%s'", claims.Username, m.ModType, m.Name, m.SourceType, m.Version),
	}, r)
}

// createMod installs the mod (blocking) and notifies the engine
func (a *Api) createMod(m *types.Mod) error {
	err := a.app.ModManager.Create(m)
	if err != nil {
		return err
	}

	a.notifyModsChanged()
	return nil
}

func (a *Api) notifyModsChanged() {
	err := a.app.Bridge.ExecuteCommand(command.COMMAND_NOTIFY_MODS_CHANGED, nil, nil, time.Second*5)
	if err != nil {
		// ignore error, just log
		logrus.WithError(err).Warn("mods updated notification failed")
	}
}

// createPresetMod returns a handler that installs a preconfigured mod
func (a *Api) createPresetMod(create func() (*types.Mod, error)) func(http.ResponseWriter, *http.Request, *types.Claims) {
	return func(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
		m, err := create()
		if err == nil && m != nil {
			a.logModCreate(claims, m, r)
		}
		Send(w, m, err)
	}
}

func (a *Api) CreateMTUIMod(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	a.createPresetMod(a.app.CreateMTUIMod)(w, r, claims)
}

func (a *Api) CreateBeerchatMod(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	a.createPresetMod(a.app.CreateBeerchatMod)(w, r, claims)
}

func (a *Api) CreateMapserverMod(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	a.createPresetMod(a.app.CreateMapserverMod)(w, r, claims)
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

	// create log entry
	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "mods",
		Message:  fmt.Sprintf("User '%s' updates the metadata of  %s '%s' (%s)", claims.Username, m.ModType, m.Name, m.SourceType),
	}, r)

	Send(w, m, nil)

	go func() {
		err = a.app.Repos.ModRepo.Update(m)
		if err != nil {
			logrus.WithError(err).Error("mod update")
		}
	}()
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

	// create log entry
	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "mods",
		Message:  fmt.Sprintf("User '%s' deletes the %s '%s' (%s)", claims.Username, m.ModType, m.Name, m.SourceType),
	}, r)

	go func() {
		err = a.app.ModManager.Remove(m)
		if err != nil {
			logrus.WithError(err).Error("mod delete")
			return
		}

		// send notification to engine
		err = a.app.Bridge.ExecuteCommand(command.COMMAND_NOTIFY_MODS_CHANGED, nil, nil, time.Second*5)
		if err != nil {
			// ignore error, just log
			logrus.WithError(err).Warn("mods updated notification failed")
		}
	}()
}

func (a *Api) ModsCheckUpdates(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	Send(w, true, a.app.ModManager.CheckUpdates())
}

func (a *Api) ModsValidate(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	ad, err := depanalyzer.AnalyzeDeps(
		path.Join(a.app.WorldDir, "worldmods"),
		path.Join(a.app.WorldDir, "game/mods"),
	)
	Send(w, ad, err)
}
