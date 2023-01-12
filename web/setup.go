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
	r.Use(loggingMiddleware)

	api := NewApi(a)

	err := api.Setup()
	if err != nil {
		return err
	}

	r.HandleFunc("/api/version", api.GetVersion)

	r.HandleFunc("/api/features", api.GetFeatures).Methods(http.MethodGet)
	r.HandleFunc("/api/feature", api.SecurePriv("server", api.SetFeature)).Methods(http.MethodPost)

	r.HandleFunc("/api/login", api.DoLogout).Methods(http.MethodDelete)
	r.HandleFunc("/api/login", api.DoLogin).Methods(http.MethodPost)
	r.HandleFunc("/api/login", api.GetLogin).Methods(http.MethodGet)

	r.HandleFunc("/api/onboard", api.GetOnboardStatus).Methods(http.MethodGet)
	r.HandleFunc("/api/onboard", api.CreateOnboardUser).Methods(http.MethodPost)

	r.HandleFunc("/api/geoip/{ip}", api.SecurePriv("ban", api.QueryGeoip)).Methods(http.MethodGet)

	r.HandleFunc("/api/changepw", api.Secure(api.ChangePassword)).Methods(http.MethodPost)

	r.HandleFunc("/api/player/info/{playername}", api.SecurePriv("interact", api.GetPlayerInfo)).Methods(http.MethodGet)
	r.HandleFunc("/api/player/search", api.SecurePriv("interact", api.SearchPlayer)).Methods(http.MethodPost)
	r.HandleFunc("/api/player/count", api.SecurePriv("interact", api.CountPlayer)).Methods(http.MethodPost)

	r.HandleFunc("/api/log/search", api.SecurePriv("ban", api.SearchLogs)).Methods(http.MethodPost)
	r.HandleFunc("/api/log/count", api.SecurePriv("ban", api.CountLogs)).Methods(http.MethodPost)
	r.HandleFunc("/api/log/events/{category}", api.SecurePriv("ban", api.GetLogEvents)).Methods(http.MethodGet)

	r.HandleFunc("/api/areas", api.Feature("areas", api.Secure(api.GetAreas))).Methods(http.MethodGet)
	r.HandleFunc("/api/areas/{playername}", api.Feature("areas", api.Secure(api.GetOwnedAreas))).Methods(http.MethodGet)

	r.HandleFunc("/api/skin", api.Feature("skinsdb", api.Secure(api.GetSkin))).Methods(http.MethodGet)
	r.HandleFunc("/api/skin", api.Feature("skinsdb", api.Secure(api.SetSkin))).Methods(http.MethodPost)

	r.HandleFunc("/api/metric_types", api.Feature("monitoring", api.GetMetricTypes)).Methods(http.MethodGet)
	r.HandleFunc("/api/metric_types/{name}", api.Feature("monitoring", api.GetMetricType)).Methods(http.MethodGet)
	r.HandleFunc("/api/metrics/search", api.Feature("monitoring", api.SearchMetrics)).Methods(http.MethodPost)
	r.HandleFunc("/api/metrics/count", api.Feature("monitoring", api.CountMetrics)).Methods(http.MethodPost)

	r.HandleFunc("/api/mail/list", api.Feature("mail", api.Secure(api.GetMails))).Methods(http.MethodGet)
	r.HandleFunc("/api/mail/{sender}/{time}", api.Feature("mail", api.Secure(api.DeleteMail))).Methods(http.MethodDelete)
	r.HandleFunc("/api/mail/{sender}/{time}/read", api.Feature("mail", api.Secure(api.MarkRead))).Methods(http.MethodPost)
	r.HandleFunc("/api/mail/checkrecipient/{recipient}", api.Feature("mail", api.Secure(api.CheckValidRecipient))).Methods(http.MethodGet)
	r.HandleFunc("/api/mail/send/{recipient}", api.Feature("mail", api.Secure(api.SendMail))).Methods(http.MethodPost)
	r.HandleFunc("/api/mail/contacts", api.Feature("mail", api.Secure(api.GetContacts))).Methods(http.MethodGet)

	r.HandleFunc("/api/xban/status", api.Feature("xban", api.SecurePriv("ban", api.GetBanDBStatus))).Methods(http.MethodGet)
	r.HandleFunc("/api/xban/records/{playername}", api.Feature("xban", api.SecurePriv("ban", api.GetBanRecord))).Methods(http.MethodGet)
	r.HandleFunc("/api/xban/records", api.Feature("xban", api.SecurePriv("ban", api.GetBannedRecords))).Methods(http.MethodGet)
	r.HandleFunc("/api/xban/ban", api.Feature("xban", api.SecurePriv("ban", api.BanPlayer))).Methods(http.MethodPost)
	r.HandleFunc("/api/xban/tempban", api.Feature("xban", api.SecurePriv("ban", api.TempBanPlayer))).Methods(http.MethodPost)
	r.HandleFunc("/api/xban/unban", api.Feature("xban", api.SecurePriv("ban", api.UnbanPlayer))).Methods(http.MethodPost)
	r.HandleFunc("/api/xban/cleanup", api.Feature("xban", api.SecurePriv("server", api.CleanupBanDB))).Methods(http.MethodPost)

	r.HandleFunc("/api/bridge/execute_chatcommand", api.Feature("shell", api.Secure(api.ExecuteChatcommand))).Methods(http.MethodPost)
	r.HandleFunc("/api/bridge/lua", api.Feature("luashell", api.SecurePriv("server", api.ExecuteLua))).Methods(http.MethodPost)
	r.HandleFunc("/api/bridge", api.CheckApiKey(a.Bridge.HandlePost)).Methods(http.MethodPost)
	r.HandleFunc("/api/bridge", api.CheckApiKey(a.Bridge.HandleGet)).Methods(http.MethodGet)

	r.HandleFunc("/api/media/index.mth", api.Feature("mediaserver", a.Mediaserver.ServeHTTPIndex)).Methods(http.MethodPost)
	r.HandleFunc("/api/media/stats", api.Feature("mediaserver", api.GetMediaStats)).Methods(http.MethodGet)
	r.HandleFunc("/api/media/{hash}", a.Mediaserver.ServeHTTPFetch).Methods(http.MethodGet)
	r.HandleFunc("/api/media/scan", api.Feature("mediaserver", api.SecurePriv("server", api.ScanMedia))).Methods(http.MethodPost)

	r.HandleFunc("/api/mods/scan", api.Feature("modmanagement", api.SecurePriv("server", api.ScanWorldDir))).Methods(http.MethodPost)
	r.HandleFunc("/api/mods", api.Feature("modmanagement", api.SecurePriv("server", api.GetMods))).Methods(http.MethodGet)
	r.HandleFunc("/api/mods/{id}/update/{version}", api.Feature("modmanagement", api.SecurePriv("server", api.UpdateModVersion))).Methods(http.MethodPost)
	r.HandleFunc("/api/mods", api.Feature("modmanagement", api.SecurePriv("server", api.CreateMod))).Methods(http.MethodPost)
	r.HandleFunc("/api/mods/{id}", api.Feature("modmanagement", api.SecurePriv("server", api.DeleteMod))).Methods(http.MethodDelete)
	r.HandleFunc("/api/mods/{id}/status", api.Feature("modmanagement", api.SecurePriv("server", api.ModStatus))).Methods(http.MethodGet)

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
	http.HandleFunc("/api/ws", api.Websocket)

	return nil
}
