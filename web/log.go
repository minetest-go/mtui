package web

import (
	"encoding/json"
	"mtui/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Api) SearchLogs(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	s := &types.LogSearch{}
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	list, err := a.app.Repos.LogRepository.Search(s)
	Send(w, list, err)
}

func (a *Api) CountLogs(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	s := &types.LogSearch{}
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	count, err := a.app.Repos.LogRepository.Count(s)
	Send(w, count, err)
}

func (a *Api) GetLogEvents(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	category := vars["category"]

	list, err := a.app.Repos.LogRepository.GetEvents(types.LogCategory(category))
	Send(w, list, err)
}
