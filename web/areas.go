package web

import (
	"mtui/areasparser"
	"mtui/types"
	"net/http"
	"path"

	"github.com/gorilla/mux"
)

func (a *Api) GetAreas(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	areas, err := areasparser.ParseFile(path.Join(a.app.WorldDir, "areas.json"))
	Send(w, areas, err)
}

func (a *Api) GetOwnedAreas(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	playername := vars["playername"]
	areas, err := areasparser.ParseFile(path.Join(a.app.WorldDir, "areas.json"))
	owned_areas := make([]*areasparser.Area, 0)
	for _, a := range areas {
		if a.Owner == playername {
			owned_areas = append(owned_areas, a)
		}
	}
	Send(w, owned_areas, err)
}
