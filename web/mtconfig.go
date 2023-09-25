package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"mtui/minetestconfig"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

func getSettingTypes(worlddir string) (minetestconfig.SettingTypes, error) {
	modst, err := minetestconfig.GetAllSettingTypes(path.Join(worlddir, "worldmods"))
	if err != nil {
		return nil, fmt.Errorf("could not get settingtypes for worldmods dir: %v", err)
	}

	gamest, err := minetestconfig.GetAllSettingTypes(path.Join(worlddir, "game/mods"))
	if err != nil {
		return nil, fmt.Errorf("could not get settingtypes for game/mods dir: %v", err)
	}

	serversettings, err := minetestconfig.GetServerSettingTypes()
	if err != nil {
		return nil, fmt.Errorf("could not get settingtypes: %v", err)
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

	return sts, nil
}

var mtconfig_mutex = sync.RWMutex{}

func readMTConfig(worlddir string, sts minetestconfig.SettingTypes) (minetestconfig.Settings, error) {
	mtconfig_mutex.RLock()
	defer mtconfig_mutex.RUnlock()

	mtconfig_file := os.Getenv("MINETEST_CONFIG")
	data, err := os.ReadFile(mtconfig_file)
	if err != nil {
		return nil, fmt.Errorf("error reading config from '%s': %v", mtconfig_file, err)
	}

	s := minetestconfig.Settings{}
	err = s.Read(bytes.NewReader(data), sts)
	return s, err
}

func writeMTConfig(cfg minetestconfig.Settings, sts minetestconfig.SettingTypes) error {
	mtconfig_mutex.Lock()
	defer mtconfig_mutex.Unlock()

	mtconfig_file := os.Getenv("MINETEST_CONFIG")
	f, err := os.OpenFile(mtconfig_file, os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		return fmt.Errorf("could not open minetest config file '%s': %v", mtconfig_file, err)
	}
	defer f.Close()

	err = cfg.Write(f, sts)
	if err != nil {
		return fmt.Errorf("could not write minetest config file '%s': %v", mtconfig_file, err)
	}

	return nil
}

var runtime_set_allowed_types = map[string]bool{
	"string": true,
	"bool":   true,
	"int":    true,
	"float":  true,
	"enum":   true,
}

func (a *Api) GetMTConfig(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	sts, err := getSettingTypes(a.app.WorldDir)
	if err != nil {
		Send(w, 500, err)
		return
	}

	s, err := readMTConfig(a.app.WorldDir, sts)
	Send(w, s, err)
}

func (a *Api) GetSettingTypes(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	sts, err := getSettingTypes(a.app.WorldDir)
	Send(w, sts, err)
}

func (a *Api) SetMTConfig(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	key := vars["key"]

	s := &minetestconfig.Setting{}
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		SendError(w, 500, err.Error())
		return
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
	cfg[key] = s

	err = writeMTConfig(cfg, sts)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	st := sts[key]
	if st == nil {
		// default to string type
		st = &minetestconfig.SettingType{Type: "string"}
	}

	// create log entry
	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "settings",
		Message:  fmt.Sprintf("User '%s' changed the setting '%s' of type '%s' to '%s'", claims.Username, key, st.Type, s.Value),
	}, r)

	if strings.HasPrefix(key, "secure.") {
		// skip runtime change for secure settings
		Send(w, true, nil)
		return
	}

	if runtime_set_allowed_types[st.Type] {
		go func() {
			// set in engine
			lua := fmt.Sprintf("minetest.settings:set(\"%s\", \"%s\")", key, s.Value)
			req := &command.LuaRequest{Code: lua}
			resp := &command.LuaResponse{}
			err = a.app.Bridge.ExecuteCommand(command.COMMAND_LUA, req, resp, time.Second*2)
			if err != nil {
				// just log error
				logrus.WithFields(logrus.Fields{
					"key": key,
					"err": err,
				}).Warn("could not apply runtime-setting")
			}
		}()
	}
	Send(w, true, nil)
}

func (a *Api) DeleteMTConfig(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	key := vars["key"]

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
	delete(cfg, key)

	err = writeMTConfig(cfg, sts)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	st := sts[key]
	if st == nil {
		st = &minetestconfig.SettingType{Type: "string"}
	}

	// create log entry
	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "settings",
		Message:  fmt.Sprintf("User '%s' removed the setting '%s' of type '%s'", claims.Username, key, st.Type),
	}, r)

	if strings.HasPrefix(key, "secure.") {
		// skip runtime change for secure settings
		Send(w, true, nil)
		return
	}

	if runtime_set_allowed_types[st.Type] {
		go func() {
			// remove in engine
			lua := fmt.Sprintf("minetest.settings:remove(\"%s\")", key)
			req := &command.LuaRequest{Code: lua}
			resp := &command.LuaResponse{}
			err = a.app.Bridge.ExecuteCommand(command.COMMAND_LUA, req, resp, time.Second*2)
			if err != nil {
				// just log error
				logrus.WithFields(logrus.Fields{
					"key": key,
					"err": err,
				}).Warn("could not remove runtime-setting")
			}
		}()
	}
	Send(w, true, nil)
}
