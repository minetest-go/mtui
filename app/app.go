package app

import (
	"fmt"
	"mtui/bridge"
	"mtui/db"
	"mtui/eventbus"
	"mtui/mail"
	"mtui/modmanager"
	"mtui/types"
	"os"

	"github.com/minetest-go/mtdb"
)

type App struct {
	DBContext  *mtdb.Context
	WorldDir   string
	Repos      *db.Repositories
	ModManager *modmanager.ModManager
	Bridge     *bridge.Bridge
	WSEvents   *eventbus.EventBus
	Mail       *mail.Mail
	Config     *types.Config
}

func Create(world_dir string) (*App, error) {
	dbctx, err := mtdb.New(world_dir)
	if err != nil {
		return nil, err
	}

	admin_username := os.Getenv("ADMIN_USERNAME")
	admin_password := os.Getenv("ADMIN_PASSWORD")
	if admin_username != "" && admin_password != "" {
		fmt.Printf("Creating admin user '%s'\n", admin_username)
		err = CreateAdminUser(dbctx, admin_username, admin_password)
		if err != nil {
			return nil, err
		}
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

	app := &App{
		WorldDir:   world_dir,
		DBContext:  dbctx,
		ModManager: modmanager.New(world_dir, repos.ModRepo),
		Repos:      repos,
		Bridge:     bridge.New(),
		WSEvents:   eventbus.NewEventBus(),
		Mail:       mail.New(world_dir),
		Config:     cfg,
	}

	return app, nil
}
