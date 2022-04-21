package web

import (
	"mtadmin/db"
	"net/http"

	"github.com/gorilla/mux"
)

func Setup(repos *db.Repositories) error {
	r := mux.NewRouter()

	http.Handle("/", r)

	return nil
}
