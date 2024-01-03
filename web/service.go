package web

import (
	"encoding/json"
	"fmt"
	"mtui/dockerservice"
	"mtui/types"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type ServiceApi struct {
	service     *dockerservice.DockerService
	imageMap    map[string]string
	api         *Api
	servicename string
}

func CreateServiceApi(service *dockerservice.DockerService, api *Api, r *mux.Router, servicename string, imageMap map[string]string) {
	if service == nil {
		return
	}

	sa := &ServiceApi{
		service:     service,
		imageMap:    imageMap,
		api:         api,
		servicename: servicename,
	}

	sr := r.PathPrefix(fmt.Sprintf("/%s", servicename)).Subrouter()
	sr.HandleFunc("/versions", api.SecurePriv(types.PRIV_SERVER, sa.GetVersions)).Methods(http.MethodGet)
	sr.HandleFunc("/stats", api.SecurePriv(types.PRIV_SERVER, sa.GetStats)).Methods(http.MethodGet)
	sr.HandleFunc("", sa.GetStatus).Methods(http.MethodGet)
	sr.HandleFunc("", api.SecurePriv(types.PRIV_SERVER, sa.Create)).Methods(http.MethodPost)
	sr.HandleFunc("", api.SecurePriv(types.PRIV_SERVER, sa.Remove)).Methods(http.MethodDelete)
	sr.HandleFunc("/start", api.SecurePriv(types.PRIV_SERVER, sa.Start)).Methods(http.MethodPost)
	sr.HandleFunc("/stop", api.SecurePriv(types.PRIV_SERVER, sa.Stop)).Methods(http.MethodPost)
	sr.HandleFunc("/logs/{since}/{until}", api.SecurePriv(types.PRIV_SERVER, sa.GetLogs)).Methods(http.MethodGet)
}

func (a *ServiceApi) GetVersions(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	SendJson(w, a.imageMap)
}

func (a *ServiceApi) GetStatus(w http.ResponseWriter, r *http.Request) {
	s, err := a.service.Status()
	Send(w, s, err)
}

func (a *ServiceApi) GetStats(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	s, err := a.service.Stats()
	Send(w, s, err)
}

type CreateRequest struct {
	Version string `json:"version"`
}

func (a *ServiceApi) Create(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	cer := &CreateRequest{}
	err := json.NewDecoder(r.Body).Decode(cer)
	if err != nil {
		SendError(w, 500, fmt.Sprintf("json error: %v", err))
		return
	}

	image := a.imageMap[cer.Version]
	if image == "" {
		SendError(w, 404, fmt.Sprintf("unknown version: %s", cer.Version))
		return
	}

	err = a.service.Create(image)
	Send(w, true, err)

	// create log entry
	a.api.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    a.servicename,
		Message:  fmt.Sprintf("User '%s' created the service with version '%s'", claims.Username, cer.Version),
	}, r)
}

func (a *ServiceApi) Start(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	err := a.service.Start()
	Send(w, true, err)

	// create log entry
	a.api.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    a.servicename,
		Message:  fmt.Sprintf("User '%s' started the service", claims.Username),
	}, r)
}

func (a *ServiceApi) Stop(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	// remove stats
	current_stats.Store(nil)

	err := a.service.Stop()
	Send(w, true, err)

	// create log entry
	a.api.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    a.servicename,
		Message:  fmt.Sprintf("User '%s' stopped the service", claims.Username),
	}, r)
}

func (a *ServiceApi) Remove(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	err := a.service.Remove()
	Send(w, true, err)

	// create log entry
	a.api.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    a.servicename,
		Message:  fmt.Sprintf("User '%s' removed the service", claims.Username),
	}, r)
}

func (a *ServiceApi) GetLogs(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
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

	slogs, err := a.service.GetLogs(since, until)
	Send(w, slogs, err)
}
