package db

import (
	"database/sql"
	"mtadmin/worldconfig"
	"path"

	_ "github.com/lib/pq"
	"github.com/minetest-go/mtdb"
	_ "modernc.org/sqlite"
)

type Repositories struct {
	Auth  mtdb.AuthRepository
	Privs *mtdb.PrivRepository
}

func CreateRepositories(world_dir string) (*Repositories, error) {

	wc, err := worldconfig.Parse(path.Join(world_dir, "world.mt"))
	if err != nil {
		return nil, err
	}
	repos := &Repositories{}

	switch wc[worldconfig.CONFIG_AUTH_BACKEND] {
	case worldconfig.BACKEND_SQLITE3:
		auth_db, err := sql.Open("sqlite", path.Join(world_dir, "auth.sqlite"))
		if err != nil {
			return nil, err
		}
		mtdb.EnableWAL(auth_db)
		mtdb.MigrateAuthDB(auth_db, mtdb.DATABASE_SQLITE)

		repos.Auth = mtdb.NewAuthRepository(auth_db, mtdb.DATABASE_SQLITE)
		repos.Privs = mtdb.NewPrivilegeRepository(auth_db, mtdb.DATABASE_SQLITE)
	case worldconfig.BACKEND_POSTGRES:
		auth_db, err := sql.Open("postgres", wc[worldconfig.CONFIG_PSQL_AUTH_CONNECTION])
		if err != nil {
			return nil, err
		}
		repos.Auth = mtdb.NewAuthRepository(auth_db, mtdb.DATABASE_POSTGRES)
		repos.Privs = mtdb.NewPrivilegeRepository(auth_db, mtdb.DATABASE_POSTGRES)
	}

	return repos, nil
}
