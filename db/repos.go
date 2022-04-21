package db

import (
	_ "modernc.org/sqlite"
)

type Repositories struct {
	Auth       *AuthRepository
	Map        *MapRepository
	ModStorage *ModStorageRepository
	Players    *PlayerRepository
}

func CreateRepositories(world_dir string) (*Repositories, error) {
	var err error
	repos := &Repositories{}

	repos.Auth, err = NewAuthRepository(world_dir + "/auth.sqlite")
	if err != nil {
		panic(err)
	}

	repos.Map, err = NewMapRepository(world_dir + "/map.sqlite")
	if err != nil {
		panic(err)
	}

	repos.ModStorage, err = NewModStorageRepository(world_dir + "/mod_storage.sqlite")
	if err != nil {
		panic(err)
	}

	repos.Players, err = NewPlayerRepository(world_dir + "/players.sqlite")
	if err != nil {
		panic(err)
	}

	return repos, nil
}
