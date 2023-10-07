package web

import (
	"encoding/json"
	"fmt"
	"mtui/types"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// map versions with full image urls (in case the registry gets switched in the future)
var VersionImageMapping = map[string]string{
	"5.6.0": "registry.gitlab.com/minetest/minetest/server:5.6.0",
	"5.7.0": "registry.gitlab.com/minetest/minetest/server:5.7.0",
}

func (a *Api) GetEngineVersions(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	SendJson(w, VersionImageMapping)
}

func (a *Api) GetEngineStatus(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	s, err := a.app.ServiceEngine.Status()
	Send(w, s, err)
}

type CreateEngineRequest struct {
	Version string `json:"version"`
}

func (a *Api) CreateEngine(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	cer := &CreateEngineRequest{}
	err := json.NewDecoder(r.Body).Decode(cer)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("json error: %v", err))
		return
	}

	image := VersionImageMapping[cer.Version]
	if image == "" {
		SendError(w, 404, fmt.Sprintf("unknown version: %s", cer.Version))
		return
	}

	err = a.app.ServiceEngine.Create(image)
	Send(w, true, err)

	// create log entry
	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "engine",
		Message:  fmt.Sprintf("User '%s' created the minetest engine with version '%s'", claims.Username, cer.Version),
	}, r)
}

func (a *Api) StartEngine(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	err := a.app.ServiceEngine.Start()
	Send(w, true, err)

	// create log entry
	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "engine",
		Message:  fmt.Sprintf("User '%s' started the minetest engine", claims.Username),
	}, r)
}

func (a *Api) StopEngine(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	// remove stats
	current_stats.Store(nil)

	err := a.app.ServiceEngine.Stop()
	Send(w, true, err)

	// create log entry
	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "engine",
		Message:  fmt.Sprintf("User '%s' stopped the minetest engine", claims.Username),
	}, r)
}

func (a *Api) RemoveEngine(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	err := a.app.ServiceEngine.Remove()
	Send(w, true, err)

	// create log entry
	a.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "engine",
		Message:  fmt.Sprintf("User '%s' removed the minetest engine", claims.Username),
	}, r)
}

func (a *Api) GetEngineLogs(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)

	since_millis, err := strconv.ParseInt(vars["since"], 10, 64)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("invalid since format: %s", vars["since"]))
		return
	}
	since := time.UnixMilli(since_millis)

	until_millis, err := strconv.ParseInt(vars["until"], 10, 64)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("invalid until format: %s", vars["until"]))
		return
	}
	until := time.UnixMilli(until_millis)

	slogs, err := a.app.ServiceEngine.GetLogs(since, until)
	Send(w, slogs, err)
}
