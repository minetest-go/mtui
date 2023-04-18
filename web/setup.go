package web

import (
	"mtui/app"
	"mtui/public"
	"mtui/types"
	"net/http"
	"os"

	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
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
	r.HandleFunc("/api/feature", api.SecurePriv(types.PRIV_SERVER, api.SetFeature)).Methods(http.MethodPost)

	r.HandleFunc("/api/login", api.DoLogout).Methods(http.MethodDelete)
	r.HandleFunc("/api/login", api.DoLogin).Methods(http.MethodPost)
	r.HandleFunc("/api/login", api.GetLogin).Methods(http.MethodGet)

	r.HandleFunc("/api/onboard", api.GetOnboardStatus).Methods(http.MethodGet)
	r.HandleFunc("/api/onboard", api.CreateOnboardUser).Methods(http.MethodPost)

	r.HandleFunc("/api/geoip/{ip}", api.SecurePriv(types.PRIV_BAN, api.QueryGeoip)).Methods(http.MethodGet)

	r.HandleFunc("/api/changepw", api.Secure(api.ChangePassword)).Methods(http.MethodPost)

	r.HandleFunc("/api/oauth_app", api.SecurePriv(types.PRIV_SERVER, api.GetOauthApps)).Methods(http.MethodGet)
	r.HandleFunc("/api/oauth_app", api.SecurePriv(types.PRIV_SERVER, api.SetOauthApp)).Methods(http.MethodPost)
	r.HandleFunc("/api/oauth_app/{id}", api.SecurePriv(types.PRIV_SERVER, api.GetOauthAppByID)).Methods(http.MethodGet)
	r.HandleFunc("/api/oauth_app/{id}", api.SecurePriv(types.PRIV_SERVER, api.DeleteOauthApp)).Methods(http.MethodDelete)

	r.HandleFunc("/api/player/info/{playername}", api.SecurePriv(types.PRIV_INTERACT, api.GetPlayerInfo)).Methods(http.MethodGet)
	r.HandleFunc("/api/player/search", api.SecurePriv(types.PRIV_INTERACT, api.SearchPlayer)).Methods(http.MethodPost)
	r.HandleFunc("/api/player/count", api.SecurePriv(types.PRIV_INTERACT, api.CountPlayer)).Methods(http.MethodPost)

	r.HandleFunc("/api/log/search", api.SecurePriv(types.PRIV_BAN, api.SearchLogs)).Methods(http.MethodPost)
	r.HandleFunc("/api/log/count", api.SecurePriv(types.PRIV_BAN, api.CountLogs)).Methods(http.MethodPost)
	r.HandleFunc("/api/log/events/{category}", api.SecurePriv(types.PRIV_BAN, api.GetLogEvents)).Methods(http.MethodGet)

	r.HandleFunc("/api/areas", api.Feature(types.PRIV_AREAS, api.Secure(api.GetAreas))).Methods(http.MethodGet)
	r.HandleFunc("/api/areas/{playername}", api.Feature(types.PRIV_AREAS, api.Secure(api.GetOwnedAreas))).Methods(http.MethodGet)

	r.HandleFunc("/api/skin", api.Feature(types.PRIV_SKINSDB, api.Secure(api.GetSkin))).Methods(http.MethodGet)
	r.HandleFunc("/api/skin", api.Feature(types.PRIV_SKINSDB, api.Secure(api.SetSkin))).Methods(http.MethodPost)

	r.HandleFunc("/api/metric_types", api.Feature(types.PRIV_MONITORING, api.GetMetricTypes)).Methods(http.MethodGet)
	r.HandleFunc("/api/metric_types/{name}", api.Feature(types.PRIV_MONITORING, api.GetMetricType)).Methods(http.MethodGet)
	r.HandleFunc("/api/metrics/search", api.Feature(types.PRIV_MONITORING, api.SearchMetrics)).Methods(http.MethodPost)
	r.HandleFunc("/api/metrics/count", api.Feature(types.PRIV_MONITORING, api.CountMetrics)).Methods(http.MethodPost)

	r.HandleFunc("/api/mail/folder/inbox", api.Feature(types.PRIV_MAIL, api.Secure(api.GetInbox))).Methods(http.MethodGet)
	r.HandleFunc("/api/mail/folder/outbox", api.Feature(types.PRIV_MAIL, api.Secure(api.GetOutbox))).Methods(http.MethodGet)
	r.HandleFunc("/api/mail/folder/drafts", api.Feature(types.PRIV_MAIL, api.Secure(api.GetDrafts))).Methods(http.MethodGet)
	r.HandleFunc("/api/mail", api.Feature(types.PRIV_MAIL, api.Secure(api.SendMail))).Methods(http.MethodPost)
	r.HandleFunc("/api/mail/{id}", api.Feature(types.PRIV_MAIL, api.Secure(api.DeleteMail))).Methods(http.MethodDelete)
	r.HandleFunc("/api/mail/{id}/read", api.Feature(types.PRIV_MAIL, api.Secure(api.MarkRead))).Methods(http.MethodPost)
	r.HandleFunc("/api/mail/{id}/unread", api.Feature(types.PRIV_MAIL, api.Secure(api.MarkUnread))).Methods(http.MethodPost)
	r.HandleFunc("/api/mail/checkrecipient/{recipient}", api.Feature(types.PRIV_MAIL, api.Secure(api.CheckValidRecipient))).Methods(http.MethodGet)

	r.HandleFunc("/api/xban/status", api.Feature(types.PRIV_XBAN, api.SecurePriv(types.PRIV_BAN, api.GetBanDBStatus))).Methods(http.MethodGet)
	r.HandleFunc("/api/xban/records/{playername}", api.Feature(types.PRIV_XBAN, api.SecurePriv(types.PRIV_BAN, api.GetBanRecord))).Methods(http.MethodGet)
	r.HandleFunc("/api/xban/records", api.Feature(types.PRIV_XBAN, api.SecurePriv(types.PRIV_BAN, api.GetBannedRecords))).Methods(http.MethodGet)
	r.HandleFunc("/api/xban/ban", api.Feature(types.PRIV_XBAN, api.SecurePriv(types.PRIV_BAN, api.BanPlayer))).Methods(http.MethodPost)
	r.HandleFunc("/api/xban/tempban", api.Feature(types.PRIV_XBAN, api.SecurePriv(types.PRIV_BAN, api.TempBanPlayer))).Methods(http.MethodPost)
	r.HandleFunc("/api/xban/unban", api.Feature(types.PRIV_XBAN, api.SecurePriv(types.PRIV_BAN, api.UnbanPlayer))).Methods(http.MethodPost)
	r.HandleFunc("/api/xban/cleanup", api.Feature(types.PRIV_XBAN, api.SecurePriv(types.PRIV_SERVER, api.CleanupBanDB))).Methods(http.MethodPost)

	r.HandleFunc("/api/bridge/execute_chatcommand", api.Feature(types.FEATURE_SHELL, api.Secure(api.ExecuteChatcommand))).Methods(http.MethodPost)
	r.HandleFunc("/api/bridge/lua", api.Feature(types.FEATURE_LUASHELL, api.SecurePriv(types.PRIV_SERVER, api.ExecuteLua))).Methods(http.MethodPost)
	r.HandleFunc("/api/bridge", api.CheckApiKey(a.Bridge.HandlePost)).Methods(http.MethodPost)
	r.HandleFunc("/api/bridge", api.CheckApiKey(a.Bridge.HandleGet)).Methods(http.MethodGet)

	r.HandleFunc("/api/media/index.mth", api.Feature(types.FEATURE_MEDIASERVER, a.Mediaserver.ServeHTTPIndex)).Methods(http.MethodPost)
	r.HandleFunc("/api/media/stats", api.Feature(types.FEATURE_MEDIASERVER, api.GetMediaStats)).Methods(http.MethodGet)
	r.HandleFunc("/api/media/{hash}", a.Mediaserver.ServeHTTPFetch).Methods(http.MethodGet)
	r.HandleFunc("/api/media/scan", api.Feature(types.FEATURE_MEDIASERVER, api.SecurePriv(types.PRIV_SERVER, api.ScanMedia))).Methods(http.MethodPost)

	r.HandleFunc("/api/mods/scan", api.Feature(types.FEATURE_MODMANAGEMENT, api.SecurePriv(types.PRIV_SERVER, api.ScanWorldDir))).Methods(http.MethodPost)
	r.HandleFunc("/api/mods", api.Feature(types.FEATURE_MODMANAGEMENT, api.SecurePriv(types.PRIV_SERVER, api.GetMods))).Methods(http.MethodGet)
	r.HandleFunc("/api/mods/{id}/update/{version}", api.Feature(types.FEATURE_MODMANAGEMENT, api.SecurePriv(types.PRIV_SERVER, api.UpdateModVersion))).Methods(http.MethodPost)
	r.HandleFunc("/api/mods", api.Feature(types.FEATURE_MODMANAGEMENT, api.SecurePriv(types.PRIV_SERVER, api.CreateMod))).Methods(http.MethodPost)
	r.HandleFunc("/api/mods/{id}", api.Feature(types.FEATURE_MODMANAGEMENT, api.SecurePriv(types.PRIV_SERVER, api.DeleteMod))).Methods(http.MethodDelete)
	r.HandleFunc("/api/mods/{id}/status", api.Feature(types.FEATURE_MODMANAGEMENT, api.SecurePriv(types.PRIV_SERVER, api.ModStatus))).Methods(http.MethodGet)

	// OAuth
	api.app.OAuthServer.SetAllowGetAccessRequest(true)
	api.app.OAuthServer.SetAllowedGrantType(oauth2.Implicit, oauth2.AuthorizationCode, oauth2.ClientCredentials)
	api.app.OAuthServer.SetClientInfoHandler(server.ClientFormHandler)
	api.app.OAuthServer.UserAuthorizationHandler = api.OAuthAuthHandler
	api.app.OAuthServer.SetUserAuthorizationHandler(api.OauthUserAuthorizationHandler)

	r.HandleFunc("/authorize", api.OAuthAuthorizeHandler)
	r.HandleFunc("/token", api.OAuthTokenHandler)

	// static files
	if os.Getenv("WEBDEV") == "true" {
		logrus.WithFields(logrus.Fields{"dir": "public"}).Info("Using live mode")
		fs := http.FileServer(http.FS(os.DirFS("public")))
		r.PathPrefix("/").HandlerFunc(fs.ServeHTTP)

	} else {
		logrus.Info("Using embed mode")
		r.PathPrefix("/").Handler(statigz.FileServer(public.Webapp, brotli.AddEncoding))
	}

	http.Handle("/", r)
	http.HandleFunc("/api/ws", api.Websocket)

	return nil
}
