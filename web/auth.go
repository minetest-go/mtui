package web

import (
	"mtui/types"
	"net/http"

	"github.com/gorilla/mux"
)

func (a *Api) GetAuth(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	playername := vars["playername"]

	if playername != claims.Username && !claims.HasPriv("ban") {
		// only players with the "ban" priv can look up other players auth
		SendError(w, 403, "unauthoized")
		return
	}

	auth, err := a.app.DBContext.Auth.GetByUsername(playername)
	Send(w, auth, err)
}
