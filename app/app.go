package app

import (
	"mtui/bridge"
	"mtui/config"
	"mtui/db"
	"mtui/eventbus"
	"mtui/mail"
	"mtui/modmanager"
	"path"

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
	Config     *config.Config
}

func Create(world_dir string) (*App, error) {
	config_path := path.Join(world_dir, "mtui.json")
	cfg, err := config.Parse(config_path)
	if err != nil {
		//no config, create a default config
		cfg, err = config.WriteDefault(config_path)
		if err != nil {
			return nil, err
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

	app := &App{
		WorldDir:   world_dir,
		DBContext:  dbctx,
		ModManager: modmanager.New(world_dir),
		Repos:      repos,
		Bridge:     bridge.New(),
		WSEvents:   eventbus.NewEventBus(),
		Mail:       mail.New(world_dir),
		Config:     cfg,
	}

	return app, nil
}
