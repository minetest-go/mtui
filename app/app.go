package app

import "mtadmin/db"

type App struct {
	Repos    *db.Repositories
	WorldDir string
}

func Create(world_dir string) (*App, error) {
	repos, err := db.CreateRepositories(world_dir)
	if err != nil {
		return nil, err
	}

	app := &App{
		WorldDir: world_dir,
		Repos:    repos,
	}

	return app, nil
}
