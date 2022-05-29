package web

import (
	"fmt"
	"mtadmin/app"
	"mtadmin/public"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

func Setup(a *app.App) error {
	r := mux.NewRouter()

	// rest api
	SetupBrowse(r, a)

	api := NewApi(a)
	r.HandleFunc("/api/login", api.Login).Methods(http.MethodPost)
	r.HandleFunc("/api/bridge", api.BridgeRx).Methods(http.MethodPost)
	r.HandleFunc("/api/bridge", api.BridgeTx).Methods(http.MethodGet)

	// static files
	if os.Getenv("WEBDEV") == "true" {
		fmt.Println("using live mode")
		fs := http.FileServer(http.FS(os.DirFS("public")))
		r.PathPrefix("/").HandlerFunc(fs.ServeHTTP)

	} else {
		fmt.Println("using embed mode")
		r.PathPrefix("/").Handler(statigz.FileServer(public.Webapp, brotli.AddEncoding))
	}

	http.Handle("/", r)
	return nil
}
