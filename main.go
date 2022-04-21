package main

import (
	"database/sql"
	"errors"
	"mtadmin/worldconfig"
	"os"

	_ "modernc.org/sqlite"
)

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
		auth_db_filename := world_dir + "/auth.sqlite"
		_, err = os.Stat(auth_db_filename)
		if errors.Is(err, os.ErrNotExist) {
			// TODO: migrate database and set journal-mode if it does not exist
			panic("auth db does not exist")
		}
		db_, err = sql.Open("sqlite", auth_db_filename)
		if err != nil {
			panic(err)
		}
	default:
		panic("unsupported backend: " + auth_backend)
	}

	_, err = db_.Exec("pragma journal_mode = wal;")
	if err != nil {
		panic(err)
	}
}
