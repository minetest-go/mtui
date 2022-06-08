package app

import (
	"mtui/bridge"
	"mtui/db"
	"mtui/eventbus"
	"mtui/mail"

	"github.com/minetest-go/mtdb"
)

type App struct {
	DBContext *mtdb.Context
	WorldDir  string
	Repos     *db.Repositories
	Bridge    *bridge.Bridge
	WSEvents  *eventbus.EventBus
	Mail      *mail.Mail
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
		Mail:      mail.New(world_dir),
	}

	return app, nil
}
