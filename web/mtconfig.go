package web

import (
	"encoding/json"
	"fmt"
	"mtui/bridge"
	"mtui/luanticonfig"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var runtime_set_allowed_types = map[string]bool{
	"string": true,
	"bool":   true,
	"int":    true,
	"float":  true,
	"enum":   true,
}

func (a *Api) GetMTConfig(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	sts, err := a.app.GetSettingTypes()
	if err != nil {
		Send(w, 500, err)
		return
	}

	s, err := a.app.ReadMTConfig(sts)
	Send(w, s, err)
}

func (a *Api) GetSettingTypes(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	sts, err := a.app.GetSettingTypes()
	Send(w, sts, err)
}

func (a *Api) SetMTConfig(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	key := vars["key"]

	s := &luanticonfig.Setting{}
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	_, err = a.setMTSetting(key, s, claims, r)
	Send(w, true, err)
}

func (a *Api) DeleteMTConfig(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	err := a.deleteMTSetting(vars["key"], claims, r)
	Send(w, true, err)
}

// setMTSetting writes the setting to the config file and applies it at runtime if possible, returns the setting type
func (a *Api) setMTSetting(key string, s *luanticonfig.Setting, claims *types.Claims, r *http.Request) (*luanticonfig.SettingType, error) {
	sts, err := a.app.GetSettingTypes()
	if err != nil {
		return nil, err
	}

	cfg, err := a.app.ReadMTConfig(sts)
	if err != nil {
		return nil, err
	}
	cfg[key] = s

	err = a.app.WriteMTConfig(cfg, sts)
	if err != nil {
		return nil, err
	}

	st := settingTypeOrDefault(sts, key)

	// create log entry
	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "settings",
		Message:  fmt.Sprintf("User '%s' changed the setting '%s' of type '%s' to '%s'", claims.Username, key, st.Type, s.Value),
	}, r)

	if !strings.HasPrefix(key, "secure.") && runtime_set_allowed_types[st.Type] {
		// set in engine, skipped for secure settings
		go a.applyRuntimeSetting(key, fmt.Sprintf("minetest.settings:set(%s, %s)", bridge.QuoteLuaString(key), bridge.QuoteLuaString(s.Value)))
	}
	return st, nil
}

// deleteMTSetting removes the setting from the config file and at runtime if possible
func (a *Api) deleteMTSetting(key string, claims *types.Claims, r *http.Request) error {
	sts, err := a.app.GetSettingTypes()
	if err != nil {
		return err
	}

	cfg, err := a.app.ReadMTConfig(sts)
	if err != nil {
		return err
	}
	delete(cfg, key)

	err = a.app.WriteMTConfig(cfg, sts)
	if err != nil {
		return err
	}

	st := settingTypeOrDefault(sts, key)

	// create log entry
	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "settings",
		Message:  fmt.Sprintf("User '%s' removed the setting '%s' of type '%s'", claims.Username, key, st.Type),
	}, r)

	if !strings.HasPrefix(key, "secure.") && runtime_set_allowed_types[st.Type] {
		// remove in engine, skipped for secure settings
		go a.applyRuntimeSetting(key, fmt.Sprintf("minetest.settings:remove(%s)", bridge.QuoteLuaString(key)))
	}
	return nil
}

func settingTypeOrDefault(sts luanticonfig.SettingTypes, key string) *luanticonfig.SettingType {
	st := sts[key]
	if st == nil {
		// default to string type
		st = &luanticonfig.SettingType{Type: "string"}
	}
	return st
}

func (a *Api) applyRuntimeSetting(key, lua string) {
	req := &command.LuaRequest{Code: lua}
	resp := &command.LuaResponse{}
	err := a.app.Bridge.ExecuteCommand(command.COMMAND_LUA, req, resp, time.Second*2)
	if err != nil {
		// just log error
		logrus.WithFields(logrus.Fields{
			"key": key,
			"err": err,
		}).Warn("could not apply runtime-setting")
	}
}
