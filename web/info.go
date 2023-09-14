package web

import (
	"net/http"
)

type AppInfo struct {
	Version    string `json:"version"`
	Servername string `json:"servername"`
}

func (a *Api) GetAppInfo(w http.ResponseWriter, r *http.Request) {
	ai := &AppInfo{
		Version:    a.app.Version,
		Servername: a.app.Config.Servername,
	}

	SendJson(w, ai)
}
