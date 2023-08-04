package web

import (
	"net/http"
	"os"
)

type AppInfo struct {
	Version    string `json:"version"`
	Servername string `json:"servername"`
}

func (a *Api) GetAppInfo(w http.ResponseWriter, r *http.Request) {
	ai := &AppInfo{
		Version:    a.app.Version,
		Servername: os.Getenv("SERVER_NAME"),
	}

	SendJson(w, ai)
}
