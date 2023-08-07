package web

import (
	"mtui/api/cdb"
	"mtui/types"
	"net/http"

	"github.com/gorilla/mux"
)

var cdbcli = cdb.New()

func (a *Api) GetCDBPackages(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	pkgtype := vars["type"]
	packages, err := cdbcli.SearchPackages(&cdb.PackageQuery{
		Type: cdb.PackageType(pkgtype),
	})
	Send(w, packages, err)
}
