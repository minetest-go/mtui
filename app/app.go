package app

import (
	"database/sql"
	"errors"
	"fmt"
	"mtui/bridge"
	"mtui/db"
	"mtui/eventbus"
	"mtui/mail"
	"mtui/mediaserver"
	"mtui/minetestconfig"
	"mtui/modmanager"
	"mtui/types"
	"os"
	"path"

	oautherrors "github.com/go-oauth2/oauth2/v4/errors"
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
	MTConfig      minetestconfig.MinetestConfig
}

const default_world_mt_content = `
mod_storage_backend = sqlite3
auth_backend = sqlite3
player_backend = sqlite3
backend = sqlite3
gameid = minetest_game
`

func Create(world_dir string) (*App, error) {

	// check world.mt file and fall back to defaults
	world_mt_file := path.Join(world_dir, "world.mt")
	_, err := os.Stat(world_mt_file)
	if errors.Is(err, os.ErrNotExist) {
		err = os.WriteFile(world_mt_file, []byte(default_world_mt_content), 0644)
		if err != nil {
			return nil, fmt.Errorf("could not write world.mt defaults to '%s'", world_mt_file)
		}
	}

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
	oauth_srv.SetInternalErrorHandler(func(err error) (re *oautherrors.Response) {
		logrus.WithFields(logrus.Fields{"err": re}).Error("Internal error")
		return
	})

	oauth_srv.SetResponseErrorHandler(func(re *oautherrors.Response) {
		logrus.WithFields(logrus.Fields{"err": re}).Error("Response error")
	})

	// Settings
	mtcfg := minetestconfig.MinetestConfig{}
	mtconfig_file := os.Getenv("MINETEST_CONFIG")
	if mtconfig_file != "" {
		f, err := os.Open(mtconfig_file)
		if err != nil {
			return nil, fmt.Errorf("could not open minetest config file '%s': %v", mtconfig_file, err)
		}

		err = mtcfg.Read(f)
		if err != nil {
			return nil, fmt.Errorf("could not parse minetest config file '%s': %v", mtconfig_file, err)
		}

		logrus.WithFields(logrus.Fields{
			"filename": mtconfig_file,
			"entries":  len(mtcfg),
		}).Info("Read minetest config")
	}

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
		Mail:          mail.New(dbctx),
		Config:        cfg,
		Mediaserver:   mediaserver.New(),
		GeoipResolver: NewGeoipResolver(world_dir),
		OAuthMgr:      oauth_mgr,
		OAuthServer:   oauth_srv,
		Version:       Version,
		MTConfig:      mtcfg,
	}

	return app, nil
}
