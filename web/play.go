package web

import (
	"mtui/public"
	"net/http"
)

func (a *Api) GetPlayPage(w http.ResponseWriter, r *http.Request) {
	data, err := public.Webapp.ReadFile("play.html")
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}
	w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
	w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
	w.Header().Add("Content-Type", "text/html")
	w.Write(data)
}
