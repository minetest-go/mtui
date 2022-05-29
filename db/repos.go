package db

import (
	"database/sql"

	"github.com/minetest-go/mtdb"
	_ "modernc.org/sqlite"
)

type Repositories struct {
	Auth  mtdb.AuthRepository
	Privs *mtdb.PrivRepository
}

func CreateRepositories(world_dir string) (*Repositories, error) {
	var err error
	repos := &Repositories{}

	auth_db, err := sql.Open("sqlite", world_dir+"/auth.sqlite")
	if err != nil {
		return nil, err
	}

	repos.Auth = mtdb.NewAuthRepository(auth_db, mtdb.DATABASE_SQLITE)
	repos.Privs = mtdb.NewPrivilegeRepository(auth_db, mtdb.DATABASE_SQLITE)

	return repos, nil
}
