package web

import (
	"encoding/json"
	"fmt"
	"mtui/api/cdb"
	"mtui/types"
	"net/http"
	"sync"

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

// simple, aggressive package- and dependency-cache
var cdb_package_cache = map[string]*cdb.PackageDetails{}
var cdb_package_dependency_cache = map[string]cdb.PackageDependency{}
var cdb_lock = sync.RWMutex{}

func (a *Api) GetCDBPackage(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	author := vars["author"]
	name := vars["name"]
	key := fmt.Sprintf("%s/%s", author, name)

	cdb_lock.RLock()
	details := cdb_package_cache[key]
	cdb_lock.RUnlock()

	if details != nil {
		Send(w, details, nil)
		return
	}

	details, err := cdbcli.GetDetails(vars["author"], vars["name"])
	cdb_lock.Lock()
	cdb_package_cache[key] = details
	cdb_lock.Unlock()
	Send(w, details, err)
}

func (a *Api) GetCDBPackageDependencies(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	author := vars["author"]
	name := vars["name"]
	key := fmt.Sprintf("%s/%s", author, name)

	cdb_lock.RLock()
	deps := cdb_package_dependency_cache[key]
	cdb_lock.RUnlock()

	if deps != nil {
		Send(w, deps, nil)
		return
	}

	deps, err := cdbcli.GetDependencies(author, name)
	cdb_lock.Lock()
	cdb_package_dependency_cache[key] = deps
	cdb_lock.Unlock()
	Send(w, deps, err)
}
