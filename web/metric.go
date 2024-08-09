package web

import (
	"encoding/json"
	"mtui/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Api) SearchMetrics(w http.ResponseWriter, r *http.Request) {
	s := &types.MetricSearch{}
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	list, err := a.app.Repos.MetricRepository.Search(s)
	Send(w, list, err)
}

func (a *Api) CountMetrics(w http.ResponseWriter, r *http.Request) {
	s := &types.MetricSearch{}
	err := json.NewDecoder(r.Body).Decode(s)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	count, err := a.app.Repos.MetricRepository.Count(s)
	Send(w, count, err)
}

func (a *Api) GetMetricType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	list, err := a.app.Repos.MetricTypeRepository.GetByName(name)
	Send(w, list, err)
}

func (a *Api) GetMetricTypes(w http.ResponseWriter, r *http.Request) {
	list, err := a.app.Repos.MetricTypeRepository.GetAll()
	Send(w, list, err)
}
