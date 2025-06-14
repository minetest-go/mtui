package web

import (
	"mtui/app"
	"mtui/public"
	"mtui/types"
	"net/http"
	"os"
	"time"

	"github.com/dchest/captcha"
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
	r.HandleFunc("/api/maintenance", api.SecurePriv(types.PRIV_SERVER, api.GetMaintenanceMode)).Methods(http.MethodGet)
	r.HandleFunc("/api/maintenance", api.SecurePriv(types.PRIV_SERVER, api.EnableMaintenanceMode)).Methods(http.MethodPut)
	r.HandleFunc("/api/maintenance", api.SecurePriv(types.PRIV_SERVER, api.DisableMaintenanceMode)).Methods(http.MethodDelete)
	r.HandleFunc("/api/stats", api.OptionalSecure(api.GetStats)).Methods(http.MethodGet)
	r.HandleFunc("/api/login", api.GetLogin).Methods(http.MethodGet)

	// captcha store and http handler
	captcha.SetCustomStore(captcha.NewMemoryStore(50, 10*time.Minute))
	r.PathPrefix("/api/captcha/").Handler(captcha.Server(350, 250))

	fbr := r.PathPrefix("/api/filebrowser").Subrouter()
	fbr.Use(SecureHandler(api.PrivCheck(types.PRIV_SERVER)))
	fbr.HandleFunc("/size", api.Secure(api.GetDirectorySize)).Methods(http.MethodGet)
	fbr.HandleFunc("/browse", api.Secure(api.BrowseFolder)).Methods(http.MethodGet)
	fbr.HandleFunc("/zip", api.Secure(api.DownloadZip)).Methods(http.MethodGet)
	fbr.HandleFunc("/zip", api.Secure(api.UploadZip)).Methods(http.MethodPost)
	fbr.HandleFunc("/unzip", api.Secure(api.UnzipFile)).Methods(http.MethodPost)
	fbr.HandleFunc("/mkdir", api.Secure(api.Mkdir)).Methods(http.MethodPost)
	fbr.HandleFunc("/file", api.Secure(api.DownloadFile)).Methods(http.MethodGet)
	fbr.HandleFunc("/file", api.Secure(api.UploadFile)).Methods(http.MethodPost)
	fbr.HandleFunc("/file", api.Secure(api.DeleteFile)).Methods(http.MethodDelete)
	fbr.HandleFunc("/file", api.Secure(api.AppendFile)).Methods(http.MethodPut)
	fbr.HandleFunc("/rename", api.Secure(api.RenameFile)).Methods(http.MethodPost)

	// backup-restore job
	apibj := r.PathPrefix("/api/backup-restore").Subrouter()
	apibj.HandleFunc("", api.GetBackupRestoreJobInfo).Methods(http.MethodGet)
	apibj.HandleFunc("/create", api.SecurePriv(types.PRIV_SERVER, api.CreateBackupRestoreJob)).Methods(http.MethodPost)

	r.HandleFunc("/api/appinfo", api.GetAppInfo)
	r.HandleFunc("/api/themes", api.SecurePriv(types.PRIV_SERVER, api.GetThemes))

	// maintenance mode middleware enabled routes below
	apir := r.PathPrefix("/api").Subrouter()
	apir.Use(MaintenanceModeCheck(a.MaintenanceMode))

	apir.HandleFunc("/healthcheck", api.HealthCheck)

	apir.HandleFunc("/features", api.GetFeatures).Methods(http.MethodGet)
	apir.HandleFunc("/feature", api.SecurePriv(types.PRIV_SERVER, api.SetFeature)).Methods(http.MethodPost)

	apir.HandleFunc("/login", api.DoLogout).Methods(http.MethodDelete)
	apir.HandleFunc("/login", api.DoLogin).Methods(http.MethodPost)
	apir.HandleFunc("/token", api.Feature(types.FEATURE_API, api.SecurePriv(types.PRIV_INTERACT, api.CreateToken))).Methods(http.MethodPost)
	apir.HandleFunc("/loginadmin/{username}", api.AdminLogin).Methods(http.MethodGet)

	apir.HandleFunc("/signup", api.Feature(types.FEATURE_SIGNUP, api.Signup))
	apir.HandleFunc("/signup/captcha", api.Feature(types.FEATURE_SIGNUP, api.SignupCaptcha))

	apir.HandleFunc("/onboard", api.GetOnboardStatus).Methods(http.MethodGet)
	apir.HandleFunc("/onboard", api.CreateOnboardUser).Methods(http.MethodPost)

	apir.HandleFunc("/geoip/{ip}", api.SecurePriv(types.PRIV_BAN, api.QueryGeoip)).Methods(http.MethodGet)

	apir.HandleFunc("/changepw", api.Secure(api.ChangePassword)).Methods(http.MethodPost)

	apir.HandleFunc("/uimod/storage/{key}", api.SecurePriv(types.PRIV_SERVER, api.GetMtUIStorage)).Methods(http.MethodGet)
	apir.HandleFunc("/uimod/storage/{key}", api.SecurePriv(types.PRIV_SERVER, api.SetMtUIStorage)).Methods(http.MethodPost)
	apir.HandleFunc("/uimod/priv_info", api.SecurePriv(types.PRIV_INTERACT, api.GetMTUIPrivInfo)).Methods(http.MethodGet)
	apir.HandleFunc("/uimod/chatcommand_info", api.SecurePriv(types.PRIV_INTERACT, api.GetMTUIChatcommandInfo)).Methods(http.MethodGet)

	apir.HandleFunc("/player/skin/{playername}", api.GetPlayerSkin).Methods(http.MethodGet)
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

	chapi := apir.PathPrefix("/chat").Subrouter()
	chapi.Use(SecureHandler(api.FeatureCheck("chat"), api.PrivCheck("shout")))
	chapi.HandleFunc("/{channel}/latest", api.Secure(api.GetLatestChatLogs)).Methods(http.MethodGet)
	chapi.HandleFunc("/{channel}/{from}/{to}", api.Secure(api.GetChatLogs)).Methods(http.MethodGet)
	chapi.HandleFunc("", api.SecurePriv("shout", api.SendChat)).Methods(http.MethodPost)

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

	atmr := apir.PathPrefix("/atm").Subrouter()
	atmr.Use(SecureHandler(api.FeatureCheck(types.FEATURE_ATM)))
	atmr.HandleFunc("/balance/{name}", api.Secure(api.GetATMBalance)).Methods(http.MethodGet)
	atmr.HandleFunc("/transfer", api.Secure(api.ATMTransfer)).Methods(http.MethodPost)

	meser := apir.PathPrefix("/mesecons").Subrouter()
	meser.Use(SecureHandler(api.FeatureCheck(types.FEATURE_MESECONS)))
	meser.HandleFunc("", api.Secure(api.GetMeseconsControls)).Methods(http.MethodGet)
	meser.HandleFunc("", api.Secure(api.SetMeseconsControl)).Methods(http.MethodPost)
	meser.HandleFunc("/{poskey}", api.Secure(api.DeleteMeseconsControl)).Methods(http.MethodDelete)
	meser.HandleFunc("/luacontroller/get", api.Secure(api.GetLuacontroller)).Methods(http.MethodPost)
	meser.HandleFunc("/luacontroller/set", api.Secure(api.SetLuacontroller)).Methods(http.MethodPost)
	meser.HandleFunc("/luacontroller/digiline_send", api.Secure(api.LuacontrollerDigilineSend)).Methods(http.MethodPost)

	msr := apir.PathPrefix("/media").Subrouter()
	msr.Use(SecureHandler(api.FeatureCheck(types.FEATURE_MEDIASERVER)))
	msr.HandleFunc("/index.mth", a.Mediaserver.ServeHTTPIndex).Methods(http.MethodPost, http.MethodGet)
	msr.HandleFunc("/stats", api.GetMediaStats).Methods(http.MethodGet)
	msr.HandleFunc("/{hash}", a.Mediaserver.ServeHTTPFetch).Methods(http.MethodGet)
	msr.HandleFunc("/scan", api.SecurePriv(types.PRIV_SERVER, api.ScanMedia)).Methods(http.MethodPost)

	prr := apir.PathPrefix("/wasm").Subrouter()
	prr.Use(SecureHandler(api.FeatureCheck(types.FEATURE_MINETEST_WEB)))
	prr.HandleFunc("/proxy", api.Secure(api.HandleProxy))
	prr.HandleFunc("/joinpassword", api.Secure(api.RequestJoinPassword)).Methods(http.MethodGet)

	modsapi := apir.PathPrefix("/mods").Subrouter()
	modsapi.Use(SecureHandler(api.FeatureCheck(types.FEATURE_MODMANAGEMENT), api.PrivCheck("server")))
	modsapi.HandleFunc("", api.Secure(api.GetMods)).Methods(http.MethodGet)
	modsapi.HandleFunc("", api.Secure(api.CreateMod)).Methods(http.MethodPost)
	modsapi.HandleFunc("/validate", api.Secure(api.ModsValidate)).Methods(http.MethodGet)
	modsapi.HandleFunc("/checkupdates", api.Secure(api.ModsCheckUpdates)).Methods(http.MethodPost)
	modsapi.HandleFunc("/create_mtui", api.Secure(api.CreateMTUIMod)).Methods(http.MethodPost)
	modsapi.HandleFunc("/create_beerchat", api.Secure(api.CreateBeerchatMod)).Methods(http.MethodPost)
	modsapi.HandleFunc("/create_mapserver", api.Secure(api.CreateMapserverMod)).Methods(http.MethodPost)
	modsapi.HandleFunc("/{id}", api.Secure(api.GetMod)).Methods(http.MethodGet)
	modsapi.HandleFunc("/{id}", api.Secure(api.UpdateMod)).Methods(http.MethodPost)
	modsapi.HandleFunc("/{id}/update/{version}", api.Secure(api.UpdateModVersion)).Methods(http.MethodPost)
	modsapi.HandleFunc("/{id}", api.Secure(api.DeleteMod)).Methods(http.MethodDelete)

	cdbapi := apir.PathPrefix("/cdb").Subrouter()
	cdbapi.Use(SecureHandler(api.FeatureCheck(types.FEATURE_MODMANAGEMENT), api.PrivCheck("server")))
	cdbapi.HandleFunc("/search", api.Secure(api.SearchCDBPackages)).Methods(http.MethodPost)
	cdbapi.HandleFunc("/resolve", api.Secure(api.ResolveCDBPackageDependencies)).Methods(http.MethodPost)
	cdbapi.HandleFunc("/detail/{author}/{name}", api.Secure(api.GetCDBPackage)).Methods(http.MethodGet)
	cdbapi.HandleFunc("/detail/{author}/{name}/dependencies", api.Secure(api.GetCDBPackageDependencies)).Methods(http.MethodGet)

	cfgr := apir.PathPrefix("/mtconfig").Subrouter()
	cfgr.Use(SecureHandler(api.FeatureCheck(types.FEATURE_MINETEST_CONFIG)))
	cfgr.HandleFunc("/settingtypes", api.SecurePriv(types.PRIV_SERVER, api.GetSettingTypes))
	cfgr.HandleFunc("/settings", api.SecurePriv(types.PRIV_SERVER, api.GetMTConfig)).Methods(http.MethodGet)
	cfgr.HandleFunc("/settings/{key}", api.SecurePriv(types.PRIV_SERVER, api.SetMTConfig)).Methods(http.MethodPost)
	cfgr.HandleFunc("/settings/{key}", api.SecurePriv(types.PRIV_SERVER, api.DeleteMTConfig)).Methods(http.MethodDelete)

	if api.app.ServiceEngine != nil {
		servapi := apir.PathPrefix("/service").Subrouter()
		servapi.Use(SecureHandler(api.FeatureCheck(types.FEATURE_DOCKER)))

		CreateServiceApi(api.app.ServiceEngine, api, servapi, "engine", types.EngineServiceImages)
		CreateServiceApi(api.app.ServiceMatterbridge, api, servapi, "matterbridge", types.MatterbridgeServiceImages)
		CreateServiceApi(api.app.ServiceMapserver, api, servapi, "mapserver", types.MapserverServiceImages)

		if api.app.Config.DockerAutoInstallEngine {
			status, err := api.app.ServiceEngine.Status()
			if err == nil && status != nil && !status.Created {
				// silently install default engine
				go api.app.ServiceEngine.Create(types.EngineServiceImages[types.EngineServiceLatest])
			}
		}
	}

	// index.html or /
	r.HandleFunc("/", api.GetIndex)
	r.HandleFunc("/index.html", api.GetIndex)

	// static files
	var fsh http.Handler
	if a.Config.Webdev {
		logrus.WithFields(logrus.Fields{"dir": "public"}).Info("Using live mode")
		fs := http.FileServer(http.FS(os.DirFS("public")))
		fsh = fs
	} else {
		logrus.Info("Using embed mode")
		fsh = statigz.FileServer(public.Webapp, brotli.AddEncoding)
	}

	// set additional headers for wasm env
	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cross-Origin-Embedder-Policy", "credentialless")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Resource-Policy", "cross-origin")
		fsh.ServeHTTP(w, r)
	})

	http.Handle("/", r)

	return nil
}
