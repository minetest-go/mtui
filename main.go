package main

import (
	"database/sql"
	"mtadmin/worldconfig"
	"os"

	_ "modernc.org/sqlite"
)

const AUTH_MIGRATE = `
pragma journal_mode = wal;
CREATE TABLE if not exists auth (id INTEGER PRIMARY KEY AUTOINCREMENT,name VARCHAR(32) UNIQUE,password VARCHAR(512),last_login INTEGER);
CREATE TABLE if not exists user_privileges (id INTEGER,privilege VARCHAR(32),PRIMARY KEY (id, privilege)CONSTRAINT fk_id FOREIGN KEY (id) REFERENCES auth (id) ON DELETE CASCADE);
`

func main() {
	world_dir := os.Getenv("WORLD_DIR")

	cfg, err := worldconfig.Parse(world_dir + "/world.mt")
	if err != nil {
		panic(err)
	}

	var db_ *sql.DB
	auth_backend := cfg[worldconfig.CONFIG_AUTH_BACKEND]
	switch auth_backend {
	case worldconfig.BACKEND_SQLITE3:
		db_, err = sql.Open("sqlite", world_dir+"/auth.sqlite")
		if err != nil {
			panic(err)
		}
	default:
		panic("unsupported backend: " + auth_backend)
	}

	_, err = db_.Exec(AUTH_MIGRATE)
	if err != nil {
		panic(err)
	}

}
