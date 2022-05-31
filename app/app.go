package app

import (
	"github.com/minetest-go/mtdb"
)

type App struct {
	DBContext *mtdb.Context
	WorldDir  string
}

func Create(world_dir string) (*App, error) {
	dbctx, err := mtdb.New(world_dir)
	if err != nil {
		return nil, err
	}

	app := &App{
		WorldDir:  world_dir,
		DBContext: dbctx,
	}

	return app, nil
}
