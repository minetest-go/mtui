package web

import (
	"encoding/json"
	"fmt"
	"mtui/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Api) GetOauthApps(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	list, err := a.app.Repos.OauthAppRepo.GetAll()
	Send(w, list, err)
}

func (a *Api) GetOauthAppByID(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	app, err := a.app.Repos.OauthAppRepo.GetByID(vars["id"])
	Send(w, app, err)
}

func (a *Api) SetOauthApp(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	app := &types.OauthApp{}
	err := json.NewDecoder(r.Body).Decode(app)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	err = a.app.Repos.OauthAppRepo.Set(app)
	Send(w, app, err)

	// create log entry
	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "oauth",
		Message:  fmt.Sprintf("User '%s' updates the oauth app '%s'", claims.Username, app.ID),
	}, r)
}

func (a *Api) DeleteOauthApp(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	err := a.app.Repos.OauthAppRepo.Delete(vars["id"])
	Send(w, true, err)

	// create log entry
	a.app.CreateUILogEntry(&types.Log{
		Username: claims.Username,
		Event:    "oauth",
		Message:  fmt.Sprintf("User '%s' deletes the oauth app '%s'", claims.Username, vars["id"]),
	}, r)
}
