package web

import (
	"html/template"
	"mtui/public"
	"mtui/types"
	"net/http"
	"os"
)

type IndexModel struct {
	BootstrapCSSUrl string
	ServerName      string
}

func (a *Api) GetIndex(w http.ResponseWriter, r *http.Request) {
	data, err := public.Webapp.ReadFile("index.html")
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	t, err := template.New("").Parse(string(data))
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	entry, err := a.app.Repos.ConfigRepo.GetByKey(types.ConfigThemeKey)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	css_url := public.ThemeMap["default"]
	if entry != nil && public.ThemeMap[entry.Value] != "" {
		css_url = public.ThemeMap[entry.Value]
	}

	m := &IndexModel{
		BootstrapCSSUrl: css_url,
		ServerName:      os.Getenv("SERVER_NAME"),
	}
	t.Execute(w, m)
}
