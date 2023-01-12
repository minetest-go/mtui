package web

import (
	"mtui/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Api) QueryGeoip(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	res := a.app.GeoipResolver.Resolve(vars["ip"])
	SendJson(w, res)
}
