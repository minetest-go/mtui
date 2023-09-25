package web

import (
	"encoding/json"
	"fmt"
	"mtui/minetestconfig"
	"mtui/minetestconfig/depanalyzer"
	"mtui/types"
	"net/http"
	"path"
	"strings"

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

	// create log entry
	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "mods",
		Message:  fmt.Sprintf("User '%s' updated the %s '%s' (%s) to version '%s'", claims.Username, m.ModType, m.Name, m.SourceType, m.Version),
	}, r)
}

func (a *Api) CreateMod(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	m := &types.Mod{}
	err := json.NewDecoder(r.Body).Decode(m)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	Send(w, m, a.app.ModManager.Create(m))

	// create log entry
	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "mods",
		Message:  fmt.Sprintf("User '%s' creates the %s '%s' (%s) in version '%s'", claims.Username, m.ModType, m.Name, m.SourceType, m.Version),
	}, r)
}

func (a *Api) CreateMTUIMod(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	m := &types.Mod{
		Name:       "mtui",
		ModType:    types.ModTypeMod,
		SourceType: types.SourceTypeGIT,
		URL:        "https://github.com/minetest-go/mtui_mod.git",
		Branch:     "refs/heads/master",
	}
	err := a.app.ModManager.Create(m)
	if err != nil {
		Send(w, 500, err)
		return
	}

	// create log entry
	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "mods",
		Message:  fmt.Sprintf("User '%s' creates the %s '%s' (%s) in version '%s'", claims.Username, m.ModType, m.Name, m.SourceType, m.Version),
	}, r)

	// settings for mtui key/url
	if a.app.Config.DockerHostname == "" {
		return
	}

	for _, fname := range []string{types.FEATURE_DOCKER, types.FEATURE_MINETEST_CONFIG} {
		feature, err := a.app.Repos.FeatureRepository.GetByName(fname)
		if err != nil {
			SendError(w, 500, err.Error())
			return
		}
		if !feature.Enabled {
			return
		}
	}

	sts, err := getSettingTypes(a.app.WorldDir)
	if err != nil {
		Send(w, 500, err)
		return
	}

	cfg, err := readMTConfig(a.app.WorldDir, sts)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	cfg["mtui.url"] = &minetestconfig.Setting{
		Value: fmt.Sprintf("http://%s:8080", a.app.Config.DockerHostname),
	}
	cfg["mtui.key"] = &minetestconfig.Setting{
		Value: a.app.Config.APIKey,
	}

	http_mods := cfg["secure.http_mods"]
	if http_mods == nil || http_mods.Value == "" {
		// create new
		http_mods = &minetestconfig.Setting{
			Value: "mtui",
		}
		cfg["secure.http_mods"] = http_mods
	} else {
		// append if not in list
		is_in_http_list := false
		for _, mod := range strings.Split(http_mods.Value, ",") {
			if strings.TrimSpace(mod) == "mtui" {
				is_in_http_list = true
				break
			}
		}
		if !is_in_http_list {
			http_mods.Value += ",mtui"
		}
	}

	err = writeMTConfig(cfg, sts)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	Send(w, m, nil)
}

func (a *Api) UpdateMod(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
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

	err = json.NewDecoder(r.Body).Decode(m)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	Send(w, m, a.app.Repos.ModRepo.Update(m))

	// create log entry
	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "mods",
		Message:  fmt.Sprintf("User '%s' updates the metadata of  %s '%s' (%s)", claims.Username, m.ModType, m.Name, m.SourceType),
	}, r)
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

	// create log entry
	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "mods",
		Message:  fmt.Sprintf("User '%s' deletes the %s '%s' (%s)", claims.Username, m.ModType, m.Name, m.SourceType),
	}, r)
}

func (a *Api) ModsCheckUpdates(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	err := a.app.ModManager.CheckUpdates()
	Send(w, true, err)
}

func (a *Api) ModsValidate(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	ad, err := depanalyzer.AnalyzeDeps(
		path.Join(a.app.WorldDir, "worldmods"),
		path.Join(a.app.WorldDir, "game/mods"),
	)
	Send(w, ad, err)
}
