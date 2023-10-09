package web

import (
	"encoding/json"
	"fmt"
	"mtui/minetestconfig"
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

	s := &minetestconfig.Setting{}
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	sts, err := a.app.GetSettingTypes()
	if err != nil {
		Send(w, 500, err)
		return
	}

	cfg, err := a.app.ReadMTConfig(sts)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	cfg[key] = s

	err = a.app.WriteMTConfig(cfg, sts)
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
	a.app.CreateUILogEntry(&types.Log{
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

	sts, err := a.app.GetSettingTypes()
	if err != nil {
		Send(w, 500, err)
		return
	}

	cfg, err := a.app.ReadMTConfig(sts)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	delete(cfg, key)

	err = a.app.WriteMTConfig(cfg, sts)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	st := sts[key]
	if st == nil {
		st = &minetestconfig.SettingType{Type: "string"}
	}

	// create log entry
	a.app.CreateUILogEntry(&types.Log{
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
