package app

import (
	"database/sql"
	"mtui/bridge"
	"mtui/db"
	"mtui/eventbus"
	"mtui/mail"
	"mtui/mediaserver"
	"mtui/modmanager"
	"mtui/types"
	"os"

	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/golang-jwt/jwt/v4"
	"github.com/minetest-go/mtdb"
	"github.com/sirupsen/logrus"
)

var Version string

type App struct {
	DBContext     *mtdb.Context
	DB            *sql.DB
	WorldDir      string
	Repos         *db.Repositories
	ModManager    *modmanager.ModManager
	Bridge        *bridge.Bridge
	WSEvents      *eventbus.EventBus
	Mail          *mail.Mail
	Config        *types.Config
	Mediaserver   *mediaserver.MediaServer
	GeoipResolver *GeoipResolver
	Version       string
	OAuthMgr      *manage.Manager
	OAuthServer   *server.Server
}

func Create(world_dir string) (*App, error) {
	dbctx, err := mtdb.New(world_dir)
	if err != nil {
		return nil, err
	}

	db_, err := db.Init(world_dir)
	if err != nil {
		return nil, err
	}

	err = db.Migrate(db_)
	if err != nil {
		return nil, err
	}

	repos := db.NewRepositories(db_)

	// admin user
	admin_username := os.Getenv("ADMIN_USERNAME")
	admin_password := os.Getenv("ADMIN_PASSWORD")
	if admin_username != "" && admin_password != "" {
		logrus.WithFields(logrus.Fields{"admin_user": admin_username}).Info("Creating admin-user")
		err = CreateAdminUser(dbctx, admin_username, admin_password)
		if err != nil {
			return nil, err
		}
	}

	// features
	err = PopulateFeatures(repos.FeatureRepository)
	if err != nil {
		return nil, err
	}

	cfg := &types.Config{
		CookieDomain: os.Getenv("COOKIE_DOMAIN"),
		CookieSecure: os.Getenv("COOKIE_SECURE") == "true",
		CookiePath:   os.Getenv("COOKIE_PATH"),
		APIKey:       os.Getenv("API_KEY"),
	}

	// config defaults
	if cfg.CookieDomain == "" {
		cfg.CookieDomain = "127.0.0.1"
	}
	if cfg.CookiePath == "" {
		cfg.CookiePath = "/"
	}

	jwtKey, err := repos.ConfigRepo.GetByKey(types.ConfigJWTKey)
	if err != nil {
		return nil, err
	}
	if jwtKey == nil {
		// create key
		jwtKey = &types.ConfigEntry{
			Key:   types.ConfigJWTKey,
			Value: randSeq(16),
		}
		err = repos.ConfigRepo.Set(jwtKey)
		if err != nil {
			return nil, err
		}
	}
	cfg.JWTKey = jwtKey.Value

	if cfg.APIKey == "" {
		apiKey, err := repos.ConfigRepo.GetByKey(types.ConfigApiKey)
		if err != nil {
			return nil, err
		}
		if apiKey == nil {
			// create key
			apiKey = &types.ConfigEntry{
				Key:   types.ConfigApiKey,
				Value: randSeq(16),
			}
			err = repos.ConfigRepo.Set(apiKey)
			if err != nil {
				return nil, err
			}
		}
		cfg.APIKey = apiKey.Value
	}

	// oauth setup

	oauth_mgr := manage.NewDefaultManager()
	oauth_mgr.MustTokenStorage(store.NewMemoryTokenStore())
	oauth_mgr.MapClientStorage(&db.OAuthAppStore{Repo: repos.OauthAppRepo})
	oauth_mgr.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte(cfg.JWTKey), jwt.SigningMethodHS512))

	oauth_srv := server.NewDefaultServer(oauth_mgr)
	oauth_srv.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		logrus.WithFields(logrus.Fields{"err": re}).Error("Internal error")
		return
	})

	oauth_srv.SetResponseErrorHandler(func(re *errors.Response) {
		logrus.WithFields(logrus.Fields{"err": re}).Error("Response error")
	})

	if Version == "" {
		Version = "DEV"
	}

	app := &App{
		WorldDir:      world_dir,
		DBContext:     dbctx,
		DB:            db_,
		ModManager:    modmanager.New(world_dir, repos.ModRepo),
		Repos:         repos,
		Bridge:        bridge.New(),
		WSEvents:      eventbus.NewEventBus(),
		Mail:          mail.New(world_dir),
		Config:        cfg,
		Mediaserver:   mediaserver.New(),
		GeoipResolver: NewGeoipResolver(world_dir),
		OAuthMgr:      oauth_mgr,
		OAuthServer:   oauth_srv,
		Version:       Version,
	}

	return app, nil
}
