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

	// always on api
	r.HandleFunc("/api/export", api.SecurePriv(types.PRIV_SERVER, api.Export)).Methods(http.MethodGet)
	r.HandleFunc("/api/maintenance", api.SecurePriv(types.PRIV_SERVER, api.GetMaintenanceMode)).Methods(http.MethodGet)
	r.HandleFunc("/api/maintenance", api.SecurePriv(types.PRIV_SERVER, api.EnableMaintenanceMode)).Methods(http.MethodPut)
	r.HandleFunc("/api/maintenance", api.SecurePriv(types.PRIV_SERVER, api.DisableMaintenanceMode)).Methods(http.MethodDelete)
	r.HandleFunc("/api/stats", api.OptionalSecure(api.GetStats)).Methods(http.MethodGet)

	// maintenance mode middleware enabled routes below
	apir := r.PathPrefix("/api").Subrouter()
	apir.Use(MaintenanceModeCheck(a.MaintenanceMode))

	apir.HandleFunc("/appinfo", api.GetAppInfo)
	apir.HandleFunc("/themes", api.SecurePriv(types.PRIV_SERVER, api.GetThemes))

	apir.HandleFunc("/features", api.GetFeatures).Methods(http.MethodGet)
	apir.HandleFunc("/feature", api.SecurePriv(types.PRIV_SERVER, api.SetFeature)).Methods(http.MethodPost)

	apir.HandleFunc("/login", api.DoLogout).Methods(http.MethodDelete)
	apir.HandleFunc("/login", api.DoLogin).Methods(http.MethodPost)
	apir.HandleFunc("/login", api.GetLogin).Methods(http.MethodGet)

	apir.HandleFunc("/onboard", api.GetOnboardStatus).Methods(http.MethodGet)
	apir.HandleFunc("/onboard", api.CreateOnboardUser).Methods(http.MethodPost)

	apir.HandleFunc("/geoip/{ip}", api.SecurePriv(types.PRIV_BAN, api.QueryGeoip)).Methods(http.MethodGet)

	apir.HandleFunc("/changepw", api.Secure(api.ChangePassword)).Methods(http.MethodPost)

	apir.HandleFunc("/oauth_app", api.SecurePriv(types.PRIV_SERVER, api.GetOauthApps)).Methods(http.MethodGet)
	apir.HandleFunc("/oauth_app", api.SecurePriv(types.PRIV_SERVER, api.SetOauthApp)).Methods(http.MethodPost)
	apir.HandleFunc("/oauth_app/{id}", api.SecurePriv(types.PRIV_SERVER, api.GetOauthAppByID)).Methods(http.MethodGet)
	apir.HandleFunc("/oauth_app/{id}", api.SecurePriv(types.PRIV_SERVER, api.DeleteOauthApp)).Methods(http.MethodDelete)

	apir.HandleFunc("/player/info/{playername}", api.SecurePriv(types.PRIV_INTERACT, api.GetPlayerInfo)).Methods(http.MethodGet)
	apir.HandleFunc("/player/search", api.SecurePriv(types.PRIV_INTERACT, api.SearchPlayer)).Methods(http.MethodPost)
	apir.HandleFunc("/player/count", api.SecurePriv(types.PRIV_INTERACT, api.CountPlayer)).Methods(http.MethodPost)

	apir.HandleFunc("/log/search", api.SecurePriv(types.PRIV_BAN, api.SearchLogs)).Methods(http.MethodPost)
	apir.HandleFunc("/log/count", api.SecurePriv(types.PRIV_BAN, api.CountLogs)).Methods(http.MethodPost)
	apir.HandleFunc("/log/events/{category}", api.SecurePriv(types.PRIV_BAN, api.GetLogEvents)).Methods(http.MethodGet)

	apir.HandleFunc("/areas", api.Feature(types.PRIV_AREAS, api.Secure(api.GetAreas))).Methods(http.MethodGet)
	apir.HandleFunc("/areas/{playername}", api.Feature(types.PRIV_AREAS, api.Secure(api.GetOwnedAreas))).Methods(http.MethodGet)

	apir.HandleFunc("/skin/{id}", api.Feature(types.FEATURE_SKINSDB, api.Secure(api.GetSkin))).Methods(http.MethodGet)
	apir.HandleFunc("/skin/{id}", api.Feature(types.FEATURE_SKINSDB, api.Secure(api.SetSkin))).Methods(http.MethodPost)
	apir.HandleFunc("/skin/{id}", api.Feature(types.FEATURE_SKINSDB, api.Secure(api.RemoveSkin))).Methods(http.MethodDelete)

	apir.HandleFunc("/metric_types", api.Feature(types.FEATURE_MONITORING, api.GetMetricTypes)).Methods(http.MethodGet)
	apir.HandleFunc("/metric_types/{name}", api.Feature(types.FEATURE_MONITORING, api.GetMetricType)).Methods(http.MethodGet)
	apir.HandleFunc("/metrics/search", api.Feature(types.FEATURE_MONITORING, api.SearchMetrics)).Methods(http.MethodPost)
	apir.HandleFunc("/metrics/count", api.Feature(types.FEATURE_MONITORING, api.CountMetrics)).Methods(http.MethodPost)

	acfgr := apir.PathPrefix("/config").Subrouter()
	acfgr.Use(SecureHandler(api.PrivCheck(types.PRIV_SERVER)))
	acfgr.HandleFunc("/{key}", api.Secure(api.GetConfig)).Methods(http.MethodGet)
	acfgr.HandleFunc("/{key}", api.Secure(api.SetConfig)).Methods(http.MethodPost)

	mr := apir.PathPrefix("/mail").Subrouter()
	mr.Use(SecureHandler(api.PrivCheck(types.PRIV_INTERACT), api.FeatureCheck(types.FEATURE_MAIL)))
	mr.HandleFunc("/folder/inbox", api.Secure(api.GetInbox)).Methods(http.MethodGet)
	mr.HandleFunc("/folder/outbox", api.Secure(api.GetOutbox)).Methods(http.MethodGet)
	mr.HandleFunc("/folder/drafts", api.Secure(api.GetDrafts)).Methods(http.MethodGet)
	mr.HandleFunc("/contacts", api.Secure(api.GetContacts)).Methods(http.MethodGet)
	mr.HandleFunc("", api.Secure(api.SendMail)).Methods(http.MethodPost)
	mr.HandleFunc("/{id}", api.Secure(api.DeleteMail)).Methods(http.MethodDelete)
	mr.HandleFunc("/{id}/read", api.Secure(api.MarkRead)).Methods(http.MethodPost)
	mr.HandleFunc("/{id}/unread", api.Secure(api.MarkUnread)).Methods(http.MethodPost)
	mr.HandleFunc("/checkrecipient/{recipient}", api.Secure(api.CheckValidRecipient)).Methods(http.MethodGet)

	xbanr := apir.PathPrefix("/xban").Subrouter()
	xbanr.Use(SecureHandler(api.PrivCheck(types.PRIV_BAN), api.FeatureCheck(types.FEATURE_XBAN)))
	xbanr.HandleFunc("/status", api.GetBanDBStatus).Methods(http.MethodGet)
	xbanr.HandleFunc("/records/{playername}", api.GetBanRecord).Methods(http.MethodGet)
	xbanr.HandleFunc("/records", api.GetBannedRecords).Methods(http.MethodGet)
	xbanr.HandleFunc("/ban", api.BanPlayer).Methods(http.MethodPost)
	xbanr.HandleFunc("/tempban", api.TempBanPlayer).Methods(http.MethodPost)
	xbanr.HandleFunc("/unban", api.UnbanPlayer).Methods(http.MethodPost)
	xbanr.HandleFunc("/cleanup", api.CleanupBanDB).Methods(http.MethodPost)

	apir.HandleFunc("/bridge/execute_chatcommand", api.Feature(types.FEATURE_SHELL, api.Secure(api.ExecuteChatcommand))).Methods(http.MethodPost)
	apir.HandleFunc("/bridge/lua", api.Feature(types.FEATURE_LUASHELL, api.SecurePriv(types.PRIV_SERVER, api.ExecuteLua))).Methods(http.MethodPost)
	apir.HandleFunc("/bridge", api.CheckApiKey(a.Bridge.HandlePost)).Methods(http.MethodPost)
	apir.HandleFunc("/bridge", api.CheckApiKey(a.Bridge.HandleGet)).Methods(http.MethodGet)

	msr := apir.PathPrefix("/media").Subrouter()
	msr.Use(SecureHandler(api.FeatureCheck(types.FEATURE_MEDIASERVER)))
	msr.HandleFunc("/index.mth", a.Mediaserver.ServeHTTPIndex).Methods(http.MethodPost)
	msr.HandleFunc("/stats", api.GetMediaStats).Methods(http.MethodGet)
	msr.HandleFunc("/{hash}", a.Mediaserver.ServeHTTPFetch).Methods(http.MethodGet)
	msr.HandleFunc("/scan", api.SecurePriv(types.PRIV_SERVER, api.ScanMedia)).Methods(http.MethodPost)

	modsapi := apir.PathPrefix("/mods").Subrouter()
	modsapi.Use(SecureHandler(api.FeatureCheck(types.FEATURE_MODMANAGEMENT), api.PrivCheck("server")))
	modsapi.HandleFunc("", api.Secure(api.GetMods)).Methods(http.MethodGet)
	modsapi.HandleFunc("", api.Secure(api.CreateMod)).Methods(http.MethodPost)
	modsapi.HandleFunc("/validate", api.Secure(api.ModsValidate)).Methods(http.MethodGet)
	modsapi.HandleFunc("/{id}/update/{version}", api.Secure(api.UpdateModVersion)).Methods(http.MethodPost)
	modsapi.HandleFunc("/{id}", api.Secure(api.DeleteMod)).Methods(http.MethodDelete)
	modsapi.HandleFunc("/{id}/status", api.Secure(api.ModStatus)).Methods(http.MethodGet)

	cdbapi := apir.PathPrefix("/cdb").Subrouter()
	cdbapi.Use(SecureHandler(api.FeatureCheck(types.FEATURE_MODMANAGEMENT), api.PrivCheck("server")))
	cdbapi.HandleFunc("/search", api.Secure(api.SearchCDBPackages)).Methods(http.MethodPost)
	cdbapi.HandleFunc("/detail/{author}/{name}", api.Secure(api.GetCDBPackage)).Methods(http.MethodGet)
	cdbapi.HandleFunc("/detail/{author}/{name}/dependencies", api.Secure(api.GetCDBPackageDependencies)).Methods(http.MethodGet)

	cfgr := apir.PathPrefix("/mtconfig").Subrouter()
	cfgr.Use(SecureHandler(api.FeatureCheck(types.FEATURE_MINETEST_CONFIG)))
	cfgr.HandleFunc("/settingtypes", api.SecurePriv(types.PRIV_SERVER, api.GetSettingTypes))
	cfgr.HandleFunc("/settings", api.SecurePriv(types.PRIV_SERVER, api.GetMTConfig)).Methods(http.MethodGet)
	cfgr.HandleFunc("/settings/{key}", api.SecurePriv(types.PRIV_SERVER, api.SetMTConfig)).Methods(http.MethodPost)
	cfgr.HandleFunc("/settings/{key}", api.SecurePriv(types.PRIV_SERVER, api.DeleteMTConfig)).Methods(http.MethodDelete)

	servapi := apir.PathPrefix("/service").Subrouter()
	servapi.Use(SecureHandler(api.FeatureCheck(types.FEATURE_DOCKER)))
	servapi.HandleFunc("/engine/versions", api.SecurePriv(types.PRIV_SERVER, api.GetEngineVersions)).Methods(http.MethodGet)
	servapi.HandleFunc("/engine", api.SecurePriv(types.PRIV_SERVER, api.GetEngineStatus)).Methods(http.MethodGet)
	servapi.HandleFunc("/engine", api.SecurePriv(types.PRIV_SERVER, api.CreateEngine)).Methods(http.MethodPost)
	servapi.HandleFunc("/engine", api.SecurePriv(types.PRIV_SERVER, api.RemoveEngine)).Methods(http.MethodDelete)
	servapi.HandleFunc("/engine/start", api.SecurePriv(types.PRIV_SERVER, api.StartEngine)).Methods(http.MethodPost)
	servapi.HandleFunc("/engine/stop", api.SecurePriv(types.PRIV_SERVER, api.StopEngine)).Methods(http.MethodPost)
	servapi.HandleFunc("/engine/logs/{since}/{until}", api.SecurePriv(types.PRIV_SERVER, api.GetEngineLogs)).Methods(http.MethodGet)

	// OAuth
	api.app.OAuthServer.SetAllowGetAccessRequest(true)
	api.app.OAuthServer.SetAllowedGrantType(oauth2.Implicit, oauth2.AuthorizationCode, oauth2.ClientCredentials)
	api.app.OAuthServer.SetClientInfoHandler(server.ClientFormHandler)
	api.app.OAuthServer.UserAuthorizationHandler = api.OAuthAuthHandler
	api.app.OAuthServer.SetUserAuthorizationHandler(api.OauthUserAuthorizationHandler)

	apir.HandleFunc("/authorize", api.OAuthAuthorizeHandler)
	apir.HandleFunc("/token", api.OAuthTokenHandler)

	// index.html or /
	r.HandleFunc("/", api.GetIndex)
	r.HandleFunc("/index.html", api.GetIndex)

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

	return nil
}
