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
	"mtui/modmanager"
	"mtui/types"
	"os"
	"path"
	"sync/atomic"

	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/minetest-go/mtdb"
	"github.com/sirupsen/logrus"
)

var Version string

type App struct {
	DBContext       *mtdb.Context
	DB              *sql.DB
	WorldDir        string
	Repos           *db.Repositories
	ModManager      *modmanager.ModManager
	Bridge          *bridge.Bridge
	WSEvents        *eventbus.EventBus
	Mail            *mail.Mail
	Config          *types.Config
	Mediaserver     *mediaserver.MediaServer
	GeoipResolver   *GeoipResolver
	Version         string
	OAuthMgr        *manage.Manager
	OAuthServer     *server.Server
	MaintenanceMode *atomic.Bool // database detached, for backup and restores
}

const default_world_mt_content = `
mod_storage_backend = sqlite3
auth_backend = sqlite3
player_backend = sqlite3
backend = sqlite3
gameid = game
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

	cfg := types.NewConfig()
	app := &App{
		WorldDir:        world_dir,
		Bridge:          bridge.New(),
		WSEvents:        eventbus.NewEventBus(),
		Config:          cfg,
		Mediaserver:     mediaserver.New(),
		GeoipResolver:   NewGeoipResolver(path.Join(world_dir, "mmdb")),
		Version:         Version,
		MaintenanceMode: &atomic.Bool{},
	}

	err = app.AttachDatabase()
	if err != nil {
		return nil, fmt.Errorf("could not attach database: %v", err)
	}

	// admin user
	if cfg.AdminUsername != "" && cfg.AdminPassword != "" {
		logrus.WithFields(logrus.Fields{"admin_user": cfg.AdminUsername}).Info("Creating admin-user")
		err = CreateAdminUser(app.DBContext, cfg.AdminUsername, cfg.AdminPassword)
		if err != nil {
			return nil, err
		}
	}

	// features
	err = PopulateFeatures(app.Repos.FeatureRepository, cfg.EnabledFeatures)
	if err != nil {
		return nil, err
	}

	// config defaults
	if cfg.CookieDomain == "" {
		cfg.CookieDomain = "127.0.0.1"
	}
	if cfg.CookiePath == "" {
		cfg.CookiePath = "/"
	}

	if cfg.JWTKey == "" {
		jwtKey, err := app.Repos.ConfigRepo.GetByKey(types.ConfigJWTKey)
		if err != nil {
			return nil, err
		}
		if jwtKey == nil {
			// create key
			jwtKey = &types.ConfigEntry{
				Key:   types.ConfigJWTKey,
				Value: randSeq(16),
			}
			err = app.Repos.ConfigRepo.Set(jwtKey)
			if err != nil {
				return nil, err
			}
		}
		cfg.JWTKey = jwtKey.Value
	}

	if cfg.APIKey == "" {
		apiKey, err := app.Repos.ConfigRepo.GetByKey(types.ConfigApiKey)
		if err != nil {
			return nil, err
		}
		if apiKey == nil {
			// create key
			apiKey = &types.ConfigEntry{
				Key:   types.ConfigApiKey,
				Value: randSeq(16),
			}
			err = app.Repos.ConfigRepo.Set(apiKey)
			if err != nil {
				return nil, err
			}
		}
		cfg.APIKey = apiKey.Value
	}

	if Version == "" {
		Version = "DEV"
	}

	return app, nil
}
