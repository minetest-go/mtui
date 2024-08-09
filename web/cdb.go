package web

import (
	"encoding/json"
	"mtui/api/cdb"
	"mtui/types"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var cached_cdbcli = cdb.NewCachedClient(cdb.New(), time.Hour*6)

func (a *Api) SearchCDBPackages(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	q := &cdb.PackageQuery{}
	err := json.NewDecoder(r.Body).Decode(q)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	res, err := cached_cdbcli.SearchPackages(q)
	Send(w, res, err)
}

func (a *Api) GetCDBPackage(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	author := vars["author"]
	name := vars["name"]

	res, err := cached_cdbcli.GetDetails(author, name)
	Send(w, res, err)
}

func (a *Api) GetCDBPackageDependencies(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	author := vars["author"]
	name := vars["name"]

	res, err := cached_cdbcli.GetDependencies(author, name)
	Send(w, res, err)
}

type ResolveCDBPackageDepsRequest struct {
	Package          string   `json:"package"`
	InstalledMods    []string `json:"installed_mods"`
	SelectedPackages []string `json:"selected_packages"`
}

func (a *Api) ResolveCDBPackageDependencies(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	rr := &ResolveCDBPackageDepsRequest{}
	err := json.NewDecoder(r.Body).Decode(rr)
	if err != nil {
		SendError(w, 500, err)
		return
	}

	rd, err := cdb.ResolveDependencies(cached_cdbcli, rr.Package, rr.SelectedPackages, rr.InstalledMods)
	Send(w, rd, err)
}
