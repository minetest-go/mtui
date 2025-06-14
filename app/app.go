package app

import (
	"database/sql"
	"errors"
	"fmt"
	"mtui/bridge"
	"mtui/db"
	"mtui/dockerservice"
	"mtui/eventbus"
	"mtui/mail"
	"mtui/mediaserver"
	"mtui/modmanager"
	"mtui/types"
	"os"
	"path"
	"sync/atomic"

	cache "github.com/Code-Hex/go-generics-cache"
	"github.com/minetest-go/mtdb"
	"gorm.io/gorm"
)

var Version string

type App struct {
	DBContext           *mtdb.Context
	DB                  *sql.DB
	G                   *gorm.DB
	WorldDir            string
	Repos               *db.Repositories
	ModManager          *modmanager.ModManager
	Bridge              *bridge.Bridge
	WSEvents            *eventbus.EventBus
	Mail                *mail.Mail
	Config              *types.Config
	Mediaserver         *mediaserver.MediaServer
	GeoipResolver       GeoIPResolver
	Version             string
	maintenanceMode     *atomic.Bool // database detached, for backup and restores
	ServiceEngine       *dockerservice.DockerService
	ServiceMatterbridge *dockerservice.DockerService
	ServiceMapserver    *dockerservice.DockerService
	offline_xban_cache  *cache.Cache[string, *types.XBanEntry]
}

const default_world_mt_content = `
mod_storage_backend = sqlite3
auth_backend = sqlite3
player_backend = sqlite3
backend = sqlite3
gameid = game
`

// Creates a new application context
func Create(cfg *types.Config) (*App, error) {

	// check world.mt file and fall back to defaults
	world_mt_file := path.Join(cfg.WorldDir, "world.mt")
	_, err := os.Stat(world_mt_file)
	if errors.Is(err, os.ErrNotExist) {
		err = os.WriteFile(world_mt_file, []byte(default_world_mt_content), 0644)
		if err != nil {
			return nil, fmt.Errorf("could not write world.mt defaults to '%s'", world_mt_file)
		}
	}

	app := &App{
		WorldDir:           cfg.WorldDir,
		Bridge:             bridge.New(),
		WSEvents:           eventbus.NewEventBus(),
		Config:             cfg,
		Mediaserver:        mediaserver.New(),
		GeoipResolver:      NewGeoIPResolver(path.Join(cfg.WorldDir, "mmdb"), cfg.GeoIPAPI),
		Version:            Version,
		maintenanceMode:    &atomic.Bool{},
		offline_xban_cache: cache.New[string, *types.XBanEntry](),
	}

	if app.Version == "" {
		app.Version = "DEV"
	}

	err = app.AttachDatabase()
	if err != nil {
		return nil, fmt.Errorf("could not attach database: %v", err)
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
				Value: RandSeq(16),
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
				Value: RandSeq(16),
			}
			err = app.Repos.ConfigRepo.Set(apiKey)
			if err != nil {
				return nil, err
			}
		}
		cfg.APIKey = apiKey.Value
	}

	// docker services, if available
	app.SetupServices()

	// (re-)install mtui mod if specified
	if cfg.InstallMtuiMod {
		_, err := app.CreateMTUIMod()
		if err != nil {
			return nil, fmt.Errorf("could not install mtui mod: %v", err)
		}
	}

	if Version == "" {
		Version = "DEV"
	}

	return app, nil
}
