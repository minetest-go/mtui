package web

import (
	"html/template"
	"mtui/public"
	"net/http"
)

type IndexModel struct {
	ServerName string
	Webdev     bool
}

func (a *Api) GetIndex(w http.ResponseWriter, r *http.Request) {
	data, err := public.Webapp.ReadFile("index.html")
	if err != nil {
		SendError(w, 500, err)
		return
	}

	t, err := template.New("").Parse(string(data))
	if err != nil {
		SendError(w, 500, err)
		return
	}

	m := &IndexModel{
		ServerName: a.app.Config.Servername,
		Webdev:     a.app.Config.Webdev,
	}

	w.Header().Set("Cross-Origin-Embedder-Policy", "credentialless")
	w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
	w.Header().Set("Cross-Origin-Resource-Policy", "cross-origin")
	t.Execute(w, m)
}
