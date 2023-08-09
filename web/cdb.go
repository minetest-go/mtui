package web

import (
	"encoding/json"
	"mtui/api/cdb"
	"mtui/types"
	"net/http"

	"github.com/gorilla/mux"
)

var cdbcli = cdb.New()

func (a *Api) SearchCDBPackages(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	q := &cdb.PackageQuery{}
	err := json.NewDecoder(r.Body).Decode(q)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	packages, err := cdbcli.SearchPackages(q)
	Send(w, packages, err)
}

func (a *Api) GetCDBPackage(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)

	details, err := cdbcli.GetDetails(vars["author"], vars["name"])
	Send(w, details, err)
}
