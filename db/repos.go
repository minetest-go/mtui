package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type Repositories struct {
	Auth       *AuthRepository
	Privs      *UserPrivilegeRepository
	Map        *MapRepository
	ModStorage *ModStorageRepository
	Players    *PlayerRepository
}

func CreateRepositories(world_dir string) (*Repositories, error) {
	var err error
	repos := &Repositories{}

	auth_db, err := sql.Open("sqlite", world_dir+"/auth.sqlite")
	if err != nil {
		return nil, err
	}

	repos.Auth = NewAuthRepository(auth_db)
	repos.Privs = NewUserPrivilegeRepository(auth_db)

	repos.Map, err = NewMapRepository(world_dir + "/map.sqlite")
	if err != nil {
		return nil, err
	}

	repos.ModStorage, err = NewModStorageRepository(world_dir + "/mod_storage.sqlite")
	if err != nil {
		return nil, err
	}

	repos.Players, err = NewPlayerRepository(world_dir + "/players.sqlite")
	if err != nil {
		return nil, err
	}

	return repos, nil
}
