package web

import (
	"encoding/json"
	"fmt"
	"mtui/api/cdb"
	"mtui/types"
	"net/http"
	"time"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/gorilla/mux"
)

var cdbcli = cdb.New()

// simple, aggressive query-, package- and dependency-cache
var cdb_search_cache = cache.New[string, []*cdb.Package]()
var cdb_package_cache = cache.New[string, *cdb.PackageDetails]()
var cdb_package_dependency_cache = cache.New[string, cdb.PackageDependency]()

func (a *Api) SearchCDBPackages(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	q := &cdb.PackageQuery{}
	err := json.NewDecoder(r.Body).Decode(q)
	if err != nil {
		SendError(w, 500, err.Error())
		return
	}

	key := q.Params().Encode()
	packages, ok := cdb_search_cache.Get(key)
	if !ok {
		packages, err = cdbcli.SearchPackages(q)
		cdb_search_cache.Set(key, packages, cache.WithExpiration(time.Hour*6))
	}
	Send(w, packages, err)
}

func (a *Api) GetCDBPackage(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	author := vars["author"]
	name := vars["name"]
	key := fmt.Sprintf("%s/%s", author, name)

	var err error
	details, ok := cdb_package_cache.Get(key)
	if !ok {
		details, err = cdbcli.GetDetails(vars["author"], vars["name"])
		cdb_package_cache.Set(key, details, cache.WithExpiration(6*time.Hour))
	}
	Send(w, details, err)
}

func (a *Api) GetCDBPackageDependencies(w http.ResponseWriter, r *http.Request, claims *types.Claims) {
	vars := mux.Vars(r)
	author := vars["author"]
	name := vars["name"]
	key := fmt.Sprintf("%s/%s", author, name)

	var err error
	deps, ok := cdb_package_dependency_cache.Get(key)
	if !ok {
		deps, err = cdbcli.GetDependencies(vars["author"], vars["name"])
		cdb_package_dependency_cache.Set(key, deps, cache.WithExpiration(6*time.Hour))
	}
	Send(w, deps, err)
}
