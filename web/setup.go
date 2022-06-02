package web

import (
	"fmt"
	"mtui/app"
	"mtui/public"
	"mtui/types"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

func Setup(a *app.App) error {
	r := mux.NewRouter()

	api := NewApi(a)
	r.HandleFunc("/api/login", api.DoLogout).Methods(http.MethodDelete)
	r.HandleFunc("/api/login", api.DoLogin).Methods(http.MethodPost)
	r.HandleFunc("/api/login", api.GetLogin).Methods(http.MethodGet)
	r.HandleFunc("/api/bridge/execute_chatcommand", Secure(api.ExecuteChatcommand)).Methods(http.MethodPost)
	r.HandleFunc("/api/bridge", CheckApiKey(os.Getenv("APIKEY"), a.Bridge.HandlePost)).Methods(http.MethodPost)
	r.HandleFunc("/api/bridge", CheckApiKey(os.Getenv("APIKEY"), a.Bridge.HandleGet)).Methods(http.MethodGet)

	c := a.Bridge.AddHandler()
	go func() {
		for {
			cmd := <-c
			payload, err := types.ParseCommand(cmd)
			if err != nil {
				fmt.Printf("Payload error: %s\n", err.Error())
				return
			}
			switch data := payload.(type) {
			case *types.StatsCommand:
				fmt.Printf("Stats: uptime=%f, max_lag=%f, tod=%f\n", data.Uptime, data.MaxLag, data.TimeOfDay)
			}
		}
	}()

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
