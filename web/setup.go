package web

import (
	"fmt"
	"mtui/app"
	"mtui/public"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"
)

func Setup(a *app.App) error {
	r := mux.NewRouter()

	api := NewApi(a)

	err := api.Setup()
	if err != nil {
		return err
	}

	r.HandleFunc("/api/ws", api.Websocket)

	r.HandleFunc("/api/features", api.GetFeatures).Methods(http.MethodGet)

	r.HandleFunc("/api/login", api.DoLogout).Methods(http.MethodDelete)
	r.HandleFunc("/api/login", api.DoLogin).Methods(http.MethodPost)
	r.HandleFunc("/api/login", api.GetLogin).Methods(http.MethodGet)

	r.HandleFunc("/api/changepw", api.Secure(api.ChangePassword)).Methods(http.MethodPost)

	r.HandleFunc("/api/auth/{playername}", api.Secure(api.GetAuth)).Methods(http.MethodGet)

	r.HandleFunc("/api/playerinfo/{playername}", api.GetPlayerInfo).Methods(http.MethodGet)

	r.HandleFunc("/api/areas", api.Secure(api.GetAreas)).Methods(http.MethodGet)
	r.HandleFunc("/api/areas/{playername}", api.Secure(api.GetOwnedAreas)).Methods(http.MethodGet)

	r.HandleFunc("/api/skin", api.Secure(api.GetSkin)).Methods(http.MethodGet)
	r.HandleFunc("/api/skin", api.Secure(api.SetSkin)).Methods(http.MethodPost)

	r.HandleFunc("/api/mail/list", api.Secure(api.GetMails)).Methods(http.MethodGet)
	r.HandleFunc("/api/mail/{sender}/{time}", api.Secure(api.DeleteMail)).Methods(http.MethodDelete)
	r.HandleFunc("/api/mail/{sender}/{time}/read", api.Secure(api.MarkRead)).Methods(http.MethodPost)
	r.HandleFunc("/api/mail/checkrecipient/{recipient}", api.Secure(api.CheckValidRecipient)).Methods(http.MethodGet)
	r.HandleFunc("/api/mail/send/{recipient}", api.Secure(api.SendMail)).Methods(http.MethodPost)
	r.HandleFunc("/api/mail/contacts", api.Secure(api.GetContacts)).Methods(http.MethodGet)

	r.HandleFunc("/api/bridge/execute_chatcommand", api.Secure(api.ExecuteChatcommand)).Methods(http.MethodPost)
	r.HandleFunc("/api/bridge/lua", api.SecurePriv("server", api.ExecuteLua)).Methods(http.MethodPost)
	r.HandleFunc("/api/bridge", api.CheckApiKey(a.Bridge.HandlePost)).Methods(http.MethodPost)
	r.HandleFunc("/api/bridge", api.CheckApiKey(a.Bridge.HandleGet)).Methods(http.MethodGet)

	r.HandleFunc("/api/mods/scan", api.SecurePriv("server", api.ScanWorldDir)).Methods(http.MethodPost)
	r.HandleFunc("/api/mods", api.SecurePriv("server", api.GetMods)).Methods(http.MethodGet)
	r.HandleFunc("/api/mods/{id}/update/{version}", api.SecurePriv("server", api.UpdateModVersion)).Methods(http.MethodPost)
	r.HandleFunc("/api/mods", api.SecurePriv("server", api.CreateMod)).Methods(http.MethodPost)
	r.HandleFunc("/api/mods/{id}", api.SecurePriv("server", api.DeleteMod)).Methods(http.MethodDelete)
	r.HandleFunc("/api/mods/{id}/status", api.SecurePriv("server", api.ModStatus)).Methods(http.MethodGet)

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
