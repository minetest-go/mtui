package web

import (
	"bytes"
	"fmt"
	"io"
	"mtui/minetestconfig"
	"mtui/types"
	"mtui/types/command"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

func (a *Api) GetMTConfig(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	Send(w, a.app.MTConfig, nil)
}

func writeMTConfig(cfg minetestconfig.Settings) error {
	mtconfig_file := os.Getenv("MINETEST_CONFIG")
	f, err := os.OpenFile(mtconfig_file, os.O_RDWR, 0755)
	if err != nil {
		return fmt.Errorf("could not open minetest config file '%s': %v", mtconfig_file, err)
	}

	err = cfg.Write(f)
	if err != nil {
		return fmt.Errorf("could not write minetest config file '%s': %v", mtconfig_file, err)
	}

	return nil
}

func (a *Api) SetMTConfig(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	key := vars["key"]

	buf := bytes.NewBuffer([]byte{})
	_, err := io.Copy(buf, r.Body)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	value := buf.String()

	// update setting
	a.app.MTConfig[key] = value

	// write back to disk
	err = writeMTConfig(a.app.MTConfig)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if strings.HasPrefix(key, "secure.") {
		// skip runtime change for secure settings
		Send(w, true, nil)
		return
	}

	// set in engine
	lua := fmt.Sprintf("minetest.settings:set(\"%s\", \"%s\")", key, value)
	req := &command.LuaRequest{Code: lua}
	resp := &command.LuaResponse{}
	err = a.app.Bridge.ExecuteCommand(command.COMMAND_LUA, req, resp, time.Second*5)
	Send(w, resp, err)
}

func (a *Api) DeleteMTConfig(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	key := vars["key"]
	value := vars["value"]

	// update setting
	a.app.MTConfig[key] = value

	// write back to disk
	err := writeMTConfig(a.app.MTConfig)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	if strings.HasPrefix(key, "secure.") {
		// skip runtime change for secure settings
		Send(w, true, nil)
		return
	}

	// remove in engine
	lua := fmt.Sprintf("minetest.settings:remove(\"%s\")", key)
	req := &command.LuaRequest{Code: lua}
	resp := &command.LuaResponse{}
	err = a.app.Bridge.ExecuteCommand(command.COMMAND_LUA, req, resp, time.Second*5)
	Send(w, resp, err)
}
