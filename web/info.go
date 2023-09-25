package web

import (
	"net/http"
)

type AppInfo struct {
	Version        string `json:"version"`
	Servername     string `json:"servername"`
	InstallMtuiMod bool   `json:"install_mtui_mod"`
}

func (a *Api) GetAppInfo(w http.ResponseWriter, r *http.Request) {
	ai := &AppInfo{
		Version:        a.app.Version,
		Servername:     a.app.Config.Servername,
		InstallMtuiMod: a.app.Config.InstallMtuiMod,
	}

	SendJson(w, ai)
}
