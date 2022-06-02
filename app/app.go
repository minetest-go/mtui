package app

import (
	"mtui/bridge"
	"mtui/db"
	"mtui/eventbus"

	"github.com/minetest-go/mtdb"
)

type App struct {
	DBContext *mtdb.Context
	WorldDir  string
	Repos     *db.Repositories
	Bridge    *bridge.Bridge
	WSEvents  *eventbus.EventBus
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

	app := &App{
		WorldDir:  world_dir,
		DBContext: dbctx,
		Repos:     repos,
		Bridge:    bridge.New(),
		WSEvents:  eventbus.NewEventBus(),
	}

	return app, nil
}
