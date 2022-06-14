package web

import (
	"mtui/areasparser"
	"mtui/types"
	"net/http"
	"path"
)

func (a *Api) GetAreas(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	areas, err := areasparser.ParseFile(path.Join(a.app.WorldDir, "areas.json"))
	Send(w, areas, err)
}
