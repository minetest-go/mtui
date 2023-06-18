package web

import (
	"mtui/types"
	"net/http"
)

func (a *Api) GetSettings(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	Send(w, a.app.Settings, nil)
}

//TODO: update, write back, set
