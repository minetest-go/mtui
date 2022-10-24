package web

import "net/http"

func (a *Api) GetVersion(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(a.app.Version))
}
